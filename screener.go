package mujica

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// ScreenerField is a numeric field conditions test, such as P/E or market cap.
type ScreenerField struct {
	ID, Label string
	Decimals  int
}

// ScreenerOp compares a field with a value. The zero value is ScreenerAbove.
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

// ScreenerRow is a symbol and its field values by field ID; missing values fail every test.
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
	Table      DataTableState[string]
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
	Status   DataStatus
	Error    string
	Disabled bool
}

// ScreenerView is a built Screener and this frame's events.
type ScreenerView struct {
	Element *ui.Element
	// Changed reports edited conditions; Save asks the caller to store them.
	Changed, Save bool
	Table         DataTableView
}

// Screener builds field, comparison and value conditions and lists the rows that pass them.
func Screener(c *ui.Context, s *ScreenerState, fields []ScreenerField, rows []ScreenerRow, opts ScreenerOptions) ScreenerView {
	if s == nil || len(fields) == 0 {
		panic("mujica: Screener needs a state and fields")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	before := fmt.Sprint(s.Conditions)
	var v ScreenerView
	v.Element = ui.Column(c).Gap(u * 3).MinHeight(0).Disabled(opts.Disabled)
	v.Element.Children(func() {
		fieldItems := make([]SelectOption[string], len(fields))
		for i, f := range fields {
			fieldItems[i] = SelectOption[string]{Value: f.ID, Label: f.Label}
		}
		opItems := make([]SelectOption[ScreenerOp], len(screenerOps))
		for i, o := range screenerOps {
			opItems[i] = SelectOption[ScreenerOp]{Value: ScreenerOp(i), Label: o + " " + msg(c, fmt.Sprintf("screener.op_%d", i))}
		}
		ui.Column(c).Gap(u*2).Padding(u*3).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
			Role(ui.RoleGroup).Label(msg(c, "screener.conditions")).Children(func() {
			if len(s.Conditions) == 0 {
				ui.Text(c, msg(c, "screener.none")).TextColor(k.TextMuted)
			}
			for i := 0; i < len(s.Conditions); i++ {
				cd := &s.Conditions[i]
				f := screenerField(fields, cd.Field)
				removed := false
				ui.Row(c).Key(i).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					Select(c, &cd.Field, fieldItems, SelectOptions{Label: msg(c, "screener.field"), Disabled: opts.Disabled})
					Select(c, &cd.Op, opItems, SelectOptions{Label: msg(c, "screener.compare"), Disabled: opts.Disabled}).Element.MinWidth(120)
					NumberInput(c, &cd.Value, NumberInputOptions{Step: math.Pow(10, -float64(f.Decimals)), Precision: f.Decimals,
						Label: msg(c, "screener.value"), Disabled: opts.Disabled}).Element.Width(150)
					x := iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "screener.remove"), i+1))
					removed = x.Clicked() && !x.IsDisabled()
				})
				if removed {
					s.Conditions = slices.Delete(slices.Clone(s.Conditions), i, i+1)
					i--
				}
			}
			ui.Row(c).Gap(u * 2).Children(func() {
				if Button(c, msg(c, "screener.add"), ButtonOptions{Variant: Secondary, Icon: icons.Must("plus"), Disabled: opts.Disabled}).Clicked() {
					s.Conditions = append(s.Conditions, ScreenerCondition{Field: fields[0].ID})
				}
				v.Save = Button(c, msg(c, "screener.save"), ButtonOptions{Variant: Ghost, Disabled: opts.Disabled}).Clicked()
			})
		})
		matched := screenerFilter(s, rows, opts.Version)
		ui.Text(c, fmt.Sprintf(msg(c, "screener.matches"), len(matched), len(rows))).FontSize(FontSize(c, theme.CaptionSize)).
			TextColor(k.TextMuted).Role(ui.RoleStatus)
		cols := []DataColumn[ScreenerRow]{{ID: "symbol", Title: msg(c, "screener.symbol"), Sortable: true,
			Compare: func(a, b ScreenerRow) int { return cmp.Compare(a.Symbol, b.Symbol) },
			Text:    func(r ScreenerRow) string { return r.Symbol }}}
		for _, f := range fields {
			cols = append(cols, DataColumn[ScreenerRow]{ID: f.ID, Title: f.Label, Align: ui.End, Sortable: true,
				Compare: func(a, b ScreenerRow) int { return cmp.Compare(a.Values[f.ID], b.Values[f.ID]) },
				Text:    func(r ScreenerRow) string { return statisticFormat(r.Values[f.ID], f.Decimals) }})
		}
		v.Table = DataTable(c, &s.Table, matched, DataTableOptions[ScreenerRow, string]{Columns: cols, Key: func(r ScreenerRow) string { return r.Symbol },
			Status: opts.Status, Error: opts.Error, Empty: msg(c, "screener.empty")})
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

// screenerFilter keeps passing rows, refiltering only when rows, version or conditions change.
func screenerFilter(s *ScreenerState, rows []ScreenerRow, version int) []ScreenerRow {
	key := fmt.Sprint(chartID(rows), version, s.Conditions)
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
