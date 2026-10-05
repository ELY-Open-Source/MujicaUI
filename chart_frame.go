package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChartFrame is what every chart shares.
type ChartFrame struct {
	// Label names the chart for assistive technology; required.
	Label string
	// Height is the plot height in DIPs; 0 means 240.
	Height float32
	// Status shows loading or failure in the plot's place.
	Status DataStatus
	Error  string
	// Empty replaces the built-in "No data" text.
	Empty string
}

// ChartResult is a built chart and this frame's events.
type ChartResult struct {
	Element *ui.Element
	// Point is the hovered or focused mark, -1 for none.
	Point          int
	changed, retry bool
	toggled        string
	value          bool
}

// Changed reports that the readout moved to another mark.
func (r ChartResult) Changed() bool { return r.changed }

// Retried reports a click on a failed load's retry.
func (r ChartResult) Retried() bool { return r.retry }

// ValueChanged reports a changed selection, path, year or pause.
func (r ChartResult) ValueChanged() bool { return r.value }

// Toggled returns the legend series switched in this frame.
func (r ChartResult) Toggled() (string, bool) { return r.toggled, r.toggled != "" }

// chartKey is a legend entry: name, color and shape.
type chartKey struct {
	Name  string
	Color ui.Color
	Shape int
}

// chartSpec is one chart's drawing and behavior.
type chartSpec struct {
	// n counts the marks keys step through; 0 is empty.
	n    int
	keys []chartKey
	// toggle lets legend entries hide and show series.
	toggle bool
	draw   func(p *ui.Painter, r ui.Rect, active int)
	// pick returns the mark at (x, y), or -1.
	pick    func(r ui.Rect, x, y float32) int
	readout func(i int) string
	// move steps from mark i; nil steps in order.
	move func(i int, key ui.Key) int
	// activate runs on Enter, Space or a click.
	keyActivate func(i int)
	activate    func(i int)
}

type chartState struct {
	focus, last int
	hidden      map[string]bool
}

// chart is a chart being built.
type chart struct {
	c     *ui.Context
	frame ChartFrame
	st    *chartState
	res   ChartResult
}

// newChart starts a chart; ask it which series are hidden.
func newChart(c *ui.Context, frame ChartFrame) *chart {
	if frame.Label == "" || frame.Height < 0 {
		panic("mujica: chart needs a Label and a non-negative Height")
	}
	if frame.Height == 0 {
		frame.Height = 240
	}
	ch := &chart{c: c, frame: frame}
	ch.res.Element = ui.Column(c).Gap(Density(c).Unit() * 2).MinWidth(0)
	ch.st = ui.Local(ch.res.Element, "mujica.chart", func() chartState { return chartState{focus: -1, last: -1, hidden: map[string]bool{}} })
	return ch
}

// hidden reports whether the legend switched series name off.
func (ch *chart) hidden(name string) bool { return ch.st.hidden[name] }

// build lays out legend, plot and readout; returns the plot.
func (ch *chart) build(spec chartSpec) (ChartResult, *ui.Element) {
	c, k, st := ch.c, Tokens(ch.c), ch.st
	u := Density(c).Unit()
	var plot *ui.Element
	ch.res.Element.Children(func() {
		if len(spec.keys) > 0 {
			ch.legend(spec)
		}
		box := ui.Box(c).Height(ch.frame.Height).FillWidth().MinWidth(0)
		if ch.frame.Status == DataLoading || ch.frame.Status == DataFailed || spec.n == 0 {
			st.focus = -1
			box.Children(func() {
				ui.Column(c).Fill().Children(func() {
					ch.res.retry = dataPlaceholder(c, ch.frame.Status, ch.frame.Error, ch.frame.Empty, msg(c, "chart.empty"))
				})
			})
			ch.res.Point = -1
			ui.Box(c).Height(max(u*5, FontSize(c, theme.CaptionSize)*theme.LineHeight))
			return
		}
		plot = box.Focusable().Label(ch.frame.Label).Role(ui.RoleGroup)
		internal.Ring(plot, k.Focus, theme.ControlRadius)
		st.focus = min(st.focus, spec.n-1)
		ch.input(plot, spec)
		active := -1
		if x, y, over := plot.PointerPosition(); over && !plot.IsDisabled() {
			b := plot.Bounds()
			active = spec.pick(ui.Rect{W: b.W, H: b.H}, x, y)
		}
		if active < 0 && plot.Focused() {
			active = st.focus
		}
		ch.res.Point = active
		plot.Draw(func(p *ui.Painter, r ui.Rect) { spec.draw(p, r, active) })
		text, col := msg(c, "chart.hint"), k.TextMuted
		if active >= 0 {
			text, col = spec.readout(active), k.Text
		}
		ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize)).TextColor(col).Role(ui.RoleStatus).SingleLine().Height(max(u*5, FontSize(c, theme.CaptionSize)*theme.LineHeight))
	})
	ch.res.changed = ch.res.Point != st.last
	st.last = ch.res.Point
	return ch.res, plot
}

