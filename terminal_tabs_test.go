package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type terminalTabsProbe struct {
	tabs            []TerminalTab
	active          string
	disabled        bool
	changed         int
	renamed, closed []string
}

func (p *terminalTabsProbe) view(c *ui.Context) {
	Use(c, Settings{})
	r := TerminalTabs(c, &p.tabs, &p.active, TerminalTabsOptions{Disabled: p.disabled}, func(id string) {
		ui.Text(c, "Content of "+id)
	})
	if r.Changed {
		p.changed++
	}
	if r.Renamed != "" {
		p.renamed = append(p.renamed, r.Renamed)
	}
	if r.Closed != "" {
		p.closed = append(p.closed, r.Closed)
	}
}

func newTerminalTabsProbe() *terminalTabsProbe {
	return &terminalTabsProbe{tabs: []TerminalTab{{ID: "a", Title: "zsh"}, {ID: "b", Title: "build", Busy: true}, {ID: "c", Title: "logs"}}, active: "a"}
}

func TestTerminalTabsSwitch(t *testing.T) {
	p := newTerminalTabsProbe()
	tt := ui.NewTester(p.view, 600, 400)
	if err := tt.Click("logs"); err != nil {
		t.Fatal(err)
	}
	if p.active != "c" || !tt.HasText("Content of c") || p.changed != 1 {
		t.Fatalf("active %q changed %d", p.active, p.changed)
	}
	tt.Key(0, ui.KeyLeft)
	if p.active != "b" {
		t.Fatalf("Left gave %q", p.active)
	}
}

func TestTerminalTabsRename(t *testing.T) {
	p := newTerminalTabsProbe()
	tt := ui.NewTester(p.view, 600, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyF2)
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("server")
	tt.Key(0, ui.KeyEnter)
	if p.tabs[0].Title != "server" || len(p.renamed) != 1 || p.renamed[0] != "a" {
		t.Fatalf("title %q renamed %v", p.tabs[0].Title, p.renamed)
	}
	// Escape cancels a rename.
	tt.Key(0, ui.KeyF2)
	tt.Type("x")
	tt.Key(0, ui.KeyEscape)
	if p.tabs[0].Title != "server" || !tt.HasText("server") {
		t.Fatalf("Escape kept %q", p.tabs[0].Title)
	}
}

func TestTerminalTabsCloseIdleAndBusy(t *testing.T) {
	p := newTerminalTabsProbe()
	tt := ui.NewTester(p.view, 600, 400)
	if err := tt.Click("Close zsh"); err != nil {
		t.Fatal(err)
	}
	if len(p.tabs) != 2 || p.active != "b" || p.closed[0] != "a" {
		t.Fatalf("after closing idle: tabs %v active %q", p.tabs, p.active)
	}
	if err := tt.Click("Close build"); err != nil {
		t.Fatal(err)
	}
	if len(p.tabs) != 2 || !tt.HasText("Close “build”?") {
		t.Fatalf("busy tab closed without asking: %v", p.tabs)
	}
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if len(p.tabs) != 2 {
		t.Fatal("Cancel closed the tab")
	}
	_ = tt.Click("Close build")
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if len(p.tabs) != 1 || p.active != "c" || len(p.closed) != 2 {
		t.Fatalf("after confirm: tabs %v active %q closed %v", p.tabs, p.active, p.closed)
	}
}

func TestTerminalTabsDisabled(t *testing.T) {
	p := newTerminalTabsProbe()
	p.disabled = true
	tt := ui.NewTester(p.view, 600, 400)
	_ = tt.Click("logs")
	_ = tt.Click("Close zsh")
	if p.active != "a" || len(p.tabs) != 3 {
		t.Fatalf("disabled tabs changed: active %q tabs %d", p.active, len(p.tabs))
	}
}

func TestTerminalTabsSnapshot(t *testing.T) {
	snapshot(t, "terminal_tabs", 460, 120, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", newTerminalTabsProbe)
		ui.Column(c).Padding(12).Children(func() {
			TerminalTabs(c, &(*p).tabs, &(*p).active, TerminalTabsOptions{}, func(id string) {
				codeMono(c, "~/court $ make seal")
			})
		})
	})
}
