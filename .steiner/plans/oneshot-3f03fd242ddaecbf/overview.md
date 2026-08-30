# Create terminal Markdown viewer CLI

## Request
Create a new Go command-line Markdown viewer. Reuse the terminal renderer used by sibling `../steiner`, provide a sensible CLI, add a README, and add a Makefile with build output in `bin/` and a test target. Keep it small, follow Go practices, and provide good test coverage. Do not add an `AGENTS.md` file.

## Overview
Create a `markdown-viewer` executable that reads Markdown from one optional file path or standard input, renders it through `charm.land/glamour/v2`, and writes terminal-ready text to standard output. It will offer narrow, explicit rendering controls: `--style` selects a Glamour built-in style and `--width` controls word wrapping. The program will use a testable internal package for argument parsing, input selection, renderer errors, and output writing; the command entrypoint stays thin.

Project scaffolding will pin the known renderer version, set the minimum Go version its module requires, ignore generated binaries, and supply Make targets for build, test, formatting, and aggregate checks. README will document local build, direct execution, stdin use, options, and output behavior.

## Key Decisions
- **D1:** Use `charm.land/glamour/v2 v2.0.1`, matching `../steiner` and current official release. It provides all rendering behavior without custom terminal rendering code.
- **D2:** Name the executable and CLI usage `markdown-viewer`; accept at most one positional `FILE`, with omitted `FILE` or `-` meaning standard input. Options must precede `FILE`, matching Go's standard `flag` parser behavior.
- **D3:** Provide `--style` (default `dark`) and `--width` (default `80`, must be positive). Glamour does not detect terminal width or TTY state, so explicit options keep output predictable and avoid another dependency.
- **D4:** Keep command wiring in `cmd/markdown-viewer/main.go` and behavior in `internal/viewer`. Inject readers, writers, and a small renderer function into unexported execution logic so unit tests can cover renderer and output failures without changing the public CLI API.
- **D5:** Require Go 1.25.8 or newer because Glamour v2.0.1 declares that minimum. Pin module path to `github.com/luispabon/terminal-markdown-viewer`, based on the configured Git remote.
- **D6:** Use exit status 0 for success and help, 2 for argument or usage errors, and 1 for input, rendering, or output failures. Help text goes to stdout; diagnostics go to stderr.
- **D7:** Do not add `AGENTS.md`, CI configuration, terminal-size detection, paging, interactive controls, recursive directory input, or escape stripping. They exceed the requested first release.

## Tradeoffs
- Automatic TTY detection and current terminal width would improve defaults but need more terminal-specific code and testing. Explicit style and width flags retain a small, portable CLI; users can choose `--style notty` for redirected output.
- A root `main.go` would reduce file count, but a small `internal/viewer` package isolates testable behavior while leaving the executable entrypoint idiomatic.
- Supporting arbitrary Glamour JSON styles or style-file paths is deferred. Built-in style selection covers the intended use without exposing renderer-specific file semantics.

## Scope Boundaries
In scope: one-file or stdin input, built-in Glamour styling, configurable positive wrapping width, clear usage and error handling, unit tests, Go module metadata, README, `.gitignore`, and Make targets.

Out of scope: an interactive TUI, pager integration, live reload, Markdown editing, multiple file concatenation, directory traversal, terminal-size probing, TTY auto-detection, CI workflows, lint-tool configuration, and `AGENTS.md`.

## Verification Strategy
| Command | Purpose | Cost | Fix mode |
| --- | --- | --- | --- |
| `gofmt -w cmd/markdown-viewer/main.go internal/viewer/*.go` | Format changed Go files | Cheap | Safe scoped auto-fix |
| `go vet ./...` | Static checks | Cheap | No auto-fix |
| `go build ./...` | Compile all packages | Cheap | No auto-fix |
| `go test ./...` | Full unit suite | Cheap | No auto-fix |
| `go test -race ./...` | Detect races | Medium | No auto-fix |
| `go mod tidy` twice | Dependency hygiene and idempotence | Cheap | `go mod tidy` is safe module-managed fix |
| `make build` | Build `bin/markdown-viewer` | Cheap | No auto-fix |
| `make test` | Project test target | Cheap | No auto-fix |
| `make check` | Format validation, vet, build, and tests | Cheap | No auto-fix |

`goimports`, `golangci-lint`, and `govulncheck` are installed locally but are not plan requirements because this new project will not add their configuration. The Makefile's aggregate `check` target will run formatting validation, vet, build, and tests using Go toolchain commands only.

## Decision Log
- **Research decision:** External research was required because Glamour v2 API, style names, release version, and toolchain requirements can change. Findings are persisted in `research.md`.
- **Assumption A1:** The configured Git remote supplies the intended module path: `github.com/luispabon/terminal-markdown-viewer`.
- **Assumption A2:** `markdown-viewer` is an acceptable executable name because the user did not prescribe one; it matches the repository name and keeps usage clear.
- **Assumption A3:** Defaulting to Glamour's `dark` style and width 80 is appropriate for a terminal viewer. Redirected output remains available through explicit `--style notty`.
- **Advisor sanity check:** Applied all material findings. Steps no longer invoke tidy or build before Go source exists; the aggregate `check` target contradiction is removed; an unexported renderer seam permits render-failure coverage; exit destinations and statuses are fixed by D6; and dependency hygiene now requires two successful tidy runs rather than an ineffective diff over initially untracked files.
