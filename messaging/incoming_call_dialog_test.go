package messaging

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestIncomingCallDefaultFocusAndEscape(t *testing.T) {
	open := false
	accepted, declined := 0, 0
	opts := IncomingCallOptions{OnAccept: func(bool) { accepted++ }, OnDecline: func() { declined++ }}
	tt := incomingTester(&open, opts, IncomingCall{Caller: "Uika"})
	if err := tt.Click("Ring"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Decline") || !tt.HasText("Uika") {
		t.Fatalf("Decline not focused: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if declined != 1 || open {
		t.Fatalf("Enter on the default declined %d", declined)
	}
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if declined != 2 || accepted != 0 || open {
		t.Fatalf("Escape declined %d", declined)
	}
}

func TestIncomingCallFocusAcceptAndVideo(t *testing.T) {
	open := true
	var videos []bool
	opts := IncomingCallOptions{FocusAccept: true, OnAccept: func(v bool) { videos = append(videos, v) }, OnDecline: func() {}}
	tt := incomingTester(&open, opts, IncomingCall{Caller: "Sakiko", Video: true})
	if !tt.Focused("Accept") {
		t.Fatal("Accept not focused")
	}
	if err := tt.Click("With video"); err != nil {
		t.Fatal(err)
	}
	if len(videos) != 1 || !videos[0] || open {
		t.Fatalf("accepted %v", videos)
	}
	mediaPanics(t, "no caller", func(c *ui.Context) { IncomingCallDialog(c, &open, IncomingCall{}, opts) })
}

func TestIncomingCallSnapshot(t *testing.T) {
	open := true
	opts := IncomingCallOptions{OnAccept: func(bool) {}, OnDecline: func() {}}
	snapshot(t, "incoming_call_dialog", 520, 400, func(c *ui.Context) {
		IncomingCallDialog(c, &open, IncomingCall{Caller: "Sakiko Togawa", Video: true}, opts)
	})
}

func TestIncomingCallDisabledAncestorCloses(t *testing.T) {
	open, disabled, events := true, false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			IncomingCallDialog(c, &open, IncomingCall{Caller: "Caller"}, IncomingCallOptions{
				OnAccept: func(bool) { events++ }, OnDecline: func() { events++ },
			})
		})
	}, 400, 300)
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if open || events != 0 || tt.HasText("Decline") {
		t.Fatalf("disabled dialog: open=%v events=%d", open, events)
	}
}

func TestIncomingCallSmallWindowKeyboardActions(t *testing.T) {
	open, accepted := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IncomingCallDialog(c, &open, IncomingCall{Caller: "Caller", Video: true}, IncomingCallOptions{
			OnAccept: func(bool) { accepted++ }, OnDecline: func() {},
		})
	}, 280, 220)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	r, ok := tt.Find("Accept")
	if !ok || r.W < 28 || r.H < 28 || r.X < 0 || r.Y < 0 || r.X+r.W > 280 || r.Y+r.H > 220 {
		t.Fatalf("accept bounds %+v", r)
	}
	tt.Key(0, ui.KeyEnter)
	if open || accepted != 1 {
		t.Fatalf("accept callbacks=%d open=%v", accepted, open)
	}
}
