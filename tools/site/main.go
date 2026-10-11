// Command site builds the static MujicaUI documentation site into site/dist.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const repoURL = "https://git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI"

// component is one documented component and everything its page needs.
type component struct {
	ID, Slug, Name, ZhName, Category string
	Doc                              string // raw Markdown
	Light, Dark                      string // absolute PNG paths, empty when missing
}

type category struct {
	Name  string
	Items []*component
}

func main() {
	root := flag.String("root", ".", "repository root")
	out := flag.String("out", "site/dist", "output directory")
	flag.Parse()
	n, err := build(*root, *out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("site: %d component pages written to %s\n", n, *out)
}

// build reads the repo at root and writes the site to out, returning the component page count.
func build(root, out string) (int, error) {
	comps, err := loadComponents(root)
	if err != nil {
		return 0, err
	}
	if err := os.RemoveAll(out); err != nil {
		return 0, err
	}
	s := &site{root: root, out: out, comps: comps, cats: groupCategories(comps)}
	if err := s.write(); err != nil {
		return 0, err
	}
	return len(comps), nil
}

var (
	docNameRe   = regexp.MustCompile(`^(\d{3})-([a-z0-9-]+)\.md$`)
	entryIDRe   = regexp.MustCompile(`\bID:\s*"([^"]+)"`)
	fieldRe     = regexp.MustCompile(`\b(Name|ZhName|Category):\s*("(?:[^"\\]|\\.)*"|[A-Za-z_][A-Za-z0-9_]*)`)
	constRe     = regexp.MustCompile(`(?m)^\s*(?:const\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=\s*"([^"]*)"`)
	titleSplitR = regexp.MustCompile(`^#\s+(?:\d{3}\s+)?(\S+)\s*(.*)$`)
)

type entry struct{ Name, ZhName, Category string }

func loadComponents(root string) ([]*component, error) {
	entries, err := loadGallery(root)
	if err != nil {
		return nil, err
	}
	snaps, err := findSnapshots(root)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, "docs", "components")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var comps []*component
	for _, f := range files {
		g := docNameRe.FindStringSubmatch(f.Name())
		if g == nil {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		c := &component{ID: g[1], Slug: strings.TrimSuffix(f.Name(), ".md"), Doc: string(raw), Category: "Other"}
		if t := titleSplitR.FindStringSubmatch(firstHeading(c.Doc)); t != nil {
			c.Name, c.ZhName = t[1], t[2]
		}
		if e, ok := entries[c.ID]; ok {
			c.Name, c.ZhName, c.Category = e.Name, e.ZhName, e.Category
		}
		if c.Name == "" {
			return nil, fmt.Errorf("site: %s has no heading and no gallery entry", f.Name())
		}
		key := normKey(c.Name)
		c.Light, c.Dark = snaps[key+"-light"], snaps[key+"-dark"]
		comps = append(comps, c)
	}
	if len(comps) == 0 {
		return nil, fmt.Errorf("site: no component docs in %s", dir)
	}
	return comps, nil
}

func firstHeading(doc string) string {
	for _, l := range strings.Split(doc, "\n") {
		if strings.HasPrefix(l, "# ") {
			return strings.TrimSpace(l)
		}
	}
	return ""
}

// loadGallery scans examples/gallery/pages_*.go for page entries without importing the package.
func loadGallery(root string) (map[string]entry, error) {
	files, err := filepath.Glob(filepath.Join(root, "examples", "gallery", "pages_*.go"))
	if err != nil {
		return nil, err
	}
	consts := map[string]string{}
	var srcs []string
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		src := string(b)
		srcs = append(srcs, src)
		for _, m := range constRe.FindAllStringSubmatch(src, -1) {
			consts[m[1]] = m[2]
		}
	}
	out := map[string]entry{}
	for _, src := range srcs {
		locs := entryIDRe.FindAllStringSubmatchIndex(src, -1)
		for i, loc := range locs {
			end := len(src)
			if i+1 < len(locs) {
				end = locs[i+1][0]
			}
			id := src[loc[2]:loc[3]]
			var e entry
			for _, f := range fieldRe.FindAllStringSubmatch(src[loc[1]:end], -1) {
				v := f[2]
				if strings.HasPrefix(v, `"`) {
					if u, err := strconv.Unquote(v); err == nil {
						v = u
					}
				} else if cv, ok := consts[v]; ok {
					v = cv
				} else {
					continue // a variable such as Code: fooCode is not a field we want
				}
				switch f[1] {
				case "Name":
					if e.Name == "" {
						e.Name = v
					}
				case "ZhName":
					if e.ZhName == "" {
						e.ZhName = v
					}
				case "Category":
					if e.Category == "" {
						e.Category = v
					}
				}
			}
			if e.Name != "" && e.Category != "" {
				out[id] = e
			}
		}
	}
	return out, nil
}

// findSnapshots indexes every PNG under any testdata/snapshots/darwin directory by normalized base name.
func findSnapshots(root string) (map[string]string, error) {
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if p != root && (strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist") {
			return filepath.SkipDir
		}
		if filepath.ToSlash(p) != "testdata/snapshots/darwin" && !strings.HasSuffix(filepath.ToSlash(p), "/testdata/snapshots/darwin") {
			return nil
		}
		pngs, err := filepath.Glob(filepath.Join(p, "*.png"))
		if err != nil {
			return err
		}
		sort.Strings(pngs)
		for _, f := range pngs {
			base := strings.TrimSuffix(filepath.Base(f), ".png")
			if strings.HasSuffix(base, ".got") {
				continue
			}
			i := strings.LastIndex(base, "-")
			if i < 0 {
				continue
			}
			key := normKey(base[:i]) + base[i:]
			if _, ok := out[key]; !ok {
				out[key] = f
			}
		}
		return filepath.SkipDir
	})
	return out, err
}

// normKey folds snake_case and CamelCase names to one key, so OAuthButtons meets oauth_buttons.
func normKey(s string) string {
	return strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(s))
}

func groupCategories(comps []*component) []*category {
	idx := map[string]*category{}
	var cats []*category
	for _, c := range comps {
		cat, ok := idx[c.Category]
		if !ok {
			cat = &category{Name: c.Category}
			idx[c.Category] = cat
			cats = append(cats, cat)
		}
		cat.Items = append(cat.Items, c)
	}
	return cats
}
