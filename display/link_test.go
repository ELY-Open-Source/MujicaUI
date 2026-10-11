package display

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestLinkOpensURL(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Link(c, "Charter", LinkOptions{URL: "https://example.com/charter"})
	}, 300, 100)
	if err := tt.Click("Charter"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	want := []string{"https://example.com/charter", "https://example.com/charter"}
	if got := tt.OpenedURLs(); !slices.Equal(got, want) {
		t.Fatalf("opened %v", got)
	}
}

func TestLinkAction(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Link(c, "Show details", LinkOptions{}).Clicked() {
			n++
		}
	}, 300, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != 1 || len(tt.OpenedURLs()) != 0 {
		t.Fatalf("clicks %d, opened %v", n, tt.OpenedURLs())
	}
}

func TestLinkInvalidURLIsNotOpened(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Link(c, "Broken", LinkOptions{URL: "https://"}).Clicked() {
			n++
		}
	}, 400, 100)
	if err := tt.Click("Broken"); err != nil {
		t.Fatal(err)
	}
	if n != 1 || len(tt.OpenedURLs()) != 0 {
		t.Fatalf("clicks %d, opened %v", n, tt.OpenedURLs())
	}
	if !tt.HasText("This link cannot be opened: https://") {
		t.Fatalf("texts %v", tt.Texts())
	}
}

func TestLinkCallerError(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		Link(c, "手册", LinkOptions{URL: "https://example.com", Error: "没有可以打开此链接的应用"})
	}, 400, 100)
	tt.Frame()
	if !tt.HasText("没有可以打开此链接的应用") {
		t.Fatalf("texts %v", tt.Texts())
	}
}

func TestLinkCheck(t *testing.T) {
	for raw, want := range map[string][2]bool{
		"https://example.com": {true, true},
		"https://":            {true, false},
		"mailto:a@b.c":        {true, true},
		"/notes/42":           {false, true},
		"?tab=info":           {false, true},
		"http://[::1":         {false, false},
	} {
		ext, ok := linkCheck(raw)
		if ext != want[0] || ok != want[1] {
			t.Errorf("%q: external %v ok %v", raw, ext, ok)
		}
	}
}

func TestLinkSnapshot(t *testing.T) {
	snapshot(t, "link", 420, 140, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			Link(c, "Read the court charter", LinkOptions{URL: "https://example.com"})
			Link(c, "Show the guest list", LinkOptions{})
			Link(c, "Broken link", LinkOptions{URL: "https://"})
		})
	})
}
