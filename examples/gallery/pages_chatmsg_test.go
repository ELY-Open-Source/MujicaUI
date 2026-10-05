package main

import (
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

// chatmsgDemoTester runs the demo of page id and returns its tester and event log.
func chatmsgDemoTester(t *testing.T, id string, mode mujica.Mode) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(chatmsgPages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("no page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{Mode: mode})
		ui.Scroll(c).Fill().Padding(24).Gap(16).Children(func() {
			chatmsgPages[i].Demo(c, func(s string) { *events = append(*events, s) })
		})
	}, 900, 800)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	return tt, events
}

func TestChatmsgDemosRender(t *testing.T) {
	if len(chatmsgPages) != 20 {
		t.Fatalf("%d chat message pages, want 20", len(chatmsgPages))
	}
	dir := os.Getenv("CHATMSG_DEMO_DIR")
	for _, p := range chatmsgPages {
		for _, m := range []mujica.Mode{mujica.Light, mujica.Dark} {
			tt, _ := chatmsgDemoTester(t, p.ID, m)
			if len(tt.Texts()) == 0 {
				t.Errorf("%s: demo shows nothing", p.ID)
			}
			if dir == "" {
				continue
			}
			f, err := os.Create(filepath.Join(dir, p.ID+"-"+m.String()+".png"))
			if err != nil {
				t.Fatal(err)
			}
			if err := png.Encode(f, tt.Image()); err != nil {
				t.Fatal(err)
			}
			f.Close()
		}
	}
}

func TestChatmsgDemosLogEvents(t *testing.T) {
	for id, click := range map[string]string{
		"202": "Add message", "207": "Copy code", "214": "Regenerate", "216": "Next branch", "217": "Retry",
	} {
		tt, events := chatmsgDemoTester(t, id, mujica.Light)
		if err := tt.Click(click); err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if len(*events) == 0 {
			t.Errorf("%s: clicking %q logged nothing", id, click)
		}
	}
}
