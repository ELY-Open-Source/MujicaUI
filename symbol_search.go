package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

type SymbolSearchState struct {
	Query, Market string
	focusFirst    bool
	cursor        int
}

// SymbolSearchOptions configure a SymbolSearch.
type SymbolSearchOptions struct {
	// Markets are the filter choices; nil hides the filter.
	Markets []string
	// Recent shows while the query is empty.
	Recent []MarketSymbol
	// Loading and Error are the caller's search for the query.
	Loading  bool
	Error    string
	Disabled bool
}

// SymbolSearchView is a built SymbolSearch and this frame's events.
type SymbolSearchView struct {
	Element *ui.Element
	Changed bool
	// Chosen is the symbol picked when Picked is set.
	Chosen MarketSymbol
	Picked bool
	Retry  bool
}

func SymbolSearch(c *ui.Context, s *SymbolSearchState, results []MarketSymbol, opts SymbolSearchOptions) SymbolSearchView {
	if s == nil {
		panic("mujica: SymbolSearch needs a state")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var v SymbolSearchView
	v.Element = ui.Column(c).Gap(u * 2).Disabled(opts.Disabled)
	opts.Disabled = v.Element.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	v.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			in := SearchInput(c, &s.Query, SearchInputOptions{Placeholder: msg(c, "symbol_search.placeholder"), Disabled: opts.Disabled})
			in.Element.Grow(1)
			v.Changed = in.Changed()
			if in.Input.Shortcut(0, ui.KeyDown) && !in.Input.IsDisabled() {
				s.focusFirst = true
			}
			if opts.Markets != nil {
				items := []SelectOption[string]{{Value: "", Label: msg(c, "symbol_search.all_markets")}}
				for _, m := range opts.Markets {
					items = append(items, SelectOption[string]{Value: m, Label: m})
				}
				r := Select(c, &s.Market, items, SelectOptions{Label: msg(c, "symbol_search.market"), Disabled: opts.Disabled})
				r.Element.MinWidth(140)
				v.Changed = v.Changed || r.Changed()
			}
		})
		list, title := results, ""
		switch {
		case opts.Error != "":
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
				ui.Text(c, opts.Error).Role(ui.RoleStatus).Grow(1)
				v.Retry = Button(c, msg(c, "symbol_search.retry"), ButtonOptions{Variant: Secondary}).Clicked()
			})
			return
		case opts.Loading && len(results) == 0:
			Spinner(c, SpinnerOptions{Label: msg(c, "symbol_search.searching"), ShowLabel: true})
			return
		case s.Query == "":
			list, title = opts.Recent, msg(c, "symbol_search.recent")
		case len(results) == 0:
			ui.Text(c, msg(c, "symbol_search.no_results")).TextColor(k.TextMuted).Role(ui.RoleStatus)
			return
		}
		if title != "" && len(list) > 0 {
			ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
		rows := ui.Local(v.Element, "mujica.symbol.rows", func() ui.ListState { return ui.ListState{} })
		rows.Key = func(i int) any { return list[i].Market + ":" + list[i].Code }
		rows.Label = func(i int) string { return symbolBadgeName(list[i]) }
		rows.Selected = &s.cursor
		if s.focusFirst {
			s.cursor = 0
			rows.ScrollTo(0, ui.Start)
		}
		restore := dataSelectTheme(c, 0)
		resultList := ui.List(c, rows, len(list), func(i int) {
			if symbolSearchRow(c, list[i], i == s.cursor) {
				s.cursor = i
				v.Chosen, v.Picked = list[i], true
			}
		}).Height(320).MaxHeight(320)
		restore()
		internal.Ring(resultList, k.Focus, theme.ControlRadius)
		if s.focusFirst || v.Picked {
			resultList.Focus()
		}
		if !resultList.IsDisabled() && resultList.Submitted() && s.cursor >= 0 && s.cursor < len(list) {
			v.Chosen, v.Picked = list[s.cursor], true
		}

		s.focusFirst = false
	})
	if v.Element.IsDisabled() {
		v.Changed, v.Picked, v.Retry = false, false, false
	}
	if v.Picked {
		slog.Debug("mujica: symbol_search pick", "code", v.Chosen.Code)
	}
	return v
}

func symbolSearchRow(c *ui.Context, sym MarketSymbol, selected bool) bool {
	k := Tokens(c)
	d := Density(c)
	b := ui.ButtonBase(c).Key(sym.Market+":"+sym.Code).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 3).
		Justify(ui.Start).Radius(theme.ControlRadius).Label(symbolBadgeName(sym)).Role(ui.RoleListItem)
	internal.Ring(b, k.Focus, theme.ControlRadius)

	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		ui.Text(c, sym.Code).Font(theme.MonoFont).FontWeight(700).Width(80).SingleLine()
		ui.Text(c, sym.Name).SingleLine().Grow(1).MinWidth(0).TextColor(k.Text)
		if selected {
			ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
		}
		if sym.Market != "" {
			ui.Text(c, sym.Market).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
	})
	return b.Clicked() && !b.IsDisabled()
}
