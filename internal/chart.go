package internal

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// Linear maps domain D0..D1 onto pixels R0..R1.
type Linear struct{ D0, D1, R0, R1 float64 }

// Map returns the pixel of v.
func (s Linear) Map(v float64) float32 {
	if s.D1 == s.D0 {
		return float32((s.R0 + s.R1) / 2)
	}
	return float32(s.R0 + (v-s.D0)/(s.D1-s.D0)*(s.R1-s.R0))
}

// Invert returns the value at pixel px.
func (s Linear) Invert(px float32) float64 {
	if s.R1 == s.R0 {
		return s.D0
	}
	return s.D0 + (float64(px)-s.R0)/(s.R1-s.R0)*(s.D1-s.D0)
}

// Nice widens lo..hi to round steps of about n ticks.
func Nice(lo, hi float64, n int) (nlo, nhi, step float64) {
	if n < 1 || math.IsNaN(lo) || math.IsNaN(hi) || lo > hi {
		panic("internal: bad Nice domain")
	}
	if lo == hi {
		lo, hi = lo-1, hi+1
	}
	raw := (hi - lo) / float64(n)
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	step = 10 * mag
	for _, m := range []float64{1, 2, 5} {
		if m*mag >= raw {
			step = m * mag
			break
		}
	}
	return math.Floor(lo/step) * step, math.Ceil(hi/step) * step, step
}

// Ticks lists the multiples of step from lo to hi.
func Ticks(lo, hi, step float64) []float64 {
	var out []float64
	for i := math.Ceil(lo/step - 1e-9); i*step <= hi+step*1e-9; i++ {
		out = append(out, i*step)
	}
	return out
}

// FormatTick writes v with as many decimals as step needs.
func FormatTick(v, step float64) string {
	d := max(0, int(-math.Floor(math.Log10(step)+1e-9)))
	if math.Abs(v) < step*1e-9 {
		v = 0
	}
	return strconv.FormatFloat(v, 'f', d, 64)
}

// timeSteps are the tick spacings a time axis chooses from.
var timeSteps = []struct {
	d      time.Duration
	layout string
}{
	{time.Second, "15:04:05"}, {5 * time.Second, "15:04:05"}, {15 * time.Second, "15:04:05"}, {30 * time.Second, "15:04:05"},
	{time.Minute, "15:04"}, {5 * time.Minute, "15:04"}, {15 * time.Minute, "15:04"}, {30 * time.Minute, "15:04"},
	{time.Hour, "15:04"}, {3 * time.Hour, "15:04"}, {6 * time.Hour, "Jan 2 15:04"}, {12 * time.Hour, "Jan 2 15:04"},
	{24 * time.Hour, "Jan 2"}, {7 * 24 * time.Hour, "Jan 2"}, {30 * 24 * time.Hour, "Jan 2006"}, {365 * 24 * time.Hour, "2006"},
}

// TimeTicks returns about n round times and their layout.
func TimeTicks(lo, hi time.Time, n int) ([]time.Time, string) {
	if n < 1 || hi.Before(lo) {
		panic("internal: bad TimeTicks domain")
	}
	span := hi.Sub(lo) / time.Duration(n)
	step := timeSteps[len(timeSteps)-1]
	for _, s := range timeSteps {
		if s.d >= span {
			step = s
			break
		}
	}
	var out []time.Time
	for t := lo.Truncate(step.d); !t.After(hi); t = t.Add(step.d) {
		if !t.Before(lo) {
			out = append(out, t)
		}
	}
	return out, step.layout
}

// Band splits R0..R1 into N bands, each Pad empty.
type Band struct {
	N      int
	R0, R1 float32
	Pad    float32
}

// Step is the room one band takes.
func (b Band) Step() float32 { return (b.R1 - b.R0) / float32(max(b.N, 1)) }

// Width is a band's drawn width.
func (b Band) Width() float32 { return b.Step() * (1 - b.Pad) }

// Start is where band i is drawn from.
func (b Band) Start(i int) float32 { return b.R0 + b.Step()*float32(i) + b.Step()*b.Pad/2 }

// Center is the middle of band i.
func (b Band) Center(i int) float32 { return b.R0 + b.Step()*(float32(i)+0.5) }

// Palette returns six series colors at 3:1 against Surface.
func Palette(k theme.Tokens) []ui.Color {
	out := []ui.Color{k.AccentText, k.Info, k.Success, k.Warning, k.Ornament, k.TextMuted}
	for i, col := range out {
		for t := float32(0.1); theme.Contrast(col, k.Surface) < 3; t += 0.1 {
			col = out[i].Mix(k.Text, t)
		}
		out[i] = col
	}
	return out
}

// Ramp is the sequential color for t in 0..1.
func Ramp(k theme.Tokens, t float64) ui.Color {
	return k.Surface.Mix(k.AccentText, float32(0.12+0.88*max(0, min(t, 1))))
}

