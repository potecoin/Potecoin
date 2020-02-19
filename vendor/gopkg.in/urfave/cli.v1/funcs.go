package cli

// BashCompleteFunc is an action to execute when the bash-completion flag is set
type BashCompleteFunc func(*Context)

// BeforeFunc is an action to execute before any subcomptcds are run, but after
// the context is ready if a non-nil error is returned, no subcomptcds are run
type BeforeFunc func(*Context) error

// AfterFunc is an action to execute after any subcomptcds are run, but after the
// subcomptcd has finished it is run even if Action() panics
type AfterFunc func(*Context) error

// ActionFunc is the action to execute when no subcomptcds are specified
type ActionFunc func(*Context) error

// ComptcdNotFoundFunc is executed if the proper comptcd cannot be found
type ComptcdNotFoundFunc func(*Context, string)

// OnUsageErrorFunc is executed if an usage error occurs. This is useful for displaying
// customized usage error messages.  This function is able to replace the
// original error messages.  If this function is not set, the "Incorrect usage"
// is displayed and the execution is interrupted.
type OnUsageErrorFunc func(context *Context, err error, isSubcomptcd bool) error

// FlagStringFunc is used by the help generation to display a flag, which is
// expected to be a single line.
type FlagStringFunc func(Flag) string
