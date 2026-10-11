package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// feedbackSection titles a part of a demo.
func feedbackSection(c *ui.Context, title string, fn func()) {
	ui.Column(c).Gap(8).FillWidth().Children(func() {
		ui.Text(c, title).FontSize(theme.CaptionSize).TextColor(core.Tokens(c).TextMuted)
		fn()
	})
}

var feedbackSeverities = []core.Severity{core.SeverityInfo, core.SeveritySuccess, core.SeverityWarning, core.SeverityError}

func alertDemo(c *ui.Context, log func(string)) {
	type state struct {
		sev           int
		title, action bool
		closable      bool
		open          bool
	}
	s := ui.Local(c.Root(), "091", func() state { return state{sev: 3, title: true, action: true, closable: true, open: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Segmented(c, &s.sev, "Info", "Success", "Warning", "Error")
			ui.Row(c).Gap(16).Children(func() {
				ui.Checkbox(c, &s.title, "Title")
				ui.Checkbox(c, &s.action, "Action")
				ui.Checkbox(c, &s.closable, "Closable")
			})
		})
		feedbackSection(c, "Live", func() {
			opts := feedback.AlertOptions{Severity: feedbackSeverities[s.sev]}
			if s.title {
				opts.Title = "Import failed"
			}
			if s.action {
				opts.Action = "Retry"
			}
			var open *bool
			if s.closable {
				open = &s.open
			}
			r := feedback.Alert(c, open, "The archive could not be read.", opts)
			if r.Action {
				log("alert: action Retry")
			}
			if r.Closed {
				log("alert: closed")
			}
			if s.closable && !s.open && input.Button(c, "Show again", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				s.open = true
				log("alert: reopened")
			}
		})
		feedbackSection(c, "Severities", func() {
			for _, sev := range feedbackSeverities {
				feedback.Alert(c, nil, "Severity "+sev.String()+" uses its own icon, bar and name.", feedback.AlertOptions{Severity: sev}).Element.Key(sev)
			}
		})
	})
}

func toastDemo(c *ui.Context, log func(string)) {
	type state struct {
		sev     int
		action  bool
		seconds float64
		n       int
		note    string
	}
	s := ui.Local(c.Root(), "092", func() state { return state{seconds: 5} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Segmented(c, &s.sev, "Info", "Success", "Warning", "Error")
			ui.Checkbox(c, &s.action, "Undo action")
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Duration %.0f s (errors stay until closed)", s.seconds))
				ui.Slider(c, &s.seconds, 1, 10).Width(200)
			})
		})
		feedbackSection(c, "Actions", func() {
			ui.Row(c).Gap(12).Wrap().Children(func() {
				if input.Button(c, "Show toast", input.ButtonOptions{}).Clicked() {
					s.n++
					m := feedback.ToastMessage{ID: fmt.Sprint(s.n), Severity: feedbackSeverities[s.sev], Message: fmt.Sprintf("Message %d", s.n)}
					if s.action {
						m.Action = "Undo"
					}
					if m.Severity != core.SeverityError {
						m.Duration = time.Duration(s.seconds * float64(time.Second))
					}
					feedback.ShowToast(c, m)
					log("toast: show " + m.ID)
				}
				if input.Button(c, "Show five at once", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					for range 5 {
						s.n++
						feedback.ShowToast(c, feedback.ToastMessage{ID: fmt.Sprint(s.n), Message: fmt.Sprintf("Queued %d", s.n)})
					}
					log("toast: five queued, three visible")
				}
				if input.Button(c, "Repeat same ID", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					s.n++
					feedback.ShowToast(c, feedback.ToastMessage{ID: "sync", Message: fmt.Sprintf("Synced %d times", s.n)})
					log("toast: dedupe id sync")
				}
			})
			ui.Text(c, "Type here while toasts appear; the focus stays:")
			ui.TextInput(c, &s.note).Width(320)
		})
	})
	if r := feedback.Toaster(c); r.Action != "" {
		log("toast: action " + r.Action)
	} else if r.Closed != "" {
		log("toast: closed " + r.Closed)
	}
}

