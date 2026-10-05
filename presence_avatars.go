package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// Collaborator is someone working on the same document.
type Collaborator struct {
	ID, Name string
	Image    *ui.Bitmap
	// Away marks someone present but idle.
	Away bool
}

// PresenceAvatarsOptions configure PresenceAvatars.
type PresenceAvatarsOptions struct {
	// Max avatars show before the "+N" button; zero is 4.
	Max int
	// Following identifies the followed collaborator, empty for none.
	Following string
	// Size is the avatar diameter; zero is 28.
	Size     float32
	Disabled bool
}

// PresenceAvatarsResult is the avatars and the follow request this frame.
type PresenceAvatarsResult struct {
	Element *ui.Element
	// Follow identifies a collaborator; Unfollow requests stopping.
	Follow   string
	Unfollow bool
}

// PresenceAvatars displays collaborators, overflow count and follow actions.
func PresenceAvatars(c *ui.Context, people []Collaborator, opts PresenceAvatarsOptions) PresenceAvatarsResult {
	if opts.Max < 0 || opts.Size < 0 {
		panic("mujica: PresenceAvatars Max and Size must not be negative")
	}
	limit, size := opts.Max, opts.Size
	if limit == 0 {
		limit = 4
	}
	if size == 0 {
		size = 28
	}
	u := Density(c).Unit()
	var res PresenceAvatarsResult
	res.Element = ui.Row(c).AlignItems(ui.Center).Gap(u).Role(ui.RoleGroup).Label(fmt.Sprintf(msg(c, "presence_avatars.label"), len(people))).
		FocusGroup(ui.Horizontal).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	open := ui.Local(res.Element, "presence_avatars.open", func() bool { return false })
	rows := ui.Local(res.Element, "presence_avatars.rows", func() ui.ListState { return ui.ListState{} })
	*open = *open && len(people) > limit && !disabled
	pick := func(p Collaborator) {
		if disabled {
			return
		}
		if p.ID == opts.Following {
			res.Unfollow = true
		} else {
			res.Follow = p.ID
		}
		*open = false
		slog.Debug("mujica: presence_avatars follow", "id", p.ID, "unfollow", res.Unfollow)
	}
	res.Element.Children(func() {
		for i, p := range people {
			if i >= limit {
				break
			}
			if presenceAvatar(c, p, size, p.ID == opts.Following).Clicked() {
				pick(p)
			}
		}
		if extra := len(people) - limit; extra > 0 {
			more := ui.ButtonBase(c).Size(max(size+4, theme.MinHitSize), max(size+4, theme.MinHitSize)).Radius((size + 4) / 2).Label(fmt.Sprintf(msg(c, "presence_avatars.more"), extra))
			presenceStyle(c, more, size, false)
			more.Children(func() {
				ui.Text(c, fmt.Sprintf("+%d", extra)).FontSize(FontSize(c, theme.CaptionSize)).FontWeight(600)
			})
			if more.Clicked() && !disabled {
				*open = !*open
			}
			Popover(c, more, open, PopoverOptions{Label: msg(c, "presence_avatars.others")}, func() {
				rows.Key = func(i int) any { return people[limit+i].ID }
				ui.List(c, rows, extra, func(i int) {
					p := people[limit+i]
					ui.Row(c).Key(p.ID).Gap(u * 2).AlignItems(ui.Center).Children(func() {
						if presenceAvatar(c, p, 24, p.ID == opts.Following).Clicked() {
							pick(p)
						}
						ui.Text(c, p.Name).SingleLine()
					})
				}).Width(260).Height(min(240, float32(extra)*Density(c).ControlHeight()))
			})
		}
	})
	return res
}

// presenceAvatar builds a collaborator button with presence status.
func presenceAvatar(c *ui.Context, p Collaborator, size float32, following bool) *ui.Element {
	if p.ID == "" {
		panic("mujica: Collaborator needs an ID")
	}
	k := Tokens(c)
	key := "presence_avatars.follow"
	if following {
		key = "presence_avatars.unfollow"
	}
	b := ui.ButtonBase(c).Key(p.ID).Size(max(size+4, theme.MinHitSize), max(size+4, theme.MinHitSize)).Radius((size + 4) / 2).Label(fmt.Sprintf(msg(c, key), p.Name)).Tooltip(p.Name)
	if p.Away {
		b.Description(msg(c, "presence_avatars.away"))
	}
	presenceStyle(c, b, size, following)
	return b.Children(func() {
		Avatar(c, p.Name, AvatarOptions{Image: p.Image, Size: size})
		dot := ui.Box(c).Size(9, 9).Radius(4.5).Absolute().Right(0).Bottom(0).Border(1.5, k.Surface).Role(ui.RoleNone)
		if p.Away {
			dot.Background(k.Surface).Border(1.5, k.TextMuted)
		} else {
			dot.Background(k.Success)
		}
	})
}

// presenceStyle styles hover, focus and following rings.
func presenceStyle(c *ui.Context, b *ui.Element, size float32, following bool) {
	k := Tokens(c)
	b.Shrink(0).Background(k.Surface).FocusRing(false)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if following {
			p.Stroke(r, k.AccentText, r.W/2, 2)
		}
		if b.FocusVisible() {
			internal.FocusRing(p, r, r.W/2, k.Focus)
		}
	})
}
