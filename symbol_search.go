package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// SymbolSearchState is the query and the market filter, empty for all markets.
type SymbolSearchState struct {
	Query, Market string
	focusFirst    bool
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
	// Changed asks the caller to search again for State.Query and State.Market.
	Changed bool
	// Chosen is the symbol picked when Picked is set.
	Chosen MarketSymbol
	Picked bool
	Retry  bool
}

// SymbolSearch finds symbols by code or name, filtered by market, with recent picks; the caller supplies results.
func SymbolSearch(c *ui.Context, s *SymbolSearchState, results []MarketSymbol, opts SymbolSearchOptions) SymbolSearchView {
	if s == nil {
		panic("mujica: SymbolSearch needs a state")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var v SymbolSearchView
	v.Element = ui.Column(c).Gap(u * 2).Disabled(opts.Disabled)
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
		case s.Query == "":
			list, title = opts.Recent, msg(c, "symbol_search.recent")
		case opts.Loading && len(results) == 0:
			Spinner(c, SpinnerOptions{Label: msg(c, "symbol_search.searching"), ShowLabel: true})
			return
		case len(results) == 0:
			ui.Text(c, msg(c, "symbol_search.no_results")).TextColor(k.TextMuted).Role(ui.RoleStatus)
			return
		}
		if title != "" && len(list) > 0 {
			ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
		ui.Scroll(c).MaxHeight(320).Children(func() {
			ui.Column(c).Gap(2).FocusGroup(ui.Vertical).Role(ui.RoleList).Children(func() {
				for i, sym := range list {
					if symbolSearchRow(c, sym, i == 0 && s.focusFirst) {
						v.Chosen, v.Picked = sym, true
					}
				}
			})
		})
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

func symbolSearchRow(c *ui.Context, sym MarketSymbol, focus bool) bool {
	k := Tokens(c)
	d := Density(c)
	b := ui.ButtonBase(c).Key(sym.Market+":"+sym.Code).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 3).
		Justify(ui.Start).Radius(theme.ControlRadius).Label(symbolBadgeName(sym)).Role(ui.RoleListItem)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if focus {
		b.Focus()
	}
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		ui.Text(c, sym.Code).Font(theme.MonoFont).FontWeight(700).Width(80).SingleLine()
		ui.Text(c, sym.Name).SingleLine().Grow(1).MinWidth(0).TextColor(k.Text)
		if sym.Market != "" {
			ui.Text(c, sym.Market).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
	})
	return b.Clicked() && !b.IsDisabled()
}
