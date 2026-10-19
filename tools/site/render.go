package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type site struct {
	root, out string
	comps     []*component
	cats      []*category
}

// extraPages are repo documents published as their own pages.
var extraPages = []struct{ Slug, Src, Title string }{
	{"readme", "README.md", "README"},
	{"compat", "docs/compat.md", "兼容与平台依赖"},
	{"upstream", "docs/upstream.md", "上游归并"},
}

var compLinkRe = regexp.MustCompile(`^(?:\.\./|docs/components/|\./)?(\d{3}-[a-z0-9-]+)\.md(#.*)?$`)

// link rewrites repository-relative hrefs to site routes or repository URLs.
func (s *site) link(from string) func(string) string {
	return func(href string) string {
		if strings.Contains(href, "://") || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "mailto:") {
			return href
		}
		if g := compLinkRe.FindStringSubmatch(href); g != nil {
			return "/c/" + g[1] + "/" + g[2]
		}
		resolved := path.Clean(path.Join(path.Dir(from), href))
		for _, p := range extraPages {
			if resolved == p.Src {
				return "/" + p.Slug + "/"
			}
		}
		if resolved == "docs/components" {
			return "/#components"
		}
		return repoURL + "/src/branch/main/" + resolved
	}
}

func (s *site) write() error {
	if err := copyDir(filepath.Join(s.root, "site"), s.out, func(name string) bool {
		return strings.HasSuffix(name, ".css") || strings.HasSuffix(name, ".js")
	}); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(s.out, "img"), 0o755); err != nil {
		return err
	}
	if err := s.writeFile("index.html", s.index()); err != nil {
		return err
	}
	if err := s.writeFile("404.html", s.page("页面不存在", "", `<article class="doc"><h1>页面不存在</h1><p>找不到这个地址。<a href="/">返回组件目录</a>。</p></article>`)); err != nil {
		return err
	}
	for _, p := range extraPages {
		raw, err := os.ReadFile(filepath.Join(s.root, p.Src))
		if err != nil {
			return err
		}
		body := `<article class="doc">` + markdown{link: s.link(p.Src)}.render(string(raw)) + `</article>`
		if err := s.writeFile(p.Slug+"/index.html", s.page(p.Title, "", body)); err != nil {
			return err
		}
	}
	for _, c := range s.comps {
		for _, src := range []string{c.Light, c.Dark} {
			if src != "" {
				if err := copyFile(src, filepath.Join(s.out, "img", filepath.Base(src))); err != nil {
					return err
				}
			}
		}
		if err := s.writeFile("c/"+c.Slug+"/index.html", s.componentPage(c)); err != nil {
			return err
		}
	}
	return nil
}

func (s *site) writeFile(rel, content string) error {
	p := filepath.Join(s.out, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(content), 0o644)
}

const themeBoot = `<script>try{var t=localStorage.getItem("mujica-theme");if(t)document.documentElement.dataset.theme=t}catch(e){}</script>`

func (s *site) page(title, current, main string) string {
	e := html.EscapeString
	var b strings.Builder
	b.WriteString(`<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">`)
	b.WriteString(`<meta name="color-scheme" content="light dark"><title>` + e(title) + ` · MujicaUI</title>`)
	b.WriteString(`<link rel="stylesheet" href="/style.css">` + themeBoot + `</head><body>`)
	b.WriteString(`<a class="skip" href="#main">跳到正文</a><header class="top"><a class="brand" href="/"><span class="lozenge" aria-hidden="true"></span>MujicaUI</a><nav class="toplinks">`)
	b.WriteString(`<a href="/#components">组件</a><a href="/readme/">README</a><a href="/compat/">兼容</a><a href="/upstream/">上游</a><a href="` + repoURL + `" rel="noopener">Git</a></nav>`)
	b.WriteString(`<button type="button" class="theme" aria-label="切换浅色或深色主题">主题</button></header>`)
	if current != "" {
		b.WriteString(`<div class="layout"><button type="button" class="navtoggle" aria-expanded="false" aria-controls="sidenav">组件目录</button>`)
		b.WriteString(s.sidebar(current))
		b.WriteString(`<main id="main">` + main + `</main></div>`)
	} else {
		b.WriteString(`<main id="main" class="solo">` + main + `</main>`)
	}
	b.WriteString(`<footer class="foot"><div class="rule" aria-hidden="true"><span></span></div><p>MujicaUI 的组件是原生 Go（MyGo）控件，在桌面窗口里由 GPU 自绘。本站图片是这些控件渲染出的快照，浏览器不运行任何组件代码。</p>`)
	b.WriteString(`<p>MujicaUI components are native Go (MyGo) widgets. The images here are rendered snapshots; this site runs no component code. Source: <a href="` + repoURL + `" rel="noopener">github.com/ELY-Open-Source/MujicaUI</a></p></footer>`)
	b.WriteString(`<script src="/app.js"></script></body></html>`)
	return b.String()
}

func (s *site) sidebar(current string) string {
	e := html.EscapeString
	var b strings.Builder
	b.WriteString(`<nav class="side" id="sidenav" aria-label="组件目录">`)
	for _, cat := range s.cats {
		open := ""
		for _, c := range cat.Items {
			if c.Slug == current {
				open = " open"
			}
		}
		fmt.Fprintf(&b, `<details%s><summary>%s <span class="count">%d</span></summary><ul>`, open, e(cat.Name), len(cat.Items))
		for _, c := range cat.Items {
			aria := ""
			if c.Slug == current {
				aria = ` aria-current="page"`
			}
			fmt.Fprintf(&b, `<li><a href="/c/%s/"%s>%s <span class="zh">%s</span></a></li>`, c.Slug, aria, e(c.Name), e(c.ZhName))
		}
		b.WriteString(`</ul></details>`)
	}
	b.WriteString(`</nav>`)
	return b.String()
}

