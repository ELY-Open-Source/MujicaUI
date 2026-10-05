package main

import (
	"strings"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

const overtureID = "overture"

// categoryZh names each PRD group in Chinese.
var categoryZh = map[string]string{
	"Layout": "布局与桌面", "Display": "基础展示", "Inputs": "基础输入", "Choices": "选择控件",
	"Forms": "表单组合", "Date & Time": "日期与时间", "Navigation": "导航", "Data": "内容与数据",
	"Overlays": "浮层与菜单", "Feedback": "反馈与状态",
}

// matches reports whether p fits the search query in English or Chinese.
func matches(p page, q string) bool {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return true
	}
	for _, s := range []string{p.ID, p.Name, p.ZhName, p.Category, categoryZh[p.Category]} {
		if strings.Contains(strings.ToLower(s), q) {
			return true
		}
	}
	return false
}

func (g *gallery) categoryName(cat string) string {
	zh, ok := categoryZh[cat]
	if !ok {
		panic("gallery: unknown category " + cat)
	}
	return g.t(cat, zh)
}

func (g *gallery) pageName(p page) string { return g.t(p.Name, p.ZhName) }

// sections groups the pages that match the query by category, in PRD order.
func (g *gallery) sections() []mujica.SidebarSection {
	out := []mujica.SidebarSection{{Title: g.t("Overture", "序章"), Items: []mujica.SidebarItem{
		{ID: overtureID, Label: g.t("Overture", "序章"), Icon: icons.Must("house")},
	}}}
	index := map[string]int{}
	for _, p := range g.pages {
		if !matches(p, g.query) {
			continue
		}
		i, ok := index[p.Category]
		if !ok {
			i = len(out)
			index[p.Category] = i
			out = append(out, mujica.SidebarSection{Title: g.categoryName(p.Category)})
		}
		out[i].Items = append(out[i].Items, mujica.SidebarItem{ID: p.ID, Label: p.ID + "  " + g.pageName(p)})
	}
	return out
}

func (g *gallery) sidebar(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(8).Padding(8, 8, 0, 8).Children(func() {
		before := g.query
		mujica.SearchInput(c, &g.query, mujica.SearchInputOptions{Placeholder: g.t("Search · 搜索", "搜索 · Search")})
		if g.query != before {
			g.logFor("shell")("search " + g.query)
		}
	})
	sections := g.sections()
	if len(sections) == 1 && g.query != "" {
		ui.Text(c, g.t("No component matches.", "没有匹配的组件。")).Padding(8, 16)
	}
	if mujica.Sidebar(c, &g.selected, sections, mujica.SidebarOptions{}).Changed() {
		g.logFor("shell")("open " + g.selected)
	}
}

func (g *gallery) commandPalette(c *ui.Context) {
	cmds := []mujica.Command{
		{ID: overtureID, Label: g.t("Overture", "序章"), Group: g.t("Go to", "前往")},
		{ID: "theme", Label: g.t("Cycle theme", "切换主题"), Group: g.t("Window", "窗口"), Keywords: []string{"dark", "light", "深色", "浅色"}},
		{ID: "density", Label: g.t("Cycle density", "切换密度"), Group: g.t("Window", "窗口"), Keywords: []string{"compact", "紧凑"}},
		{ID: "locale", Label: g.t("Switch to Chinese", "切换到英文"), Group: g.t("Window", "窗口"), Keywords: []string{"language", "语言"}},
	}
	for _, p := range g.pages {
		cmds = append(cmds, mujica.Command{ID: "page:" + p.ID, Label: p.ID + " " + g.pageName(p), Group: g.categoryName(p.Category),
			Keywords: []string{p.Name, p.ZhName, p.Category, categoryZh[p.Category]}})
	}
	id, ok := mujica.CommandPalette(c, &g.palette, cmds, mujica.CommandPaletteOptions{}).Chosen()
	if !ok {
		return
	}
	g.logFor("shell")("command " + id)
	switch {
	case id == overtureID:
		g.selected = overtureID
	case id == "theme":
		g.mode = (g.mode + 1) % 3
		g.overture.switches++
	case id == "density":
		g.density = (g.density + 1) % len(densities)
		g.overture.switches++
	case id == "locale":
		g.zh = !g.zh
	case strings.HasPrefix(id, "page:"):
		g.selected = strings.TrimPrefix(id, "page:")
	default:
		panic("gallery: unknown command " + id)
	}
}
