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
	docsPageID  = regexp.MustCompile(`ID: "(\d{3})"`)
)

// docsLines normalizes code into its non-empty, non-comment lines without whitespace.
func docsLines(code string) []string {
	var out []string
	for _, l := range strings.Split(code, "\n") {
		s := strings.TrimSpace(l)
		if s != "" && !strings.HasPrefix(s, "//") {
			out = append(out, docsSpace.ReplaceAllString(s, ""))
		}
	}
	return out
}

// docsContains reports whether want appears in order inside have.
func docsContains(have, want []string) bool {
	i := 0
	for _, h := range have {
		if i < len(want) && h == want[i] {
			i++
		}
	}
	return i == len(want)
}

func readGlob(t *testing.T, pattern string) map[string]string {
	t.Helper()
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		out[f] = string(b)
	}
	return out
}

// TestDocSnippetsAreCompiledExamples keeps each doc's first Go block, in order, inside one Example.
func TestDocSnippetsAreCompiledExamples(t *testing.T) {
	var examples [][]string
	for _, src := range readGlob(t, "example_*_test.go") {
		for _, fn := range strings.Split(src, "\nfunc Example")[1:] {
			examples = append(examples, docsLines(fn))
		}
	}
	for d, src := range readGlob(t, "docs/components/*.md") {
		m := docsGoBlock.FindStringSubmatch(src)
		if m == nil {
			t.Errorf("%s has no Go example", d)
			continue
		}
		want, found := docsLines(m[1]), false
		for _, ex := range examples {
			if docsContains(ex, want) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s: its Go block is not, in order, inside any compiled Example", d)
		}
	}
}

// TestDocsMatchGalleryPages wants exactly one doc and one gallery page per component ID.
func TestDocsMatchGalleryPages(t *testing.T) {
	pages := map[string]int{}
	for _, src := range readGlob(t, "examples/gallery/pages_*.go") {
		for _, m := range docsPageID.FindAllStringSubmatch(src, -1) {
			pages[m[1]]++
		}
	}
	docs := map[string]int{}
	for d := range readGlob(t, "docs/components/*.md") {
		docs[filepath.Base(d)[:3]]++
	}
	for id, n := range pages {
		if n != 1 || docs[id] != 1 {
			t.Errorf("component %s: %d gallery pages, %d docs; want 1 and 1", id, n, docs[id])
		}
	}
	for id, n := range docs {
		if pages[id] == 0 {
			t.Errorf("doc %s (%d files) has no gallery page", id, n)
		}
	}
}