func progressDemo(c *ui.Context, log func(string)) {
	type state struct {
		value, meter     float64
		indet, showValue bool
		ring             int
	}
	s := ui.Local(c.Root(), "093", func() state { return state{value: 0.62, meter: 72, showValue: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Segmented(c, &s.ring, "Linear", "Ring")
			ui.Row(c).Gap(16).Children(func() {
				if ui.Checkbox(c, &s.indet, "Indeterminate").Changed() {
					log(fmt.Sprintf("progress: indeterminate %v", s.indet))
				}
				ui.Checkbox(c, &s.showValue, "Show value")
			})
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Value")
				if ui.Slider(c, &s.value, 0, 1).Width(240).Changed() {
					log(fmt.Sprintf("progress: value %.2f", s.value))
				}
			})
		})
		feedbackSection(c, "Progress", func() {
			feedback.Progress(c, s.value, feedback.ProgressOptions{
				Shape: feedback.ProgressShape(s.ring), Indeterminate: s.indet, Label: "Uploading", ShowValue: s.showValue,
			})
		})
		feedbackSection(c, "Meter with thresholds 70 and 90", func() {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Level")
				if ui.Slider(c, &s.meter, 0, 100).Width(240).Changed() {
					log(fmt.Sprintf("meter: value %.0f", s.meter))
				}
			})
			v := float64(int(s.meter))
			feedback.Meter(c, v, feedback.MeterOptions{Max: 100, Label: "Disk", ShowValue: true, Levels: &ui.MeterLevels{Warning: 70, Critical: 90}})
			feedback.Meter(c, 100-v, feedback.MeterOptions{Max: 100, Label: "Battery (low is bad)", ShowValue: true, Levels: &ui.MeterLevels{Warning: 30, Critical: 10}})
		})
	})
}

func spinnerDemo(c *ui.Context, log func(string)) {
	type state struct {
		size              float64
		showLabel, active bool
	}
	s := ui.Local(c.Root(), "094", func() state { return state{size: 20, showLabel: true, active: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Size %.0f", s.size))
				ui.Slider(c, &s.size, 12, 48).Width(200)
			})
			ui.Row(c).Gap(16).Children(func() {
				ui.Checkbox(c, &s.showLabel, "Show label")
				if ui.Checkbox(c, &s.active, "Loading (unchecked stops the animation)").Changed() {
					log(fmt.Sprintf("spinner: loading %v", s.active))
				}
			})
		})
		feedbackSection(c, "Live", func() {
			if s.active {
				core.Spinner(c, core.SpinnerOptions{Size: float32(s.size), Label: "Syncing", ShowLabel: s.showLabel})
			} else {
				ui.Text(c, "Done")
			}
		})
		feedbackSection(c, "In a button", func() {
			input.Button(c, "Saving", input.ButtonOptions{Loading: true})
		})
	})
}

func skeletonDemo(c *ui.Context, log func(string)) {
	type state struct {
		lines   float64
		loading bool
	}
	s := ui.Local(c.Root(), "095", func() state { return state{lines: 3, loading: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Lines %.0f", s.lines))
				ui.Slider(c, &s.lines, 1, 6).Width(200)
			})
			if ui.Checkbox(c, &s.loading, "Loading").Changed() {
				log(fmt.Sprintf("skeleton: loading %v", s.loading))
			}
		})
		feedbackSection(c, "Card (same size loading and loaded)", func() {
			ui.Row(c).Gap(12).Padding(16).Width(420).Background(core.Tokens(c).Surface).Radius(theme.CardRadius).Children(func() {
				if s.loading {
					feedback.Skeleton(c, feedback.SkeletonOptions{Shape: feedback.SkeletonCircle})
					ui.Column(c).Grow(1).Children(func() {
						feedback.Skeleton(c, feedback.SkeletonOptions{Lines: int(s.lines)})
					})
					return
				}
				ui.Box(c).Size(40, 40).Radius(20).Background(core.Tokens(c).Selection)
				ui.Column(c).Grow(1).Children(func() {
					for i := range int(s.lines) {
						ui.Text(c, fmt.Sprintf("Loaded line %d", i+1)).LineHeight(theme.LineHeight)
					}
				})
			})
		})
		feedbackSection(c, "Rect", func() {
			feedback.Skeleton(c, feedback.SkeletonOptions{Shape: feedback.SkeletonRect, Width: 420, Height: 96})
		})
	})
}

func emptyStateDemo(c *ui.Context, log func(string)) {
	type state struct {
		kind         int
		query        string
		icon, action bool
	}
	s := ui.Local(c.Root(), "096", func() state { return state{kind: 1, query: "quarterly", icon: true, action: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Segmented(c, &s.kind, "No data", "No results")
			ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
				ui.Checkbox(c, &s.icon, "Icon")
				ui.Checkbox(c, &s.action, "Action")
				if s.kind == 1 {
					ui.TextInput(c, &s.query).Width(200)
				}
			})
		})
		feedbackSection(c, "Live", func() {
			opts := feedback.EmptyStateOptions{Kind: feedback.EmptyKind(s.kind), NoIcon: !s.icon}
			if s.kind == 1 {
				opts.Query = s.query
			}
			if s.action {
				opts.Action = []string{"New project", "Clear search"}[s.kind]
			}
			if feedback.EmptyState(c, opts).Action {
				log("empty state: action " + opts.Action)
				s.query = ""
			}
		})
	})
}

