package mujica_test

import (
	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleAlert() {
	open := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.Alert(c, &open, "The archive could not be read.", mujica.AlertOptions{
			Severity: mujica.SeverityError,
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
		mujica.Use(c, mujica.Settings{})
		if mujica.Button(c, "Delete note", mujica.ButtonOptions{Variant: mujica.Danger}).Clicked() {
			mujica.ShowToast(c, mujica.ToastMessage{ID: "note-1", Message: "Note deleted", Action: "Undo"})
		}
		if r := mujica.Toaster(c); r.Action == "note-1" {
			// restore the note
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProgress() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			mujica.Progress(c, 0.62, mujica.ProgressOptions{Label: "Uploading", ShowValue: true})
			mujica.Progress(c, 0, mujica.ProgressOptions{Shape: mujica.ProgressRing, Indeterminate: true, Label: "Connecting"})
			mujica.Meter(c, 93, mujica.MeterOptions{
				Max: 100, Label: "Disk", ShowValue: true,
				Levels: &ui.MeterLevels{Warning: 70, Critical: 90},
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSpinner() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Spinner(c, mujica.SpinnerOptions{Size: 20, Label: "Syncing", ShowLabel: true})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSkeleton() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			mujica.Skeleton(c, mujica.SkeletonOptions{Shape: mujica.SkeletonCircle})
			ui.Column(c).Grow(1).Children(func() {
				mujica.Skeleton(c, mujica.SkeletonOptions{Lines: 2})
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEmptyState() {
	query := "quarterly"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.EmptyState(c, mujica.EmptyStateOptions{
			Kind:   mujica.EmptyNoResults,
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
		mujica.Use(c, mujica.Settings{})
		r := mujica.Result(c, mujica.ResultOptions{
			Kind:        mujica.ResultForbidden,
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
		mujica.Use(c, mujica.Settings{})
		ui.Row(c).Gap(16).Children(func() {
			mujica.StatusIndicator(c, mujica.StatusSuccess, mujica.StatusIndicatorOptions{Label: "API online"})
			mujica.StatusIndicator(c, mujica.StatusOffline, mujica.StatusIndicatorOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleNotificationCenter() {
	items := []mujica.Notification{
		{ID: "a", Title: "Build failed", Body: "3 tests failing", Time: "2 min ago", Action: "View log"},
		{ID: "b", Title: "Backup complete", Time: "Yesterday", Read: true},
	}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.NotificationCenter(c, &items, mujica.NotificationCenterOptions{})
		if r.Action == "a" {
			// open the build log
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLoadingOverlay() {
	loading := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.LoadingOverlay(c, loading, mujica.LoadingOverlayOptions{Cancel: "Cancel"}, func() {
			ui.Text(c, "Quarterly report")
			mujica.Button(c, "Export", mujica.ButtonOptions{})
		})
		if r.Canceled {
			loading = false
		}
	}
	ui.Render(view, 400, 300, 1)
}
