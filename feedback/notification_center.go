package feedback

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Notification is one in-app message. The caller owns and persists the list.
type Notification struct {
	// ID is required and unique within the list.
	ID    string
	Title string
	Body  string
	// Time shows as given, as "5 min ago".
	Time string
	Read bool
	// Action labels a button; empty shows none.
	Action string
}

// NotificationCenterOptions configure a NotificationCenter.
type NotificationCenterOptions struct {
	// Title uses the built-in "Notifications" when empty.
	Title string
}

// NotificationCenterResult reports the frame's events by notification ID.
type NotificationCenterResult struct {
	Element *ui.Element
	Action  string
	// Toggled changed its Read state; Opened was clicked and is now read.
	Toggled string
	Opened  string
	Removed string
	// MarkedAllRead and ClearedAll report the header buttons.
	MarkedAllRead bool
	ClearedAll    bool
}

// NotificationCenter lists *items and edits them in place: read, unread, remove, clear.
func NotificationCenter(c *ui.Context, items *[]Notification, opts NotificationCenterOptions) NotificationCenterResult {
	if items == nil {
		panic("mujica: NotificationCenter needs items")
	}
	seen := map[string]bool{}
	unread := 0
	for _, n := range *items {
		if n.ID == "" || seen[n.ID] {
			panic(fmt.Sprintf("mujica: notification ID %q empty or repeated", n.ID))
		}
		seen[n.ID] = true
		if !n.Read {
			unread++
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	title := cmp.Or(opts.Title, core.Msg(c, "notification_center.title"))
	var res NotificationCenterResult
	panel := ui.Column(c).FillWidth().Shrink(1).MinHeight(0).Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Role(ui.RoleGroup).Label(title).Clip()
	panel.Children(func() {
		// Wraps so title and actions stay whole when text is enlarged.
		ui.Row(c).Wrap().AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(u*2).Padding(u*3, u*3, u*3, u*4).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
				ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.Text).SingleLine()
				if unread > 0 {
					ui.Text(c, fmt.Sprintf(core.Msg(c, "notification_center.unread"), unread)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.AccentText).SingleLine()
				}
			})
			ui.Row(c).AlignItems(ui.Center).Gap(u * 2).Children(func() {
				if input.Button(c, core.Msg(c, "notification_center.mark_all_read"), input.ButtonOptions{Variant: input.Ghost, Disabled: unread == 0}).Clicked() {
					for i := range *items {
						(*items)[i].Read = true
					}
					res.MarkedAllRead = true
					slog.Debug("mujica: notification_center mark all read", "count", unread)
				}
				if input.Button(c, core.Msg(c, "notification_center.clear_all"), input.ButtonOptions{Variant: input.Ghost, Disabled: len(*items) == 0}).Clicked() {
					slog.Debug("mujica: notification_center clear all", "count", len(*items))
					*items = (*items)[:0]
					res.ClearedAll = true
				}
			})
		})
		// Ornament hairline under the panel title.
		ui.Box(c).FillWidth().Height(1).Background(k.Ornament).Role(ui.RoleNone)
		if len(*items) == 0 && !res.ClearedAll {
			EmptyState(c, EmptyStateOptions{
				Title:       core.Msg(c, "notification_center.empty"),
				Description: core.Msg(c, "notification_center.empty_hint"),
				Icon:        icons.Must("bell"),
			})
			return
		}
		list := ui.Scroll(c).Shrink(1).FillWidth().Role(ui.RoleList).FocusGroup(ui.Vertical | ui.Horizontal)
		list.Dividers(theme.BorderWidth, k.Border)
		list.Children(func() {
			for i := range *items {
				notificationRow(c, k, u, &(*items)[i], &res)
			}
		})
	})
	if res.Removed != "" {
		for i, n := range *items {
			if n.ID == res.Removed {
				*items = append((*items)[:i], (*items)[i+1:]...)
				break
			}
		}
	}
	res.Element = panel
	return res
}

func notificationRow(c *ui.Context, k theme.Tokens, u float32, n *Notification, res *NotificationCenterResult) {
	name := n.Title
	if !n.Read {
		name = core.Msg(c, "notification_center.unread_mark") + ": " + n.Title
	}
	row := ui.Row(c).Key(n.ID).FillWidth().AlignItems(ui.Start).Gap(u*2).Padding(u*3, u*3, u*3, u*2).
		Role(ui.RoleListItem)
	if row.Hovered() {
		row.Background(k.SurfaceHover)
	}
	row.Children(func() {
		// Unread mark: a diamond and bold title, not color alone.
		ui.Box(c).Size(12, 20).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			if !n.Read {
				p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, 8), k.AccentText)
			}
		})
		open := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Grow(1).Shrink(1).Column().AlignItems(ui.Start).Justify(ui.Start).Gap(u).
			Radius(theme.ControlRadius).Label(name)
		internal.Ring(open, k.Focus, theme.ControlRadius)
		open.Children(func() {
			ui.Row(c).FillWidth().Gap(u * 2).Children(func() {
				t := ui.Text(c, n.Title).TextColor(k.Text).Grow(1).Shrink(1)
				if !n.Read {
					t.Bold()
				}
				if n.Time != "" {
					ui.Text(c, n.Time).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				}
			})
			if n.Body != "" {
				ui.Text(c, n.Body).TextColor(k.TextMuted).LineHeight(theme.LineHeight)
			}
		})
		if open.Clicked() && !n.Read {
			n.Read, res.Opened = true, n.ID
			slog.Debug("mujica: notification_center opened", "id", n.ID)
		}
		ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
			if n.Action != "" && input.Button(c, n.Action, input.ButtonOptions{Variant: input.Outline}).Clicked() {
				res.Action = n.ID
				slog.Debug("mujica: notification_center action", "id", n.ID)
			}
			toggle, icon := "notification_center.mark_read", icons.Must("check")
			if n.Read {
				toggle, icon = "notification_center.mark_unread", icons.Must("circle")
			}
			if notificationIconButton(c, k, u, icon, core.Msg(c, toggle)).Clicked() {
				n.Read, res.Toggled = !n.Read, n.ID
				slog.Debug("mujica: notification_center toggled", "id", n.ID, "read", n.Read)
			}
			if notificationIconButton(c, k, u, icons.Must("x"), core.Msg(c, "notification_center.clear")).Clicked() {
				res.Removed = n.ID
			}
		})
	})
	if row.Shortcut(0, ui.KeyDelete) || row.Shortcut(0, ui.KeyBackspace) {
		res.Removed = n.ID
	}
	if res.Removed == n.ID {
		slog.Debug("mujica: notification_center removed", "id", n.ID)
	}
}

func notificationIconButton(c *ui.Context, k theme.Tokens, u float32, icon *ui.SVG, label string) *ui.Element {
	b := input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icon, Label: label})
	b.TextColor(k.TextMuted).Padding(0, u).Tooltip(label)
	return b
}
