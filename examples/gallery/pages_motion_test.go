package main

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestMotionPagesControlsAndEvents(t *testing.T) {
	controls := []string{"Confirm", "Lines", "Tool", "Empty", "Empty", "Disabled", "Loading", "Loading", "Loading", "Loading", "Loading", "Loading", "Add 500", "Append", "Thinking", "Replay", "Reverse", "Expanded", "Show", "Tick"}
	if len(motionPages) != 20 {
		t.Fatalf("motion pages %d", len(motionPages))
	}
	for i, p := range motionPages {
		t.Run(p.ID, func(t *testing.T) {
			if p.Category != "Agent Ops & Motion" || p.Code == "" || p.Keys == "" {
				t.Fatal("incomplete page metadata")
			}
			var events []string
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ui.Column(c).Padding(16).Children(func() { p.Demo(c, func(s string) { events = append(events, s) }) })
			}, 1000, 700)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.SetDark(true)
			events = nil
			if err := tt.Click(controls[i]); err != nil {
				t.Fatal(err)
			}
			if len(events) == 0 {
				t.Fatalf("%s emitted no event", controls[i])
			}
			if p.ID == "294" && !strings.Contains(strings.Join(events, " "), "Finished") {
				t.Fatalf("completion events %q", events)
			}
		})
	}
}

func TestMotionTypewriterRestartReplays(t *testing.T) {
	finished := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		motionTypewriterDemo(c, func(s string) {
			if s == "Finished" {
				finished++
			}
		})
	}, 1000, 700)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	before := finished
	_ = tt.Click("Restart")
	_ = tt.Click("Restart")
	if finished != before+2 {
		t.Fatalf("restart completions %d, want %d", finished, before+2)
	}
}
