package mujica

import (
	"fmt"
	"log/slog"

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
	Element *ui.Element
	// Clicked is the ID of the action item clicked, or empty.
	Clicked string
}

const (
	statusBarProgressWidth float32 = 64
	statusBarIcon          float32 = 16
)

// StatusBar is a window's bottom strip with left and right groups of status items.
func StatusBar(c *ui.Context, opts StatusBarOptions) StatusBarResult {
	k := Tokens(c)
	d := Density(c)
	statusBarCheck(opts)
	gap := d.Unit() * 2
	e := ui.Row(c).MinHeight(max(theme.MinHitSize, d.ControlHeight()-8)).Shrink(0).PaddingX(d.Unit()*2).Gap(gap).
		Background(k.Surface).BorderWidth(theme.BorderWidth, 0, 0, 0).BorderColor(k.Border).
		Label(msg(c, "status_bar.label"))
	was := e.Bounds()
	internal.Remeasure(c, e, was.W, was.H)
	hidden := statusBarHide(c, opts, was.W-d.Unit()*4, gap)
	open := ui.Local(e, "status_bar.open", func() bool { return false })
	*open = *open && len(hidden) > 0
	// focused is the bar item that had the focus last frame, for moving it to "+N" when that item folds.
	focused := ui.Local(e, "status_bar.focused", func() string { return "" })
	wasFocused := *focused
	*focused = ""
	var res StatusBarResult
	item := func(it StatusItem) bool {
		clicked, hasFocus := statusBarItem(c, it)
		if clicked {
			res.Clicked = it.ID
			*open = false
			slog.Debug("mujica: status_bar clicked", "id", it.ID)
		}
		return hasFocus
	}
	group := func(key string, items []StatusItem) {
		ui.Row(c).Key(key).Gap(gap).Shrink(0).Dividers(theme.BorderWidth, k.Border).Children(func() {
			for _, it := range items {
				if !hidden[it.ID] && item(it) {
					*focused = it.ID
				}
			}
		})
	}
	e.Children(func() {
		group("status_bar.left", opts.Left)
		ui.Spacer(c).Key("status_bar.spacer")
		if len(hidden) > 0 {
			more := statusBarMore(c, len(hidden))
			if wasFocused != "" && hidden[wasFocused] {
				more.Focus()
				slog.Debug("mujica: status_bar focus folded", "id", wasFocused)
			}
			if more.Clicked() {
				*open = !*open
				slog.Debug("mujica: status_bar hidden items", "open", *open)
			}
			if *open {
				ui.PopoverBase(c, more, open, func(panel *ui.Element) {
					overlayStyle(c, panel, overlayEnter(c, panel))
					panel.Margin(d.Unit(), 0, 0, 0).Padding(d.Unit()).Gap(d.Unit()).AlignItems(ui.Stretch)
					for _, it := range statusBarAll(opts) {
						if hidden[it.ID] {
							item(it)
						}
					}
				})
			}
		}
		group("status_bar.right", opts.Right)
	})
	res.Element = e
	return res
}

// statusBarMore is the "+N" button that opens the hidden items.
func statusBarMore(c *ui.Context, n int) *ui.Element {
	k := Tokens(c)
	b := ui.ButtonBase(c).Key("status_bar.more").MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).
		PaddingX(Density(c).Unit()).Radius(theme.ControlRadius).TextColor(k.Text).
		Label(fmt.Sprintf(msg(c, "status_bar.hidden"), n))
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Pressed() {
		b.Background(k.Selection)
	} else if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	return b.Children(func() {
		ui.Text(c, fmt.Sprintf("+%d", n)).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
	})
}

func statusBarAll(opts StatusBarOptions) []StatusItem {
	return append(append([]StatusItem{}, opts.Left...), opts.Right...)
}

func statusBarCheck(opts StatusBarOptions) {
	seen := map[string]bool{}
	for _, it := range statusBarAll(opts) {
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

// statusBarHide returns the secondary items to hide, last first, so the rest fit in room DIPs.
func statusBarHide(c *ui.Context, opts StatusBarOptions, room, gap float32) map[string]bool {
	hidden := map[string]bool{}
	if room <= 0 {
		// Not laid out yet: show all; Remeasure rebuilds once measured.
		return hidden
	}
	u := Density(c).Unit()
	group := func(items []StatusItem) float32 {
		var w float32
		n := 0
		for _, it := range items {
			if !hidden[it.ID] {
				w += statusBarWidth(c, it)
				n++
			}
		}
		return w + gap*float32(max(n-1, 0))
	}
	// The bar's children: left group, spacer, the "+N" button when shown, right group.
	need := func() float32 {
		w := group(opts.Left) + group(opts.Right) + 2*gap
		if len(hidden) > 0 {
			tw, _ := c.MeasureText(0, ui.Span{Text: fmt.Sprintf("+%d", len(hidden)), Size: FontSize(c, theme.CaptionSize)})
			w += max(theme.MinHitSize, tw+2*u) + gap
		}
		return w
	}
	all := statusBarAll(opts)
	for i := len(all) - 1; i >= 0 && need() > room; i-- {
		if all[i].Secondary {
			hidden[all[i].ID] = true
		}
	}
	return hidden
}

// statusBarWidth is the width an item takes, as statusBarItem lays it out.
func statusBarWidth(c *ui.Context, it StatusItem) float32 {
	u := Density(c).Unit()
	var parts []float32
	if it.Icon != nil {
		parts = append(parts, statusBarIcon)
	}
	if it.Text != "" {
		tw, _ := c.MeasureText(0, ui.Span{Text: it.Text, Size: FontSize(c, theme.CaptionSize)})
		parts = append(parts, tw)
	}
	if it.ShowProgress {
		parts = append(parts, statusBarProgressWidth)
	}
	w := 2 * u
	for i, p := range parts {
		w += p
		if i > 0 {
			w += u
		}
	}
	if it.Action {
		w = max(w, theme.MinHitSize)
	}
	return w
}

// statusBarItem builds one item and reports whether its action was clicked and whether it has the focus.
func statusBarItem(c *ui.Context, it StatusItem) (clicked, focused bool) {
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
			ui.Icon(c, it.Icon).Size(statusBarIcon, statusBarIcon)
		}
		if it.Text != "" {
			ui.Text(c, it.Text).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
		}
		if it.ShowProgress {
			statusBarProgress(c, it.Progress)
		}
	})
	return it.Action && e.Clicked(), e.FocusWithin()
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
