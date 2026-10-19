package feedback_test

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/feedback"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func ExampleAlert() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.Alert(c, &open, "The archive could not be read.", feedback.AlertOptions{
			Severity: core.SeverityError,
			Title:    "Import failed",
			Action:   "Retry",
		})
		if r.Action {
			// retry the import
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToaster() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Delete note", input.ButtonOptions{Variant: input.Danger}).Clicked() {
			feedback.ShowToast(c, feedback.ToastMessage{ID: "note-1", Message: "Note deleted", Action: "Undo"})
		}
		if r := feedback.Toaster(c); r.Action == "note-1" {
			// restore the note
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProgress() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			feedback.Progress(c, 0.62, feedback.ProgressOptions{Label: "Uploading", ShowValue: true})
			feedback.Progress(c, 0, feedback.ProgressOptions{Shape: feedback.ProgressRing, Indeterminate: true, Label: "Connecting"})
			feedback.Meter(c, 93, feedback.MeterOptions{
				Max: 100, Label: "Disk", ShowValue: true,
				Levels: &ui.MeterLevels{Warning: 70, Critical: 90},
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSkeleton() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			feedback.Skeleton(c, feedback.SkeletonOptions{Shape: feedback.SkeletonCircle})
			ui.Column(c).Grow(1).Children(func() {
				feedback.Skeleton(c, feedback.SkeletonOptions{Lines: 2})
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEmptyState() {
	query := "quarterly"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.EmptyState(c, feedback.EmptyStateOptions{
			Kind:   feedback.EmptyNoResults,
			Query:  query,
			Action: "Clear search",
		})
		if r.Action {
			query = ""
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleResult() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.Result(c, feedback.ResultOptions{
			Kind:        feedback.ResultForbidden,
			Description: "You need the Editor role to change billing.",
			Actions:     []string{"Request access", "Go back"},
		})
		switch r.Action {
		case 0: // send the request
		case 1: // navigate back
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatusIndicator() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Gap(16).Children(func() {
			feedback.StatusIndicator(c, feedback.StatusSuccess, feedback.StatusIndicatorOptions{Label: "API online"})
			feedback.StatusIndicator(c, feedback.StatusOffline, feedback.StatusIndicatorOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleNotificationCenter() {
	items := []feedback.Notification{
		{ID: "a", Title: "Build failed", Body: "3 tests failing", Time: "2 min ago", Action: "View log"},
		{ID: "b", Title: "Backup complete", Time: "Yesterday", Read: true},
	}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.NotificationCenter(c, &items, feedback.NotificationCenterOptions{})
		if r.Action == "a" {
			// open the build log
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLoadingOverlay() {
	loading := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := feedback.LoadingOverlay(c, loading, feedback.LoadingOverlayOptions{Cancel: "Cancel"}, func() {
			ui.Text(c, "Quarterly report")
			input.Button(c, "Export", input.ButtonOptions{})
		})
		if r.Canceled {
			loading = false
		}
	}
	ui.Render(view, 400, 300, 1)
}
