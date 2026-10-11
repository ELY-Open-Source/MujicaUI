package chat

import (
	"cmp"
	"log/slog"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChatConversation is one saved conversation, as the caller stores it.
type ChatConversation struct {
	ID, Title string
	Updated   time.Time
	Pinned    bool
}

// ConversationListOptions configure a ConversationList.
type ConversationListOptions struct {
	// Now dates the groups; zero uses the window clock.
	Now time.Time
	// Label names the list for assistive technology.
	Label    string
	Disabled bool
}

// ConversationListState is the chosen conversation and the list's place.
type ConversationListState struct {
	// Selected is the chosen conversation's ID; empty for none.
	Selected string
	rows     ui.ListState
	cursor   int
	plan     conversationListPlan
}

// conversationListPlan holds grouped rows, rebuilt on new data or day.
type conversationListPlan struct {
	src   []ChatConversation
	day   core.Date
	rows  []conversationListRow
	built bool
}

// conversationListRow is a header (item < 0) or item index.
type conversationListRow struct {
	group, item int
}

// ConversationListView is a built ConversationList and this frame's events.
type ConversationListView struct {
	Element   *ui.Element
	changed   bool
	submitted bool
}

// Changed reports a newly chosen conversation, in State.Selected.
func (v ConversationListView) Changed() bool { return v.changed }

// Submitted reports Enter or a double click on the choice.
func (v ConversationListView) Submitted() bool { return v.submitted }

// ConversationList groups conversations by recency, building only visible rows.
func ConversationList(c *ui.Context, s *ConversationListState, items []ChatConversation, opts ConversationListOptions, row func(conv ChatConversation)) ConversationListView {
	if s == nil {
		panic("mujica: ConversationList needs a state")
	}
	now := opts.Now
	if now.IsZero() {
		now = c.Now()
	}
	conversationListPlanRows(&s.plan, items, core.DateOf(now))
	rows := s.plan.rows
	s.cursor = slices.IndexFunc(rows, func(r conversationListRow) bool { return r.item >= 0 && items[r.item].ID == s.Selected })
	before := s.cursor
	k := core.Tokens(c)
	d := core.Density(c)
	cursor := &s.cursor
	if opts.Disabled {
		// A disabled list moves a copy of the cursor.
		cp := s.cursor
		cursor = &cp
	}
	s.rows.Key = func(i int) any {
		if r := rows[i]; r.item >= 0 {
			return items[r.item].ID
		}
		return conversationListGroups[rows[i].group]
	}
	s.rows.Header = func(i int) bool { return rows[i].item < 0 }
	s.rows.Label = func(i int) string { return conversationListTitle(c, items, rows[i]) }
	s.rows.Selected = cursor
	var list *ui.Element
	wrap := ui.Column(c).MinHeight(0).Disabled(opts.Disabled)
	wrap.Children(func() {
		restore := data.DataSelectTheme(c, 0)
		list = ui.List(c, &s.rows, len(rows), func(i int) {
			r := rows[i]
			if r.item < 0 {
				conversationListHeader(c, k, core.Msg(c, conversationListGroups[r.group]))
				return
			}
			conv := items[r.item]
			ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(8, 8).Shrink(0).Draw(func(p *ui.Painter, b ui.Rect) {
					if conv.ID == s.Selected {
						p.FillPath(internal.Diamond(b.X+b.W/2, b.Y+b.H/2, 8), k.AccentText)
					}
				})
				if row != nil {
					ui.Row(c).Grow(1).MinWidth(0).Children(func() { row(conv) })
					return
				}
				ui.Text(c, conversationListTitle(c, items, r)).SingleLine().Grow(1).MinWidth(0)
			})
		}).Grow(1).MinHeight(0).Padding(d.Unit()).Background(k.Surface).Border(theme.BorderWidth, k.Border).
			Radius(theme.ControlRadius).Label(opts.Label)
		restore()
		internal.Ring(list, k.Focus, theme.ControlRadius)
	})
	v := ConversationListView{Element: wrap}
	if !wrap.IsDisabled() {
		if s.cursor != before && s.cursor >= 0 && rows[s.cursor].item >= 0 {
			s.Selected, v.changed = items[rows[s.cursor].item].ID, true
			slog.Debug("mujica: conversation_list select", "id", s.Selected)
		}
		v.submitted = list.Submitted() && s.Selected != ""
	}
	return v
}

func conversationListTitle(c *ui.Context, items []ChatConversation, r conversationListRow) string {
	if r.item < 0 {
		return core.Msg(c, conversationListGroups[r.group])
	}
	return cmp.Or(items[r.item].Title, core.Msg(c, "conversation_list.untitled"))
}

func conversationListHeader(c *ui.Context, k theme.Tokens, title string) {
	d := core.Density(c)
	ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).AlignItems(ui.Center).Background(k.Surface).Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Ornament.Mix(k.Text, 0.35)).SingleLine()
	})
}

// conversationListPlanRows regroups items when the slice or the day changed.
func conversationListPlanRows(p *conversationListPlan, items []ChatConversation, today core.Date) {
	same := p.built && p.day == today && len(p.src) == len(items) && (len(items) == 0 || &p.src[0] == &items[0])
	if same {
		return
	}
	p.src, p.day, p.built = items, today, true
	order := make([]int, len(items))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int { return items[b].Updated.Compare(items[a].Updated) })
	groups := [5][]int{}
	for _, i := range order {
		g := conversationListGroup(items[i], today)
		groups[g] = append(groups[g], i)
	}
	p.rows = p.rows[:0]
	for g, idx := range groups {
		if len(idx) == 0 {
			continue
		}
		p.rows = append(p.rows, conversationListRow{group: g, item: -1})
		for _, i := range idx {
			p.rows = append(p.rows, conversationListRow{item: i})
		}
	}
	slog.Debug("mujica: conversation_list regroup", "count", len(items))
}

var conversationListGroups = [5]string{
	"conversation_list.pinned", "conversation_list.today", "conversation_list.yesterday",
	"conversation_list.week", "conversation_list.older",
}

// conversationListGroup returns the recency group index of conv.
func conversationListGroup(conv ChatConversation, today core.Date) int {
	day := core.DateOf(conv.Updated)
	switch {
	case conv.Pinned:
		return 0
	case day.Compare(today) >= 0:
		return 1
	case day == today.AddDays(-1):
		return 2
	case day.Compare(today.AddDays(-7)) >= 0:
		return 3
	}
	return 4
}
