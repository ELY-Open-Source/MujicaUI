package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestCapabilityCardsExample(t *testing.T) {
	got := ""
	tt := chatTester(900, 400, func(c *ui.Context) {
		if v := CapabilityCards(c, chatCards(), CapabilityCardsOptions{}); v.Example != "" {
			got = v.Example
		}
	})
	if !tt.HasText("Banquets, seating and schedules.") {
		t.Fatal("description missing")
	}
	if err := tt.Click("“Plan a ball for 200 guests”"); err != nil {
		t.Fatal(err)
	}
	if got != "Plan a ball for 200 guests" {
		t.Fatalf("example %q", got)
	}
	// Tab moves from the clicked example to the next card's.
	got = ""
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if got != "Explain the minuet" {
		t.Fatalf("keys chose %q", got)
	}
}

func TestCapabilityCardsColumnsAndDisabled(t *testing.T) {
	got := ""
	var g *ui.Element
	tt := chatTester(900, 600, func(c *ui.Context) {
		v := CapabilityCards(c, chatCards(), CapabilityCardsOptions{Columns: 1, Disabled: true})
		g = v.Element
		if v.Example != "" {
			got = v.Example
		}
	})
	_ = tt.Click("“Explain the minuet”")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if got != "" {
		t.Fatalf("disabled cards chose %q", got)
	}
	a, _ := tt.Find("Write")
	b, _ := tt.Find("Plan")
	if b.Y <= a.Y || g == nil {
		t.Fatalf("one column should stack cards: %v %v", a, b)
	}
}

func TestCapabilityCardsNeedTitles(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		CapabilityCards(c, []CapabilityCard{{}}, CapabilityCardsOptions{})
	}, 100, 100, 1)
}

func TestCapabilityCardsSnapshot(t *testing.T) {
	snapshot(t, "capability_cards", 860, 200, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { CapabilityCards(c, chatCards(), CapabilityCardsOptions{}) })
	})
}
