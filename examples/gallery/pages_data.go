package main

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

const dataCategory = "Data"

var dataPages = []page{
	{ID: "069", Name: "Card", ZhName: "卡片", Category: dataCategory, Demo: dataCardDemo, Code: dataCardCode,
		Keys: "A card takes no focus itself; Tab moves through the controls inside it."},
	{ID: "070", Name: "List", ZhName: "列表", Category: dataCategory, Demo: dataListDemo, Code: dataListCode,
		Keys: "Up/Down (arrows in grid mode), Home/End and PageUp/PageDown move the choice; typing jumps by label; Enter submits; Cmd/Ctrl-click and Shift-click choose several; Cmd/Ctrl+A chooses all."},
	{ID: "071", Name: "DataTable", ZhName: "数据表格", Category: dataCategory, Demo: dataTableDemo, Code: dataTableCode,
		Keys: "Up/Down, Home/End move the choice; Enter submits; Shift and Cmd/Ctrl extend it; click a header to sort, again to reverse; drag header edges to resize and headers to move columns."},
	{ID: "072", Name: "Tree", ZhName: "树", Category: dataCategory, Demo: dataTreeDemo, Code: dataTreeCode,
		Keys: "Up/Down move; Right opens or goes to the first child; Left closes or goes to the parent; Option/Shift with Right/Left opens or closes all inside; typing jumps by label."},
	{ID: "073", Name: "TreeTable", ZhName: "树形表格", Category: dataCategory, Demo: dataTreeTableDemo, Code: dataTreeTableCode,
		Keys: "As Tree for rows; click a header to sort siblings; drag header edges to resize and headers to move columns."},
	{ID: "074", Name: "Accordion", ZhName: "手风琴", Category: dataCategory, Demo: dataAccordionDemo, Code: dataAccordionCode,
		Keys: "Tab reaches a header; Up/Down, Home/End move between headers, skipping disabled ones; Enter or Space opens and closes."},
	{ID: "075", Name: "Collapsible", ZhName: "折叠面板", Category: dataCategory, Demo: dataCollapsibleDemo, Code: dataCollapsibleCode,
		Keys: "Enter or Space on the trigger opens and closes; closed content leaves the Tab order."},
	{ID: "076", Name: "Timeline", ZhName: "时间线", Category: dataCategory, Demo: dataTimelineDemo, Code: dataTimelineCode,
		Keys: "Display only: no keys."},
	{ID: "077", Name: "DescriptionList", ZhName: "描述列表", Category: dataCategory, Demo: dataDescriptionDemo, Code: dataDescriptionCode,
		Keys: "Tab reaches each copy button; Enter or Space copies its value."},
	{ID: "078", Name: "Statistic", ZhName: "统计数值", Category: dataCategory, Demo: dataStatisticDemo, Code: dataStatisticCode,
		Keys: "Display only: no keys."},
	{ID: "079", Name: "ImageViewer", ZhName: "图片查看器", Category: dataCategory, Demo: dataViewerDemo, Code: dataViewerCode,
		Keys: "+ and - zoom, 0 fits; Left/Right switch images; Shift+arrows pan; the wheel zooms; drag pans; Escape closes."},
	{ID: "080", Name: "Carousel", ZhName: "走马灯", Category: dataCategory, Demo: dataCarouselDemo, Code: dataCarouselCode,
		Keys: "With focus inside: Left/Right step, Home/End jump; arrows move among the position dots; Enter or Space activates buttons."},
}

// dataSection lays out a titled block of a demo.
func dataSection(c *ui.Context, title string, fn func()) {
	ui.Column(c).Gap(12).Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size))
		fn()
	})
}

// dataStates lays labeled variants side by side.
func dataStates(c *ui.Context, fn func(state func(name string, body func()))) {
	k := core.Tokens(c)
	ui.Row(c).Gap(24).Wrap().AlignItems(ui.Start).Children(func() {
		fn(func(name string, body func()) {
			ui.Column(c).Gap(6).AlignItems(ui.Start).Children(func() {
				ui.Text(c, name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
				body()
			})
		})
	})
}

func dataControls(c *ui.Context, fn func()) {
	ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Children(fn)
}

var dataStatuses = []string{"Ready", "Loading", "Loading more", "Failed", "Empty"}

// dataStatus maps a status choice to the status and whether rows show.
func dataStatus(s string) (data.DataStatus, bool) {
	switch s {
	case "Ready":
		return data.DataReady, true
	case "Loading":
		return data.DataLoading, false
	case "Loading more":
		return data.DataLoadingMore, true
	case "Failed":
		return data.DataFailed, true
	case "Empty":
		return data.DataReady, false
	}
	panic("gallery: unknown status " + s)
}
