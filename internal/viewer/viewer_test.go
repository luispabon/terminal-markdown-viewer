package viewer

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInputSources(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		stdin    string
		fileData string
	}{
		{name: "stdin", stdin: "# stdin"},
		{name: "explicit stdin", args: []string{"-"}, stdin: "# dash"},
		{name: "file", args: []string{"doc.md"}, fileData: "# file"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := tc.args
			stdin := strings.NewReader(tc.stdin)
			if tc.fileData != "" {
				path := filepath.Join(t.TempDir(), "doc.md")
				if err := os.WriteFile(path, []byte(tc.fileData), 0600); err != nil {
					t.Fatal(err)
				}
				args = []string{path}
			}
			var stdout, stderr bytes.Buffer
			status := run(args, stdin, &stdout, &stderr, func(style string, width int, md string) (string, error) {
				return "rendered:" + md, nil
			})
			if status != 0 {
				t.Fatalf("status = %d, stderr = %q", status, stderr.String())
			}
			if expected := "rendered:" + tc.stdin + tc.fileData; stdout.String() != expected {
				t.Fatalf("stdout = %q, want %q", stdout.String(), expected)
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr = %q", stderr.String())
			}
		})
	}
}

func TestRunForwardsStyleAndWidth(t *testing.T) {
	var gotStyle string
	var gotWidth int
	var stdout, stderr bytes.Buffer
	status := run([]string{"--style", "light", "--width", "37"}, strings.NewReader("text"), &stdout, &stderr,
		func(style string, width int, md string) (string, error) {
			gotStyle, gotWidth = style, width
			return "ok", nil
		})
	if status != 0 || stdout.String() != "ok" {
		t.Fatalf("status = %d, stdout = %q, stderr = %q", status, stdout.String(), stderr.String())
	}
	if gotStyle != "light" || gotWidth != 37 {
		t.Fatalf("renderer got style %q and width %d", gotStyle, gotWidth)
	}
}

func TestRunDefaults(t *testing.T) {
	var gotStyle string
	var gotWidth int
	var stdout, stderr bytes.Buffer
	status := run(nil, strings.NewReader("text"), &stdout, &stderr,
		func(style string, width int, md string) (string, error) {
			gotStyle, gotWidth = style, width
			return "ok", nil
		})
	if status != 0 || stdout.String() != "ok" {
		t.Fatalf("status = %d, stdout = %q, stderr = %q", status, stdout.String(), stderr.String())
	}
	if gotStyle != "dark" || gotWidth != 80 {
		t.Fatalf("renderer got style %q and width %d", gotStyle, gotWidth)
	}
}

func TestRunHelp(t *testing.T) {
	for _, arg := range []string{"--help", "-h"} {
		t.Run(arg, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := run([]string{arg}, strings.NewReader("ignored"), &stdout, &stderr, failingRenderer)
			if status != 0 {
				t.Fatalf("status = %d", status)
			}
			for _, text := range []string{"markdown-viewer [options] [FILE]", "--style", "--width", "Options must precede FILE"} {
				if !strings.Contains(stdout.String(), text) {
					t.Errorf("help lacks %q: %q", text, stdout.String())
				}
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr = %q", stderr.String())
			}
		})
	}
}

func TestRunArgumentErrors(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{name: "zero width", args: []string{"--width", "0"}},
		{name: "negative width", args: []string{"--width", "-1"}},
		{name: "two files", args: []string{"one.md", "two.md"}},
		{name: "unknown option", args: []string{"--unknown"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := run(tc.args, strings.NewReader("ignored"), &stdout, &stderr, failingRenderer)
			if status != 2 {
				t.Fatalf("status = %d, want 2", status)
			}
			if stderr.Len() == 0 {
				t.Fatal("expected stderr diagnostic")
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
		})
	}
}

func TestRunReadErrors(t *testing.T) {
	var stdout, stderr bytes.Buffer
	status := run([]string{"missing.md"}, strings.NewReader("ignored"), &stdout, &stderr, failingRenderer)
	if status != 1 || !strings.Contains(stderr.String(), "missing.md") {
		t.Fatalf("status = %d, stderr = %q", status, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	status = run(nil, errorReader{}, &stdout, &stderr, failingRenderer)
	if status != 1 || !strings.Contains(stderr.String(), "stdin") {
		t.Fatalf("status = %d, stderr = %q", status, stderr.String())
	}
}

func TestRunRendererError(t *testing.T) {
	wantErr := errors.New("renderer failed")
	var stdout, stderr bytes.Buffer
	status := run(nil, strings.NewReader("text"), &stdout, &stderr, func(string, int, string) (string, error) {
		return "", wantErr
	})
	if status != 1 || !strings.Contains(stderr.String(), wantErr.Error()) {
		t.Fatalf("status = %d, stderr = %q", status, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunWriterError(t *testing.T) {
	var stderr bytes.Buffer
	status := run(nil, strings.NewReader("text"), failingWriter{}, &stderr, func(string, int, string) (string, error) {
		return "rendered", nil
	})
	if status != 1 || !strings.Contains(stderr.String(), "write output") {
		t.Fatalf("status = %d, stderr = %q", status, stderr.String())
	}
}

func TestRenderMarkdown(t *testing.T) {
	output, err := render("dark", 80, "# Heading\n\nSome **bold** text.\n\n- first item\n")
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"Heading", "bold", "first", "item"} {
		if !strings.Contains(output, text) {
			t.Errorf("rendered output lacks %q: %q", text, output)
		}
	}
}

func TestRenderInvalidStyle(t *testing.T) {
	_, err := render("bogus", 80, "text")
	if err == nil {
		t.Fatal("expected invalid style error")
	}
}

func failingRenderer(string, int, string) (string, error) {
	return "", errors.New("unexpected renderer call")
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

var _ io.Reader = errorReader{}
var _ io.Writer = failingWriter{}
