# terminal-markdown-viewer

`markdown-viewer` is a command-line Markdown viewer. It renders Markdown to terminal-ready text with [Glamour](https://github.com/charmbracelet/glamour).

## Prerequisite

Go 1.25.8 or newer.

## Build and test

```sh
make build
```

This produces `bin/markdown-viewer`.

```sh
make test
```

This runs the unit tests. `make check` runs formatting validation, `go vet`, a build, and the tests.

## Usage

```text
markdown-viewer [options] [FILE]
```

Render Markdown from a file:

```sh
./bin/markdown-viewer README.md
```

With no `FILE`, read Markdown from standard input:

```sh
cat README.md | ./bin/markdown-viewer
```

Use `-` to read from standard input explicitly:

```sh
./bin/markdown-viewer -
```

Options must precede `FILE`.

Options:

- `--style`: rendering style. Default: `dark`.
- `--width`: word-wrap width. Default: `80`. Must be a positive integer.
- `-h`, `--help`: show help.

Built-in styles are `ascii`, `dark`, `light`, `notty`, `pink`, `dracula`, and `tokyo-night`. Use `--style notty` for redirected or piped output when plain text without styling is wanted.

Exit statuses:

- `0`: success or help shown.
- `1`: input, rendering, or output failure.
- `2`: invalid usage.
