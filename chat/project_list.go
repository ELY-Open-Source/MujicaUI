package chat

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ChatProject holds conversations; Folder groups projects, empty for none.
type ChatProject struct {
	ID, Name, Folder string
}

// ProjectListOptions configure a ProjectList.
type ProjectListOptions struct {
	Label    string
	Disabled bool
}

// ProjectListState is the chosen project and open folders and projects.
type ProjectListState struct {
	// Selected is the chosen project's ID; empty for none.
	Selected string
	closed   map[string]bool
	open     map[string]bool
}

// ProjectListView is a built ProjectList and this frame's events.
type ProjectListView struct {
	Element ui.Element
	// Changed reports a newly chosen project, in State.Selected.
	Changed bool
	// Opened is a conversation the user opened, empty for none.
	Opened string
	// Moved asks the caller to move MoveID into MoveTo.
	Moved          bool
	MoveID, MoveTo string
}

// ProjectList shows project folders and conversations; drops move them.
func ProjectList(c *ui.Context, s *ProjectListState, projects []ChatProject, conversations func(projectID string) []ChatConversation, opts ProjectListOptions) ProjectListView {
	if s == nil || conversations == nil {
		panic("mujica: ProjectList needs a state and conversations")
	}
	if s.closed == nil {
		s.closed, s.open = map[string]bool{}, map[string]bool{}
	}
	d := core.Density(c)
	col := ui.Column(c).Gap(2).FocusGroup(ui.Vertical).Role(ui.RoleTree).Label(opts.Label).Disabled(opts.Disabled)
	v := ProjectListView{Element: col}
	col.Children(func() {
		seen := map[string]bool{}
		for _, p := range projects {
			if p.Folder == "" || seen[p.Folder] {
				continue
			}
			seen[p.Folder] = true
			open := !s.closed[p.Folder]
			ui.Column(c.Key("folder:" + p.Folder)).Children(func() { projectListFolder(c, p.Folder, &open) })
			s.closed[p.Folder] = !open
			if !open {
				continue
			}
			for _, q := range projects {
				if q.Folder == p.Folder {
					projectListProject(c, s, q, projects, conversations, d.Unit()*6, &v)
				}
			}
		}
		for _, p := range projects {
			if p.Folder == "" {
				projectListProject(c, s, p, projects, conversations, 0, &v)
			}
		}
	})
	if col.IsDisabled() {
		v.Changed, v.Opened, v.Moved = false, "", false
	}
	switch {
	case v.Moved:
		slog.Debug("mujica: project_list move", "conversation", v.MoveID, "project", v.MoveTo)
	case v.Changed:
		slog.Debug("mujica: project_list select", "project", s.Selected)
	}
	return v
}

// projectListFolder is a folder heading on a CollapsibleBase trigger.
func projectListFolder(c *ui.Context, title string, open *bool) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	parts := ui.CollapsibleBase(c, open)
	tr := parts.Trigger.MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).AlignItems(ui.Center).
		Radius(theme.ControlRadius).Label(title)
	internal.Ring(tr, k.Focus, theme.ControlRadius)
	if tr.Hovered() {
		tr.Background(k.SurfaceHover)
	}
	chevron := "chevron-right"
	if *open {
		chevron = "chevron-down"
	}
	tr.Children(func() {
		ui.Icon(c, icons.Must(chevron)).Size(14, 14).TextColor(k.TextMuted)
		ui.Icon(c, icons.Must("folder")).Size(16, 16).TextColor(k.Ornament.Mix(k.Text, 0.3))
		ui.Text(c, title).Font(theme.SerifFont).SingleLine()
	})
	return tr
}

