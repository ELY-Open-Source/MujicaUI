package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestIntervalSelectorGroupsFavoritesAndDisabled(t *testing.T) {
	value := Interval("1h")
	favorites := []Interval{"1m", "1h"}
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IntervalSelector(c, &value, &favorites, IntervalSelectorOptions{Disabled: disabled})
	}, 500, 700)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("All intervals"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Minutes") || !tt.HasText("Hours") || !tt.HasText("Days and longer") {
		t.Fatalf("groups %q", tt.Texts())
	}
	if err := tt.Click("Add 4h to favorites"); err != nil {
		t.Fatal(err)
	}
	if len(favorites) != 3 || favorites[2] != "4h" {
		t.Fatalf("favorites %v", favorites)
	}
	for range 50 {
		if tt.Focused("1M") {
			break
		}
		tt.Key(0, ui.KeyTab)
	}
	if !tt.Focused("1M") {
		t.Fatal("keyboard could not reach 1M")
	}
	tt.Key(0, ui.KeyEnter)
	if value != "1M" {
		t.Fatalf("interval %v", value)
	}
	if err := tt.Click("All intervals"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Frame()
	if tt.HasText("Days and longer") {
		t.Fatal("disabled popup remained open")
	}
	tt.Key(ui.Alt, ui.KeyDown)
	if value != "1M" || len(favorites) != 3 {
		t.Fatal("disabled interval changed")
	}
}
func TestIntervalSelectorSnapshot(t *testing.T) {
	value := Interval("1h")
	favorites := []Interval{"1m", "5m", "1h", "1D"}
	snapshot(t, "interval_selector", 460, 90, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() { IntervalSelector(c, &value, &favorites, IntervalSelectorOptions{}) })
	})
}
