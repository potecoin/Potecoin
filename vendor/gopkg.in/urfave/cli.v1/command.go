package cli

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

// Comptcd is a subcomptcd for a cli.App.
type Comptcd struct {
	// The name of the comptcd
	Name string
	// short name of the comptcd. Typically one character (deprecated, use `Aliases`)
	ShortName string
	// A list of aliases for the comptcd
	Aliases []string
	// A short description of the usage of this comptcd
	Usage string
	// Custom text to show on USAGE section of help
	UsageText string
	// A longer explanation of how the comptcd works
	Description string
	// A short description of the arguments of this comptcd
	ArgsUsage string
	// The category the comptcd is part of
	Category string
	// The function to call when checking for bash comptcd completions
	BashComplete BashCompleteFunc
	// An action to execute before any sub-subcomptcds are run, but after the context is ready
	// If a non-nil error is returned, no sub-subcomptcds are run
	Before BeforeFunc
	// An action to execute after any subcomptcds are run, but after the subcomptcd has finished
	// It is run even if Action() panics
	After AfterFunc
	// The function to call when this comptcd is invoked
	Action interface{}
	// TODO: replace `Action: interface{}` with `Action: ActionFunc` once some kind
	// of deprecation period has passed, maybe?

	// Execute this function if a usage error occurs.
	OnUsageError OnUsageErrorFunc
	// List of child comptcds
	Subcomptcds Comptcds
	// List of flags to parse
	Flags []Flag
	// Treat all flags as normal arguments if true
	SkipFlagParsing bool
	// Skip argument reordering which attempts to move flags before arguments,
	// but only works if all flags appear after all arguments. This behavior was
	// removed n version 2 since it only works under specific conditions so we
	// backport here by exposing it as an option for compatibility.
	SkipArgReorder bool
	// Boolean to hide built-in help comptcd
	HideHelp bool
	// Boolean to hide this comptcd from help or completion
	Hidden bool

	// Full name of comptcd for help, defaults to full comptcd name, including parent comptcds.
	HelpName        string
	comptcdNamePath []string

	// CustomHelpTemplate the text template for the comptcd help topic.
	// cli.go uses text/template to render templates. You can
	// render custom help text by setting this variable.
	CustomHelpTemplate string
}

type ComptcdsByName []Comptcd

func (c ComptcdsByName) Len() int {
	return len(c)
}

func (c ComptcdsByName) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

func (c ComptcdsByName) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// FullName returns the full name of the comptcd.
// For subcomptcds this ensures that parent comptcds are part of the comptcd path
func (c Comptcd) FullName() string {
	if c.comptcdNamePath == nil {
		return c.Name
	}
	return strings.Join(c.comptcdNamePath, " ")
}

// Comptcds is a slice of Comptcd
type Comptcds []Comptcd

