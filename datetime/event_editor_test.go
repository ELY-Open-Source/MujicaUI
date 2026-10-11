package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

type editorProbe struct {
	ev             CalendarEvent
	saves, cancels int
	disabled       bool
}

func (p *editorProbe) tester() *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Scroll(c).Fill().Padding(16).Children(func() {
			r := EventEditor(c, &p.ev, EventEditorOptions{Location: calTokyo, Disabled: p.disabled})
			if r.Saved() {
				p.saves++
			}
			if r.Canceled() {
				p.cancels++
			}
		})
	}, 900, 1400)
}

func TestEventEditorValidatesAndSaves(t *testing.T) {
	p := &editorProbe{ev: CalendarEvent{ID: "n", Start: calAt(t, 9, 0, 14), End: calAt(t, 9, 0, 13)}}
	tt := p.tester()
	_ = tt.Click("Save")
	if p.saves != 0 || !tt.HasText("Enter a title.") || !tt.HasText("The end must be after the start.") {
		t.Fatalf("invalid draft saved: %v", tt.Texts())
	}
	_ = tt.Click("Title")
	tt.Type("Masquerade")
	_ = tt.Click("Location")
	tt.Type("Hall")
	_ = tt.Click("Guests")
	tt.Type("Mutsumi")
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("15 minutes before")
	_ = tt.Click("2026-10-13 09:00 JST")
	_ = tt.Click(datetimeTestLabel(oct(14)))
	_ = tt.Click("Time: Hours")
	tt.Key(0, ui.KeyUp)
	_ = tt.Click("Confirm")
	if p.ev.Title != "" {
		t.Fatal("the caller's event changed before Save")
	}
	_ = tt.Click("Save")
	ev := p.ev
	if p.saves != 1 || ev.Title != "Masquerade" || ev.Location != "Hall" || len(ev.Attendees) != 1 || ev.Attendees[0].Name != "Mutsumi" ||
		len(ev.Reminders) != 1 || !ev.End.Equal(calAt(t, 10, 0, 14)) {
		t.Fatalf("saved %+v", ev)
	}
}

func TestEventEditorAllDayRepeatCancel(t *testing.T) {
	p := &editorProbe{ev: calEvent("e", 14, 9, 0, 60)}
	tt := p.tester()
	_ = tt.Click("All day")
	_ = tt.Click("Does not repeat")
	_ = tt.Click("Weekly")
	_ = tt.Click("Save")
	if p.saves != 1 || !p.ev.AllDay || p.ev.Repeat.Unit != RepeatWeekly || p.ev.Repeat.Weekdays[0] != time.Wednesday ||
		!p.ev.Start.Equal(calAt(t, 0, 0, 14)) || !p.ev.End.Equal(calAt(t, 0, 0, 15)) {
		t.Fatalf("saved %+v", p.ev)
	}
	_ = tt.Click("Title")
	tt.Type(" draft")
	_ = tt.Click("Cancel")
	if p.cancels != 1 || p.ev.Title != "Event e" {
		t.Fatalf("cancel kept %q", p.ev.Title)
	}
	tt.Frame()
	_ = tt.Click("Save")
	if p.ev.Title != "Event e" {
		t.Fatal("canceled draft came back")
	}
}

func TestEventEditorDisabled(t *testing.T) {
	p := &editorProbe{ev: calEvent("e", 14, 9, 0, 60), disabled: true}
	tt := p.tester()
	_ = tt.Click("Title")
	tt.Type("x")
	_ = tt.Click("Save")
	_ = tt.Click("Cancel")
	if p.saves != 0 || p.cancels != 0 || p.ev.Title != "Event e" {
		t.Fatal("disabled editor acted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without Location")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		EventEditor(c, &p.ev, EventEditorOptions{})
	}, 10, 10)
}

func TestEventEditorSnapshot(t *testing.T) {
	ev := calRichEvent()
	snapshot(t, "event_editor", 760, 1100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { EventEditor(c, &ev, EventEditorOptions{Location: calTokyo}) })
	})
}
