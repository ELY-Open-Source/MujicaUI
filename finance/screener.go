package finance

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chart"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type ScreenerField struct {
	ID, Label string
	Decimals  int
}

type ScreenerOp int

const (
	ScreenerAbove ScreenerOp = iota
	ScreenerAtLeast
	ScreenerBelow
	ScreenerAtMost
	ScreenerEqual
)

var screenerOps = [5]string{">", "≥", "<", "≤", "="}

// ScreenerCondition is one test: Field Op Value.
type ScreenerCondition struct {
	Field string
	Op    ScreenerOp
	Value float64
}

type ScreenerRow struct {
	Symbol, Name string
	Values       map[string]float64
}

// ScreenerMatch reports whether row passes every condition.
func ScreenerMatch(row ScreenerRow, conds []ScreenerCondition) bool {
	for _, cd := range conds {
		v, ok := row.Values[cd.Field]
		if !ok || math.IsNaN(v) {
			return false
		}
		pass := [5]bool{v > cd.Value, v >= cd.Value, v < cd.Value, v <= cd.Value, v == cd.Value}
		if cd.Op < ScreenerAbove || cd.Op > ScreenerEqual {
			panic("mujica: ScreenerOp out of range")
		}
		if !pass[cd.Op] {
			return false
		}
	}
	return true
}

// ScreenerState is the conditions and the result table's state.
type ScreenerState struct {
	Conditions []ScreenerCondition
	Table      data.DataTableState[string]
	cache      screenerCache
}

type screenerCache struct {
	key  string
	rows []ScreenerRow
}

// ScreenerOptions configure a Screener.
type ScreenerOptions struct {
	// Version must change when rows change in place.
	Version  int
	Status   data.DataStatus
	Error    string
	Disabled bool
}

// ScreenerView is a built Screener and this frame's events.
type ScreenerView struct {
	Element       ui.Element
	Changed, Save bool
	Table         data.DataTableView
}

