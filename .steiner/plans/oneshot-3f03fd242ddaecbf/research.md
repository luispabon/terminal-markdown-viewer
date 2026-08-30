## Question
What current official Glamour v2 API and terminal behavior should govern a small Go Markdown viewer CLI?

## Findings
- `charm.land/glamour/v2 v2.0.1` is current and is the version used by the sibling `../steiner` project.
- The required API is `glamour.NewTermRenderer(options ...TermRendererOption) (*TermRenderer, error)` followed by `(*TermRenderer).Render(string) (string, error)`.
- `glamour.WithWordWrap(int)`, `glamour.WithPreservedNewLines()`, and `glamour.WithStandardStyle(string)` are available in v2.0.1.
- Supported built-in style names are `ascii`, `dark`, `light`, `notty`, `pink`, `dracula`, and `tokyo-night`; `dark` is Glamour's default.
- Glamour does not inspect terminal capabilities or width. Its default wrapping width is 80. It does not automatically change styles for piped output.
- Glamour v2.0.1 declares `go 1.25.8` in its module, so this project must require Go 1.25.8 or newer.
- Renderer creation and rendering can both return errors and must be checked.

## Implications
The CLI can be dependency-light: expose an explicit `--width` flag with a default of 80 and an explicit `--style` flag with a default of `dark`. Do not add terminal-size or TTY detection in the initial release. Document `--style notty` for plain output when redirecting or piping. Validate a positive width before calling Glamour; pass style validation through Glamour and report its error.

## Risks and Uncertainties
- The `notty` style removes color but Glamour source still emits OSC 8 hyperlinks for Markdown links. Byte-clean piped output would need escape stripping, which is outside this simple viewer's scope.
- The plan avoids terminal capability detection to avoid an extra dependency and platform-specific behavior. Users choose style and width explicitly.

## Sources
- https://pkg.go.dev/charm.land/glamour/v2
- https://raw.githubusercontent.com/charmbracelet/glamour/v2.0.1/glamour.go
- https://raw.githubusercontent.com/charmbracelet/glamour/v2.0.1/styles/styles.go
- https://raw.githubusercontent.com/charmbracelet/glamour/v2.0.1/README.md
- https://raw.githubusercontent.com/charmbracelet/glamour/v2.0.1/UPGRADE_GUIDE_V2.md
- https://proxy.golang.org/charm.land/glamour/v2/@latest

## Open Questions
None material to the initial CLI.
