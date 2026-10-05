package display

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TagOptions configure a Tag; zero is a plain label.
type TagOptions struct {
	// Icon names an icon shown before the label.
	Icon string
	// Clickable makes the label a button reported by TagView.Clicked.
	Clickable bool
	// Closable adds a remove button reported by TagView.Closed.
	Closable bool
	Disabled bool
}

// TagView is a shown Tag and this frame's events.
type TagView struct {
	Element *ui.Element
	// Clicked reports the label was activated, not the remove button.
	Clicked bool
	// Closed reports the remove button, Delete or Backspace.
	Closed bool
}

// Tag shows a label; click and remove are separate stops.
func Tag(c *ui.Context, label string, opts TagOptions) TagView {
	if label == "" {
		panic("mujica: Tag needs a label")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var v TagView
	box := ui.Row(c).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Radius(theme.ControlRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border).Shrink(0).Disabled(opts.Disabled)
	v.Element = box
	box.Children(func() {
		var main *ui.Element
		if opts.Clickable {
			main = ui.ButtonBase(c).Label(label)
			internal.Ring(main, k.Focus, theme.ControlRadius)
		} else {
			main = ui.Row(c).AlignItems(ui.Center)
		}
		main.Gap(u*1.5).Padding(0, u*2).FillHeight().Radius(theme.ControlRadius).TextColor(k.Text)
		if opts.Clickable && !opts.Disabled {
			if main.Pressed() {
				main.Background(k.Selection)
			} else if main.Hovered() {
				main.Background(k.SurfaceHover)
			}
			main.Transition(core.ColorFade)
		}
		if opts.Closable {
			main.Padding(0, u, 0, u*2)
		}
		main.Children(func() {
			if opts.Icon != "" {
				Icon(c, opts.Icon, IconOptions{Size: 14})
			}
			ui.Text(c, label).FontSize(core.FontSize(c, theme.BodySize-1)).SingleLine()
		})
		if opts.Clickable && main.Clicked() {
			v.Clicked = true
		}
		if opts.Closable {
			if opts.Clickable && !main.IsDisabled() && (main.Shortcut(0, ui.KeyDelete) || main.Shortcut(0, ui.KeyBackspace)) {
				v.Closed = true
			}
			if tagClose(c, k, label).Clicked() {
				v.Closed = true
			}
		}
	})
	if v.Clicked {
		slog.Debug("mujica: tag click", "label", label)
	}
	if v.Closed {
		slog.Debug("mujica: tag close", "label", label)
	}
	return v
}

func tagClose(c *ui.Context, k theme.Tokens, label string) *ui.Element {
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).
		Radius(theme.ControlRadius).Label(fmt.Sprintf(core.Msg(c, "tag.remove"), label))
	fg := k.TextMuted
	if b.Pressed() {
		b.Background(k.Selection)
		fg = k.Text
	} else if b.Hovered() {
		b.Background(k.SurfaceHover)
		fg = k.Text
	}
	internal.Ring(b, k.Focus, theme.ControlRadius)
	b.Transition(core.ColorFade)
	b.TextColor(fg).Children(func() {
		Icon(c, "x", IconOptions{Size: 12})
	})
	return b
}
