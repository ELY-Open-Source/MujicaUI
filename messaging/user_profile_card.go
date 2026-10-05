package messaging

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/display"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// UserProfile is the person a UserProfileCard shows.
type UserProfile struct {
	// Name is required.
	Name     string
	Image    *ui.Bitmap
	Presence PresenceStatus
	// Status is a custom status line, as "In a rehearsal".
	Status string
	// Role, as "Admin"; empty shows none.
	Role string
	// About is a short description.
	About string
}

// UserProfileCardOptions configure a UserProfileCard.
type UserProfileCardOptions struct {
	// Actions builds buttons at the bottom, as Message or Call.
	Actions func()
	// Width is the card's width; zero means 300.
	Width float32
}

// UserProfileCard displays identity, presence, role and action slots.
func UserProfileCard(c *ui.Context, prof UserProfile, opts UserProfileCardOptions) *ui.Element {
	if prof.Name == "" {
		panic("mujica: UserProfileCard needs a Name")
	}
	if opts.Width < 0 {
		panic("mujica: UserProfileCard Width must not be negative")
	}
	w := opts.Width
	if w == 0 {
		w = 300
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	card := ui.Column(c).Width(w).Padding(u*4).Gap(u*3).Radius(theme.CardRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Role(ui.RoleGroup).Label(core.Msg(c, "user_profile_card.label") + ": " + prof.Name)
	card.Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Shrink(0).Children(func() {
				display.Avatar(c, prof.Name, display.AvatarOptions{Image: prof.Image, Size: 56})
				ui.Box(c).Absolute().Right(-1).Bottom(-1).Padding(2).Radius(8).Background(k.Surface).Children(func() {
					OnlineStatus(c, prof.Presence, OnlineStatusOptions{HideText: true, Size: 12})
				})
			})
			ui.Column(c).Gap(u).MinWidth(0).Grow(1).Children(func() {
				ui.Text(c, prof.Name).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).SingleLine()
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Wrap().Children(func() {
					OnlineStatus(c, prof.Presence, OnlineStatusOptions{})
					if prof.Role != "" {
						display.Badge(c, display.BadgeOptions{Text: prof.Role, Tone: display.BadgeNeutral}, nil)
					}
				})
			})
		})
		if prof.Status != "" {
			ui.Text(c, prof.Status).MaxLines(2)
		}
		if prof.About != "" {
			ui.Text(c, prof.About).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).MaxLines(4)
		}
		if opts.Actions != nil {
			ui.Box(c).Height(1).Background(k.Border)
			ui.Row(c).Gap(u * 2).Wrap().Children(opts.Actions)
		}
	})
	return card
}
