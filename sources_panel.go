package mujica

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// OpenResult is a list's element and the 1-based item the user opened this frame.
type OpenResult struct {
	Element *ui.Element
	opened  int
}

// Opened returns the number of the source opened.
func (r OpenResult) Opened() (int, bool) { return r.opened, r.opened > 0 }

// SourcesPanel lists sources in citation order, numbered from 1; the Active one scrolls into view.
// It grows to fill its container, which needs a height.
func SourcesPanel(c *ui.Context, sources []Source, st *CitationState) OpenResult {
	if st == nil || st.Active < 0 || st.Active > len(sources) {
		panic("mujica: SourcesPanel needs a state whose Active is 0 or a listed source")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r OpenResult
	r.Element = ui.Column(c).Gap(u * 2).Grow(1).MinHeight(0).Role(ui.RoleList).Label(msg(c, "sources_panel.label"))
	shown := ui.Local(r.Element, "active", func() int { return 0 })
	r.Element.Children(func() {
		ui.Text(c, fmt.Sprintf("%s (%d)", msg(c, "sources_panel.label"), len(sources))).Font(theme.SerifFont).
			FontSize(FontSize(c, theme.H3Size)).FontWeight(600)
		if len(sources) == 0 {
			ui.Text(c, msg(c, "sources_panel.empty")).TextColor(k.TextMuted)
			return
		}
		ui.Scroll(c).Grow(1).MinHeight(0).Children(func() {
			ui.Column(c).Gap(u * 2).Children(func() {
				for i, src := range sources {
					card := SourceCard(c, i+1, src, st).Role(ui.RoleListItem)
					if card.Clicked() && !card.IsDisabled() {
						r.opened = i + 1
					}
					if st.Active == i+1 && *shown != st.Active {
						card.ScrollIntoView()
					}
				}
			})
		})
	})
	*shown = st.Active
	return r
}
