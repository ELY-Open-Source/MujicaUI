package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFixture(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildFixture(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "docs/components/001-oauth-buttons.md", "# 001 OAuthButtons 第三方登录按钮\n\n## 参数\n\n"+
		"| 字段 | 类型 |\n| --- | --- |\n| `Value` | `*T` |\n| `Pipe` | `a \\| b` |\n\n"+
		"- 一项，见 [兼容](../compat.md)\n- 二项\n\n```go\nx := a < b\n```\n")
	writeFixture(t, root, "examples/gallery/pages_fixture.go", "package main\n\nconst fixCategory = \"Fixtures\"\n\n"+
		"var fixPages = []page{\n\t{ID: \"001\", Name: \"OAuthButtons\", ZhName: \"第三方登录按钮\", Category: fixCategory, Code: fixCode[\"x\"]},\n}\n")
	writeFixture(t, root, "pkg/auth/testdata/snapshots/darwin/oauth_buttons-light.png", "png")
	writeFixture(t, root, "pkg/auth/testdata/snapshots/darwin/oauth_buttons-dark.png", "png")
	writeFixture(t, root, "README.md", "# MujicaUI\n\n见 [兼容](docs/compat.md)。\n")
	writeFixture(t, root, "docs/compat.md", "# 兼容\n")
	writeFixture(t, root, "docs/upstream.md", "# 上游\n")
	writeFixture(t, root, "site/style.css", "body{}")
	writeFixture(t, root, "site/app.js", "")

	out := filepath.Join(root, "site", "dist")
	n, err := build(root, out)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("pages = %d, want 1", n)
	}
	raw, err := os.ReadFile(filepath.Join(out, "c", "001-oauth-buttons", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	page := string(raw)
	for _, want := range []string{
		"<th>字段</th><th>类型</th>",
		"<td><code>Value</code></td><td><code>*T</code></td>",
		"<td><code>a | b</code></td>",
		`<a href="/compat/">兼容</a>`,
		`<code class="lang-go">x := a &lt; b</code>`,
		`class="copy"`,
		`<li><a href="/c/001-oauth-buttons/" aria-current="page">OAuthButtons`,
		`src="/img/oauth_buttons-light.png"`,
		`src="/img/oauth_buttons-dark.png"`,
		"Fixtures",
	} {
		if !strings.Contains(page, want) {
			t.Errorf("component page missing %q", want)
		}
	}
	for _, rel := range []string{"index.html", "404.html", "readme/index.html", "compat/index.html", "upstream/index.html", "img/oauth_buttons-dark.png", "style.css", "app.js"} {
		if _, err := os.Stat(filepath.Join(out, rel)); err != nil {
			t.Errorf("missing %s: %v", rel, err)
		}
	}
	index, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(index), `"n":"OAuthButtons","z":"第三方登录按钮"`) {
		t.Error("index search data missing the component")
	}
}
