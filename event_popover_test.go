package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func calRichEvent() CalendarEvent {
	ev := calEvent("r", 14, 19, 0, 120)
	ev.Title, ev.Location, ev.Color = "Masquerade rehearsal", "Rehearsal hall", EventInfo
	ev.Attendees = calAttendees[:2]
	ev.Reminders = []time.Duration{15 * time.Minute, time.Hour}
	ev.Repeat = Recurrence{Unit: RepeatWeekly, Interval: 1, Weekdays: []time.Weekday{time.Wednesday}}
	return ev
}

type popoverProbe struct {
	open           bool
	opts           EventPopoverOptions
	edits, deletes int
}

func (p *popoverProbe) tester() *ui.Tester {
	p.opts.Location = calTokyo
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(16).Children(func() {
			b := Button(c, "Details", ButtonOptions{Variant: Secondary})
			if b.Clicked() {
				p.open = true
			}
			r := EventPopover(c, b, &p.open, calRichEvent(), p.opts)
			if r.Edit() {
				p.edits++
			}
			if r.Delete() {
				p.deletes++
			}
		})
	}, 600, 500)
}

func TestEventPopoverSummaryAndActions(t *testing.T) {
	p := &popoverProbe{}
	tt := p.tester()
	_ = tt.Click("Details")
	for _, s := range []string{"Masquerade rehearsal", "Wednesday, October 14, 2026 · 19:00 – 21:00", "Rehearsal hall",
		"Weekly on Wednesday", "2 guests: Sakiko Togawa, Mutsumi Wakaba", "15 minutes before, 1 hour before"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %v", s, tt.Texts())
		}
	}
	_ = tt.Click("Edit")
	if p.edits != 1 || p.open {
		t.Fatal("Edit should report and close")
	}
	_ = tt.Click("Details")
	tt.Key(0, ui.KeyTab)
	_ = tt.Click("Delete")
	if p.deletes != 1 || p.open {
		t.Fatal("Delete should report and close")
	}
	_ = tt.Click("Details")
	tt.Key(0, ui.KeyEscape)
	if p.open {
		t.Fatal("Escape should close")
	}
}

func TestEventPopoverReadOnlyAndDisabled(t *testing.T) {
	p := &popoverProbe{opts: EventPopoverOptions{ReadOnly: true}}
	tt := p.tester()
	_ = tt.Click("Details")
	if tt.HasText("Edit") || tt.HasText("Delete") || !tt.HasText("Rehearsal hall") {
		t.Fatal("read-only popover shows actions")
	}
	p.opts.ReadOnly, p.opts.Disabled = false, true
	tt.Frame()
	if p.open || tt.HasText("Rehearsal hall") {
		t.Fatal("disabling should close the popover")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without Location")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		open := true
		EventPopover(c, ui.Box(c), &open, calRichEvent(), EventPopoverOptions{})
	}, 10, 10)
}

func TestEventPopoverSnapshot(t *testing.T) {
	snapshot(t, "event_popover", 460, 400, func(c *ui.Context) {
		ui.Row(c).Padding(16).Children(func() {
			b := Button(c, "Details", ButtonOptions{Variant: Secondary})
			// Open once the anchor is laid out.
			frames := ui.Local(c.Root(), "frames", func() int { return 0 })
			*frames++
			if *frames == 1 {
				c.Invalidate()
			}
			open := *frames > 1
			EventPopover(c, b, &open, calRichEvent(), EventPopoverOptions{Location: calTokyo})
		})
	})
}
