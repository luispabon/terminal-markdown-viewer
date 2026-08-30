# Review: terminal Markdown viewer

## Scope and inputs

Reviewed `overview.md`, `plan.yaml`, and `execution.md` against committed implementation: module metadata, Makefile, README, `cmd/markdown-viewer/main.go`, and `internal/viewer` source and tests.

## Review status

`pass_with_notes`

Plan adherence is complete. The CLI accepts an optional file or stdin, supports the planned style and width flags, uses Glamour v2, and preserves the specified exit-code and stream behavior in ordinary error paths. No `AGENTS.md` or out-of-scope runtime features were added.

## Findings

### Blocking

None.

### Non-blocking

- **N1, closed-pipe exit status:** When stdout is closed by a downstream consumer, such as `markdown-viewer big.md | head -1`, the process receives SIGPIPE and exits `141`, rather than the README's documented `1` for output failure. This is conventional Unix behavior, is recorded in `execution.md`, and does not affect ordinary surfaced write errors, which report an error and exit `1`.

### Informational

- **I1, help output errors:** `printUsage` ignores write errors, so `--help` can return `0` with a failing stdout.
- **I2, direct test gap:** The explicit `-` stdin read-error message path has no direct test. It shares the tested input helper with the default-stdin path.

## Verification reruns

All passed:

- `go vet ./...`
- `go build ./...`
- `go test ./...`
- `go test -race ./...`
- `make check`
- `make build`
- `gofmt -l cmd/markdown-viewer internal/viewer` produced no output.
- Ran `go mod tidy` twice with no changes to `go.mod` or `go.sum`.

Behavior checks also passed for stdin, `-`, file input, help output, invalid width and argument forms, invalid style, missing files, directories, and every README-listed style.

## Residual risks and closeout

N1 is the only non-blocking contract caveat. I1 and I2 are low-risk edge cases and need no scope expansion. Worktree is clean after review checks. Branch push and PR creation are engine-owned closeout actions.

## Advisor sanity check

Advisor recommends `pass_with_notes`. It found no material blocker: SIGPIPE is a conventional signal path, and the help-write and direct-test observations are low-risk. Artifacts are ready for closeout.
