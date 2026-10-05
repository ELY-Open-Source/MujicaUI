package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestSendButtonSendStopDisabled(t *testing.T) {
	opts := SendButtonOptions{Shortcut: "Enter"}
	sends, stops := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := SendButton(c, opts)
		if r.Sent() {
			sends++
		}
		if r.Stopped() {
			stops++
		}
	}, 200, 100)
	_ = tt.Click("Send")
	opts.Streaming = true
	tt.Frame()
	_ = tt.Click("Stop generating")
	opts = SendButtonOptions{Disabled: true}
	tt.Frame()
	_ = tt.Click("Send")
	opts.Streaming = true
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if sends != 1 || stops != 2 {
		t.Fatalf("sends %d stops %d (disabled blocks send, not stop)", sends, stops)
	}
}

func TestSendButtonSnapshot(t *testing.T) {
	snapshot(t, "send_button", 200, 70, func(c *ui.Context) {
		ui.Row(c).Padding(12).Gap(12).Children(func() {
			SendButton(c, SendButtonOptions{})
			SendButton(c, SendButtonOptions{Streaming: true})
			SendButton(c, SendButtonOptions{Disabled: true})
		})
	})
}
