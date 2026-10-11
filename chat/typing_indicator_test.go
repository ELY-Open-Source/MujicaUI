package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestTypingIndicatorMotionAndStatic(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			TypingIndicator(c, "Herald")
		}, 300, 100)
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
		tt.Frame()
		if _, ok := tt.Find("Herald is typing…"); !ok {
			t.Fatalf("reduce %v: not named", reduce)
		}
		shown := 0
		for _, s := range tt.Texts() {
			if s == "Herald is typing…" {
				shown++
			}
		}
		// The status label always; the static text only under reduced motion.
		if want := map[bool]int{false: 1, true: 2}[reduce]; shown != want {
			t.Errorf("reduce %v: %d labels, want %d", reduce, shown, want)
		}
	}
}

func TestTypingIndicatorNeedsName(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without a name")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		TypingIndicator(c, "")
	}, 100, 100)
}

func TestTypingIndicatorSnapshot(t *testing.T) {
	snapshot(t, "typing_indicator", 240, 60, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() { TypingIndicator(c, "Herald") })
	})
}
