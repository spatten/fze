package fze

type RunnerOptions struct {
	EmacsServer      string
	Multi            bool
	SwitchToTmuxPane bool
}

type pathArg struct {
	path       string
	lineNumber string
}
