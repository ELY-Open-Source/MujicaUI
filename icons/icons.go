// Package icons embeds the Lucide icons MujicaUI draws.
package icons

import (
	"embed"
	"path"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"
)

//go:embed svg/*.svg
var files embed.FS

var set = load()

func load() map[string]*ui.SVG {
	entries, err := files.ReadDir("svg")
	if err != nil {
		panic(err)
	}
	m := make(map[string]*ui.SVG, len(entries))
	for _, e := range entries {
		data, err := files.ReadFile(path.Join("svg", e.Name()))
		if err != nil {
			panic(err)
		}
		m[strings.TrimSuffix(e.Name(), ".svg")] = ui.MustParseSVG(data)
	}
	return m
}

// Get returns the icon named name.
func Get(name string) (*ui.SVG, bool) {
	s, ok := set[name]
	return s, ok
}

// Must returns the icon named name and panics when it is missing.
func Must(name string) *ui.SVG {
	s, ok := set[name]
	if !ok {
		panic("icons: no icon " + name)
	}
	return s
}

// Names lists every icon name in order.
func Names() []string {
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	slices.Sort(names)
	return names
}
