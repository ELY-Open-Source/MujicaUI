package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// CollapsibleOptions configure a Collapsible.
type CollapsibleOptions struct {
	// Title is the trigger's text; required unless Trigger is set.
	Title string
	// Trigger builds custom trigger content after the arrow.
	Trigger  func()
	Disabled bool
}

// CollapsibleView is a built Collapsible and what happened to it in this frame.
type CollapsibleView struct {
	*ui.Element
	changed bool
}

// Changed reports that the trigger opened or closed the content.
func (v CollapsibleView) Changed() bool { return v.changed }

// Collapsible shows content below a trigger while *open; closed content is not built.
func Collapsible(c *ui.Context, open *bool, opts CollapsibleOptions, content func()) CollapsibleView {
	if opts.Title == "" && opts.Trigger == nil {
		panic("mujica: Collapsible needs Title or Trigger")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var v CollapsibleView
	v.Element = ui.Column(c).Shrink(0)
	v.Children(func() {
		was := *open
		p := ui.CollapsibleBase(c, open)
		tr := p.Trigger
		if opts.Disabled {
			*open = was
		}
		v.changed = *open != was
		tr.Disabled(opts.Disabled).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(u*2).
			MinHeight(Density(c).ControlHeight()).Padding(0, u*3, 0, u*2).Radius(theme.ControlRadius)
		if opts.Title != "" {
			tr.Label(opts.Title)
		}
		internal.Ring(tr, k.Focus, theme.ControlRadius)
		if tr.Hovered() && !opts.Disabled {
			tr.Background(k.SurfaceHover)
		}
		prog := tr.Animate("mujica.open", collapsibleTarget(*open), Motion(c, theme.StateDuration))
		tr.Children(func() {
			ui.Icon(c, icons.Must("chevron-right")).Size(16, 16).Rotate(90 * prog).TextColor(k.TextMuted)
			if opts.Trigger != nil {
				opts.Trigger()
			} else {
				ui.Text(c, opts.Title).FontWeight(600).SingleLine()
			}
		})
		if panel := collapsiblePanel(c, *open, prog, content); panel != nil {
			panel.Padding(u*2, u*2, u*2, u*8).Gap(u * 2)
		}
	})
	if v.changed {
		slog.Debug("mujica: collapsible toggle", "open", *open)
	}
	return v
}

func collapsibleTarget(open bool) float32 {
	if open {
		return 1
	}
	return 0
}

// collapsiblePanel builds fn while open or animating, clipped to its progress.
func collapsiblePanel(c *ui.Context, open bool, prog float32, fn func()) *ui.Element {
	if !open && prog == 0 {
		return nil
	}
	clip := ui.Column(c).Shrink(0)
	var inner *ui.Element
	clip.Children(func() {
		inner = ui.Column(c).Shrink(0)
		inner.Children(fn)
	})
	if prog < 1 {
		// The content's height from the last frame, grown by progress.
		clip.ClipY().Height(prog * inner.Bounds().H)
	}
	return inner
}