func projectListProject(c *ui.Context, s *ProjectListState, p ChatProject, all []ChatProject, conversations func(string) []ChatConversation, indent float32, v *ProjectListView) {
	k := core.Tokens(c)
	d := core.Density(c)
	convs := conversations(p.ID)
	chosen := s.Selected == p.ID
	var b ui.Element
	line := ui.Row(c.Key("project:" + p.ID)).AlignItems(ui.Center)
	line.Children(func() { b = ui.ButtonBase(c).Grow(1).MinWidth(0).MinHeight(d.ControlHeight()) })
	b.Padding(0, d.Unit()*2, 0, d.Unit()*2+indent).Padding(0, d.Unit()*2, 0, d.Unit()*2+indent).
		Gap(d.Unit() * 2).Justify(ui.Start).Radius(theme.ControlRadius).Role(ui.RoleTreeItem).Label(p.Name).Selected(chosen)
	b.Background(ui.Transparent).TextColor(k.Text)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	_, over := ui.DragOver[ConversationDrag](b)
	switch {
	case chosen:
		b.Background(k.Selection).TextColor(k.AccentText)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	if over {
		b.DrawOver(func(pt *ui.Painter, r ui.Rect) { internal.Frame(pt, r, 1, 5, k.AccentText) })
	}
	if drop, ok := ui.Drop[ConversationDrag](b); ok && !b.IsDisabled() {
		v.Moved, v.MoveID, v.MoveTo = true, drop.ID, p.ID
	}
	if b.Clicked() && !chosen {
		s.Selected, v.Changed = p.ID, true
	}
	open := s.open[p.ID]
	if !b.IsDisabled() && len(convs) > 0 {
		switch {
		case b.Shortcut(0, ui.KeyRight):
			s.open[p.ID] = true
		case b.Shortcut(0, ui.KeyLeft):
			s.open[p.ID] = false
		}
	}
	if len(convs) > 0 {
		b.Description(core.Msg(c, map[bool]string{true: "project_list.expanded", false: "project_list.collapsed"}[open]))
	}
	b.Children(func() {
		ui.Box(c).Size(8, 8).Shrink(0).Draw(func(pt *ui.Painter, r ui.Rect) {
			if chosen {
				pt.FillPath(internal.Diamond(r.X+4, r.Y+4, 8), k.AccentText)
			}
		})
		ui.Text(c, p.Name).SingleLine().Grow(1).MinWidth(0)
		ui.Text(c, fmt.Sprint(len(convs))).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	})
	if len(convs) > 0 {
		chevron := map[bool]string{true: "chevron-down", false: "chevron-right"}[open]
		line.Children(func() {
			t := core.IconAction(c, icons.Must(chevron), fmt.Sprintf(core.Msg(c, map[bool]string{true: "project_list.collapse", false: "project_list.expand"}[open]), p.Name))
			if t.Clicked() && !t.IsDisabled() {
				s.open[p.ID] = !open
			}
		})
	}
	if !open {
		return
	}
	for _, conv := range convs {
		projectListConversation(c, conv, p, all, indent+d.Unit()*6, v)
	}
}

func projectListConversation(c *ui.Context, conv ChatConversation, home ChatProject, all []ChatProject, indent float32, v *ProjectListView) {
	k := core.Tokens(c)
	d := core.Density(c)
	ui.Row(c.Key("conv:" + conv.ID)).AlignItems(ui.Center).Children(func() {
		b := ui.ButtonBase(c).Grow(1).MinWidth(0).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2, 0, d.Unit()*2+indent).
			Justify(ui.Start).Radius(theme.ControlRadius).Label(conv.Title).Drag(ConversationDrag{ID: conv.ID})
		internal.Ring(b, k.Focus, theme.ControlRadius)
		if b.Dragging() {
			b.Opacity(0.4)
		}
		if b.Hovered() {
			b.Background(k.SurfaceHover)
		}
		b.Children(func() {
			ui.Text(c, cmp.Or(conv.Title, core.Msg(c, "conversation_list.untitled"))).SingleLine().TextColor(k.TextMuted)
		})
		if b.Clicked() {
			v.Opened = conv.ID
		}
		conversationItemMenu(c, fmt.Sprintf(core.Msg(c, "project_list.move_to"), conv.Title), func(m *overlay.PopupMenu) {
			m.Group(core.Msg(c, "project_list.move_group"), func() {
				for _, p := range all {
					if p.ID != home.ID && m.Item(p.Name, overlay.PopupMenuItemOptions{}) {
						v.Moved, v.MoveID, v.MoveTo = true, conv.ID, p.ID
					}
				}
			})
		})
	})
}
