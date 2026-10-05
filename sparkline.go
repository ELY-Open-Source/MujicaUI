package mujica

import (
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// SparklineKind is how a Sparkline draws its values.
type SparklineKind int

const (
	SparklineLine SparklineKind = iota
	SparklineBar
	SparklineArea
)

// SparklineOptions configure a Sparkline; Frame.Height 0 means 28 DIPs.
type SparklineOptions struct {
	Frame ChartFrame
	Kind  SparklineKind
	// Width is the plot width; 0 means 96 DIPs.
	Width  float32
	Format func(v float64) string
}

// Sparkline is an inline chart without axes.
func Sparkline(c *ui.Context, values []float64, opts SparklineOptions) ChartResult {
	f := opts.Frame
	if f.Label == "" || f.Height < 0 || opts.Width < 0 || opts.Kind > SparklineArea {
		panic("mujica: Sparkline needs a Label, sizes ≥ 0 and a known Kind")
	}
	if f.Height == 0 {
		f.Height = 28
	}
	w := opts.Width
	if w == 0 {
		w = 96
	}
	k := Tokens(c)
	u := Density(c).Unit()
	col := internal.Palette(k)[0]
	size := FontSize(c, theme.CaptionSize)
	lo, hi := math.Inf(1), math.Inf(-1)
	last := -1
	for i, v := range values {
		if !math.IsNaN(v) {
			lo, hi, last = min(lo, v), max(hi, v), i
		}
	}
	res := ChartResult{Point: -1}
	res.Element = ui.Row(c).Gap(u * 2).AlignItems(ui.Center)
	st := ui.Local(res.Element, "mujica.sparkline", func() chartState { return chartState{focus: -1, last: -1} })
	res.Element.Children(func() {
		if f.Status == DataLoading || f.Status == DataFailed || last < 0 {
			box := ui.Row(c).Size(w, f.Height).AlignItems(ui.Center).Gap(u)
			box.Children(func() {
				switch f.Status {
				case DataLoading:
					Spinner(c, SpinnerOptions{Size: 14, Label: msg(c, "list.loading")})
				case DataFailed:
					ui.Text(c, dataError(c, f.Error)).FontSize(size).TextColor(k.Danger).SingleLine().Grow(1).MinWidth(0)
					res.retry = iconAction(c, icons.Must("refresh-cw"), msg(c, "list.retry")).Clicked()
				default:
					empty := f.Empty
					if empty == "" {
						empty = msg(c, "chart.empty")
					}
					ui.Text(c, empty).FontSize(size).TextColor(k.TextMuted).SingleLine()
				}
			})
			st.focus = -1
			return
		}
		if lo == hi {
			lo, hi = lo-1, hi+1
		}
		if opts.Kind == SparklineBar {
			lo, hi = min(lo, 0), max(hi, 0)
		}
		plot := ui.Box(c).Size(w, f.Height).Shrink(0).Focusable().Label(f.Label).Role(ui.RoleGroup)
		internal.Ring(plot, k.Focus, theme.ControlRadius)
		n := len(values)
		x := func(r ui.Rect, i int) float32 { return r.X + 3 + (r.W-6)*float32(i)/float32(max(n-1, 1)) }
		ys := func(r ui.Rect) internal.Linear {
			return internal.Linear{D0: lo, D1: hi, R0: float64(r.Y + r.H - 3), R1: float64(r.Y + 3)}
		}
		if !plot.IsDisabled() {
			to := st.focus
			switch {
			case plot.Shortcut(0, ui.KeyLeft):
				to = max(st.focus-1, 0)
			case plot.Shortcut(0, ui.KeyRight):
				to = min(st.focus+1, n-1)
			case plot.Shortcut(0, ui.KeyHome):
				to = 0
			case plot.Shortcut(0, ui.KeyEnd):
				to = n - 1
			case plot.Shortcut(0, ui.KeyEscape):
				to = -1
			}
			if to != st.focus {
				st.focus = to
				if to >= 0 {
					c.Announce(chartValue(c, values[to], opts.Format))
				}
				slog.Debug("mujica: sparkline focus", "label", f.Label, "mark", to)
			}
		}
		active := -1
		if px, _, over := plot.PointerPosition(); over && !plot.IsDisabled() {
			active = max(0, min(n-1, int(math.Round(float64((px-3)/(w-6)*float32(max(n-1, 1)))))))
		} else if plot.Focused() {
			active = st.focus
		}
		res.Point = active
		plot.Draw(func(p *ui.Painter, r ui.Rect) {
			s := ys(r)
			if opts.Kind == SparklineBar {
				bw := max((r.W-6)/float32(n)-1, 1)
				for i, v := range values {
					if math.IsNaN(v) {
						continue
					}
					a, b := s.Map(0), s.Map(v)
					cx := x(r, i)
					p.Fill(ui.Rect{X: cx - bw/2, Y: min(a, b), W: bw, H: max(abs32(a-b), 1)}, col, 0)
				}
			} else {
				chartRuns(values, func(i int) float32 { return x(r, i) }, func(idx []int) {
					var line ui.Path
					for j, i := range idx {
						px, py := internal.Snap(p, x(r, i)), internal.Snap(p, s.Map(values[i]))
						if j == 0 {
							line.MoveTo(px, py)
						} else {
							line.LineTo(px, py)
						}
					}
					if opts.Kind == SparklineArea {
						fill := line
						fill.LineTo(internal.Snap(p, x(r, idx[len(idx)-1])), r.Y+r.H-3).LineTo(internal.Snap(p, x(r, idx[0])), r.Y+r.H-3).Close()
						p.FillPath(&fill, col.Alpha(0.22))
					}
					p.StrokePath(&line, 1.5, col)
				})
			}
			// Diamond: last value. Ring: the value being read.
			p.FillPath(internal.Diamond(internal.Snap(p, x(r, last)), internal.Snap(p, s.Map(values[last])), 7), k.Text)
			if active >= 0 && !math.IsNaN(values[active]) {
				var ring ui.Path
				ring.Circle(internal.Snap(p, x(r, active)), internal.Snap(p, s.Map(values[active])), 4)
				p.StrokePath(&ring, 1.5, k.Text)
			}
		})
		shown := last
		if active >= 0 {
			shown = active
		}
		ui.Text(c, chartValue(c, values[shown], opts.Format)).FontSize(size).TextColor(k.Text).SingleLine().Role(ui.RoleStatus)
	})
	res.changed = res.Point != st.last
	st.last = res.Point
	return res
}
