package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
	"mujicaui/theme"
)

// componentPage shows one component: demo, code, keys and its events.
func (g *gallery) componentPage(c *ui.Context, p page) {
	k := mujica.Tokens(c)
	mujica.Container(c, mujica.ContainerOptions{MaxWidth: 1080, Center: true}, func() {
		ui.Column(c).FillWidth().Gap(24).Children(func() {
			mujica.Heading(c, g.pageName(p), mujica.HeadingOptions{
				Subtitle: fmt.Sprintf("%s · %s · %s", p.ID, g.t(p.ZhName, p.Name), g.categoryName(p.Category)),
				Divider:  true,
			})
			ui.Column(c).FillWidth().Gap(16).Children(func() { p.Demo(c, g.logFor(p.ID)) })
			mujica.Grid(c, mujica.GridOptions{Columns: []ui.Track{ui.Fr(3), ui.Fr(2)}, Gap: mujica.SpaceL}, func() {
				mujica.Card(c, mujica.CardOptions{Title: g.t("Code", "代码"), Footer: func() {
					if mujica.Button(c, g.t("Copy", "复制"), mujica.ButtonOptions{Variant: mujica.Secondary, Icon: icons.Must("copy")}).Clicked() {
						c.WriteClipboard(p.Code)
						g.logFor(p.ID)("code copied")
						mujica.ShowToast(c, mujica.ToastMessage{ID: "copied", Severity: mujica.SeveritySuccess, Message: g.t("Code copied", "代码已复制")})
					}
				}}, func() {
					ui.Text(c, p.Code).Font(theme.MonoFont).FontSize(mujica.FontSize(c, theme.CaptionSize)).Selectable().FillWidth()
				})
				mujica.Card(c, mujica.CardOptions{Title: g.t("Keyboard", "键盘")}, func() {
					ui.Text(c, p.Keys).TextColor(k.TextMuted).FillWidth()
				})
			})
			g.eventCard(c, p.ID)
		})
	})
}

// eventCard lists the newest events of one page.
func (g *gallery) eventCard(c *ui.Context, id string) {
	k := mujica.Tokens(c)
	mujica.Card(c, mujica.CardOptions{Title: g.t("Events", "事件")}, func() {
		shown := 0
		for i := len(g.events) - 1; i >= 0 && shown < 8; i-- {
			e := g.events[i]
			if e.page != id {
				continue
			}
			shown++
			ui.Row(c).Gap(12).Children(func() {
				ui.Text(c, e.at.Format("15:04:05.000")).Font(theme.MonoFont).FontSize(mujica.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				ui.Text(c, e.text).Grow(1)
			})
		}
		if shown == 0 {
			ui.Text(c, g.t("Interact with the demo; its events appear here.", "操作上方示例，事件会显示在这里。")).TextColor(k.TextMuted)
		}
	})
}

// compositePage shows one composition demo above its events.
func (g *gallery) compositePage(c *ui.Context, id string, body func(c *ui.Context, log func(string))) {
	cp := compositionOf(id)
	mujica.Container(c, mujica.ContainerOptions{MaxWidth: 1080, Center: true}, func() {
		ui.Column(c).FillWidth().Gap(24).Children(func() {
			mujica.Heading(c, g.t(cp.Name, cp.ZhName), mujica.HeadingOptions{Subtitle: g.t(cp.About, cp.ZhAbout), Divider: true})
			ui.Column(c).FillWidth().Gap(16).Children(func() { body(c, g.logFor(id)) })
			g.eventCard(c, id)
		})
	})
}
