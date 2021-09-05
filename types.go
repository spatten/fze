package fze

type RunnerOptions struct {
	EmacsServer      string
	TestFilter       string
	Multi            bool
	SwitchToTmuxPane bool
}

type pathArg struct {
	path       string
	lineNumber string
}
