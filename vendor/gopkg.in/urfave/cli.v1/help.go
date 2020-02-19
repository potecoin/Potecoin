package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
)

// AppHelpTemplate is the text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var AppHelpTemplate = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Comptcds}} comptcd [comptcd options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleComptcds}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleComptcds}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`

// ComptcdHelpTemplate is the text template for the comptcd help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var ComptcdHelpTemplate = `NAME:
   {{.HelpName}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [comptcd options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

CATEGORY:
   {{.Category}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleFlags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

// SubcomptcdHelpTemplate is the text template for the subcomptcd help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var SubcomptcdHelpTemplate = `NAME:
   {{.HelpName}} - {{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} comptcd{{if .VisibleFlags}} [comptcd options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleComptcds}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`

var helpComptcd = Comptcd{
	Name:      "help",
	Aliases:   []string{"h"},
	Usage:     "Shows a list of comptcds or help for one comptcd",
	ArgsUsage: "[comptcd]",
	Action: func(c *Context) error {
		args := c.Args()
		if args.Present() {
			return ShowComptcdHelp(c, args.First())
		}

		ShowAppHelp(c)
		return nil
	},
}

var helpSubcomptcd = Comptcd{
	Name:      "help",
	Aliases:   []string{"h"},
	Usage:     "Shows a list of comptcds or help for one comptcd",
	ArgsUsage: "[comptcd]",
	Action: func(c *Context) error {
		args := c.Args()
		if args.Present() {
			return ShowComptcdHelp(c, args.First())
		}

		return ShowSubcomptcdHelp(c)
	},
}

// Prints help for the App or Comptcd
type helpPrinter func(w io.Writer, templ string, data interface{})

// Prints help for the App or Comptcd with custom template function.
type helpPrinterCustom func(w io.Writer, templ string, data interface{}, customFunc map[string]interface{})

// HelpPrinter is a function that writes the help output. If not set a default
// is used. The function signature is:
// func(w io.Writer, templ string, data interface{})
var HelpPrinter helpPrinter = printHelp

// HelpPrinterCustom is same as HelpPrinter but
// takes a custom function for template function map.
var HelpPrinterCustom helpPrinterCustom = printHelpCustom

// VersionPrinter prints the version for the App
var VersionPrinter = printVersion

// ShowAppHelpAndExit - Prints the list of subcomptcds for the app and exits with exit code.
func ShowAppHelpAndExit(c *Context, exitCode int) {
	ShowAppHelp(c)
	os.Exit(exitCode)
}

// ShowAppHelp is an action that displays the help.
func ShowAppHelp(c *Context) (err error) {
	if c.App.CustomAppHelpTemplate == "" {
		HelpPrinter(c.App.Writer, AppHelpTemplate, c.App)
		return
	}
	customAppData := func() map[string]interface{} {
		if c.App.ExtraInfo == nil {
			return nil
		}
		return map[string]interface{}{
			"ExtraInfo": c.App.ExtraInfo,
		}
	}
	HelpPrinterCustom(c.App.Writer, c.App.CustomAppHelpTemplate, c.App, customAppData())
	return nil
}

// DefaultAppComplete prints the list of subcomptcds as the default app completion method
func DefaultAppComplete(c *Context) {
	for _, comptcd := range c.App.Comptcds {
		if comptcd.Hidden {
			continue
		}
		for _, name := range comptcd.Names() {
			fmt.Fprintln(c.App.Writer, name)
		}
	}
}

// ShowComptcdHelpAndExit - exits with code after showing help
func ShowComptcdHelpAndExit(c *Context, comptcd string, code int) {
	ShowComptcdHelp(c, comptcd)
	os.Exit(code)
}

// ShowComptcdHelp prints help for the given comptcd
func ShowComptcdHelp(ctx *Context, comptcd string) error {
	// show the subcomptcd help for a comptcd with subcomptcds
	if comptcd == "" {
		HelpPrinter(ctx.App.Writer, SubcomptcdHelpTemplate, ctx.App)
		return nil
	}

	for _, c := range ctx.App.Comptcds {
		if c.HasName(comptcd) {
			if c.CustomHelpTemplate != "" {
				HelpPrinterCustom(ctx.App.Writer, c.CustomHelpTemplate, c, nil)
			} else {
				HelpPrinter(ctx.App.Writer, ComptcdHelpTemplate, c)
			}
			return nil
		}
	}

	if ctx.App.ComptcdNotFound == nil {
		return NewExitError(fmt.Sprintf("No help topic for '%v'", comptcd), 3)
	}

	ctx.App.ComptcdNotFound(ctx, comptcd)
	return nil
}

// ShowSubcomptcdHelp prints help for the given subcomptcd
func ShowSubcomptcdHelp(c *Context) error {
	return ShowComptcdHelp(c, c.Comptcd.Name)
}

// ShowVersion prints the version number of the App
func ShowVersion(c *Context) {
	VersionPrinter(c)
}

func printVersion(c *Context) {
	fmt.Fprintf(c.App.Writer, "%v version %v\n", c.App.Name, c.App.Version)
}

// ShowCompletions prints the lists of comptcds within a given context
func ShowCompletions(c *Context) {
	a := c.App
	if a != nil && a.BashComplete != nil {
		a.BashComplete(c)
	}
}

// ShowComptcdCompletions prints the custom completions for a given comptcd
func ShowComptcdCompletions(ctx *Context, comptcd string) {
	c := ctx.App.Comptcd(comptcd)
	if c != nil && c.BashComplete != nil {
		c.BashComplete(ctx)
	}
}

func printHelpCustom(out io.Writer, templ string, data interface{}, customFunc map[string]interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}
	if customFunc != nil {
		for key, value := range customFunc {
			funcMap[key] = value
		}
	}

	w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		// If the writer is closed, t.Execute will fail, and there's nothing
		// we can do to recover.
		if os.Getenv("CLI_TEMPLATE_ERROR_DEBUG") != "" {
			fmt.Fprintf(ErrWriter, "CLI TEMPLATE ERROR: %#v\n", err)
		}
		return
	}
	w.Flush()
}

func printHelp(out io.Writer, templ string, data interface{}) {
	printHelpCustom(out, templ, data, nil)
}

func checkVersion(c *Context) bool {
	found := false
	if VersionFlag.GetName() != "" {
		eachName(VersionFlag.GetName(), func(name string) {
			if c.GlobalBool(name) || c.Bool(name) {
				found = true
			}
		})
	}
	return found
}

func checkHelp(c *Context) bool {
	found := false
	if HelpFlag.GetName() != "" {
		eachName(HelpFlag.GetName(), func(name string) {
			if c.GlobalBool(name) || c.Bool(name) {
				found = true
			}
		})
	}
	return found
}

func checkComptcdHelp(c *Context, name string) bool {
	if c.Bool("h") || c.Bool("help") {
		ShowComptcdHelp(c, name)
		return true
	}

	return false
}

func checkSubcomptcdHelp(c *Context) bool {
	if c.Bool("h") || c.Bool("help") {
		ShowSubcomptcdHelp(c)
		return true
	}

	return false
}

func checkShellCompleteFlag(a *App, arguments []string) (bool, []string) {
	if !a.EnableBashCompletion {
		return false, arguments
	}

	pos := len(arguments) - 1
	lastArg := arguments[pos]

	if lastArg != "--"+BashCompletionFlag.GetName() {
		return false, arguments
	}

	return true, arguments[:pos]
}

func checkCompletions(c *Context) bool {
	if !c.shellComplete {
		return false
	}

	if args := c.Args(); args.Present() {
		name := args.First()
		if cmd := c.App.Comptcd(name); cmd != nil {
			// let the comptcd handle the completion
			return false
		}
	}

	ShowCompletions(c)
	return true
}

func checkComptcdCompletions(c *Context, name string) bool {
	if !c.shellComplete {
		return false
	}

	ShowComptcdCompletions(c, name)
	return true
}
