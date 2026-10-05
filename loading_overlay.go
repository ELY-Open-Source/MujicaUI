package mujica

import (
	"cmp"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// LoadingOverlayOptions configure a LoadingOverlay.
type LoadingOverlayOptions struct {
	// Label explains the wait; empty uses "Loading…".
	Label string
	// Cancel labels a button that reports Canceled; empty shows none.
	Cancel string
}

// LoadingOverlayResult is what a LoadingOverlay built and reported this frame.
type LoadingOverlayResult struct {
	Element *ui.Element
	// Canceled reports the cancel button or Escape while loading.
	Canceled bool
}

// LoadingOverlay builds content and, while loading, covers it with a
// spinner and label. Covered content takes no pointer, keys or focus;
// the rest of the window stays usable. Focus inside moves to the overlay.
func LoadingOverlay(c *ui.Context, loading bool, opts LoadingOverlayOptions, content func()) LoadingOverlayResult {
	if content == nil {
		panic("mujica: LoadingOverlay needs content")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	label := cmp.Or(opts.Label, msg(c, "loading_overlay.label"))
	var res LoadingOverlayResult
	box := ui.Box(c)
	was := ui.Local(box, "loading_overlay.loading", func() bool { return false })
	if loading != *was {
		slog.Debug("mujica: loading_overlay changed", "loading", loading)
		if loading {
			c.Announce(label)
		}
		*was = loading
	}
	box.Children(func() {
		inner := ui.Column(c).FillWidth().Disabled(loading)
		inner.Children(content)
		if !loading {
			return
		}
		held := inner.FocusWithin()
		scrim := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Padding(u * 4).
			Background(k.Background.Alpha(0.5)).Radius(theme.CardRadius)
		scrim.Transition(ui.ElementTransition{Duration: theme.StateDuration, Enter: &ui.Motion{}, Exit: &ui.Motion{}})
		scrim.Children(func() {
			panel := ui.Column(c).AlignItems(ui.Center).Gap(u*3).Padding(u*4, u*6).Radius(theme.PopupRadius).
				Background(k.Surface).Border(theme.BorderWidth, k.Border).Role(ui.RoleStatus).Label(label)
			if opts.Cancel == "" {
				// Without a button the panel keeps the focus in the region.
				internal.Ring(panel.Focusable(), k.Focus, theme.PopupRadius)
				if held {
					panel.Focus()
				}
			}
			panel.Children(func() {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					Spinner(c, SpinnerOptions{Size: 20, Label: label})
					ui.Text(c, label).TextColor(k.Text)
				})
				if opts.Cancel == "" {
					return
				}
				b := Button(c, opts.Cancel, ButtonOptions{Variant: Secondary})
				if held {
					b.Focus()
				}
				if b.Clicked() {
					res.Canceled = true
				}
			})
			if opts.Cancel != "" && panel.Shortcut(0, ui.KeyEscape) {
				res.Canceled = true
			}
		})
	})
	if res.Canceled {
		slog.Debug("mujica: loading_overlay canceled")
	}
	res.Element = box
	return res
}
