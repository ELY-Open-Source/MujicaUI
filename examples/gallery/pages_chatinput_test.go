package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func chatinputDemoTester(t *testing.T, id string, mode core.Mode) (*ui.Tester, *[]string) {
	t.Helper()
	var p *page
	for i := range chatinputPages {
		if chatinputPages[i].ID == id {
			p = &chatinputPages[i]
			break
		}
	}
	if p == nil {
		t.Fatalf("missing page %s", id)
	}
	events := new([]string)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: mode})
		ui.Scroll(c).Fill().Padding(24).Gap(16).Children(func() { p.Demo(c, func(s string) { *events = append(*events, s) }) })
	}, 900, 800)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	return tt, events
}

func TestChatinputDemosRender(t *testing.T) {
	if len(chatinputPages) != 20 {
		t.Fatalf("pages %d, want 20", len(chatinputPages))
	}
	for _, p := range chatinputPages {
		if p.Category != "Chat Input" || p.Code == "" || p.Keys == "" {
			t.Errorf("incomplete page %s", p.ID)
		}
		for _, mode := range []core.Mode{core.Light, core.Dark} {
			tt, _ := chatinputDemoTester(t, p.ID, mode)
			if tt.HasText("panic") {
				t.Errorf("%s: panic shown", p.ID)
			}
			if dir := os.Getenv("CHATINPUT_DEMO_DIR"); dir != "" {
				f, err := os.Create(filepath.Join(dir, p.ID+"-"+mode.String()+".png"))
				if err != nil {
					t.Fatal(err)
				}
				if err := png.Encode(f, tt.Image()); err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

func TestChatinputDemosEvents(t *testing.T) {
	for _, tc := range []struct{ id, click string }{
		{"221", "Prompt"}, {"222", "Remove guests.csv"}, {"223", "Remove image Picture 1"}, {"225", "Start voice input"},
		{"227", "Command prompt"}, {"228", "Context prompt"}, {"229", "Remove guests.csv"}, {"230", "Model"},
		{"231", "Deep research"}, {"232", "Tools, 1 on"}, {"233", "Reset Temperature"}, {"236", "Send"},
		{"237", "Source 1: The Palace of Versailles"}, {"238", "Open The Palace of Versailles"},
		{"239", "Locate last source"}, {"240", "Searching for “Versailles 1682”"},
	} {
		tt, events := chatinputDemoTester(t, tc.id, core.Light)
		if err := tt.Click(tc.click); err != nil {
			t.Fatalf("%s: %v", tc.id, err)
		}
		switch tc.id {
		case "221":
			tt.Type("hello")
			tt.Key(0, ui.KeyEnter)
		case "227":
			tt.Type("s")
			tt.Key(0, ui.KeyTab)
		case "228":
			tt.Key(0, ui.KeyEnter)
		case "230":
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
		case "232":
			if err := tt.Click("Code execution"); err != nil {
				t.Fatal(err)
			}
		}
		if len(*events) == 0 {
			t.Errorf("%s: action logged no event", tc.id)
		}
	}
}