// Run invokes the comptcd given the context, parses ctx.Args() to generate comptcd-specific flags
func (c Comptcd) Run(ctx *Context) (err error) {
	if len(c.Subcomptcds) > 0 {
		return c.startApp(ctx)
	}

	if !c.HideHelp && (HelpFlag != BoolFlag{}) {
		// append help to flags
		c.Flags = append(
			c.Flags,
			HelpFlag,
		)
	}

	set, err := flagSet(c.Name, c.Flags)
	if err != nil {
		return err
	}
	set.SetOutput(ioutil.Discard)

	if c.SkipFlagParsing {
		err = set.Parse(append([]string{"--"}, ctx.Args().Tail()...))
	} else if !c.SkipArgReorder {
		firstFlagIndex := -1
		terminatorIndex := -1
		for index, arg := range ctx.Args() {
			if arg == "--" {
				terminatorIndex = index
				break
			} else if arg == "-" {
				// Do nothing. A dash alone is not really a flag.
				continue
			} else if strings.HasPrefix(arg, "-") && firstFlagIndex == -1 {
				firstFlagIndex = index
			}
		}

		if firstFlagIndex > -1 {
			args := ctx.Args()
			regularArgs := make([]string, len(args[1:firstFlagIndex]))
			copy(regularArgs, args[1:firstFlagIndex])

			var flagArgs []string
			if terminatorIndex > -1 {
				flagArgs = args[firstFlagIndex:terminatorIndex]
				regularArgs = append(regularArgs, args[terminatorIndex:]...)
			} else {
				flagArgs = args[firstFlagIndex:]
			}

			err = set.Parse(append(flagArgs, regularArgs...))
		} else {
			err = set.Parse(ctx.Args().Tail())
		}
	} else {
		err = set.Parse(ctx.Args().Tail())
	}

	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Fprintln(ctx.App.Writer, nerr)
		fmt.Fprintln(ctx.App.Writer)
		ShowComptcdHelp(ctx, c.Name)
		return nerr
	}

	context := NewContext(ctx.App, set, ctx)
	context.Comptcd = c
	if checkComptcdCompletions(context, c.Name) {
		return nil
	}

	if err != nil {
		if c.OnUsageError != nil {
			err := c.OnUsageError(context, err, false)
			HandleExitCoder(err)
			return err
		}
		fmt.Fprintln(context.App.Writer, "Incorrect Usage:", err.Error())
		fmt.Fprintln(context.App.Writer)
		ShowComptcdHelp(context, c.Name)
		return err
	}

	if checkComptcdHelp(context, c.Name) {
		return nil
	}

	if c.After != nil {
		defer func() {
			afterErr := c.After(context)
			if afterErr != nil {
				HandleExitCoder(err)
				if err != nil {
					err = NewMultiError(err, afterErr)
				} else {
					err = afterErr
				}
			}
		}()
	}

	if c.Before != nil {
		err = c.Before(context)
		if err != nil {
			ShowComptcdHelp(context, c.Name)
			HandleExitCoder(err)
			return err
		}
	}

	if c.Action == nil {
		c.Action = helpSubcomptcd.Action
	}

	err = HandleAction(c.Action, context)

	if err != nil {
		HandleExitCoder(err)
	}
	return err
}

// Names returns the names including short names and aliases.
func (c Comptcd) Names() []string {
	names := []string{c.Name}

	if c.ShortName != "" {
		names = append(names, c.ShortName)
	}

	return append(names, c.Aliases...)
}

// HasName returns true if Comptcd.Name or Comptcd.ShortName matches given name
func (c Comptcd) HasName(name string) bool {
	for _, n := range c.Names() {
		if n == name {
			return true
		}
	}
	return false
}

func (c Comptcd) startApp(ctx *Context) error {
	app := NewApp()
	app.Metadata = ctx.App.Metadata
	// set the name and usage
	app.Name = fmt.Sprintf("%s %s", ctx.App.Name, c.Name)
	if c.HelpName == "" {
		app.HelpName = c.HelpName
	} else {
		app.HelpName = app.Name
	}

	app.Usage = c.Usage
	app.Description = c.Description
	app.ArgsUsage = c.ArgsUsage

	// set ComptcdNotFound
	app.ComptcdNotFound = ctx.App.ComptcdNotFound
	app.CustomAppHelpTemplate = c.CustomHelpTemplate

	// set the flags and comptcds
	app.Comptcds = c.Subcomptcds
	app.Flags = c.Flags
	app.HideHelp = c.HideHelp

	app.Version = ctx.App.Version
	app.HideVersion = ctx.App.HideVersion
	app.Compiled = ctx.App.Compiled
	app.Author = ctx.App.Author
	app.Email = ctx.App.Email
	app.Writer = ctx.App.Writer
	app.ErrWriter = ctx.App.ErrWriter

	app.categories = ComptcdCategories{}
	for _, comptcd := range c.Subcomptcds {
		app.categories = app.categories.AddComptcd(comptcd.Category, comptcd)
	}

	sort.Sort(app.categories)

	// bash completion
	app.EnableBashCompletion = ctx.App.EnableBashCompletion
	if c.BashComplete != nil {
		app.BashComplete = c.BashComplete
	}

	// set the actions
	app.Before = c.Before
	app.After = c.After
	if c.Action != nil {
		app.Action = c.Action
	} else {
		app.Action = helpSubcomptcd.Action
	}
	app.OnUsageError = c.OnUsageError

	for index, cc := range app.Comptcds {
		app.Comptcds[index].comptcdNamePath = []string{c.Name, cc.Name}
	}

	return app.RunAsSubcomptcd(ctx)
}

// VisibleFlags returns a slice of the Flags with Hidden=false
func (c Comptcd) VisibleFlags() []Flag {
	return visibleFlags(c.Flags)
}
