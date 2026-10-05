package mujica

import (
	"cmp"
	"fmt"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// EmptyKind tells why there is nothing to show. The zero value is EmptyNoData.
type EmptyKind int

const (
	// EmptyNoData: the collection itself is empty.
	EmptyNoData EmptyKind = iota
	// EmptyNoResults: a search or filter matched nothing.
	EmptyNoResults
)

// EmptyStateOptions configure an EmptyState. Empty texts use the kind's built-in copy.
type EmptyStateOptions struct {
	Kind        EmptyKind
	Title       string
	Description string
	// Query, for EmptyNoResults, is named in the default description.
	Query string
	// Icon replaces the kind's icon (inbox or search); NoIcon shows none.
	Icon   *ui.SVG
	NoIcon bool
	// Action labels a button, such as "Clear search"; empty shows none.
	Action string
}

// EmptyStateResult is what an EmptyState built and reported this frame.
type EmptyStateResult struct {
	Element *ui.Element
	Action  bool
}

// EmptyState explains an empty view and offers a next step.
func EmptyState(c *ui.Context, opts EmptyStateOptions) EmptyStateResult {
	if opts.NoIcon && opts.Icon != nil {
		panic("mujica: EmptyState with both Icon and NoIcon")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	title, desc, icon := opts.Title, opts.Description, opts.Icon
	switch opts.Kind {
	case EmptyNoData:
		if opts.Query != "" {
			panic("mujica: EmptyState Query needs EmptyNoResults")
		}
		title = cmp.Or(title, msg(c, "empty_state.no_data.title"))
		desc = cmp.Or(desc, msg(c, "empty_state.no_data.description"))
		if icon == nil {
			icon = icons.Must("inbox")
		}
	case EmptyNoResults:
		title = cmp.Or(title, msg(c, "empty_state.no_results.title"))
		if desc == "" && opts.Query != "" {
			desc = fmt.Sprintf(msg(c, "empty_state.no_results.query"), opts.Query)
		}
		desc = cmp.Or(desc, msg(c, "empty_state.no_results.description"))
		if icon == nil {
			icon = icons.Must("search")
		}
	default:
		panic(fmt.Sprintf("mujica: unknown EmptyKind %d", int(opts.Kind)))
	}
	var res EmptyStateResult
	res.Element = ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(u*2).Padding(u*8, u*4).Role(ui.RoleGroup).Label(title)
	res.Element.Children(func() {
		if !opts.NoIcon {
			ui.Icon(c, icon).Size(32, 32).TextColor(k.TextMuted).Margin(0, 0, u*2, 0)
		}
		ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).TextColor(k.Text).TextAlign(ui.Center)
		ui.Text(c, desc).TextColor(k.TextMuted).TextAlign(ui.Center).MaxWidth(360).LineHeight(theme.LineHeight)
		if opts.Action != "" {
			ui.Box(c).Margin(u*2, 0, 0, 0).Children(func() {
				res.Action = Button(c, opts.Action, ButtonOptions{Variant: Secondary}).Clicked()
			})
		}
	})
	return res
}
