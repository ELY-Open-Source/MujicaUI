package main

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
)

const overtureID = "overture"

// composition is one PRD §8 demo that combines components.
type composition struct {
	ID, Name, ZhName, About, ZhAbout, Icon string
}

var compositions = []composition{
	{settingsID, "Settings", "设置页", "Tabs of labelled controls with save, reset and an unsaved-changes flag.", "分页的带标签控件，可保存、重置，并标出未保存的更改。", "settings"},
	{filesID, "File list", "文件列表", "Filter, select, rename and delete files; drop or pick files to add them.", "筛选、选择、重命名与删除文件；拖入或选取文件以添加。", "folder"},
	{formID, "Form submission", "表单提交", "Validated fields, an async name check, a guarded submit and its result.", "带校验的字段、异步名称检查、防重复提交及其结果。", "file-text"},
}

func compositionOf(id string) composition {
	for _, cp := range compositions {
		if cp.ID == id {
			return cp
		}
	}
	panic("gallery: no composition " + id)
}

// categoryZh names each PRD group in Chinese.
var categoryZh = map[string]string{
	"Layout": "布局与桌面", "Display": "基础展示", "Inputs": "基础输入", "Choices": "选择控件",
	"Forms": "表单组合", "Date & Time": "日期与时间", "Navigation": "导航", "Data": "内容与数据",
	"Overlays": "浮层与菜单", "Feedback": "反馈与状态",
	"Charts I":               "通用图表一",
	"Charts II":              "通用图表二",
	"Market Charts":          "行情图表",
	"Trading":                "报价与交易",
	"Markets":                "市场与资产",
	"Chat Messages":          "对话消息",
	"Chat Input":             "对话输入",
	"Chat Sessions":          "对话会话",
	"Agent":                  "智能体",
	"Agent Ops & Motion":     "智能体运维与动效",
	"Git":                    "版本控制",
	"Terminal & Code":        "终端与代码",
	"Media":                  "媒体",
	"Projects":               "项目管理",
	"Files & Collaboration":  "文件与协作",
	"Messaging & Mail":       "通讯与邮件",
	"Calendar":               "日历",
	"Dashboards & Dev Tools": "仪表盘与开发工具",
	"Settings & Account":     "设置与账户",
	"Essentials":             "通用补充",
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
func (g *gallery) sections() []navigation.SidebarSection {
	out := []navigation.SidebarSection{{Title: g.t("Overture", "序章"), Items: []navigation.SidebarItem{
		{ID: overtureID, Label: g.t("Overture", "序章"), Icon: icons.Must("house")},
	}}}
	var comps []navigation.SidebarItem
	q := strings.ToLower(strings.TrimSpace(g.query))
	for _, cp := range compositions {
		if q == "" || strings.Contains(strings.ToLower(cp.ID+" "+cp.Name+" "+cp.ZhName), q) {
			comps = append(comps, navigation.SidebarItem{ID: cp.ID, Label: g.t(cp.Name, cp.ZhName), Icon: icons.Must(cp.Icon)})
		}
	}
	if len(comps) > 0 {
		out = append(out, navigation.SidebarSection{Title: g.t("Compositions", "组合示例"), Items: comps})
	}
	index := map[string]int{}
	for _, p := range g.pages {
		if !matches(p, g.query) {
			continue
		}
		i, ok := index[p.Category]
		if !ok {
			i = len(out)
			index[p.Category] = i
			out = append(out, navigation.SidebarSection{Title: g.categoryName(p.Category)})
		}
		out[i].Items = append(out[i].Items, navigation.SidebarItem{ID: p.ID, Label: p.ID + "  " + g.pageName(p)})
	}
	return out
}

func (g *gallery) sidebar(c *ui.Context) {
	ui.Column(c).FillWidth().Gap(8).Padding(8, 8, 0, 8).Children(func() {
		before := g.query
		input.SearchInput(c, &g.query, input.SearchInputOptions{Placeholder: g.t("Search · 搜索", "搜索 · Search")})
		if g.query != before {
			g.logFor("shell")("search " + g.query)
		}
	})
	sections := g.sections()
	if len(sections) == 1 && g.query != "" {
		ui.Text(c, g.t("No component matches.", "没有匹配的组件。")).Padding(8, 16)
	}
	if navigation.Sidebar(c, &g.selected, sections, navigation.SidebarOptions{}).Changed() {
		g.logFor("shell")("open " + g.selected)
	}
}

func (g *gallery) commandPalette(c *ui.Context) {
	cmds := []navigation.Command{
		{ID: overtureID, Label: g.t("Overture", "序章"), Group: g.t("Go to", "前往")},
		{ID: "theme", Label: g.t("Cycle theme", "切换主题"), Group: g.t("Window", "窗口"), Keywords: []string{"dark", "light", "深色", "浅色"}},
		{ID: "density", Label: g.t("Cycle density", "切换密度"), Group: g.t("Window", "窗口"), Keywords: []string{"compact", "紧凑"}},
		{ID: "locale", Label: g.t("Switch to Chinese", "切换到英文"), Group: g.t("Window", "窗口"), Keywords: []string{"language", "语言"}},
	}
	for _, cp := range compositions {
		cmds = append(cmds, navigation.Command{ID: "page:" + cp.ID, Label: g.t(cp.Name, cp.ZhName), Group: g.t("Compositions", "组合示例"),
			Keywords: []string{cp.Name, cp.ZhName, "composition", "组合"}})
	}
	for _, p := range g.pages {
		cmds = append(cmds, navigation.Command{ID: "page:" + p.ID, Label: p.ID + " " + g.pageName(p), Group: g.categoryName(p.Category),
			Keywords: []string{p.Name, p.ZhName, p.Category, categoryZh[p.Category]}})
	}
	id, ok := navigation.CommandPalette(c, &g.palette, cmds, navigation.CommandPaletteOptions{}).Chosen()
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
