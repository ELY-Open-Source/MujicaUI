package mujica

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	docsGoBlock = regexp.MustCompile("(?s)```go\n(.*?)```")
	docsSpace   = regexp.MustCompile(`\s+`)
)

// TestDocSnippetsAreCompiledExamples keeps every doc's first Go block inside a compiled Example.
func TestDocSnippetsAreCompiledExamples(t *testing.T) {
	examples, err := filepath.Glob("example_*_test.go")
	if err != nil {
		t.Fatal(err)
	}
	var all strings.Builder
	for _, f := range examples {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		all.Write(b)
	}
	compiled := docsSpace.ReplaceAllString(all.String(), "")
	docs, err := filepath.Glob("docs/components/*.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 100 {
		t.Fatalf("%d component docs, want 100", len(docs))
	}
	for _, d := range docs {
		b, err := os.ReadFile(d)
		if err != nil {
			t.Fatal(err)
		}
		m := docsGoBlock.FindSubmatch(b)
		if m == nil {
			t.Errorf("%s has no Go example", d)
			continue
		}
		for _, line := range strings.Split(string(m[1]), "\n") {
			s := strings.TrimSpace(line)
			if s == "" || strings.HasPrefix(s, "//") {
				continue
			}
			if !strings.Contains(compiled, docsSpace.ReplaceAllString(s, "")) {
				t.Errorf("%s: %q is not in any compiled Example", d, s)
			}
		}
	}
}
