package viewer

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"charm.land/glamour/v2"
)

type renderer func(style string, width int, md string) (string, error)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return run(args, stdin, stdout, stderr, render)
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, render renderer) int {
	fs := flag.NewFlagSet("markdown-viewer", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}
	style := fs.String("style", "dark", "rendering style")
	width := fs.Int("width", 80, "word wrap width")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(stdout)
			return 0
		}
		fmt.Fprintf(stderr, "markdown-viewer: %s\n", err)
		return 2
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(stderr, "markdown-viewer: at most one FILE may be supplied")
		return 2
	}
	if *width < 1 {
		fmt.Fprintln(stderr, "markdown-viewer: --width must be a positive integer")
		return 2
	}

	md, err := readInput(fs.Arg(0), stdin)
	if err != nil {
		if fs.NArg() == 0 || fs.Arg(0) == "-" {
			fmt.Fprintf(stderr, "markdown-viewer: stdin: %s\n", err)
		} else {
			fmt.Fprintf(stderr, "markdown-viewer: %s: %s\n", fs.Arg(0), err)
		}
		return 1
	}

	output, err := render(*style, *width, string(md))
	if err != nil {
		fmt.Fprintf(stderr, "markdown-viewer: render: %s\n", err)
		return 1
	}
	if n, err := io.WriteString(stdout, output); err != nil {
		fmt.Fprintf(stderr, "markdown-viewer: write output: %s\n", err)
		return 1
	} else if n != len(output) {
		fmt.Fprintf(stderr, "markdown-viewer: write output: %s\n", io.ErrShortWrite)
		return 1
	}
	return 0
}

func readInput(path string, stdin io.Reader) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(stdin)
	}
	return os.ReadFile(path)
}

func render(style string, width int, md string) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(width),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		return "", err
	}
	return r.Render(md)
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: markdown-viewer [options] [FILE]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Render Markdown from FILE or standard input. Use '-' for standard input.")
	fmt.Fprintln(w, "Options must precede FILE.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  --style string   rendering style (default \"dark\")")
	fmt.Fprintln(w, "  --width int      word wrap width (default 80)")
	fmt.Fprintln(w, "  -h, --help       show this help")
}
