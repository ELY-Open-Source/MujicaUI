package feedback_test

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/feedback"
)

func ExampleUndoToast() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.UndoToast(c, &open, "Note deleted", feedback.UndoToastOptions{
			OnUndo: func() {
				// restore the note
			},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleActivityFeed() {
	events := []feedback.Activity{
		{ID: "1", Actor: "Lin", Action: "commented on", Object: "Quarterly report", Time: time.Now()},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		feedback.ActivityFeed(c, events, feedback.ActivityFeedOptions{}).Size(380, 280)
	}
	ui.Render(view, 400, 300, 1)
}
