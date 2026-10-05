package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMessageHeaderFieldsAndStatus(t *testing.T) {
	at := time.Date(2026, 3, 1, 14, 5, 0, 0, time.UTC)
	for status, want := range map[MessageStatus]string{
		MessageSending: "Sending", MessageStreaming: "Generating", MessageSent: "Sent", MessageFailed: "Failed",
	} {
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Locale: EnUS})
			MessageHeader(c, MessageHeaderOptions{Model: "claude", Time: at, Status: status})
		}, 400, 100)
		for _, s := range []string{"claude", "14:05", want} {
			if !tt.HasText(s) {
				t.Errorf("status %v: missing %q in %q", status, s, tt.Texts())
			}
		}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		MessageHeader(c, MessageHeaderOptions{})
	}, 400, 100)
	if len(tt.Texts()) != 0 {
		t.Fatalf("zero header shows %q", tt.Texts())
	}
}

func TestMessageHeaderUnknownStatus(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an unknown status")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		MessageHeader(c, MessageHeaderOptions{Status: 99})
	}, 100, 100)
}

func TestMessageHeaderSnapshot(t *testing.T) {
	at := time.Date(2026, 3, 1, 21, 30, 0, 0, time.UTC)
	snapshot(t, "message_header", 360, 150, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			for _, s := range []MessageStatus{MessageSending, MessageStreaming, MessageSent, MessageFailed} {
				MessageHeader(c, MessageHeaderOptions{Model: "herald-2", Time: at, Status: s})
			}
		})
	})
}
