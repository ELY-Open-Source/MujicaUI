package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// componentPage shows one component: demo, code, keys and its events.
func (g *gallery) componentPage(c *ui.Context, p page) {
	k := core.Tokens(c)
	layout.Container(c, layout.ContainerOptions{MaxWidth: 1080, Center: true}, func() {
		ui.Column(c).FillWidth().Gap(24).Children(func() {
			display.Heading(c, g.pageName(p), display.HeadingOptions{
				Subtitle: fmt.Sprintf("%s · %s · %s", p.ID, g.t(p.ZhName, p.Name), g.categoryName(p.Category)),
				Divider:  true,
			})
			ui.Column(c).FillWidth().Gap(16).Children(func() { p.Demo(c, g.logFor(p.ID)) })
			layout.Grid(c, layout.GridOptions{Columns: []ui.Track{ui.Fr(3), ui.Fr(2)}, Gap: layout.SpaceL}, func() {
				data.Card(c, data.CardOptions{Title: g.t("Code", "代码"), Footer: func() {
					if input.Button(c, g.t("Copy", "复制"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("copy")}).Clicked() {
						c.WriteClipboard(p.Code)
						g.logFor(p.ID)("code copied")
						feedback.ShowToast(c, feedback.ToastMessage{ID: "copied", Severity: core.SeveritySuccess, Message: g.t("Code copied", "代码已复制")})
					}
				}}, func() {
					ui.Text(c, p.Code).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).Selectable().FillWidth()
				})
				data.Card(c, data.CardOptions{Title: g.t("Keyboard", "键盘")}, func() {
					ui.Text(c, p.Keys).TextColor(k.TextMuted).FillWidth()
				})
			})
			g.eventCard(c, p.ID)
		})
	})
}

// eventCard lists the newest events of one page.
func (g *gallery) eventCard(c *ui.Context, id string) {
	k := core.Tokens(c)
	data.Card(c, data.CardOptions{Title: g.t("Events", "事件")}, func() {
		shown := 0
		for i := len(g.events) - 1; i >= 0 && shown < 8; i-- {
			e := g.events[i]
			if e.page != id {
				continue
			}
			shown++
			ui.Row(c).Gap(12).Children(func() {
				ui.Text(c, e.at.Format("15:04:05.000")).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
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
	layout.Container(c, layout.ContainerOptions{MaxWidth: 1080, Center: true}, func() {
		ui.Column(c).FillWidth().Gap(24).Children(func() {
			display.Heading(c, g.t(cp.Name, cp.ZhName), display.HeadingOptions{Subtitle: g.t(cp.About, cp.ZhAbout), Divider: true})
			ui.Column(c).FillWidth().Gap(16).Children(func() { body(c, g.logFor(id)) })
			g.eventCard(c, id)
		})
	})
}
