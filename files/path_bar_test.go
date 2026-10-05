package files

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func pathBarTester(p *string, opts PathBarOptions, w int, went *[]string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Width(float32(w)).Children(func() {
			if r := PathBar(c, *p, opts); r.Navigated != "" {
				*went = append(*went, r.Navigated)
				*p = r.Navigated
			}
		})
	}, w+20, 300)
}

func TestPathSegments(t *testing.T) {
	segs := pathSegments("/Users/ada/court", "/", "Disk")
	want := []pathSegment{{"Disk", "/"}, {"Users", "/Users"}, {"ada", "/Users/ada"}, {"court", "/Users/ada/court"}}
	if len(segs) != 4 || segs[0] != want[0] || segs[3] != want[3] {
		t.Fatalf("segments %v", segs)
	}
	if s := pathSegments(`C:\Court\Seals`, `\`, ""); len(s) != 3 || s[2].path != `C:\Court\Seals` {
		t.Fatalf("windows segments %v", s)
	}
}

func TestPathBarClickNavigates(t *testing.T) {
	p := "/Users/ada/court/seals"
	var went []string
	tt := pathBarTester(&p, PathBarOptions{}, 600, &went)
	if err := tt.Click("ada"); err != nil {
		t.Fatal(err)
	}
	if len(went) != 1 || went[0] != "/Users/ada" || tt.HasText("seals") {
		t.Fatalf("went %v", went)
	}
	if err := tt.Click("/"); err != nil {
		t.Fatal(err)
	}
	if went[1] != "/" {
		t.Fatalf("root went %v", went)
	}
}

func TestPathBarFoldsMiddle(t *testing.T) {
	p := "/Users/ada/chronicles/of/the/royal/court/seals"
	var went []string
	tt := pathBarTester(&p, PathBarOptions{RootLabel: "Disk"}, 300, &went)
	if !tt.HasText("Disk") || !tt.HasText("seals") || tt.HasText("chronicles") {
		t.Fatalf("narrow bar shows %q", tt.Texts())
	}
	if err := tt.Click("Hidden folders"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("chronicles"); err != nil {
		t.Fatal(err)
	}
	if len(went) != 1 || went[0] != "/Users/ada/chronicles" {
		t.Fatalf("went %v", went)
	}
	if !tt.HasText("chronicles") {
		t.Fatal("shorter path still folded")
	}
}

func TestPathBarKeyboardAndDisabled(t *testing.T) {
	p := "/a/b/c"
	var went []string
	tt := pathBarTester(&p, PathBarOptions{}, 600, &went)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if len(went) != 1 || went[0] != "/a" {
		t.Fatalf("keyboard went %v", went)
	}
	p = "/a/b/c"
	went = nil
	tt = pathBarTester(&p, PathBarOptions{Disabled: true}, 600, &went)
	_ = tt.Click("a")
	if len(went) != 0 {
		t.Fatal("disabled bar navigated")
	}
}

func TestPathBarSnapshot(t *testing.T) {
	snapshot(t, "path_bar", 460, 110, func(c *ui.Context) {
		ui.Column(c).Padding(12).Gap(12).Children(func() {
			PathBar(c, "/Users/ada/court/seals", PathBarOptions{RootLabel: "Disk"})
			ui.Row(c).Width(300).Children(func() {
				PathBar(c, "/Users/ada/chronicles/of/the/royal/court/seals", PathBarOptions{RootLabel: "Disk"})
			})
		})
	})
}