// input moves the keyboard mark, and takes clicks on marks.
func (ch *chart) input(plot *ui.Element, spec chartSpec) {
	st := ch.st
	if plot.IsDisabled() {
		return
	}
	to := st.focus
	for _, key := range []ui.Key{ui.KeyLeft, ui.KeyRight, ui.KeyUp, ui.KeyDown} {
		if !plot.Shortcut(0, key) {
			continue
		}
		switch {
		case st.focus < 0:
			to = 0
		case spec.move != nil:
			to = spec.move(st.focus, key)
		case key == ui.KeyLeft || key == ui.KeyUp:
			to = max(st.focus-1, 0)
		default:
			to = min(st.focus+1, spec.n-1)
		}
	}
	if plot.Shortcut(0, ui.KeyHome) {
		to = 0
	}
	if plot.Shortcut(0, ui.KeyEnd) {
		to = spec.n - 1
	}
	if plot.Shortcut(0, ui.KeyEscape) {
		to = -1
	}
	if to != st.focus {
		st.focus = to
		if to >= 0 {
			ch.c.Announce(spec.readout(to))
		}
		slog.Debug("mujica: chart focus", "label", ch.frame.Label, "mark", to)
	}
	act := st.focus >= 0 && (plot.Shortcut(0, ui.KeyEnter) || plot.Shortcut(0, ui.KeySpace))
	if act && spec.keyActivate != nil {
		spec.keyActivate(st.focus)
	}
	if plot.Clicked() {
		plot.Focus()
		x, y, _ := plot.PointerPosition()
		b := plot.Bounds()
		if i := spec.pick(ui.Rect{W: b.W, H: b.H}, x, y); i >= 0 {
			st.focus, act = i, true
		}
	}
	if act && spec.activate != nil {
		spec.activate(st.focus)
	}
}

// legend shows the series keys; toggling ones are check boxes.
func (ch *chart) legend(spec chartSpec) {
	c, k, st := ch.c, Tokens(ch.c), ch.st
	u := Density(c).Unit()
	ui.Row(c).Wrap().Gap(u * 3).Children(func() {
		for _, key := range spec.keys {
			ui.Row(c).Key(key.Name).Children(func() {
				on := !st.hidden[key.Name]
				var e *ui.Element
				if spec.toggle {
					e = ui.CheckboxBase(c, &on).Disabled(ch.res.Element.IsDisabled()).Label(key.Name)
					internal.Ring(e, k.Focus, theme.ControlRadius)
					if e.Changed() && !e.IsDisabled() {
						st.hidden[key.Name] = !on
						ch.res.toggled = key.Name
						slog.Debug("mujica: chart toggle", "series", key.Name, "visible", on)
					}
				} else {
					e = ui.Row(c)
				}
				e.MinHeight(theme.MinHitSize).Gap(u*2).AlignItems(ui.Center).Padding(0, u).Radius(theme.ControlRadius)
				col := key.Color
				if !on {
					col = k.TextMuted
				}
				e.Children(func() {
					ui.Box(c).Size(12, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
						if on {
							internal.Marker(p, key.Shape, r.X+6, r.Y+6, 10, col)
						} else {
							p.Stroke(r, col, 2, 1)
						}
					})
					t := ui.Text(c, key.Name).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Text)
					if !on {
						t.Strikethrough().TextColor(k.TextMuted)
					}
				})
			})
		}
	})
}

// chartText draws s anchored at (x, y) by ax, ay.
func chartText(p *ui.Painter, s string, x, y, size float32, col ui.Color, ax, ay float32) (w, h float32) {
	w, h = p.MeasureText(0, ui.Span{Text: s, Size: size})
	p.Text(x-w*ax, y-h*ay, s, size, col)
	return w, h
}

// chartInset is r less room for axis labels.
func chartInset(r ui.Rect, left, bottom float32) ui.Rect {
	return ui.Rect{X: r.X + left, Y: r.Y + 8, W: max(r.W-left-12, 1), H: max(r.H-bottom-8, 1)}
}

// chartInk is Text or Surface, whichever reads better on bg.
func chartInk(k theme.Tokens, bg ui.Color) ui.Color {
	if theme.Contrast(k.Text, bg) >= theme.Contrast(k.Surface, bg) {
		return k.Text
	}
	return k.Surface
}