func (s *site) componentPage(c *component) string {
	e := html.EscapeString
	var b strings.Builder
	fmt.Fprintf(&b, `<p class="crumb"><a href="/#%s">%s</a> · #%s</p>`, e(slug(c.Category)), e(c.Category), c.ID)
	b.WriteString(`<section class="snaps" aria-label="渲染快照">`)
	for _, v := range []struct{ src, mode, label string }{{c.Light, "light", "浅色"}, {c.Dark, "dark", "深色"}} {
		fmt.Fprintf(&b, `<figure class="snap snap-%s">`, v.mode)
		if v.src != "" {
			fmt.Fprintf(&b, `<img src="/img/%s" alt="%s %s主题渲染快照" loading="lazy">`, e(filepath.Base(v.src)), e(c.Name), v.label)
		} else {
			b.WriteString(`<div class="nosnap">暂无快照</div>`)
		}
		fmt.Fprintf(&b, `<figcaption>%s · 渲染快照</figcaption></figure>`, v.label)
	}
	b.WriteString(`</section>`)
	b.WriteString(`<article class="doc">` + markdown{link: s.link("docs/components/" + c.Slug + ".md")}.render(c.Doc) + `</article>`)
	b.WriteString(s.pager(c))
	return s.page(c.Name+" "+c.ZhName, c.Slug, b.String())
}

func (s *site) pager(c *component) string {
	var prev, next *component
	for i, x := range s.comps {
		if x == c {
			if i > 0 {
				prev = s.comps[i-1]
			}
			if i+1 < len(s.comps) {
				next = s.comps[i+1]
			}
		}
	}
	var b strings.Builder
	b.WriteString(`<nav class="pager" aria-label="相邻组件">`)
	if prev != nil {
		fmt.Fprintf(&b, `<a href="/c/%s/" rel="prev">← %s</a>`, prev.Slug, html.EscapeString(prev.Name))
	}
	if next != nil {
		fmt.Fprintf(&b, `<a class="next" href="/c/%s/" rel="next">%s →</a>`, next.Slug, html.EscapeString(next.Name))
	}
	b.WriteString(`</nav>`)
	return b.String()
}

type searchItem struct {
	ID   string `json:"id"`
	Slug string `json:"s"`
	Name string `json:"n"`
	Zh   string `json:"z"`
	Cat  string `json:"c"`
}

func (s *site) index() string {
	e := html.EscapeString
	var b strings.Builder
	shots := 0
	for _, c := range s.comps {
		if c.Light != "" || c.Dark != "" {
			shots++
		}
	}
	b.WriteString(`<section class="hero"><p class="eyebrow">Court-gothic components for MyGo</p><h1>MujicaUI</h1>`)
	b.WriteString(`<div class="rule" aria-hidden="true"><span></span></div>`)
	b.WriteString(`<p class="honest"><strong>说明：</strong>这些组件是原生 Go（MyGo）控件，在桌面窗口里由 GPU 自绘。本站展示的是渲染快照图片，浏览器不运行任何组件代码。<br>`)
	b.WriteString(`<strong>Note:</strong> the components are native Go (MyGo) widgets. The images are rendered snapshots; your browser runs no component code.</p>`)
	b.WriteString(`<p class="ctas"><a class="btn" href="/readme/">开始使用</a><a class="btn ghost" href="` + repoURL + `" rel="noopener">源码</a></p></section>`)

	b.WriteString(`<section class="search" aria-label="搜索组件"><label for="q">搜索组件</label><input id="q" type="search" autocomplete="off" placeholder="名称、中文名或编号，如 Button、按钮、042"><p id="qcount" class="qcount" aria-live="polite"></p><ul id="results" class="results"></ul></section>`)

	b.WriteString(`<section id="components" class="cats"><h2>分类</h2><div class="grid">`)
	for _, cat := range s.cats {
		fmt.Fprintf(&b, `<a class="cat" href="#%s"><span class="cname">%s</span><span class="count">%d</span></a>`, e(slug(cat.Name)), e(cat.Name), len(cat.Items))
	}
	b.WriteString(`</div></section>`)
	for _, cat := range s.cats {
		fmt.Fprintf(&b, `<section class="catlist" id="%s"><h2>%s <span class="count">%d</span></h2><ul class="cards">`, e(slug(cat.Name)), e(cat.Name), len(cat.Items))
		for _, c := range cat.Items {
			fmt.Fprintf(&b, `<li><a href="/c/%s/"><span class="num">%s</span><span class="n">%s</span><span class="zh">%s</span></a></li>`, c.Slug, c.ID, e(c.Name), e(c.ZhName))
		}
		b.WriteString(`</ul></section>`)
	}
	items := make([]searchItem, 0, len(s.comps))
	for _, c := range s.comps {
		items = append(items, searchItem{c.ID, c.Slug, c.Name, c.ZhName, c.Category})
	}
	data, err := json.Marshal(items)
	if err != nil {
		panic(err) // plain strings always marshal
	}
	b.WriteString(`<script id="search-data" type="application/json">` + string(data) + `</script>`)
	return s.page("组件", "", b.String())
}

func copyDir(src, dst string, keep func(string) bool) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() || !keep(f.Name()) {
			continue
		}
		if err := copyFile(filepath.Join(src, f.Name()), filepath.Join(dst, f.Name())); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
