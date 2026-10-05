package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// StatusItemTone colors a status item's text and icon. The zero value is muted text.
type StatusItemTone int

const (
	StatusItemNeutral StatusItemTone = iota
	StatusItemSuccess
	StatusItemWarning
	StatusItemDanger
	StatusItemInfo
)

// StatusItem is one entry of a StatusBar.
type StatusItem struct {
	// ID is the item's stable key and the value StatusBarResult.Clicked reports.
	ID   string
	Text string
	Icon *ui.SVG
	Tone StatusItemTone
	// ShowProgress draws a short bar filled to Progress, from 0 to 1.
	ShowProgress bool
	Progress     float32
	// Action makes the item a button.
	Action bool
	// Label names an item without Text for assistive technology.
	Label string
	// Secondary items hide first, last first, when the bar is too narrow.
	Secondary bool
}

// StatusBarOptions configure a StatusBar.
type StatusBarOptions struct {
	Left, Right []StatusItem
}

// StatusBarResult is a StatusBar and its event in this frame.
type StatusBarResult struct {
	*ui.Element
	// Clicked is the ID of the action item clicked, or empty.
	Clicked string
}

const statusBarProgressWidth float32 = 64

// StatusBar is a window's bottom strip with left and right groups of status items.
func StatusBar(c *ui.Context, opts StatusBarOptions) StatusBarResult {
	k := Tokens(c)
	d := Density(c)
	statusBarCheck(opts)
	h := max(theme.MinHitSize, d.ControlHeight()-8)
	gap := d.Unit() * 2
	e := ui.Row(c).Height(h).Shrink(0).PaddingX(d.Unit()*2).Gap(gap).
		Background(k.Surface).BorderWidth(theme.BorderWidth, 0, 0, 0).BorderColor(k.Border).
		Label(msg(c, "status_bar.label"))
	hidden := statusBarHide(c, e, opts, e.Bounds().W-d.Unit()*4, gap)
	var res StatusBarResult
	group := func(items []StatusItem) {
		ui.Row(c).Gap(gap).Shrink(0).Dividers(theme.BorderWidth, k.Border).Children(func() {
			for _, it := range items {
				if hidden[it.ID] {
					continue
				}
				if statusBarItem(c, it) {
					res.Clicked = it.ID
					slog.Debug("mujica: status_bar clicked", "id", it.ID)
				}
			}
		})
	}
	e.Children(func() {
		group(opts.Left)
		ui.Spacer(c)
		if len(hidden) > 0 {
			var names []string
			for _, it := range append(append([]StatusItem{}, opts.Left...), opts.Right...) {
				if hidden[it.ID] {
					names = append(names, statusBarName(it))
				}
			}
			label := fmt.Sprintf(msg(c, "status_bar.hidden"), len(hidden))
			ui.Text(c, fmt.Sprintf("+%d", len(hidden))).FontSize(theme.CaptionSize).TextColor(k.TextMuted).
				Label(label).Tooltip(strings.Join(names, "\n")).Shrink(0)
		}
		group(opts.Right)
	})
	res.Element = e
	return res
}

func statusBarCheck(opts StatusBarOptions) {
	seen := map[string]bool{}
	for _, it := range append(append([]StatusItem{}, opts.Left...), opts.Right...) {
		switch {
		case it.ID == "":
			panic("mujica: StatusItem needs an ID")
		case seen[it.ID]:
			panic("mujica: duplicate StatusItem ID " + it.ID)
		case it.Progress < 0 || it.Progress > 1:
			panic("mujica: StatusItem Progress must be within 0 and 1")
		case it.Text == "" && it.Label == "":
			panic("mujica: StatusItem " + it.ID + " needs Text or Label")
		}
		seen[it.ID] = true
	}
}

func statusBarName(it StatusItem) string {
	if it.Text != "" {
		return it.Text
	}
	return it.Label
}

// statusBarHide returns the secondary items to hide so the rest fit in room DIPs.
func statusBarHide(c *ui.Context, bar *ui.Element, opts StatusBarOptions, room, gap float32) map[string]bool {
	hidden := map[string]bool{}
	if room <= 0 {
		// Not laid out yet: show all, and measure in the next frame.
		if first := ui.Local(bar, "status_bar.first", func() bool { return true }); *first {
			*first = false
			c.Invalidate()
		}
		return hidden
	}
	all := append(append([]StatusItem{}, opts.Left...), opts.Right...)
	need := func() float32 {
		var w float32
		for _, it := range all {
			if !hidden[it.ID] {
				w += statusBarWidth(c, it) + gap
			}
		}
		if len(hidden) > 0 {
			mw, _ := c.MeasureText(0, ui.Span{Text: fmt.Sprintf("+%d", len(hidden)), Size: theme.CaptionSize})
			w += mw + gap
		}
		return w
	}
	for i := len(all) - 1; i >= 0 && need() > room; i-- {
		if all[i].Secondary {
			hidden[all[i].ID] = true
		}
	}
	return hidden
}

// statusBarWidth estimates the width an item takes.
func statusBarWidth(c *ui.Context, it StatusItem) float32 {
	d := Density(c)
	w := d.Unit() * 2
	if it.Icon != nil {
		w += 14 + d.Unit()
	}
	if it.Text != "" {
		tw, _ := c.MeasureText(0, ui.Span{Text: it.Text, Size: theme.CaptionSize})
		w += tw + d.Unit()
	}
	if it.ShowProgress {
		w += statusBarProgressWidth + d.Unit()
	}
	return w
}

// statusBarItem builds one item and reports whether its action was clicked.
func statusBarItem(c *ui.Context, it StatusItem) bool {
	k := Tokens(c)
	d := Density(c)
	fg := k.TextMuted
	if it.Action {
		fg = k.Text
	}
	switch it.Tone {
	case StatusItemNeutral:
	case StatusItemSuccess:
		fg = k.Success
	case StatusItemWarning:
		fg = k.Warning
	case StatusItemDanger:
		fg = k.Danger
	case StatusItemInfo:
		fg = k.Info
	default:
		panic("mujica: invalid StatusItemTone")
	}
	var e *ui.Element
	if it.Action {
		e = ui.ButtonBase(c).Key(it.ID).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Radius(theme.ControlRadius)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		if e.Pressed() {
			e.Background(k.Selection)
		} else if e.Hovered() {
			e.Background(k.SurfaceHover)
		}
	} else {
		e = ui.Row(c).Key(it.ID).MinHeight(theme.MinHitSize)
	}
	e.PaddingX(d.Unit()).Gap(d.Unit()).TextColor(fg).Shrink(0)
	if it.Label != "" {
		e.Label(it.Label)
	}
	e.Children(func() {
		if it.Icon != nil {
			ui.Icon(c, it.Icon).Size(14, 14)
		}
		if it.Text != "" {
			ui.Text(c, it.Text).FontSize(theme.CaptionSize).SingleLine()
		}
		if it.ShowProgress {
			statusBarProgress(c, it.Progress)
		}
	})
	return it.Action && e.Clicked()
}

func statusBarProgress(c *ui.Context, v float32) {
	k := Tokens(c)
	ui.Box(c).Size(statusBarProgressWidth, 4).Shrink(0).Role(ui.RoleProgress).
		Description(fmt.Sprintf("%.0f%%", v*100)).
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, k.Border, 2)
			if v > 0 {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.W*v, r.H), H: r.H}, k.AccentText, 2)
			}
		})
}
