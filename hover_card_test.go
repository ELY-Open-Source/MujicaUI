package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

const hoverCardTestDelay = 20 * time.Millisecond

type hoverCardApp struct {
	follows, next int
}

func (a *hoverCardApp) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Column(c).Padding(40).Gap(80).AlignItems(ui.Start).Children(func() {
		b := Button(c, "Lady Sakiko", ButtonOptions{Variant: Ghost})
		HoverCard(c, b, HoverCardOptions{OpenDelay: hoverCardTestDelay, CloseDelay: hoverCardTestDelay}, func() {
			ui.Text(c, "Keyboardist of the court")
			if Button(c, "Follow", ButtonOptions{}).Clicked() {
				a.follows++
			}
		})
		if Button(c, "Next", ButtonOptions{}).Clicked() {
			a.next++
		}
	})
}

func hoverCardWait(tt *ui.Tester) {
	time.Sleep(hoverCardTestDelay + 15*time.Millisecond)
	tt.Frame()
}

func TestHoverCardPointerCanEnterCard(t *testing.T) {
	a := &hoverCardApp{}
	tt := ui.NewTester(a.view, 500, 400)
	r, _ := tt.Find("Lady Sakiko")
	tt.Move(r.X+5, r.Y+5)
	if tt.HasText("Keyboardist of the court") {
		t.Fatal("opened before the delay")
	}
	hoverCardWait(tt)
	tt.Frame()
	f, ok := tt.Find("Follow")
	if !ok {
		t.Fatalf("card did not open: %q", tt.Texts())
	}
	// Moving into the card keeps it; its button works.
	tt.Move(f.X+5, f.Y+5)
	hoverCardWait(tt)
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	if a.follows != 1 {
		t.Fatalf("follows %d", a.follows)
	}
	// Leaving both closes it after the delay, and the window stays usable.
	tt.Move(490, 390)
	tt.Frame()
	if !tt.HasText("Follow") {
		t.Fatal("closed without the close delay")
	}
	hoverCardWait(tt)
	if tt.HasText("Follow") {
		t.Fatal("card stayed open")
	}
}

func TestHoverCardKeyboardEntersAndEscapes(t *testing.T) {
	a := &hoverCardApp{}
	tt := ui.NewTester(a.view, 500, 400)
	tt.Key(0, ui.KeyTab)
	hoverCardWait(tt)
	if !tt.HasText("Keyboardist of the court") {
		t.Fatal("no card on keyboard focus")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Follow") {
		t.Fatal("Tab did not enter the card")
	}
	tt.Key(0, ui.KeyEnter)
	if a.follows != 1 {
		t.Fatalf("follows %d", a.follows)
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Follow") || !tt.Focused("Lady Sakiko") {
		t.Fatalf("Escape: card shown %v, anchor focused %v", tt.HasText("Follow"), tt.Focused("Lady Sakiko"))
	}
	hoverCardWait(tt)
	if tt.HasText("Follow") {
		t.Fatal("dismissed card came back while focused")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Next") || a.next != 0 {
		t.Fatal("Tab after the card did not reach Next")
	}
}

func TestHoverCardSnapshot(t *testing.T) {
	snapshot(t, "hover_card", 420, 240, func(c *ui.Context) {
		ui.Column(c).Padding(24).AlignItems(ui.Start).Children(func() {
			Button(c, "Lady Uika", ButtonOptions{Variant: Ghost})
		})
		ui.Overlay(c, func() {
			p := ui.Column(c).Absolute().Left(24).Top(66)
			overlayStyle(c, p, 1)
			p.Width(280).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
				ui.Text(c, "Lady Uika").FontWeight(600)
				ui.Text(c, "Vocalist and guitarist; keeps the ledger of the masquerade.").TextColor(Tokens(c).TextMuted)
				Button(c, "Follow", ButtonOptions{Variant: Outline})
			})
		})
	})
}
