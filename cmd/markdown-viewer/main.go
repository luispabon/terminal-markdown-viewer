package main

import (
	"os"

	"github.com/luispabon/terminal-markdown-viewer/internal/viewer"
)

func main() {
	os.Exit(viewer.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
