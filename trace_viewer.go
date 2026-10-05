package mujica

import (
	"cmp"
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// TraceAttr is a key and value recorded on a span.
type TraceAttr struct{ Key, Value string }

// TraceSpan is one timed call of a trace; Start and End are offsets from the trace start.
type TraceSpan struct {
	ID     string
	Parent string
	Name   string
	Start  time.Duration
	End    time.Duration
	Failed bool
	Attrs  []TraceAttr
}

// TraceViewerState holds the tree's sort, open spans, selection and index.
type TraceViewerState struct {
	Tree  TreeTableState[string]
	index traceIndex
	src   *TraceSpan
	n     int
	built bool
}

type traceIndex struct {
	byID  map[string]int
	kids  map[string][]string
	roots []string
	end   time.Duration
}

// TraceViewerOptions configure a TraceViewer.
type TraceViewerOptions struct {
	Status DataStatus
	Error  string
}

// TraceViewer shows a call tree with each span's time bar; the chosen span's attributes show below.
func TraceViewer(c *ui.Context, s *TraceViewerState, spans []TraceSpan, opts TraceViewerOptions) DataTableView {
	var first *TraceSpan
	if len(spans) > 0 {
		first = &spans[0]
	}
	if !s.built || s.src != first || s.n != len(spans) {
		s.index = traceBuild(spans)
		s.src, s.n, s.built = first, len(spans), true
		for _, id := range s.index.roots {
			s.Tree.SetOpen(id, true)
		}
		slog.Debug("mujica: trace viewer index", "spans", len(spans), "roots", len(s.index.roots))
	}
	ix := s.index
	k := Tokens(c)
	u := Density(c).Unit()
	span := func(id string) TraceSpan { return spans[ix.byID[id]] }
	cols := []TreeColumn[string]{
		{ID: "name", Title: msg(c, "trace_viewer.span"), Compare: func(a, b string) int { return cmp.Compare(span(a).Name, span(b).Name) },
			Cell: func(id string) {
				sp := span(id)
				if sp.Failed {
					ui.Icon(c, icons.Must("circle-x")).Size(14, 14).TextColor(k.Danger).Shrink(0).Margin(0, u, 0, 0)
				}
				t := ui.Text(c, sp.Name).SingleLine()
				if sp.Failed {
					t.TextColor(k.Danger).Label(sp.Name + ", " + msg(c, "trace_viewer.failed"))
				}
			}},
		{ID: "duration", Title: msg(c, "trace_viewer.duration"), Width: 96, Align: ui.End,
			Compare: func(a, b string) int { return cmp.Compare(span(a).End-span(a).Start, span(b).End-span(b).Start) },
			Cell: func(id string) {
				sp := span(id)
				ui.Text(c, traceDuration(sp.End-sp.Start)).FontFeatures("tnum").TextColor(k.TextMuted)
			}},
		{ID: "timeline", Title: msg(c, "trace_viewer.timeline"), Width: 260, Compare: func(a, b string) int { return cmp.Compare(span(a).Start, span(b).Start) },
			Cell: func(id string) {
				sp := span(id)
				col := k.AccentText
				if sp.Failed {
					col = k.Danger
				}
				ui.Box(c).Height(12).Grow(1).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H/2, W: r.W, H: 1}, k.Border, 0)
					if ix.end <= 0 {
						return
					}
					x0 := r.X + r.W*float32(float64(sp.Start)/float64(ix.end))
					w := max(2, r.W*float32(float64(sp.End-sp.Start)/float64(ix.end)))
					bar := ui.Rect{X: internal.Snap(p, x0), Y: r.Y, W: internal.Snap(p, w), H: r.H}
					if sp.Failed {
						// Failed spans are hollow, not told apart by color alone.
						p.Fill(bar, k.Surface, 2)
						p.Stroke(bar, col, 2, 1.5)
						return
					}
					p.Fill(bar, col, 2)
				})
			}},
	}
	var v DataTableView
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 2)
	col.Children(func() {
		v = TreeTable(c, &s.Tree, TreeTableOptions[string]{Columns: cols, Roots: ix.roots,
			Children: func(id string) []string { return ix.kids[id] },
			Label:    func(id string) string { return span(id).Name },
			Status:   opts.Status, Error: opts.Error, Empty: msg(c, "trace_viewer.empty")})
		v.Element.Grow(1)
		id, ok := s.Tree.Selected()
		if _, known := ix.byID[id]; !ok || !known {
			return
		}
		traceAttrs(c, span(id))
	})
	v.Element = col
	return v
}

// traceAttrs lists the chosen span's timing and attributes.
func traceAttrs(c *ui.Context, sp TraceSpan) {
	k := Tokens(c)
	u := Density(c).Unit()
	ui.Column(c).Gap(u).Padding(u*3, u*4).Background(k.Background).Radius(theme.CardRadius).Shrink(0).Children(func() {
		ui.Text(c, sp.Name).Font(theme.SerifFont).FontWeight(600)
		row := func(key, val string) {
			ui.Row(c).Gap(u * 3).Children(func() {
				ui.Text(c, key).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Width(140).Shrink(0)
				ui.Text(c, val).Font(theme.MonoFont).FontSize(FontSize(c, theme.CaptionSize)).Grow(1).MinWidth(0)
			})
		}
		row(msg(c, "trace_viewer.start"), traceDuration(sp.Start))
		row(msg(c, "trace_viewer.duration"), traceDuration(sp.End-sp.Start))
		for _, a := range sp.Attrs {
			row(a.Key, a.Value)
		}
	})
}

// traceBuild indexes spans by ID and parent; unknown parents make roots.
func traceBuild(spans []TraceSpan) traceIndex {
	ix := traceIndex{byID: make(map[string]int, len(spans)), kids: map[string][]string{}}
	for i, sp := range spans {
		if sp.ID == "" || sp.End < sp.Start || sp.Start < 0 {
			panic(fmt.Sprintf("mujica: TraceSpan %q needs an ID and 0 <= Start <= End", sp.ID))
		}
		if _, dup := ix.byID[sp.ID]; dup {
			panic("mujica: duplicate TraceSpan ID " + sp.ID)
		}
		ix.byID[sp.ID] = i
		ix.end = max(ix.end, sp.End)
	}
	for _, sp := range spans {
		if _, ok := ix.byID[sp.Parent]; ok && sp.Parent != sp.ID {
			ix.kids[sp.Parent] = append(ix.kids[sp.Parent], sp.ID)
		} else {
			ix.roots = append(ix.roots, sp.ID)
		}
	}
	seen := 0
	var walk func(ids []string)
	walk = func(ids []string) {
		for _, id := range ids {
			seen++
			walk(ix.kids[id])
		}
	}
	walk(ix.roots)
	if seen != len(spans) {
		panic("mujica: TraceSpan parents form a cycle")
	}
	return ix
}

// traceDuration prints microseconds, milliseconds or seconds.
func traceDuration(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return fmt.Sprintf("%d µs", d.Microseconds())
	case d < time.Second:
		return fmt.Sprintf("%.1f ms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2f s", d.Seconds())
}
