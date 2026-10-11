package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

// toolMenuApp holds tools and what the menu reported.
type toolMenuApp struct {
	web, code, mail bool
	authorized      bool
	changed, auth   []string
	off             bool
}

func (a *toolMenuApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(16).AlignItems(ui.Start).Disabled(a.off).Children(func() {
		r := ToolToggleMenu(c, []ToolItem{
			{ID: "web", Label: "Web search", Description: "Search the open web", On: &a.web},
			{ID: "code", Label: "Code execution", On: &a.code},
			{ID: "mail", Label: "Mail", Description: "Read your inbox", On: &a.mail, NeedsAuth: true, Authorized: a.authorized},
		})
		if id, ok := r.Changed(); ok {
			a.changed = append(a.changed, id)
		}
		if id, ok := r.Authorize(); ok {
			a.auth = append(a.auth, id)
		}
	})
}

func TestToolToggleMenuSwitchesAndAuthorization(t *testing.T) {
	a := &toolMenuApp{web: true}
	tt := ui.NewTester(a.view, 500, 500)
	if err := tt.Click("Tools, 1 on"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("Code execution")
	if !a.code || len(a.changed) != 1 || a.changed[0] != "code" {
		t.Fatalf("code %v changed %v", a.code, a.changed)
	}
	if !tt.HasText("Needs authorization") {
		t.Fatalf("auth mark missing: %q", tt.Texts())
	}
	_ = tt.Click("Authorize Mail")
	if a.mail || len(a.auth) != 1 {
		t.Fatalf("mail %v auth %v", a.mail, a.auth)
	}
	a.authorized = true
	tt.Frame()
	_ = tt.Click("Mail")
	if !a.mail {
		t.Fatal("authorized tool did not turn on")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Search the open web") {
		t.Fatal("Escape did not close the panel")
	}
}

func TestToolToggleMenuClosesWhenDisabled(t *testing.T) {
	a := &toolMenuApp{}
	tt := ui.NewTester(a.view, 500, 500)
	_ = tt.Click("Tools, 0 on")
	a.off = true
	tt.Frame()
	if tt.HasText("Search the open web") {
		t.Fatal("panel stayed open under a disabled control")
	}
}

func TestToolToggleMenuSnapshot(t *testing.T) {
	snapshotAfter(t, "tool_toggle_menu", 400, 360, func(c *ui.Context) {
		(*ui.Local(c.Root(), "a", func() *toolMenuApp { return &toolMenuApp{web: true} })).view(c)
	}, func(tt *ui.Tester) { _ = tt.Click("Tools, 1 on") })
}
