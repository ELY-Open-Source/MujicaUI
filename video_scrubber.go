package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// MediaRange is a span of a media timeline.
type MediaRange struct{ Start, End time.Duration }

// VideoScrubberOptions configure a VideoScrubber.
type VideoScrubberOptions struct {
	// Duration is the timeline's length; required.
	Duration time.Duration
	// Buffered spans show as a lighter fill.
	Buffered []MediaRange
	// Preview builds the hover slot's content for a time.
	Preview func(at time.Duration)
	// PreviewWidth, PreviewHeight size the slot; zero means 160×90.
	PreviewWidth, PreviewHeight float32
	// Step is one arrow press; zero means 5 seconds.
	Step     time.Duration
	Disabled bool
	// Label names the timeline for assistive technology.
	Label string
}

// VideoScrubber sets *pos by dragging, clicking, arrows, Home and End.
func VideoScrubber(c *ui.Context, pos *time.Duration, opts VideoScrubberOptions) ChoiceResult {
	if opts.Duration <= 0 || opts.Step < 0 || opts.PreviewWidth < 0 || opts.PreviewHeight < 0 {
		panic("mujica: VideoScrubber needs Duration > 0 and no negative Step or preview size")
	}
	if *pos < 0 || *pos > opts.Duration {
		panic("mujica: VideoScrubber position outside 0..Duration")
	}
	step := opts.Step
	if step == 0 {
		step = 5 * time.Second
	}
	pw, ph := opts.PreviewWidth, opts.PreviewHeight
	if pw == 0 {
		pw, ph = 160, 90
	}
	k := Tokens(c)
	dur := opts.Duration.Seconds()
	g, h := sliderSetup(0, dur, min(step.Seconds(), dur), nil, false, c)
	old := *pos
	var s *ui.Element
	var next float64
	wrap := ui.Box(c).MinWidth(120).Shrink(0).Disabled(opts.Disabled)
	wrap.Children(func() {
		s, next, _ = sliderControl(c, old.Seconds(), g, 0, dur)
		s.Height(h).PaddingX(sliderThumb / 2).FocusRing(false).Cursor(ui.CursorPointer).
			Label(scrubberLabel(c, opts.Label))
		px, _, over := s.PointerPosition()
		r := s.Bounds()
		at := time.Duration(float64((px-sliderThumb/2)/max(r.W-sliderThumb, 1)) * float64(opts.Duration))
		at = max(0, min(opts.Duration, at))
		if over && !s.IsDisabled() {
			scrubberHover(c, k, r, g, at, opts, pw, ph)
		}
	})
	if !s.IsDisabled() {
		*pos = max(0, min(opts.Duration, time.Duration(next*float64(time.Second)).Round(time.Millisecond)))
		if next == dur {
			*pos = opts.Duration
		}
	}
	changed := *pos != old
	if changed {
		slog.Debug("mujica: video_scrubber seek", "pos", *pos)
	}
	v, buffered, focus := pos.Seconds(), opts.Buffered, s
	s.Draw(func(p *ui.Painter, r ui.Rect) {
		y := r.Y + g.trackY
		x0 := g.x(r, 0)
		sliderTrackPaint(p, r, g, k, nil, 0, x0, x0)
		// Buffered spans fill the rail under the played part.
		for _, b := range buffered {
			b0, b1 := g.x(r, b.Start.Seconds()), g.x(r, b.End.Seconds())
			p.Fill(ui.Rect{X: b0, Y: y - sliderTrack/2, W: b1 - b0, H: sliderTrack}, k.TextMuted.Alpha(0.45), 0)
		}
		p.Fill(ui.Rect{X: x0, Y: y - sliderTrack/2, W: g.x(r, v) - x0, H: sliderTrack}, k.AccentText, 0)
		sliderThumbPaint(p, k, g.x(r, v), y, focus.Hovered(), focus.FocusVisible())
	})
	return ChoiceResult{Element: s, changed: changed}
}

func scrubberLabel(c *ui.Context, label string) string {
	if label != "" {
		return label
	}
	return msg(c, "video_scrubber.label")
}

// scrubberHover shows hover time and preview slot above the track.
func scrubberHover(c *ui.Context, k theme.Tokens, r ui.Rect, g sliderGeom, at time.Duration, opts VideoScrubberOptions, pw, ph float32) {
	x := g.x(ui.Rect{W: r.W}, at.Seconds())
	w := pw
	if opts.Preview == nil {
		w = 64
	}
	panel := ui.Column(c).Absolute().Left(max(0, min(r.W-w, x-w/2))).Bottom(r.H - g.trackY + sliderThumb).
		Width(w).Padding(4).Gap(4).AlignItems(ui.Center).PassThrough()
	overlayStyle(c, panel, 1)
	panel.Children(func() {
		if opts.Preview != nil {
			ui.Box(c).Size(pw-8, ph).Clip().Radius(theme.ControlRadius).Background(k.Background).Children(func() {
				opts.Preview(at)
			})
		}
		ui.Text(c, mediaClock(at)).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
	})
}
