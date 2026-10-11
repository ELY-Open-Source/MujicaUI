package feedback

import (
	"cmp"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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

// LoadingOverlay covers content while loading, blocking input to that region only.
func LoadingOverlay(c *ui.Context, loading bool, opts LoadingOverlayOptions, content func()) LoadingOverlayResult {
	if content == nil {
		panic("mujica: LoadingOverlay needs content")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := cmp.Or(opts.Label, core.Msg(c, "loading_overlay.label"))
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
					core.Spinner(c, core.SpinnerOptions{Size: 20, Label: label})
					ui.Text(c, label).TextColor(k.Text)
				})
				if opts.Cancel == "" {
					return
				}
				b := input.Button(c, opts.Cancel, input.ButtonOptions{Variant: input.Secondary})
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
