package feedback

import (
	"cmp"
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TraceAttr is a key and value recorded on a span.
type TraceAttr struct{ Key, Value string }

// TraceSpan records one call with trace-relative time offsets.
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
	Tree    data.TreeTableState[string]
	index   traceIndex
	src     *TraceSpan
	n       int
	built   bool
	version int
}

type traceIndex struct {
	byID  map[string]int
	kids  map[string][]string
	roots []string
	end   time.Duration
}

// TraceViewerOptions configure a TraceViewer.
type TraceViewerOptions struct {
	Status data.DataStatus
	Error  string
	// Version changes when spans are edited in place.
	Version int
}

// TraceViewer shows a span tree with time bars and attributes.
func TraceViewer(c *ui.Context, s *TraceViewerState, spans []TraceSpan, opts TraceViewerOptions) data.DataTableView {
	var first *TraceSpan
	if len(spans) > 0 {
		first = &spans[0]
	}
	if !s.built || s.src != first || s.n != len(spans) || s.version != opts.Version {
		previous := s.index.byID
		s.index = traceBuild(spans)
		s.src, s.n, s.built = first, len(spans), true
		s.version = opts.Version
		for _, id := range s.index.roots {
			if _, known := previous[id]; !known {
				s.Tree.SetOpen(id, true)
			}
		}
		slog.Debug("mujica: trace viewer index", "spans", len(spans), "roots", len(s.index.roots))
	}
	ix := s.index
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	span := func(id string) TraceSpan { return spans[ix.byID[id]] }
	cols := []data.TreeColumn[string]{
		{ID: "name", Title: core.Msg(c, "trace_viewer.span"), Compare: func(a, b string) int { return cmp.Compare(span(a).Name, span(b).Name) },
			Cell: func(id string) {
				sp := span(id)
				if sp.Failed {
					ui.Icon(c, icons.Must("circle-x")).Size(14, 14).TextColor(k.Danger).Shrink(0).Margin(0, u, 0, 0)
				}
				t := ui.Text(c, sp.Name).SingleLine()
				if sp.Failed {
					t.TextColor(k.Danger).Label(sp.Name + ", " + core.Msg(c, "trace_viewer.failed"))
				}
			}},
		{ID: "duration", Title: core.Msg(c, "trace_viewer.duration"), Width: 96, Align: ui.End,
			Compare: func(a, b string) int { return cmp.Compare(span(a).End-span(a).Start, span(b).End-span(b).Start) },
			Cell: func(id string) {
				sp := span(id)
				ui.Text(c, traceDuration(sp.End-sp.Start)).FontFeatures("tnum").TextColor(k.TextMuted)
			}},
		{ID: "timeline", Title: core.Msg(c, "trace_viewer.timeline"), Width: 260, Compare: func(a, b string) int { return cmp.Compare(span(a).Start, span(b).Start) },
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
	var v data.DataTableView
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 2)
	bounds := col.Bounds()
	internal.Remeasure(c, col, bounds.W, bounds.H)
	col.Children(func() {
		v = data.TreeTable(c, &s.Tree, data.TreeTableOptions[string]{Columns: cols, Roots: ix.roots,
			Children: func(id string) []string { return ix.kids[id] },
			Label:    func(id string) string { return span(id).Name },
			Status:   opts.Status, Error: opts.Error, Empty: core.Msg(c, "trace_viewer.empty")})
		v.Element.Grow(1)
		id, ok := s.Tree.Selected()
		if _, known := ix.byID[id]; !ok || !known {
			return
		}
		height := bounds.H
		if height == 0 {
			_, height = c.Size()
		}
		traceAttrs(c, span(id), height)
	})
	v.Element = col
	return v
}

// traceAttrs lists the chosen span's timing and attributes.
func traceAttrs(c *ui.Context, sp TraceSpan, height float32) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ui.Column(c).Gap(u).Padding(u*3, u*4).Background(k.Background).Radius(theme.CardRadius).Shrink(0).Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).MinHeight(core.Density(c).ControlHeight()).Children(func() {
			ui.Text(c, sp.Name).Font(theme.SerifFont).FontWeight(600).SingleLine().Grow(1).MinWidth(0)
			ui.Text(c, core.Msg(c, "trace_viewer.start")+" "+traceDuration(sp.Start)+" · "+core.Msg(c, "trace_viewer.duration")+" "+traceDuration(sp.End-sp.Start)).
				FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		})
		row := func(key, val string) {
			ui.Row(c).Gap(u * 3).MinHeight(core.Density(c).ControlHeight()).Children(func() {
				ui.Text(c, key).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Width(140).Shrink(0)
				ui.Text(c, val).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).Grow(1).MinWidth(0)
			})
		}
		if len(sp.Attrs) > 0 {
			ui.Column(c).Key(sp.ID).Height(min(float32(len(sp.Attrs))*core.Density(c).ControlHeight(), max(core.Density(c).ControlHeight(), min(200, height/4)))).Children(func() {
				ui.List(c, nil, len(sp.Attrs), func(i int) {
					a := sp.Attrs[i]
					row(a.Key, a.Value)
				}).Grow(1).MinHeight(0)
			})
		}
	})
}

// traceBuild indexes spans; unknown parents make roots.
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
		if _, ok := ix.byID[sp.Parent]; ok {
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
