package overlay

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

const hoverCardTestDelay = 150 * time.Millisecond

type hoverCardApp struct {
	follows, next int
	inDialog      bool
}

func (a *hoverCardApp) card(c *ui.Context) {
	b := input.Button(c, "Lady Sakiko", input.ButtonOptions{Variant: input.Ghost})
	HoverCard(c, b, HoverCardOptions{OpenDelay: hoverCardTestDelay, CloseDelay: hoverCardTestDelay}, func() {
		ui.Text(c, "Keyboardist of the court")
		if input.Button(c, "Follow", input.ButtonOptions{}).Clicked() {
			a.follows++
		}
	})
	if input.Button(c, "Next", input.ButtonOptions{}).Clicked() {
		a.next++
	}
}

func (a *hoverCardApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	if a.inDialog {
		open := true
		Dialog(c, &open, DialogOptions{Title: "Court", HideClose: true}, func() {
			ui.Column(c).Gap(80).AlignItems(ui.Start).Children(func() { a.card(c) })
		})
		return
	}
	ui.Column(c).Padding(40).Gap(80).AlignItems(ui.Start).Children(func() { a.card(c) })
}

func hoverCardShown(tt *ui.Tester) func() bool {
	return func() bool { return tt.HasText("Keyboardist of the court") }
}

// hoverCardSettle waits past both delays so any pending open or close has happened.
func hoverCardSettle(tt *ui.Tester) {
	tt.Advance(hoverCardTestDelay + 30*time.Millisecond)
}

func TestHoverCardPointerCanEnterCard(t *testing.T) {
	a := &hoverCardApp{}
	tt := ui.NewTester(a.view, 500, 400)
	r, _ := tt.Find("Lady Sakiko")
	start := time.Now()
	tt.Move(r.X+5, r.Y+5)
	if time.Since(start) < hoverCardTestDelay && tt.HasText("Keyboardist of the court") {
		t.Fatal("opened before the delay")
	}
	if !overlayUntil(tt, hoverCardShown(tt)) {
		t.Fatalf("card did not open: %q", tt.Texts())
	}
	// The pointer inside the card keeps it open past the close delay.
	k, _ := tt.Find("Keyboardist of the court")
	tt.Move(k.X+5, k.Y+5)
	hoverCardSettle(tt)
	if !tt.HasText("Follow") {
		t.Fatal("card closed with the pointer inside it")
	}
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	if a.follows != 1 || !tt.Focused("Follow") {
		t.Fatalf("follows %d, Follow focused %v", a.follows, tt.Focused("Follow"))
	}
	// With the focus inside, leaving keeps it; a click outside closes it.
	tt.Move(490, 390)
	hoverCardSettle(tt)
	if !tt.HasText("Follow") {
		t.Fatal("card closed with the focus inside it")
	}
	tt.ClickAt(490, 390)
	if tt.HasText("Follow") {
		t.Fatal("outside click did not close the card")
	}
}

func TestHoverCardCloseDelay(t *testing.T) {
	a := &hoverCardApp{}
	tt := ui.NewTester(a.view, 500, 400)
	r, _ := tt.Find("Lady Sakiko")
	tt.Move(r.X+5, r.Y+5)
	if !overlayUntil(tt, hoverCardShown(tt)) {
		t.Fatal("card did not open")
	}
	start := time.Now()
	tt.Move(490, 390)
	if time.Since(start) < hoverCardTestDelay && !tt.HasText("Follow") {
		t.Fatal("closed without the close delay")
	}
	if !overlayUntil(tt, func() bool { return !tt.HasText("Follow") }) {
		t.Fatal("card stayed open after the close delay")
	}
}

func TestHoverCardKeyboardEntersAndEscapes(t *testing.T) {
	a := &hoverCardApp{}
	tt := ui.NewTester(a.view, 500, 400)
	tt.Key(0, ui.KeyTab)
	if !overlayUntil(tt, hoverCardShown(tt)) {
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
	hoverCardSettle(tt)
	if tt.HasText("Follow") {
		t.Fatal("dismissed card came back while focused")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Next") || a.next != 0 {
		t.Fatal("Tab after the card did not reach Next")
	}
}

func TestHoverCardInsideDialog(t *testing.T) {
	a := &hoverCardApp{inDialog: true}
	tt := ui.NewTester(a.view, 600, 500)
	r, _ := tt.Find("Lady Sakiko")
	tt.Move(r.X+5, r.Y+5)
	if !overlayUntil(tt, hoverCardShown(tt)) {
		t.Fatal("card did not open in the dialog")
	}
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	if a.follows != 1 || !tt.Focused("Follow") {
		t.Fatalf("follows %d, focus kept on Follow %v", a.follows, tt.Focused("Follow"))
	}
	tt.Key(ui.Shift, ui.KeyTab)
	if !tt.Focused("Lady Sakiko") {
		t.Fatal("Shift+Tab from the card did not reach its anchor in the dialog")
	}
}

func TestHoverCardSnapshot(t *testing.T) {
	snapshotAfter(t, "hover_card", 420, 240, func(c *ui.Context) {
		ui.Column(c).Padding(24).AlignItems(ui.Start).Children(func() {
			b := input.Button(c, "Lady Uika", input.ButtonOptions{Variant: input.Ghost})
			HoverCard(c, b, HoverCardOptions{OpenDelay: time.Millisecond}, func() {
				ui.Column(c).Gap(8).AlignItems(ui.Start).Children(func() {
					ui.Text(c, "Lady Uika").FontWeight(600)
					ui.Text(c, "Vocalist and guitarist; keeps the ledger of the masquerade.").TextColor(core.Tokens(c).TextMuted)
					input.Button(c, "Follow", input.ButtonOptions{Variant: input.Outline})
				})
			})
		})
	}, func(tt *ui.Tester) {
		r, _ := tt.Find("Lady Uika")
		tt.Move(r.X+5, r.Y+5)
		if !overlayUntil(tt, func() bool { return tt.HasText("Follow") }) {
			t.Fatal("card did not open")
		}
	})
}
