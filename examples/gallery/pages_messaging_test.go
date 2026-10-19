package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func messagingDemoTester(t *testing.T, id string, mode core.Mode) (*ui.Tester, *[]string) {
	t.Helper()
	i := slices.IndexFunc(messagingPages, func(p page) bool { return p.ID == id })
	if i < 0 {
		t.Fatalf("missing page %s", id)
	}
	events := &[]string{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: mode})
		ui.Scroll(c).Fill().Padding(24).Gap(12).Children(func() {
			messagingPages[i].Demo(c, func(s string) { *events = append(*events, s) })
		})
	}, 900, 800)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	return tt, events
}

func TestMessagingDemosInteractAndLog(t *testing.T) {
	cases := []struct{ id, click, event string }{
		{"401", "random, 120 unread", "channel random"},
		{"402", "Like", "reaction true"},
		{"403", "Unread count", "unread 7"},
		{"404", "Read", "receipt 0"},
		{"405", "Sakiko", "member 1"},
		{"406", "Message", "Message"},
		{"407", "Presence", "presence 3"},
		{"408", "Status text", "status changed"},
		{"409", "Mute", "mute true"},
		{"410", "Speaking", "Speaking = false"},
		{"411", "Incoming call", "Declined"},
		{"412", "grinning face", "emoji "},
		{"413", "Unpin Uika's message", "unpin b"},
		{"414", "Reply in thread", "reply sent"},
		{"415", "Star Second act", "star 1"},
		{"416", "Download score.pdf", "download attachment 0"},
		{"417", "Attach", "attachment added"},
		{"418", "Recipients", "recipients 1"},
		{"419", "Show quoted text (2 lines)", "quote expanded true"},
		{"420", "Tomorrow", "snooze "},
	}
	if len(messagingPages) != len(cases) {
		t.Fatalf("pages=%d, cases=%d", len(messagingPages), len(cases))
	}
	for _, tc := range cases {
		for _, mode := range []core.Mode{core.Light, core.Dark} {
			t.Run(tc.id+mode.String(), func(t *testing.T) {
				tt, events := messagingDemoTester(t, tc.id, mode)
				if err := tt.Click(tc.click); err != nil {
					t.Fatal(err)
				}
				switch tc.id {
				case "403":
					tt.Command("selectAll")
					tt.Type("7")
					tt.Key(0, ui.KeyEnter)
				case "404":
					if err := tt.Click("Sent"); err != nil {
						t.Fatal(err)
					}
				case "407":
					tt.Key(0, ui.KeyDown)
					tt.Key(0, ui.KeyDown)
					tt.Key(0, ui.KeyEnter)
				case "408":
					tt.Type("Rehearsing")
				case "411":
					tt.Key(0, ui.KeyEnter)
				case "414":
					tt.Type("Second act")
					tt.Key(0, ui.KeyEnter)
				case "418":
					tt.Type("guest@example.com")
					tt.Key(0, ui.KeyEnter)
				}
				if !slices.ContainsFunc(*events, func(s string) bool { return strings.HasPrefix(s, tc.event) }) {
					t.Fatalf("want %q, events %v", tc.event, *events)
				}
			})
		}
	}
}

func TestMessagingDemoCaptures(t *testing.T) {
	dir := os.Getenv("MESSAGING_DEMO_DIR")
	if dir == "" {
		t.Skip("set MESSAGING_DEMO_DIR to export demos")
	}
	for _, p := range messagingPages {
		if p.Category != "Messaging & Mail" || p.Code == "" || p.Keys == "" {
			t.Fatalf("incomplete page %s", p.ID)
		}
		for _, mode := range []core.Mode{core.Light, core.Dark} {
			tt, _ := messagingDemoTester(t, p.ID, mode)
			f, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s-%s.png", p.ID, mode)))
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
