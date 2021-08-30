package fze

type runnerOptions struct {
	fzfOptions fzfOptions
}

type fzfOptions struct {
	multi bool
}

type openEditorArgs struct {
	path       string
	lineNumber string
}
