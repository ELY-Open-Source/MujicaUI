package media

import (
	"log/slog"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// AudioWaveformOptions configure an AudioWaveform.
type AudioWaveformOptions struct {
	// Duration is the audio's length; required when data is ready.
	Duration time.Duration
	// Position draws the played part; OnSeek reports clicks and arrows.
	Position time.Duration
	OnSeek   func(time.Duration)
	// Selection, when set, is chosen by dragging or Shift+arrows.
	Selection *MediaRange
	// Height is the chart's height; zero means 96.
	Height float32
	Status data.DataStatus
	Error  string
	Label  string
}

// waveformPeaks caches column peaks until data or width change.
type waveformPeaks struct {
	first *float32
	n     int
	cols  int
	peaks []float32
}

type waveformKey struct{}

type waveformPress struct {
	on, dragged bool
	x           float32
}

// waveformBuckets counts peak rebuilds, for tests.
var waveformBuckets int

// AudioWaveform draws samples with played part and selectable range.
func AudioWaveform(c *ui.Context, samples []float32, opts AudioWaveformOptions) core.ChoiceResult {
	h := opts.Height
	if h < 0 || opts.Position < 0 {
		panic("mujica: AudioWaveform Height and Position must not be negative")
	}
	if h == 0 {
		h = 96
	}
	k := core.Tokens(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "audio_waveform.label")
	}
	e := ui.Box(c).Height(h).MinWidth(120).Radius(theme.ControlRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Focusable().Label(label).Cursor(ui.CursorText)
	internal.Ring(e, k.Focus, theme.ControlRadius)
	v := core.ChoiceResult{Element: e}
	if mediaChartState(c, e, opts.Status, opts.Error, len(samples) == 0, "audio_waveform.empty") {
		return v
	}
	if opts.Duration <= 0 || opts.OnSeek == nil || opts.Position > opts.Duration {
		panic("mujica: AudioWaveform needs Duration > 0, OnSeek and Position <= Duration")
	}
	if s := opts.Selection; s != nil && (s.Start < 0 || s.End < s.Start || s.End > opts.Duration) {
		panic("mujica: AudioWaveform Selection outside 0..Duration")
	}
	r := e.Bounds()
	cols := int(r.W)
	cache := ui.Local(e, waveformKey{}, func() waveformPeaks { return waveformPeaks{} })
	if cols <= 0 {
		// Peaks need the width: lay out first, then build again.
		c.AnimationFrame()
	} else if cache.first != &samples[0] || cache.n != len(samples) || cache.cols != cols {
		*cache = waveformPeaks{&samples[0], len(samples), cols, waveformBucket(samples, cols)}
	}
	if !e.IsDisabled() {
		v = core.NewChoiceResult(v.Element, waveformInput(c, e, r, opts))
	}
	e.Description(core.MediaClock(opts.Position) + " / " + core.MediaClock(opts.Duration))
	peaks, pos, sel := cache.peaks, float32(opts.Position)/float32(opts.Duration), opts.Selection
	dur := float32(opts.Duration)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		mid := r.Y + r.H/2
		amp := r.H/2 - 6
		if sel != nil && sel.End > sel.Start {
			x0, x1 := r.X+r.W*float32(sel.Start)/dur, r.X+r.W*float32(sel.End)/dur
			p.Fill(ui.Rect{X: x0, Y: r.Y, W: x1 - x0, H: r.H}, k.Selection, 0)
			p.Fill(ui.Rect{X: x0, Y: r.Y, W: 1, H: r.H}, k.AccentText, 0)
			p.Fill(ui.Rect{X: x1 - 1, Y: r.Y, W: 1, H: r.H}, k.AccentText, 0)
		}
		played := r.X + r.W*pos
		for i, pk := range peaks {
			x := r.X + r.W*float32(i)/float32(len(peaks))
			col := k.TextMuted
			if x < played {
				col = k.AccentText
			}
			ph := max(1, amp*pk)
			p.Fill(ui.Rect{X: x, Y: mid - ph, W: max(1, r.W/float32(len(peaks))-0.5), H: 2 * ph}, col, 0)
		}
		px := internal.Snap(p, played)
		p.Fill(ui.Rect{X: px - 1, Y: r.Y, W: 2, H: r.H}, k.Text, 0)
		p.FillPath(internal.Diamond(px, internal.Snap(p, r.Y+5), 8), k.Text)
	})
	return v
}

// waveformBucket returns the largest magnitude in each of cols spans.
func waveformBucket(samples []float32, cols int) []float32 {
	waveformBuckets++
	peaks := make([]float32, cols)
	for i, s := range samples {
		col := i * cols / len(samples)
		peaks[col] = max(peaks[col], float32(math.Abs(float64(s))))
	}
	return peaks
}

// waveformInput seeks and selects; true means the selection changed.
func waveformInput(c *ui.Context, e ui.Element, r ui.Rect, opts AudioWaveformOptions) bool {
	at := func(x float32) time.Duration {
		f := max(0, min(1, (x-r.X)/max(1, r.W)))
		return time.Duration(float64(f) * float64(opts.Duration)).Round(time.Millisecond)
	}
	seek := func(t time.Duration) {
		t = max(0, min(opts.Duration, t))
		slog.Debug("mujica: audio_waveform seek", "pos", t)
		opts.OnSeek(t)
	}
	sel := opts.Selection
	old := MediaRange{}
	if sel != nil {
		old = *sel
	}
	px, _, _ := e.PointerPosition()
	px += r.X
	st := ui.Local(e, waveformPress{}, func() waveformPress { return waveformPress{} })
	_, _, dragging := e.Dragged()
	switch {
	case e.Pressed() && !st.on:
		*st = waveformPress{on: true, x: px}
	case e.Pressed() && sel != nil && (dragging || st.dragged) && math.Abs(float64(px-st.x)) > 3:
		st.dragged = true
		a, b := at(st.x), at(px)
		*sel = MediaRange{min(a, b), max(a, b)}
	case !e.Pressed() && st.on:
		if !st.dragged {
			seek(at(st.x))
		}
		*st = waveformPress{}
	}
	const step = 5 * time.Second
	switch {
	case e.Shortcut(0, ui.KeyLeft):
		seek(opts.Position - step)
	case e.Shortcut(0, ui.KeyRight):
		seek(opts.Position + step)
	case sel != nil && e.Shortcut(ui.Shift, ui.KeyRight):
		*sel = waveformExtend(*sel, opts.Position, step, opts.Duration)
	case sel != nil && e.Shortcut(ui.Shift, ui.KeyLeft):
		*sel = waveformExtend(*sel, opts.Position, -step, opts.Duration)
	case sel != nil && e.Shortcut(0, ui.KeyEscape):
		*sel = MediaRange{}
	}
	changed := sel != nil && *sel != old
	if changed {
		slog.Debug("mujica: audio_waveform select", "start", sel.Start, "end", sel.End)
	}
	return changed
}

// waveformExtend starts a selection at pos or moves its end.
func waveformExtend(s MediaRange, pos, d, dur time.Duration) MediaRange {
	if s.End <= s.Start {
		s = MediaRange{pos, pos}
	}
	s.End = max(s.Start, min(dur, s.End+d))
	return s
}
