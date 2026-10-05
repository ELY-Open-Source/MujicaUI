package feedback

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CostItem is the usage and cost of one call.
type CostItem struct {
	Model, Tool, Step string
	InputTokens       int
	OutputTokens      int
	Cost              float64
}

// CostGroup chooses what a CostBreakdown sums by.
type CostGroup int

const (
	CostByModel CostGroup = iota
	CostByTool
	CostByStep
)

// CostBreakdownOptions configure a CostBreakdown.
type CostBreakdownOptions struct {
	GroupBy CostGroup
	// Version changes when items are edited in place.
	Version int
	// Currency precedes amounts; empty means "$".
	Currency string
	// Decimals of amounts; zero means 2.
	Decimals int
}

// CostBreakdownState holds the table's sort and selection and the sums.
type CostBreakdownState struct {
	Table   data.DataTableState[string]
	rows    []costRow
	total   costRow
	src     *CostItem
	n       int
	group   CostGroup
	built   bool
	version int
	unnamed string
}

type costRow struct {
	name          string
	calls         int
	input, output int
	cost          float64
}

// CostBreakdown totals grouped costs; give the element a height.
func CostBreakdown(c *ui.Context, s *CostBreakdownState, items []CostItem, opts CostBreakdownOptions) data.DataTableView {
	if opts.GroupBy < CostByModel || opts.GroupBy > CostByStep || opts.Decimals < 0 {
		panic("mujica: CostBreakdown needs a known GroupBy and non-negative Decimals")
	}
	cur, dec := opts.Currency, opts.Decimals
	if cur == "" {
		cur = "$"
	}
	if dec == 0 {
		dec = 2
	}
	var first *CostItem
	if len(items) > 0 {
		first = &items[0]
	}
	unnamed := core.Msg(c, "cost_breakdown.unnamed")
	if !s.built || s.src != first || s.n != len(items) || s.group != opts.GroupBy || s.version != opts.Version || s.unnamed != unnamed {
		s.rows, s.total = costSum(items, opts.GroupBy, unnamed)
		s.src, s.n, s.group, s.built = first, len(items), opts.GroupBy, true
		s.version, s.unnamed = opts.Version, unnamed
		slog.Debug("mujica: cost breakdown sum", "items", len(items), "groups", len(s.rows))
	}
	money := func(v float64) string { return cur + strconv.FormatFloat(v, 'f', dec, 64) }
	count := func(n int) string { return strconv.Itoa(n) }
	share := func(r costRow) float64 {
		if s.total.cost == 0 {
			return 0
		}
		return r.cost / s.total.cost
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	groupTitle := core.Msg(c, []string{"cost_breakdown.model", "cost_breakdown.tool", "cost_breakdown.step"}[opts.GroupBy])
	cols := []data.DataColumn[costRow]{
		{ID: "name", Title: groupTitle, Sortable: true, Compare: func(a, b costRow) int { return cmp.Compare(a.name, b.name) }, Text: func(r costRow) string { return r.name }},
		{ID: "calls", Title: core.Msg(c, "cost_breakdown.calls"), Width: 80, Align: ui.End, Sortable: true, Compare: func(a, b costRow) int { return cmp.Compare(a.calls, b.calls) }, Text: func(r costRow) string { return count(r.calls) }},
		{ID: "input", Title: core.Msg(c, "cost_breakdown.input"), Width: 120, Align: ui.End, Sortable: true, Compare: func(a, b costRow) int { return cmp.Compare(a.input, b.input) }, Text: func(r costRow) string { return count(r.input) }},
		{ID: "output", Title: core.Msg(c, "cost_breakdown.output"), Width: 130, Align: ui.End, Sortable: true, Compare: func(a, b costRow) int { return cmp.Compare(a.output, b.output) }, Text: func(r costRow) string { return count(r.output) }},
		{ID: "cost", Title: core.Msg(c, "cost_breakdown.cost"), Width: 160, Sortable: true, Compare: func(a, b costRow) int { return cmp.Compare(a.cost, b.cost) },
			Text: func(r costRow) string { return money(r.cost) },
			Cell: func(r costRow) {
				ui.Row(c).Grow(1).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					f := share(r)
					ui.Box(c).Height(6).Grow(1).Radius(1).Background(k.Border).Role(ui.RoleNone).Draw(func(p *ui.Painter, b ui.Rect) {
						p.Fill(ui.Rect{X: b.X, Y: b.Y, W: b.W * float32(f), H: b.H}, k.AccentText, 1)
					})
					ui.Text(c, money(r.cost)).FontFeatures("tnum").Shrink(0)
					ui.Text(c, fmt.Sprintf("%.0f%%", 100*f)).FontFeatures("tnum").TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).Width(36).TextAlign(ui.End).Shrink(0)
				})
			}},
	}
	var v data.DataTableView
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 2)
	col.Children(func() {
		v = data.DataTable(c, &s.Table, s.rows, data.DataTableOptions[costRow, string]{Columns: cols, Key: func(r costRow) string { return r.name }, Empty: core.Msg(c, "cost_breakdown.empty")})
		v.Element.Grow(1)
		ui.Row(c).AlignItems(ui.Center).Gap(u*4).Padding(u*2, u*3).Background(k.Background).Radius(theme.ControlRadius).Role(ui.RoleStatus).
			Label(fmt.Sprintf(core.Msg(c, "cost_breakdown.total_name"), money(s.total.cost), s.total.calls)).Children(func() {
			ui.Text(c, core.Msg(c, "cost_breakdown.total")).Font(theme.SerifFont).FontWeight(600).Grow(1)
			ui.Text(c, fmt.Sprintf(core.Msg(c, "cost_breakdown.tokens"), s.total.input, s.total.output)).FontFeatures("tnum").TextColor(k.TextMuted)
			ui.Text(c, money(s.total.cost)).FontWeight(700).FontFeatures("tnum")
		})
	})
	v.Element = col
	return v
}

// costSum totals items per group, in order of first appearance.
func costSum(items []CostItem, by CostGroup, unnamed string) ([]costRow, costRow) {
	var rows []costRow
	at := map[string]int{}
	var total costRow
	for _, it := range items {
		if it.InputTokens < 0 || it.OutputTokens < 0 || it.Cost < 0 || math.IsNaN(it.Cost) || math.IsInf(it.Cost, 0) {
			panic("mujica: CostItem needs non-negative, finite counts and cost")
		}
		name := []string{it.Model, it.Tool, it.Step}[by]
		if name == "" {
			name = unnamed
		}
		i, ok := at[name]
		if !ok {
			i = len(rows)
			at[name] = i
			rows = append(rows, costRow{name: name})
		}
		for _, r := range []*costRow{&rows[i], &total} {
			if r.input > math.MaxInt-it.InputTokens || r.output > math.MaxInt-it.OutputTokens || math.IsInf(r.cost+it.Cost, 0) {
				panic("mujica: CostBreakdown total exceeds numeric range")
			}
			r.calls++
			r.input += it.InputTokens
			r.output += it.OutputTokens
			r.cost += it.Cost
		}
	}
	return rows, total
}
