package mujica

import (
	"cmp"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ConversationItemOptions configure a ConversationItem.
type ConversationItemOptions struct {
	// Today decides the time stamp; zero uses the clock.
	Today    Date
	Selected bool
	Disabled bool
	// Draggable lets the row be dropped on a ProjectList.
	Draggable bool
}

// ConversationDrag is the value a dragged conversation carries.
type ConversationDrag struct{ ID string }

// ConversationItemView is a built ConversationItem and its requests.
type ConversationItemView struct {
	Element *ui.Element
	// Opened, Deleted, Exported and PinToggled are requests for the caller.
	Opened, Deleted, Exported, PinToggled bool
	// Renamed reports a confirmed new title, already in conv.Title.
	Renamed bool
}

type conversationItemState struct {
	editing, fresh bool
	draft          string
}

type conversationItemKey struct{}

// ConversationItem is a conversation row with rename, pin, export, delete.
func ConversationItem(c *ui.Context, conv *ChatConversation, opts ConversationItemOptions) ConversationItemView {
	if conv == nil || conv.ID == "" {
		panic("mujica: ConversationItem needs a conversation with an ID")
	}
	k := Tokens(c)
	d := Density(c)
	today := opts.Today
	if today.IsZero() {
		today = DateOf(c.Now())
	}
	row := ui.Row(c).Key(conv.ID).FillWidth().MinHeight(d.ControlHeight()).Gap(d.Unit()).AlignItems(ui.Center).
		Radius(theme.ControlRadius).Disabled(opts.Disabled)
	v := ConversationItemView{Element: row}
	st := ui.Local(row, conversationItemKey{}, func() conversationItemState { return conversationItemState{} })
	if row.IsDisabled() {
		st.editing = false
	}
	if opts.Selected {
		row.Background(k.Selection)
	}
	title := cmp.Or(conv.Title, msg(c, "conversation_list.untitled"))
	rename := false
	row.Children(func() {
		if st.editing {
			conversationItemEditor(c, k, conv, st, &v)
		} else {
			main := ui.ButtonBase(c).Grow(1).MinWidth(0).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).
				Justify(ui.Start).Radius(theme.ControlRadius).Label(title)
			internal.Ring(main, k.Focus, theme.ControlRadius)
			if opts.Draggable {
				main.Drag(ConversationDrag{ID: conv.ID})
			}
			if main.Hovered() && !opts.Selected {
				main.Background(k.SurfaceHover)
			}
			main.Transition(colorFade)
			main.Children(func() {
				if conv.Pinned {
					ui.Icon(c, icons.Must("star")).Size(14, 14).Shrink(0).TextColor(k.AccentText).Label(msg(c, "conversation_item.pinned"))
				}
				t := ui.Text(c, title).SingleLine().Grow(1).MinWidth(0)
				if opts.Selected {
					t.TextColor(k.AccentText)
				}
				ui.Text(c, dateSeparatorWhen(c, conv.Updated, today)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Shrink(0)
			})
			if !main.IsDisabled() {
				v.Opened = main.Clicked()
				rename = main.Shortcut(0, ui.KeyF2)
				v.Deleted = main.Shortcut(0, ui.KeyDelete)
			}
		}
		conversationItemMenu(c, msg(c, "conversation_item.actions"), func(m *PopupMenu) {
			rename = m.Item(msg(c, "conversation_item.rename"), PopupMenuItemOptions{}) || rename
			pin := msg(c, "conversation_item.pin")
			if conv.Pinned {
				pin = msg(c, "conversation_item.unpin")
			}
			v.PinToggled = m.Item(pin, PopupMenuItemOptions{})
			v.Exported = m.Item(msg(c, "conversation_item.export"), PopupMenuItemOptions{})
			m.Separator()
			v.Deleted = m.Item(msg(c, "conversation_item.delete"), PopupMenuItemOptions{}) || v.Deleted
		})
	})
	if rename && !row.IsDisabled() {
		st.editing, st.fresh, st.draft = true, true, conv.Title
	}
	switch {
	case v.Opened:
		slog.Debug("mujica: conversation_item open", "id", conv.ID)
	case v.Deleted, v.Exported, v.PinToggled:
		slog.Debug("mujica: conversation_item action", "id", conv.ID, "delete", v.Deleted, "export", v.Exported, "pin", v.PinToggled)
	}
	return v
}

// conversationItemEditor edits the title in place; Escape cancels.
func conversationItemEditor(c *ui.Context, k theme.Tokens, conv *ChatConversation, st *conversationItemState, v *ConversationItemView) {
	d := Density(c)
	in := ui.TextInputBase(c, &st.draft).Grow(1).MinHeight(d.ControlHeight()-4).Padding(0, d.Unit()*2).
		Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.AccentText).
		Label(msg(c, "conversation_item.rename"))
	internal.Ring(in, k.Focus, theme.ControlRadius)
	switch {
	case st.fresh:
		st.fresh = false
		in.Focus()
	case in.IsDisabled():
	case in.Submitted():
		if t := strings.TrimSpace(st.draft); t != "" {
			v.Renamed = t != conv.Title
			conv.Title = t
		}
		st.editing = false
		slog.Debug("mujica: conversation_item rename", "id", conv.ID, "changed", v.Renamed)
	case in.Shortcut(0, ui.KeyEscape):
		st.editing = false
	}
}

// conversationItemMenu is a borderless ellipsis button opening a menu.
func conversationItemMenu(c *ui.Context, label string, build func(m *PopupMenu)) *ui.Element {
	b := iconAction(c, icons.Must("ellipsis"), label).Role(ui.RoleMenuButton).Shrink(0)
	open := ui.Local(b, conversationItemKey{}, func() bool { return false })
	if b.IsDisabled() {
		*open = false
	} else if b.Clicked() || b.Shortcut(0, ui.KeyDown) {
		*open = true
	}
	popupMenuOpen(c, b, open, build)
	return b
}
