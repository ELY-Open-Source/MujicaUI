package chat

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// OpenResult carries the list and opened 1-based source number.
type OpenResult struct {
	Element *ui.Element
	opened  int
}

// Opened returns the number of the source opened.
func (r OpenResult) Opened() (int, bool) { return r.opened, r.opened > 0 }

// SourcesPanel lists citations and scrolls to activated sources.
// It grows to fill its container, which needs a height.
func SourcesPanel(c *ui.Context, sources []Source, st *CitationState) OpenResult {
	if st == nil || st.Active < 0 || st.Active > len(sources) {
		panic("mujica: SourcesPanel needs a state whose Active is 0 or a listed source")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r OpenResult
	r.Element = ui.Column(c).Gap(u * 2).Grow(1).MinHeight(0).Role(ui.RoleList).Label(core.Msg(c, "sources_panel.label"))
	shown := ui.Local(r.Element, "active", func() int { return 0 })
	located := ui.Local(r.Element, "locate", func() uint64 { return 0 })
	list := ui.Local(r.Element, "list", func() ui.ListState { return ui.ListState{} })
	if st.Active > 0 && (*shown != st.Active || *located != st.locate) {
		list.ScrollIntoView(st.Active - 1)
		slog.Debug("mujica: sources_panel locate", "source", st.Active)
	}
	r.Element.Children(func() {
		ui.Text(c, fmt.Sprintf("%s (%d)", core.Msg(c, "sources_panel.label"), len(sources))).Font(theme.SerifFont).
			FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600)
		if len(sources) == 0 {
			ui.Text(c, core.Msg(c, "sources_panel.empty")).TextColor(k.TextMuted)
			return
		}
		ui.List(c, list, len(sources), func(i int) {
			card := SourceCard(c, i+1, sources[i], st).MinHeight(core.Density(c).ControlHeight()).Role(ui.RoleListItem)
			if card.Clicked() && !card.IsDisabled() {
				r.opened = i + 1
			}
		}).Grow(1).MinHeight(0).Gap(u * 2)
	})
	*shown = st.Active
	*located = st.locate
	return r
}
