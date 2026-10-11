package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// LinkPreview is a link's metadata, fetched by the caller.
type LinkPreview struct {
	// URL is the link; required.
	URL string
	// Title is required; Description, Site and Thumbnail may be empty.
	Title, Description, Site string
	Thumbnail                ui.ImageSource
}

// LinkPreviewCard shows a link's site, title, description and thumbnail; Clicked opens it.
func LinkPreviewCard(c *ui.Context, p LinkPreview) *ui.Element {
	if p.URL == "" || p.Title == "" {
		panic("mujica: LinkPreviewCard needs URL and Title")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	b := ui.ButtonBase(c).MaxWidth(420).Gap(u*3).Padding(u*2).AlignItems(ui.Start).Justify(ui.Start).
		Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).Label(p.Title).Tooltip(p.URL).Transition(core.ColorFade)
	internal.Ring(b, k.Focus, theme.CardRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Gap(u).Children(func() {
			site := p.Site
			if site == "" {
				site = p.URL
			}
			ui.Text(c, site).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
			ui.Text(c, p.Title).FontWeight(600).TextColor(k.AccentText).MaxLines(2)
			if p.Description != "" {
				ui.Text(c, p.Description).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).MaxLines(2)
			}
		})
		if p.Thumbnail != nil {
			ui.Image(c, p.Thumbnail).Fit(ui.Cover).Size(64, 64).Radius(theme.ControlRadius).Shrink(0)
		}
	})
	return b
}
