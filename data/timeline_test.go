package data

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var timelineItems = []TimelineItem{
	{Time: "06:00", Title: "Matins sung", Content: "The choir gathered before dawn.", Status: TimelineDone},
	{Time: "09:00", Title: "Relics moved", Status: TimelineActive,
		Content: "The reliquaries travel from the crypt to the high altar, carried in procession by four acolytes in black and wine."},
	{Time: "12:00", Title: "Bell cracked", Status: TimelineError},
	{Time: "18:00", Title: "Vespers"},
}

func TestTimelineOrderAndStatus(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(260).Children(func() { Timeline(c, timelineItems) })
	}, 300, 600)
	var ys []float32
	for _, title := range []string{"Matins sung", "Relics moved", "Bell cracked", "Vespers"} {
		r, ok := tt.Find(title)
		if !ok {
			t.Fatalf("missing %q", title)
		}
		ys = append(ys, r.Y)
	}
	for i := 1; i < len(ys); i++ {
		if ys[i] <= ys[i-1] {
			t.Fatalf("not top to bottom: %v", ys)
		}
	}
	for _, w := range []string{"Done", "Current", "Error"} {
		if !tt.HasText(w) {
			t.Fatalf("status word %q missing", w)
		}
	}
	long, _ := tt.Find(timelineItems[1].Content)
	if long.H < 3*14 || long.W > 260 {
		t.Fatalf("long content did not wrap: %v", long)
	}
}

func TestTimelineSnapshot(t *testing.T) {
	snapshot(t, "timeline", 420, 380, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { Timeline(c, timelineItems) })
	})
}