func resultDemo(c *ui.Context, log func(string)) {
	type state struct{ kind int }
	s := ui.Local(c.Root(), "097", func() state { return state{kind: 2} })
	descs := []string{"Your order has been placed.", "The payment was declined by the bank.", "You need the Editor role to change billing."}
	actions := [][]string{{"View order", "Keep shopping"}, {"Try again", "Back to cart"}, {"Request access", "Go back"}}
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Segmented(c, &s.kind, "Success", "Failure", "Forbidden")
		})
		r := feedback.Result(c, feedback.ResultOptions{Kind: feedback.ResultKind(s.kind), Description: descs[s.kind], Actions: actions[s.kind]})
		if r.Action >= 0 {
			log("result: action " + actions[s.kind][r.Action])
		}
	})
}

func statusDemo(c *ui.Context, log func(string)) {
	type state struct {
		status int
		custom bool
	}
	s := ui.Local(c.Root(), "098", func() state { return state{} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			if ui.Segmented(c, &s.status, "Success", "Warning", "Error", "Offline", "Busy").Changed() {
				log("status: " + feedback.Status(s.status).String())
			}
			ui.Checkbox(c, &s.custom, "Custom label")
		})
		feedbackSection(c, "Live", func() {
			opts := feedback.StatusIndicatorOptions{}
			if s.custom {
				opts.Label = "Database replica"
			}
			feedback.StatusIndicator(c, feedback.Status(s.status), opts)
		})
		feedbackSection(c, "All states (shape and word, not color alone)", func() {
			ui.Row(c).Gap(20).Wrap().Children(func() {
				for st := feedback.StatusSuccess; st <= feedback.StatusBusy; st++ {
					feedback.StatusIndicator(c, st, feedback.StatusIndicatorOptions{}).Key(st)
				}
			})
		})
	})
}

func notificationDemo(c *ui.Context, log func(string)) {
	type state struct {
		items []feedback.Notification
		n     int
	}
	s := ui.Local(c.Root(), "099", func() state {
		return state{items: []feedback.Notification{
			{ID: "a", Title: "Build failed", Body: "main · 3 tests failing", Time: "2 min ago", Action: "View log"},
			{ID: "b", Title: "Invite accepted", Body: "Mira joined the workspace", Time: "1 h ago"},
			{ID: "c", Title: "Backup complete", Time: "Yesterday", Read: true},
		}}
	})
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Caller-owned list", func() {
			if input.Button(c, "Add notification", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				s.n++
				id := fmt.Sprintf("new-%d", s.n)
				s.items = append([]feedback.Notification{{ID: id, Title: fmt.Sprintf("Comment %d on the design doc", s.n), Time: "now"}}, s.items...)
				log("notifications: added " + id)
			}
		})
		ui.Column(c).Width(520).Height(420).Children(func() {
			r := feedback.NotificationCenter(c, &s.items, feedback.NotificationCenterOptions{})
			for _, ev := range [][2]string{{"action", r.Action}, {"toggled", r.Toggled}, {"opened", r.Opened}, {"removed", r.Removed}} {
				if ev[1] != "" {
					log("notifications: " + ev[0] + " " + ev[1])
				}
			}
			if r.MarkedAllRead {
				log("notifications: marked all read")
			}
			if r.ClearedAll {
				log("notifications: cleared all")
			}
		})
	})
}

func loadingOverlayDemo(c *ui.Context, log func(string)) {
	type state struct {
		loading, cancel bool
		saves           int
	}
	s := ui.Local(c.Root(), "100", func() state { return state{loading: true, cancel: true} })
	ui.Column(c).Gap(16).FillWidth().Children(func() {
		feedbackSection(c, "Parameters", func() {
			ui.Row(c).Gap(16).Children(func() {
				if ui.Checkbox(c, &s.loading, "Loading").Changed() {
					log(fmt.Sprintf("overlay: loading %v", s.loading))
				}
				ui.Checkbox(c, &s.cancel, "Cancel action")
			})
		})
		feedbackSection(c, "Covered region", func() {
			opts := feedback.LoadingOverlayOptions{Label: "Generating report…"}
			if s.cancel {
				opts.Cancel = "Cancel"
			}
			r := feedback.LoadingOverlay(c, s.loading, opts, func() {
				ui.Column(c).Gap(8).Padding(16).FillWidth().AlignItems(ui.Start).Background(core.Tokens(c).Surface).Radius(theme.CardRadius).Children(func() {
					ui.Text(c, "Quarterly report").Font(theme.SerifFont).FontSize(theme.H3Size)
					ui.Text(c, fmt.Sprintf("Exported %d times.", s.saves))
					if input.Button(c, "Export", input.ButtonOptions{}).Clicked() {
						s.saves++
						log("overlay: export clicked")
					}
				})
			})
			r.Element.Width(420)
			if r.Canceled {
				s.loading = false
				log("overlay: canceled")
			}
		})
		feedbackSection(c, "Outside the region (always usable)", func() {
			if input.Button(c, "Outside button", input.ButtonOptions{Variant: input.Secondary}).AlignSelf(ui.Start).Clicked() {
				log("overlay: outside clicked")
			}
		})
	})
}
