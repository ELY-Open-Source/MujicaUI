package navigation

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Segment is one choice of a SegmentedControl.
type Segment struct {
	// Label is required: it is shown, or with HideLabel names the icon.
	Label     string
	Icon      *ui.SVG
	HideLabel bool
}

// SegmentedControlOptions configure a SegmentedControl.
type SegmentedControlOptions struct {
	// Equal gives every segment the width of the widest.
	Equal bool
	// Label names the group for assistive technology.
	Label    string
	Disabled bool
}

// SegmentedControl chooses one of segments into *selected; Changed reports a new choice.
func SegmentedControl(c *ui.Context, selected *int, segments []Segment, opts SegmentedControlOptions) ui.Element {
	if len(segments) == 0 {
		panic("mujica: SegmentedControl needs at least one segment")
	}
	if *selected < 0 || *selected >= len(segments) {
		panic("mujica: SegmentedControl selected index out of range")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	inset := float32(2)
	prev := *selected
	parts := ui.SegmentedBase(c, selected, len(segments))
	track := parts.Track.Padding(inset).Gap(inset).Radius(theme.ControlRadius+inset).
		Background(k.Background).Border(theme.BorderWidth, k.ControlBorder).AlignItems(ui.Stretch).
		MinHeight(d.ControlHeight()).Disabled(opts.Disabled)
	if opts.Equal {
		track.Grid().Columns(len(segments))
	}
	if opts.Label != "" {
		track.Label(opts.Label)
	}
	track.Children(func() {
		for i, s := range segments {
			if s.Label == "" {
				panic("mujica: Segment needs a Label")
			}
			if s.HideLabel && s.Icon == nil {
				panic("mujica: Segment with HideLabel needs an Icon")
			}
			on := i == *selected
			seg := parts.Segment(i).Padding(0, u*3).Gap(u * 1.5).Radius(theme.ControlRadius).
				MinWidth(theme.MinHitSize).Label(s.Label)
			if !opts.Disabled {
				seg.Cursor(ui.CursorPointer)
			}
			fg := k.TextMuted
			switch {
			case on:
				seg.Background(k.Accent).FontWeight(600)
				fg = k.OnAccent
			case seg.Pressed():
				seg.Background(k.Selection)
				fg = k.Text
			case seg.Hovered():
				seg.Background(k.SurfaceHover)
				fg = k.Text
			}
			internal.Ease(seg, core.Motion(c, theme.StateDuration), fg)
			internal.Ring(seg, k.Focus, theme.ControlRadius)
			if s.HideLabel {
				seg.Tooltip(s.Label)
			}
			seg.Children(func() {
				if s.Icon != nil {
					ui.Icon(c, s.Icon).Size(16, 16)
				}
				if !s.HideLabel {
					ui.Text(c, s.Label).SingleLine()
				}
			})
		}
	})
	if track.Changed() {
		slog.Debug("mujica: segmented control change", "from", prev, "to", *selected)
	}
	return track
}