func Screener(c *ui.Context, s *ScreenerState, fields []ScreenerField, rows []ScreenerRow, opts ScreenerOptions) ScreenerView {
	if s == nil || len(fields) == 0 {
		panic("mujica: Screener needs a state and fields")
	}
	seen := map[string]bool{}
	for _, f := range fields {
		if f.ID == "" || seen[f.ID] || f.Decimals < 0 || f.Decimals > 12 {
			panic("mujica: Screener fields need unique IDs and decimals in 0..12")
		}
		seen[f.ID] = true
	}
	for _, cd := range s.Conditions {
		if cd.Op < ScreenerAbove || cd.Op > ScreenerEqual || math.IsNaN(cd.Value) || math.IsInf(cd.Value, 0) {
			panic("mujica: Screener condition is invalid")
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	before := fmt.Sprint(s.Conditions)
	var v ScreenerView
	v.Element = ui.Column(c).Gap(u * 3).Grow(1).MinHeight(0).Disabled(opts.Disabled)
	opts.Disabled = v.Element.IsDisabled()
	if opts.Disabled {
		copy := *s
		copy.Conditions = slices.Clone(s.Conditions)
		s = &copy
	}
	v.Element.Children(func() {
		fieldItems := make([]input.SelectOption[string], len(fields))
		for i, f := range fields {
			fieldItems[i] = input.SelectOption[string]{Value: f.ID, Label: f.Label}
		}
		opItems := make([]input.SelectOption[ScreenerOp], len(screenerOps))
		for i, o := range screenerOps {
			opItems[i] = input.SelectOption[ScreenerOp]{Value: ScreenerOp(i), Label: o + " " + core.Msg(c, fmt.Sprintf("screener.op_%d", i))}
		}
		ui.Column(c).Gap(u*2).Padding(u*3).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
			Role(ui.RoleGroup).Label(core.Msg(c, "screener.conditions")).Children(func() {
			if len(s.Conditions) == 0 {
				ui.Text(c, core.Msg(c, "screener.none")).TextColor(k.TextMuted)
			}
			for i := 0; i < len(s.Conditions); i++ {
				cd := &s.Conditions[i]
				f := screenerField(fields, cd.Field)
				removed := false
				ui.Row(c.Key(i)).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					input.Select(c, &cd.Field, fieldItems, input.SelectOptions{Label: core.Msg(c, "screener.field"), Disabled: opts.Disabled})
					input.Select(c, &cd.Op, opItems, input.SelectOptions{Label: core.Msg(c, "screener.compare"), Disabled: opts.Disabled}).Element.MinWidth(120)
					input.NumberInput(c, &cd.Value, input.NumberInputOptions{Step: math.Pow(10, -float64(f.Decimals)), Precision: f.Decimals,
						Label: core.Msg(c, "screener.value"), Disabled: opts.Disabled}).Element.Width(150)
					x := core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "screener.remove"), i+1))
					removed = x.Clicked() && !x.IsDisabled()
				})
				if removed {
					s.Conditions = slices.Delete(slices.Clone(s.Conditions), i, i+1)
					i--
				}
			}
			ui.Row(c).Gap(u * 2).Children(func() {
				if input.Button(c, core.Msg(c, "screener.add"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("plus"), Disabled: opts.Disabled}).Clicked() {
					s.Conditions = append(s.Conditions, ScreenerCondition{Field: fields[0].ID})
				}
				v.Save = input.Button(c, core.Msg(c, "screener.save"), input.ButtonOptions{Variant: input.Ghost, Disabled: opts.Disabled}).Clicked()
			})
		})
		matched := screenerFilter(s, rows, opts.Version)
		ui.Text(c, fmt.Sprintf(core.Msg(c, "screener.matches"), len(matched), len(rows))).FontSize(core.FontSize(c, theme.CaptionSize)).
			TextColor(k.TextMuted).Role(ui.RoleStatus)
		cols := []data.DataColumn[ScreenerRow]{{ID: "symbol", Title: core.Msg(c, "screener.symbol"), Sortable: true,
			Compare: func(a, b ScreenerRow) int { return cmp.Compare(a.Symbol, b.Symbol) },
			Text:    func(r ScreenerRow) string { return r.Symbol }}}
		for _, f := range fields {
			cols = append(cols, data.DataColumn[ScreenerRow]{ID: f.ID, Title: f.Label, Align: ui.End, Sortable: true,
				Compare: func(a, b ScreenerRow) int { return cmp.Compare(a.Values[f.ID], b.Values[f.ID]) },
				Text:    func(r ScreenerRow) string { return data.StatisticFormat(r.Values[f.ID], f.Decimals) }})
		}
		v.Table = data.DataTable(c, &s.Table, matched, data.DataTableOptions[ScreenerRow, string]{Columns: cols, Key: func(r ScreenerRow) string { return r.Symbol },
			Version: opts.Version, Status: opts.Status, Error: opts.Error, Empty: core.Msg(c, "screener.empty")})
		v.Table.Element.Grow(1)
	})
	v.Changed = fmt.Sprint(s.Conditions) != before
	if v.Element.IsDisabled() {
		v.Save = false
	}
	if v.Changed || v.Save {
		slog.Debug("mujica: screener conditions", "count", len(s.Conditions), "save", v.Save)
	}
	return v
}

func screenerField(fields []ScreenerField, id string) ScreenerField {
	for _, f := range fields {
		if f.Decimals < 0 {
			panic("mujica: ScreenerField Decimals is negative")
		}
		if f.ID == id {
			return f
		}
	}
	panic("mujica: Screener condition on unknown field " + id)
}

func screenerFilter(s *ScreenerState, rows []ScreenerRow, version int) []ScreenerRow {
	key := fmt.Sprint(chart.ChartID(rows), version, s.Conditions)
	if key != s.cache.key {
		s.cache.key, s.cache.rows = key, nil
		for _, r := range rows {
			if ScreenerMatch(r, s.Conditions) {
				s.cache.rows = append(s.cache.rows, r)
			}
		}
	}
	return s.cache.rows
}
