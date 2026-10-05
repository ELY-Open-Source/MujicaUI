package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

func commandPaletteSample() []Command {
	return []Command{
		{ID: "new", Label: "New document", Group: "File", Icon: icons.Must("file")},
		{ID: "open", Label: "Open recent", Group: "File", Icon: icons.Must("folder")},
		{ID: "theme", Label: "Toggle dark mode", Group: "View", Keywords: []string{"theme", "night"}, Icon: icons.Must("moon")},
		{ID: "zoom", Label: "Zoom in", Group: "View", Icon: icons.Must("zoom-in")},
		{ID: "save", Label: "Save all", Group: "File"},
	}
}

type paletteApp struct {
	open bool
	ran  []string
}

func (a *paletteApp) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Column(c).Fill().Children(func() {
		ui.Text(c, "Editor")
		if id, ok := CommandPalette(c, &a.open, commandPaletteSample(), CommandPaletteOptions{}).Chosen(); ok {
			a.ran = append(a.ran, id)
		}
	})
}

func TestCommandPaletteFilter(t *testing.T) {
	ids := func(cs []Command) (out []string) {
		for _, c := range cs {
			out = append(out, c.ID)
		}
		return
	}
	got := ids(commandPaletteFilter(commandPaletteSample(), ""))
	// Groups keep first-appearance order: File's Save joins its group.
	want := []string{"new", "open", "save", "theme", "zoom"}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
	if got := ids(commandPaletteFilter(commandPaletteSample(), "NIGHT")); len(got) != 1 || got[0] != "theme" {
		t.Fatalf("keyword search %v", got)
	}
	if got := ids(commandPaletteFilter(commandPaletteSample(), "file save")); len(got) != 1 || got[0] != "save" {
		t.Fatalf("multi-word search %v", got)
	}
}

func TestCommandPaletteShortcutSearchRun(t *testing.T) {
	a := &paletteApp{}
	tt := ui.NewTester(a.view, 800, 600)
	tt.Key(ui.Cmd, ui.KeyK)
	if !a.open || !tt.HasText("Zoom in") {
		t.Fatalf("Cmd+K did not open: %q", tt.Texts())
	}
	tt.Type("zo")
	if tt.HasText("New document") || !tt.HasText("Zoom in") {
		t.Fatalf("search did not filter: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if a.open || len(a.ran) != 1 || a.ran[0] != "zoom" {
		t.Fatalf("Enter: open %v ran %q", a.open, a.ran)
	}
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.HasText("New document") {
		t.Fatal("query not reset on reopen")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp) // wraps to the last
	tt.Key(0, ui.KeyEnter)
	if len(a.ran) != 2 || a.ran[1] != "zoom" {
		t.Fatalf("arrows ran %q", a.ran)
	}
}

func TestCommandPaletteEmptyClickEscape(t *testing.T) {
	a := &paletteApp{open: true}
	tt := ui.NewTester(a.view, 800, 600)
	tt.Type("qqq")
	if !tt.HasText("No matching commands") {
		t.Fatalf("no empty state: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if !a.open || len(a.ran) != 0 {
		t.Fatal("Enter with no results ran something")
	}
	tt.Key(0, ui.KeyEscape)
	if a.open {
		t.Fatal("Escape did not close")
	}
	a.open = true
	tt.Frame()
	if err := tt.Click("Open recent"); err != nil {
		t.Fatal(err)
	}
	if a.open || len(a.ran) != 1 || a.ran[0] != "open" {
		t.Fatalf("click: open %v ran %q", a.open, a.ran)
	}
}

func TestCommandPaletteInvalid(t *testing.T) {
	for name, cmds := range map[string][]Command{
		"id":  {{Label: "A"}},
		"dup": {{ID: "a", Label: "A"}, {ID: "a", Label: "B"}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			open := true
			ui.Render(func(c *ui.Context) {
				Use(c, Settings{})
				CommandPalette(c, &open, cmds, CommandPaletteOptions{})
			}, 300, 100, 1)
		})
	}
}

func TestCommandPaletteSnapshot(t *testing.T) {
	open := true
	snapshot(t, "command_palette", 640, 460, func(c *ui.Context) {
		CommandPalette(c, &open, commandPaletteSample(), CommandPaletteOptions{})
	})
}
