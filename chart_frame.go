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

// ChartSpec is a chart's marks and drawing, for Chart; r is the chart box in each callback.
type ChartSpec struct {
	// Marks counts what the keys step through; 0 shows the empty state.
	Marks int
	// Legend lists series keys above the plot.
	Legend []ChartLegendItem
	// Hidden, when set, makes legend keys check boxes switching series off.
	Hidden *map[string]bool
	// Draw paints the plot; active is the read mark, or -1.
	Draw func(p *ui.Painter, r ui.Rect, active int)
	// Pick returns the mark at (x, y), or -1; required.
	Pick func(r ui.Rect, x, y float32) int
	// Readout describes mark i; it is shown, announced and the tooltip's default.
	Readout func(i int) string
	// Move steps from mark i for an arrow key; nil steps in order.
	Move func(i int, key ui.Key) int
	// Activate runs on Enter, Space or a click on a mark.
	Activate func(i int)
	// Anchor places mark i's tooltip for the keyboard; nil uses the top middle.
	Anchor func(r ui.Rect, i int) (x, y float32)
	// Tooltip writes mark i's tooltip; nil splits Readout at " · ".
	Tooltip func(i int) []ChartTooltipLine
	// Input takes the plot's other input, as a ChartAnnotation's drags.
	Input func(plot *ui.Element, r ui.Rect)
}

// Chart frames a custom chart: legend, focus, keys, readout, tooltip and empty states.
func Chart(c *ui.Context, frame ChartFrame, spec ChartSpec) ChartResult {
	if spec.Marks < 0 || spec.Draw == nil || spec.Pick == nil || spec.Readout == nil {
		panic("mujica: Chart needs Marks ≥ 0, Draw, Pick and Readout")
	}
	res, _ := newChart(c, frame).build(spec)
	return res
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

// toggles returns the legend's hidden set when series can be switched off.
func (ch *chart) toggles(on bool) *map[string]bool {
	if on {
		return &ch.st.hidden
	}
	return nil
}

// hidden reports whether the legend switched series name off.
func (ch *chart) hidden(name string) bool { return ch.st.hidden[name] }

// build lays out legend, plot and readout; returns the plot.
func (ch *chart) build(spec ChartSpec) (ChartResult, *ui.Element) {
	c, k, st := ch.c, Tokens(ch.c), ch.st
	u := Density(c).Unit()
	var plot *ui.Element
	ch.res.Element.Children(func() {
		if len(spec.Legend) > 0 {
			if name, ok := ChartLegend(c, spec.Hidden, spec.Legend, ChartLegendOptions{}).Toggled(); ok {
				ch.res.toggled = name
			}
		}
		if ch.frame.Status == DataLoading || ch.frame.Status == DataFailed || spec.Marks == 0 {
			st.focus = -1
			ch.res.retry = ChartEmptyState(c, ch.frame.Status, ChartEmptyStateOptions{Height: ch.frame.Height, Error: ch.frame.Error, Empty: ch.frame.Empty}).Changed()
			ch.res.Point = -1
			ui.Box(c).Height(u * 5)
			return
		}
		box := ui.Box(c).Height(ch.frame.Height).FillWidth().MinWidth(0)
		plot = box.Focusable().Label(ch.frame.Label).Role(ui.RoleGroup)
		internal.Ring(plot, k.Focus, theme.ControlRadius)
		st.focus = min(st.focus, spec.Marks-1)
		ch.input(plot, spec)
		if spec.Input != nil && !plot.IsDisabled() {
			b := plot.Bounds()
			spec.Input(plot, ui.Rect{W: b.W, H: b.H})
		}
		active := -1
		b := plot.Bounds()
		local := ui.Rect{W: b.W, H: b.H}
		x, y, over := plot.PointerPosition()
		if over && !plot.IsDisabled() {
			active = spec.Pick(local, x, y)
		}
		if active < 0 && plot.Focused() && st.focus >= 0 && !plot.IsDisabled() {
			active = st.focus
			x, y = local.W/2, 0
			if spec.Anchor != nil {
				x, y = spec.Anchor(local, active)
			}
		}
		ch.res.Point = active
		if active >= 0 && b.W > 0 {
			lines := chartTipLines(spec.Readout(active))
			if spec.Tooltip != nil {
				lines = spec.Tooltip(active)
			}
			ChartTooltip(c, b.X+x, b.Y+y, lines)
		}
		plot.Draw(func(p *ui.Painter, r ui.Rect) { spec.Draw(p, r, active) })
		text, col := msg(c, "chart.hint"), k.TextMuted
		if active >= 0 {
			text, col = spec.Readout(active), k.Text
		}
		ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize)).TextColor(col).Role(ui.RoleStatus).SingleLine().Height(u * 5)
	})
	ch.res.changed = ch.res.Point != st.last
	st.last = ch.res.Point
	return ch.res, plot
}

// input moves the keyboard mark, and takes clicks on marks.
func (ch *chart) input(plot *ui.Element, spec ChartSpec) {
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
		case spec.Move != nil:
			to = spec.Move(st.focus, key)
		case key == ui.KeyLeft || key == ui.KeyUp:
			to = max(st.focus-1, 0)
		default:
			to = min(st.focus+1, spec.Marks-1)
		}
	}
	if plot.Shortcut(0, ui.KeyHome) {
		to = 0
	}
	if plot.Shortcut(0, ui.KeyEnd) {
		to = spec.Marks - 1
	}
	if plot.Shortcut(0, ui.KeyEscape) {
		to = -1
	}
	if to != st.focus {
		st.focus = to
		if to >= 0 {
			ch.c.Announce(spec.Readout(to))
		}
		slog.Debug("mujica: chart focus", "label", ch.frame.Label, "mark", to)
	}
	act := st.focus >= 0 && (plot.Shortcut(0, ui.KeyEnter) || plot.Shortcut(0, ui.KeySpace))
	if plot.Clicked() {
		plot.Focus()
		x, y, _ := plot.PointerPosition()
		b := plot.Bounds()
		if i := spec.Pick(ui.Rect{W: b.W, H: b.H}, x, y); i >= 0 {
			st.focus, act = i, true
		}
	}
	if act && spec.Activate != nil {
		spec.Activate(st.focus)
	}
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
