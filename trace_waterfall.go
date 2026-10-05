package mujica

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// WaterfallSpan holds timing; empty Parent marks a root.
type WaterfallSpan struct {
	ID, Parent string
	Name       string
	Service    string
	// Start is the offset from the trace start.
	Start, Duration time.Duration
	Error           bool
}

// TraceWaterfallOptions configure a TraceWaterfall.
type TraceWaterfallOptions struct {
	// Height of the span list; zero means 320.
	Height float32
}

// traceRow is a span in tree order with its depth.
type traceRow struct {
	span  int
	depth int
}

type traceState struct {
	key      string
	rows     []traceRow
	total    time.Duration
	selected string
	services map[string]int // palette slot of each service
}

// TraceWaterfall shows timed span hierarchies and selected details.
func TraceWaterfall(c *ui.Context, spans []WaterfallSpan, opts TraceWaterfallOptions) ChoiceResult {
	if opts.Height < 0 {
		panic("mujica: TraceWaterfall Height is negative")
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = codePanel(c).Label(msg(c, "trace_waterfall.label"))
	st := ui.Local(r.Element, "mujica.trace", func() traceState { return traceState{} })
	traceLayout(st, spans)
	r.Element.Children(func() {
		if len(st.rows) == 0 {
			codeCaption(c, msg(c, "trace_waterfall.empty")).Padding(u * 4)
			return
		}
		codeHeader(c, func() {
			codeCaption(c, fmt.Sprintf(msg(c, "trace_waterfall.summary"), len(spans), agentDuration(st.total))).Grow(1)
		})
		agentRows(c, len(st.rows), func(i int) {
			row := st.rows[i]
			s := spans[row.span]
			if traceSpanRow(c, s, row.depth, st.total, st.services[s.Service], s.ID == st.selected) {
				st.selected, r.changed = s.ID, true
				slog.Debug("mujica: trace_waterfall select", "span", s.ID)
			}
		}).Height(h).FocusGroup(ui.Vertical)
		for _, s := range spans {
			if s.ID == st.selected {
				ui.Column(c).Shrink(0).Gap(u).Padding(u*2, u*3).BorderWidth(theme.BorderWidth, 0, 0, 0).BorderColor(k.Border).
					Role(ui.RoleStatus).Children(func() {
					ui.Text(c, s.Name).FontWeight(600)
					codeCaption(c, fmt.Sprintf(msg(c, "trace_waterfall.detail"), s.Service, agentDuration(s.Start), agentDuration(s.Duration)))
				})
			}
		}
	})
	return r
}

// traceLayout caches spans in depth-first time order.
func traceLayout(st *traceState, spans []WaterfallSpan) {
	key := fmt.Sprintf("%d", len(spans))
	if len(spans) > 0 {
		key += fmt.Sprintf("%p", &spans[0])
	}
	if key == st.key {
		return
	}
	st.key, st.rows, st.total, st.services = key, st.rows[:0], 0, map[string]int{}
	index := make(map[string]int, len(spans))
	kids := map[string][]int{}
	for i, s := range spans {
		if s.ID == "" || s.Name == "" || s.Start < 0 || s.Duration < 0 {
			panic(fmt.Sprintf("mujica: TraceWaterfall span %d needs an ID, a Name and non-negative times", i))
		}
		if _, dup := index[s.ID]; dup {
			panic("mujica: TraceWaterfall span IDs must be unique: " + s.ID)
		}
		index[s.ID] = i
		if _, ok := st.services[s.Service]; !ok {
			st.services[s.Service] = len(st.services)
		}
		st.total = max(st.total, s.Start+s.Duration)
	}
	for i, s := range spans {
		if _, ok := index[s.Parent]; s.Parent != "" && !ok {
			panic("mujica: TraceWaterfall span " + s.ID + " has an unknown Parent")
		}
		kids[s.Parent] = append(kids[s.Parent], i)
	}
	var walk func(parent string, depth int)
	walk = func(parent string, depth int) {
		list := kids[parent]
		slices.SortStableFunc(list, func(a, b int) int { return cmp.Compare(spans[a].Start, spans[b].Start) })
		for _, i := range list {
			st.rows = append(st.rows, traceRow{i, depth})
			walk(spans[i].ID, depth+1)
		}
	}
	walk("", 0)
	if len(st.rows) != len(spans) {
		panic("mujica: TraceWaterfall spans form a cycle")
	}
}

// traceSpanRow shows span indentation, timing, selection and duration.
func traceSpanRow(c *ui.Context, s WaterfallSpan, depth int, total time.Duration, slot int, selected bool) bool {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	row := ui.ButtonBase(c).Key(s.ID).Justify(ui.Start).Gap(u*2).MinHeight(d.ControlHeight()).Padding(0, u*2).
		Label(fmt.Sprintf("%s, %s, %s", s.Name, s.Service, agentDuration(s.Duration)))
	internal.Ring(row, k.Focus, theme.ControlRadius)
	if selected {
		row.Background(k.Selection).DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Stroke(ui.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}, k.Focus, theme.ControlRadius, 1)
		})
	} else if row.Hovered() {
		row.Background(k.SurfaceHover)
	}
	pal := internal.Palette(k)
	col := pal[slot%len(pal)]
	if s.Error {
		col = k.Danger
	}
	row.Children(func() {
		ui.Row(c).Width(220).Shrink(0).Gap(u).Padding(0, 0, 0, float32(depth)*u*4).AlignItems(ui.Center).Children(func() {
			if s.Error {
				ui.Icon(c, icons.Must("circle-x")).Size(12, 12).TextColor(k.Danger).Label(msg(c, "trace_waterfall.error"))
			}
			ui.Text(c, s.Name).FontSize(FontSize(c, theme.CaptionSize+1)).SingleLine().Shrink(1)
		})
		ui.Box(c).Grow(1).Height(14).Draw(func(p *ui.Painter, r ui.Rect) {
			span := max(float64(total), 1)
			x0 := r.X + r.W*float32(float64(s.Start)/span)
			w := max(r.W*float32(float64(s.Duration)/span), 2)
			bar := ui.Rect{X: internal.Snap(p, x0), Y: r.Y, W: w, H: r.H}
			if s.Error { // hollow, so errors differ by shape too
				p.Fill(bar, col.Alpha(0.25), 2)
				p.Stroke(bar, col, 2, 1.5)
				return
			}
			p.Fill(bar, col, 2)
		})
		ui.Text(c, agentDuration(s.Duration)).Width(64).TextAlign(ui.End).FontSize(FontSize(c, theme.CaptionSize)).
			TextColor(k.TextMuted).FontFeatures("tnum").Shrink(0)
	})
	return !row.IsDisabled() && row.Clicked()
}