// Marker draws series shape i centered at (cx, cy).
func Marker(p *ui.Painter, i int, cx, cy, size float32, col ui.Color) {
	cx, cy = Snap(p, cx), Snap(p, cy)
	h := size / 2
	var path ui.Path
	switch i % 5 {
	case 0:
		path.Circle(cx, cy, h)
	case 1:
		h *= 0.85
		path.MoveTo(cx-h, cy-h).LineTo(cx+h, cy-h).LineTo(cx+h, cy+h).LineTo(cx-h, cy+h).Close()
	case 2:
		p.FillPath(Diamond(cx, cy, size*1.2), col)
		return
	case 3:
		path.MoveTo(cx, cy-h).LineTo(cx+h, cy+h*0.8).LineTo(cx-h, cy+h*0.8).Close()
	default:
		path.MoveTo(cx-h, cy-h*0.8).LineTo(cx+h, cy-h*0.8).LineTo(cx, cy+h).Close()
	}
	p.FillPath(&path, col)
}

// Arc appends a polyline from a0 to a1, clockwise radians.
func Arc(path *ui.Path, cx, cy, r float32, a0, a1 float64, move bool) {
	n := max(2, int(math.Ceil(math.Abs(a1-a0)/(math.Pi/48))))
	for i := 0; i <= n; i++ {
		x, y := Polar(cx, cy, r, a0+(a1-a0)*float64(i)/float64(n))
		if i == 0 && move {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
}

// Polar returns the point r from center at angle a.
func Polar(cx, cy, r float32, a float64) (float32, float32) {
	return cx + r*float32(math.Sin(a)), cy - r*float32(math.Cos(a))
}

// Sector returns a ring slice; zero r0 makes a wedge.
func Sector(cx, cy, r0, r1 float32, a0, a1 float64) *ui.Path {
	var path ui.Path
	Arc(&path, cx, cy, r1, a0, a1, true)
	if r0 > 0 {
		Arc(&path, cx, cy, r0, a1, a0, false)
	} else {
		path.LineTo(cx, cy)
	}
	path.Close()
	return &path
}

// Rule draws a pixel-aligned 1 DIP horizontal line at y.
func Rule(p *ui.Painter, x0, x1, y float32, c ui.Color) {
	yy := ruleY(p, y)
	p.Line(x0, yy, x1, yy, 1, c)
}

// AxisY draws grid lines and tick labels left of plot.
func AxisY(p *ui.Painter, plot ui.Rect, s Linear, ticks []float64, label func(float64) string, size float32, k theme.Tokens) {
	for _, t := range ticks {
		y := s.Map(t)
		Rule(p, plot.X, plot.X+plot.W, y, k.Border)
		txt := label(t)
		w, h := p.MeasureText(0, ui.Span{Text: txt, Size: size})
		p.Text(plot.X-w-6, y-h/2, txt, size, k.TextMuted)
	}
}

// AxisX draws a baseline and labels centered at xs.
func AxisX(p *ui.Painter, plot ui.Rect, xs []float32, labels []string, size float32, k theme.Tokens) {
	Rule(p, plot.X, plot.X+plot.W, plot.Y+plot.H, k.ControlBorder)
	right := float32(math.Inf(-1))
	for i, x := range xs {
		w, _ := p.MeasureText(0, ui.Span{Text: labels[i], Size: size})
		at := min(max(x-w/2, plot.X-12), plot.X+plot.W+12-w)
		// Labels that would overlap the last one drawn are skipped.
		if at < right+6 {
			continue
		}
		p.Text(at, plot.Y+plot.H+6, labels[i], size, k.TextMuted)
		right = at + w
	}
}

// Nearest returns the closest point within reach, or -1.
func Nearest(n int, at func(i int) (float32, float32), x, y, reach float32, byX bool) int {
	best, dist := -1, reach*reach
	for i := range n {
		px, py := at(i)
		d := (px - x) * (px - x)
		if !byX {
			d += (py - y) * (py - y)
		}
		if d <= dist {
			best, dist = i, d
		}
	}
	return best
}

// Columns keeps first, lowest, highest and last per pixel column.
func Columns(n int, x func(i int) float32, y func(i int) float64) []int {
	var out []int
	for i := 0; i < n; {
		col := math.Floor(float64(x(i)))
		lo, hi, j := i, i, i
		for ; j < n && math.Floor(float64(x(j))) == col; j++ {
			if y(j) < y(lo) {
				lo = j
			}
			if y(j) > y(hi) {
				hi = j
			}
		}
		out = append(out, i)
		for _, k := range []int{min(lo, hi), max(lo, hi), j - 1} {
			if k != out[len(out)-1] {
				out = append(out, k)
			}
		}
		i = j
	}
	return out
}
