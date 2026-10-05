package navigation

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
)

type toolbarLog struct{ clicks []string }

func (l *toolbarLog) view(width float32) func(c *ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).Width(width).Children(func() {
			Toolbar(c, ToolbarOptions{Label: "Format"}, func() {
				ToolbarGroup(c, func() {
					for _, n := range []string{"Bold", "Italic", "Underline"} {
						if ToolbarButton(c, n, ToolbarButtonOptions{Icon: icons.Must(map[string]string{"Bold": "bold", "Italic": "italic", "Underline": "underline"}[n])}).Clicked() {
							l.clicks = append(l.clicks, n)
						}
					}
				})
				ToolbarSeparator(c)
				if ToolbarButton(c, "Copy", ToolbarButtonOptions{Icon: icons.Must("copy"), ShowLabel: true}).Clicked() {
					l.clicks = append(l.clicks, "Copy")
				}
				if ToolbarButton(c, "Delete", ToolbarButtonOptions{Icon: icons.Must("trash-2"), Disabled: true}).Clicked() {
					l.clicks = append(l.clicks, "Delete")
				}
				if ToolbarButton(c, "Share", ToolbarButtonOptions{}).Clicked() {
					l.clicks = append(l.clicks, "Share")
				}
			})
		})
	}
}

func TestToolbarClickAndArrows(t *testing.T) {
	l := &toolbarLog{}
	tt := ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); l.view(600)(c) }, 700, 100)
	if err := tt.Click("Italic"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight) // from Italic to Underline
	tt.Key(0, ui.KeyRight) // Copy
	tt.Key(0, ui.KeyRight) // Share, skipping disabled Delete
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeySpace)
	_ = tt.Click("Delete")
	if want := []string{"Italic", "Share", "Bold"}; !slices.Equal(l.clicks, want) {
		t.Fatalf("clicks %q, want %q", l.clicks, want)
	}
}

func TestToolbarOverflow(t *testing.T) {
	l := &toolbarLog{}
	// The snapshot's layout: MyGo caches the overflow chevrons by position.
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Gap(16).Children(func() { l.view(180)(c) })
	}, 520, 300)
	tt.Frame()
	if r, ok := tt.Find("Share"); ok && r.W > 0 {
		t.Fatalf("Share should overflow, at %+v", r)
	}
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	menu := tt.Menu()
	if !slices.Contains(menu, "Share") || !slices.Contains(menu, "Copy") {
		t.Fatalf("overflow menu %q", menu)
	}
	if err := tt.ChooseMenuItem("Share"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !slices.Equal(l.clicks, []string{"Share"}) {
		t.Fatalf("overflow choice clicked %q", l.clicks)
	}
}

func TestToolbarButtonNeedsLabel(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ToolbarButton(c, "", ToolbarButtonOptions{Icon: icons.Must("copy")})
	}, 200, 100, 1)
}

func TestToolbarSnapshot(t *testing.T) {
	l := &toolbarLog{}
	snapshot(t, "toolbar", 520, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			l.view(180)(c)
			l.view(420)(c)
		})
	})
}
