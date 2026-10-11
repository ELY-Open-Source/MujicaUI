package messaging

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestThreadPanelReplyAndClose(t *testing.T) {
	at := time.Date(2026, 10, 5, 20, 0, 0, 0, time.UTC)
	root := ChatEntry{ID: "r", Author: "Sakiko", Time: at, Text: "Set list for Friday?"}
	replies := []ChatEntry{{ID: "1", Author: "Uika", Time: at.Add(time.Minute), Text: "Overture first."}}
	open := true
	var sent []string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Fill().Children(func() {
			ui.Box(c).Grow(1)
			if r := ThreadPanel(c, &open, root, replies, ThreadPanelOptions{}); r.Sent != "" {
				sent = append(sent, r.Sent)
			}
		})
	}, 600, 600)
	if !tt.HasText("Set list for Friday?") || !tt.HasText("1 replies") || !tt.HasText("Overture first.") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if err := tt.Click("Reply in thread"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Then the masquerade.")
	tt.Key(0, ui.KeyEnter)
	if len(sent) != 1 || sent[0] != "Then the masquerade." {
		t.Fatalf("sent %q", sent)
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not close")
	}
}

func TestThreadPanelSnapshot(t *testing.T) {
	at := time.Date(2026, 10, 5, 20, 0, 0, 0, time.UTC)
	open := true
	root := ChatEntry{ID: "r", Author: "Sakiko", Time: at, Text: "Set list for Friday?"}
	replies := []ChatEntry{{ID: "1", Author: "Uika", Time: at.Add(time.Minute), Text: "Overture first."},
		{ID: "2", Author: "Uika", Time: at.Add(2 * time.Minute), Text: "Then the masquerade."}}
	snapshot(t, "thread_panel", 400, 560, func(c *ui.Context) {
		ui.Row(c).Fill().Children(func() { ThreadPanel(c, &open, root, replies, ThreadPanelOptions{}) })
	})
}

func TestThreadPanelFitsWindowAndDisabled(t *testing.T) {
	open, disabled := true, false
	root := ChatEntry{ID: "root", Author: "Author", Text: "Question"}
	sent := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Fill().Children(func() {
			if ThreadPanel(c, &open, root, nil, ThreadPanelOptions{Disabled: disabled}).Sent != "" {
				sent++
			}
		})
	}, 280, 400)
	r, ok := tt.Find("Thread")
	if !ok || r.X < 0 || r.X+r.W > 280 {
		t.Fatalf("thread bounds %+v", r)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Reply in thread")
	tt.Type("Blocked")
	tt.Key(0, ui.KeyEnter)
	if open || sent != 0 {
		t.Fatal("disabled thread stayed open or sent")
	}
}
