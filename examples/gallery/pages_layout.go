package main

import (
	"fmt"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/layout"
)

var layoutPages = []page{
	{ID: "001", Name: "AppShell", ZhName: "应用外壳", Category: "Layout", Demo: layoutAppShellDemo, Code: layoutAppShellCode,
		Keys: "Cmd/Ctrl+B or the toggle button folds the sidebar. Tab moves between the toggle, the sidebar and the content; arrows, Page Up/Down, Home and End scroll the focused region."},
	{ID: "002", Name: "TitleBar", ZhName: "标题栏", Category: "Layout", Demo: layoutTitleBarDemo, Code: layoutTitleBarCode,
		Keys: "Tools in the slots take Tab and Enter/Space as usual. Dragging the empty bar moves the window; double-click zooms it."},
	{ID: "003", Name: "StatusBar", ZhName: "状态栏", Category: "Layout", Demo: layoutStatusBarDemo, Code: layoutStatusBarCode,
		Keys: "Tab reaches action items and the +N button; Enter or Space activates them. +N opens the folded items, Tab moves through them, Escape closes the list. Plain status text is not a tab stop."},
	{ID: "004", Name: "Container", ZhName: "容器", Category: "Layout", Demo: layoutContainerDemo, Code: layoutContainerCode,
		Keys: "None: Container never takes the focus."},
	{ID: "005", Name: "Stack", ZhName: "堆叠", Category: "Layout", Demo: layoutStackDemo, Code: layoutStackCode,
		Keys: "None: Stack never takes the focus; Tab follows the children's order."},
	{ID: "006", Name: "Grid", ZhName: "网格", Category: "Layout", Demo: layoutGridDemo, Code: layoutGridCode,
		Keys: "None: Grid never takes the focus; Tab follows the build order."},
	{ID: "007", Name: "ScrollArea", ZhName: "滚动区域", Category: "Layout", Demo: layoutScrollAreaDemo, Code: layoutScrollAreaCode,
		Keys: "When the content overflows, Tab focuses the area; arrows scroll by a line, Page Up/Down and Space by a page, Home and End to the ends."},
	{ID: "008", Name: "SplitPane", ZhName: "分栏", Category: "Layout", Demo: layoutSplitPaneDemo, Code: layoutSplitPaneCode,
		Keys: "Tab focuses the divider. Left/Right (Up/Down when vertical) move it by 4 units, Home and End to the limits, Enter folds the collapsible pane; an arrow past its minimum folds it, the opposite arrow opens it."},
	{ID: "009", Name: "AspectRatio", ZhName: "宽高比", Category: "Layout", Demo: layoutAspectRatioDemo, Code: layoutAspectRatioCode,
		Keys: "None: AspectRatio never takes the focus."},
	{ID: "010", Name: "Divider", ZhName: "分隔线", Category: "Layout", Demo: layoutDividerDemo, Code: layoutDividerCode,
		Keys: "None: Divider and its ornament never take the focus."},
}

// layoutCheck is a check box that logs its changes.
func layoutCheck(c *ui.Context, v *bool, label string, log func(string)) {
	was := *v
	ui.Checkbox(c, v, label)
	if *v != was {
		log(fmt.Sprintf("%s: %v", label, *v))
	}
}

// layoutChoice is a labeled drop-down that logs its changes.
func layoutChoice(c *ui.Context, v *string, label string, options []string, log func(string)) {
	was := *v
	ui.Row(c).Gap(8).Children(func() { ui.Text(c, label); ui.Select(c, v, options) })
	if *v != was {
		log(fmt.Sprintf("%s: %s", label, *v))
	}
}

// layoutSlider is a labeled slider that logs whole-number changes.
func layoutSlider(c *ui.Context, v *float64, label string, lo, hi float64, log func(string)) {
	was := math.Round(*v)
	ui.Row(c).Gap(8).Children(func() {
		ui.Text(c, fmt.Sprintf("%s %.0f", label, *v))
		ui.Slider(c, v, lo, hi).Width(140)
	})
	if math.Round(*v) != was {
		log(fmt.Sprintf("%s: %.0f", label, *v))
	}
}

func layoutControls(c *ui.Context, fn func()) { ui.Row(c).Wrap().Gap(12).Children(fn) }

func layoutFrame(c *ui.Context, h float32, fn func()) *ui.Element {
	return ui.Column(c).Height(h).Border(1, core.Tokens(c).Border).Radius(4).Clip().Children(fn)
}

var layoutSpaces = []string{"Default", "None", "XS", "S", "M", "L", "XL", "XXL"}

func layoutSpace(s string) layout.Space { return layout.Space(slices.Index(layoutSpaces, s)) }

const layoutAppShellCode = `var collapsed bool
layout.AppShell(c, &collapsed, layout.AppShellOptions{
	TitleBar: func() { layout.TitleBar(c, "Court Records", layout.TitleBarOptions{}) },
	Sidebar:  func() { ui.Text(c, "Chronicles") },
	Content:  func() { ui.Text(c, "The Grand Hall") },
})`

func layoutAppShellDemo(c *ui.Context, log func(string)) {
	type state struct {
		collapsed, title, status bool
		width                    float64
		scroll                   ui.ScrollState
	}
	s := ui.Local(c.Root(), "001", func() state { return state{title: true, status: true, width: 220} })
	layoutControls(c, func() {
		layoutCheck(c, &s.title, "Title bar", log)
		layoutCheck(c, &s.status, "Status bar", log)
		layoutCheck(c, &s.collapsed, "Collapsed", log)
		layoutSlider(c, &s.width, "Sidebar width", 160, 320, log)
	})
	opts := layout.AppShellOptions{
		SidebarWidth:  float32(s.width),
		ContentScroll: &s.scroll,
		Sidebar: func() {
			for _, n := range []string{"Chronicles", "Letters", "Decrees", "Portraits", "Ledgers"} {
				ui.Text(c, n).Padding(6, 12)
			}
		},
		Content: func() {
			layout.Container(c, layout.ContainerOptions{}, func() {
				for i := range 30 {
					ui.Text(c, fmt.Sprintf("Entry %d of the court record", i+1)).Padding(4, 0)
				}
			})
		},
	}
	if s.title {
		opts.TitleBar = func() { layout.TitleBar(c, "Court Records", layout.TitleBarOptions{}) }
	}
	if s.status {
		opts.StatusBar = func() {
			layout.StatusBar(c, layout.StatusBarOptions{Left: []layout.StatusItem{{ID: "pos", Text: fmt.Sprintf("Scrolled %.0f", s.scroll.Y)}}})
		}
	}
	ui.Column(c).Margin(16, 0, 0, 0).Children(func() {
		layoutFrame(c, 360, func() {
			if layout.AppShell(c, &s.collapsed, opts).Toggled {
				log(fmt.Sprintf("sidebar collapsed: %v", s.collapsed))
			}
		})
	})
}

const layoutTitleBarCode = `layout.TitleBar(c, "Court Records", layout.TitleBarOptions{
	Trailing: func() { input.Button(c, "Share", input.ButtonOptions{Variant: input.Ghost}) },
})`

func layoutTitleBarDemo(c *ui.Context, log func(string)) {
	type state struct {
		leading, trailing bool
		title             string
	}
	s := ui.Local(c.Root(), "002", func() state { return state{trailing: true, title: "Court Records"} })
	layoutControls(c, func() {
		layoutCheck(c, &s.leading, "Leading tool", log)
		layoutCheck(c, &s.trailing, "Trailing tool", log)
		layoutChoice(c, &s.title, "Title", []string{"Court Records", "宫廷档案", "A very long title of the royal chronicle of the court"}, log)
	})
	var opts layout.TitleBarOptions
	if s.leading {
		opts.Leading = func() {
			if input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("panel-left"), Label: "Sidebar"}).Clicked() {
				log("leading tool clicked")
			}
		}
	}
	if s.trailing {
		opts.Trailing = func() {
			if input.Button(c, "Share", input.ButtonOptions{Variant: input.Ghost}).Clicked() {
				log("Share clicked")
			}
		}
	}
	ui.Column(c).Margin(16, 0, 0, 0).Width(520).Children(func() {
		layoutFrame(c, 90, func() { layout.TitleBar(c, s.title, opts) })
	})
	bar := c.TitleBar()
	ui.Text(c, fmt.Sprintf("Window controls reserve: left %.0f, right %.0f, height %.0f DIP", bar.Left, bar.Right, bar.Height)).
		TextColor(core.Tokens(c).TextMuted).Margin(8, 0, 0, 0)
}

const layoutStatusBarCode = `bar := layout.StatusBar(c, layout.StatusBarOptions{
	Left:  []layout.StatusItem{{ID: "sync", Text: "Syncing", ShowProgress: true, Progress: 0.4}},
	Right: []layout.StatusItem{{ID: "bell", Label: "Notifications", Icon: icons.Must("bell"), Action: true}},
})
if bar.Clicked == "bell" {
	opened++
}`

func layoutStatusBarDemo(c *ui.Context, log func(string)) {
	type state struct {
		width, progress float64
		showProgress    bool
		tone            string
	}
	s := ui.Local(c.Root(), "003", func() state { return state{width: 640, progress: 40, showProgress: true, tone: "Warning"} })
	layoutControls(c, func() {
		layoutSlider(c, &s.width, "Width", 240, 720, log)
		layoutCheck(c, &s.showProgress, "Progress", log)
		layoutSlider(c, &s.progress, "Percent", 0, 100, log)
		layoutChoice(c, &s.tone, "Tone", []string{"Neutral", "Success", "Warning", "Danger", "Info"}, log)
	})
	tones := map[string]layout.StatusItemTone{"Neutral": layout.StatusItemNeutral, "Success": layout.StatusItemSuccess,
		"Warning": layout.StatusItemWarning, "Danger": layout.StatusItemDanger, "Info": layout.StatusItemInfo}
	opts := layout.StatusBarOptions{
		Left: []layout.StatusItem{
			{ID: "branch", Text: "main", Icon: icons.Must("folder")},
			{ID: "sync", Text: "Syncing", ShowProgress: s.showProgress, Progress: float32(s.progress / 100)},
			{ID: "lint", Text: "2 warnings", Tone: tones[s.tone], Secondary: true, Action: true},
		},
		Right: []layout.StatusItem{
			{ID: "enc", Text: "UTF-8", Secondary: true},
			{ID: "pos", Text: "Ln 12, Col 4", Secondary: true},
			{ID: "bell", Label: "Notifications", Icon: icons.Must("bell"), Action: true},
		},
	}
	ui.Column(c).Margin(16, 0, 0, 0).Width(float32(s.width)).Children(func() {
		if id := layout.StatusBar(c, opts).Clicked; id != "" {
			log("clicked " + id)
		}
	})
	ui.Text(c, "Narrow the width: secondary items fold into a count, last first.").
		TextColor(core.Tokens(c).TextMuted).Margin(8, 0, 0, 0)
}

const layoutContainerCode = `layout.Container(c, layout.ContainerOptions{MaxWidth: 320, Center: true}, func() {
	ui.Text(c, "A centered column of prose.")
})`

func layoutContainerDemo(c *ui.Context, log func(string)) {
	type state struct {
		maxWidth                 float64
		padding                  string
		center, fill, clip, wide bool
	}
	s := ui.Local(c.Root(), "004", func() state { return state{maxWidth: 360, padding: "Default", center: true} })
	layoutControls(c, func() {
		layoutSlider(c, &s.maxWidth, "Max width", 0, 600, log)
		layoutChoice(c, &s.padding, "Padding", layoutSpaces, log)
		layoutCheck(c, &s.center, "Center", log)
		layoutCheck(c, &s.fill, "Fill", log)
		layoutCheck(c, &s.clip, "Clip", log)
		layoutCheck(c, &s.wide, "Overflowing child", log)
	})
	k := core.Tokens(c)
	ui.Column(c).Margin(16, 0, 0, 0).Children(func() {
		layoutFrame(c, 200, func() {
			layout.Container(c, layout.ContainerOptions{MaxWidth: float32(s.maxWidth), Padding: layoutSpace(s.padding),
				Center: s.center, Fill: s.fill, Clip: s.clip}, func() {
				ui.Text(c, "Of the Court").FontSize(18)
				ui.Text(c, "Prose keeps a readable measure inside a container.").TextColor(k.TextMuted)
				if s.wide {
					ui.Box(c).Size(900, 24).Background(k.Selection).Shrink(0)
				}
			}).Background(k.Surface).Border(1, k.Border)
		})
	})
}

const layoutStackCode = `layout.Stack(c, layout.StackOptions{Horizontal: true, Gap: layout.SpaceM}, func() {
	input.Button(c, "Save", input.ButtonOptions{})
	if showDraft {
		input.Button(c, "Draft", input.ButtonOptions{Variant: input.Secondary})
	}
	input.Button(c, "Cancel", input.ButtonOptions{Variant: input.Ghost})
})`

func layoutStackDemo(c *ui.Context, log func(string)) {
	type state struct {
		horizontal, wrap, hide bool
		gap, align, justify    string
	}
	s := ui.Local(c.Root(), "005", func() state {
		return state{horizontal: true, gap: "Default", align: "Auto", justify: "Start"}
	})
	aligns := []string{"Auto", "Start", "Center", "End", "Stretch"}
	justifies := map[string]ui.Align{"Start": ui.Start, "Center": ui.Center, "End": ui.End, "SpaceBetween": ui.SpaceBetween}
	layoutControls(c, func() {
		layoutCheck(c, &s.horizontal, "Horizontal", log)
		s.wrap = s.wrap && s.horizontal
		if s.horizontal {
			layoutCheck(c, &s.wrap, "Wrap", log)
		}
		layoutCheck(c, &s.hide, "Hide Draft", log)
		layoutChoice(c, &s.gap, "Gap", layoutSpaces, log)
		layoutChoice(c, &s.align, "Align", aligns, log)
		layoutChoice(c, &s.justify, "Justify", []string{"Start", "Center", "End", "SpaceBetween"}, log)
	})
	align := layout.StackAlign(slices.Index(aligns, s.align))
	ui.Column(c).Margin(16, 0, 0, 0).Width(420).Children(func() {
		layoutFrame(c, 200, func() {
			layout.Stack(c, layout.StackOptions{Horizontal: s.horizontal, Wrap: s.wrap, Gap: layoutSpace(s.gap),
				Align: align, Justify: justifies[s.justify]}, func() {
				for _, n := range []string{"Save", "Draft", "Publish", "Archive", "Cancel"} {
					if (!s.hide || n != "Draft") && input.Button(c, n, input.ButtonOptions{Variant: input.Secondary}).Clicked() {
						log(n + " clicked")
					}
				}
			}).Grow(1).Padding(12)
		})
	})
}

const layoutGridCode = `layout.Grid(c, layout.GridOptions{Columns: []ui.Track{ui.Fixed(96), ui.Fr(1)}}, func() {
	layout.GridCell(c, layout.GridCellOptions{ColumnSpan: -1}, func() { ui.Text(c, "Header") })
	ui.Text(c, "Name")
	ui.Text(c, "Ada")
})`

func layoutGridDemo(c *ui.Context, log func(string)) {
	type state struct {
		first, gap, rowGap string
		header, tall       bool
	}
	s := ui.Local(c.Root(), "006", func() state {
		return state{first: "Fixed 120", gap: "Default", rowGap: "Default", header: true, tall: true}
	})
	tracks := map[string]ui.Track{"Fixed 120": ui.Fixed(120), "Fit content": ui.FitContent(), "Fr 1": ui.Fr(1), "Fr 2": ui.Fr(2)}
	layoutControls(c, func() {
		layoutChoice(c, &s.first, "First column", []string{"Fixed 120", "Fit content", "Fr 1", "Fr 2"}, log)
		layoutChoice(c, &s.gap, "Gap", layoutSpaces, log)
		layoutChoice(c, &s.rowGap, "Row gap", layoutSpaces, log)
		layoutCheck(c, &s.header, "Header spans all", log)
		layoutCheck(c, &s.tall, "First cell spans 2 rows", log)
	})
	k := core.Tokens(c)
	cell := func(text string) {
		ui.Box(c).Center().Padding(8).Background(k.Surface).Border(1, k.Border).Radius(4).
			Children(func() { ui.Text(c, text).TextColor(k.TextMuted) })
	}
	ui.Column(c).Margin(16, 0, 0, 0).Width(520).Children(func() {
		layout.Grid(c, layout.GridOptions{Columns: []ui.Track{tracks[s.first], ui.Fr(1), ui.Fr(1)},
			Gap: layoutSpace(s.gap), RowGap: layoutSpace(s.rowGap)}, func() {
			if s.header {
				layout.GridCell(c, layout.GridCellOptions{ColumnSpan: -1}, func() { cell("Header, ColumnSpan -1") })
			}
			span := 1
			if s.tall {
				span = 2
			}
			layout.GridCell(c, layout.GridCellOptions{RowSpan: span}, func() { cell(s.first) })
			for i := range 4 {
				cell(fmt.Sprintf("Cell %d", i+1))
			}
		})
	})
}

const layoutScrollAreaCode = `var offset ui.ScrollState
layout.ScrollArea(c, &offset, layout.ScrollAreaOptions{Label: "Log"}, func() {
	for range 50 {
		ui.Text(c, "A line of the log")
	}
}).Height(200)`

func layoutScrollAreaDemo(c *ui.Context, log func(string)) {
	type state struct {
		direction    string
		offset       ui.ScrollState
		lastX, lastY float32
	}
	s := ui.Local(c.Root(), "007", func() state { return state{direction: "Vertical"} })
	dirs := map[string]layout.ScrollDirection{"Vertical": layout.ScrollVertical, "Horizontal": layout.ScrollHorizontal, "Both": layout.ScrollBoth}
	layoutControls(c, func() {
		layoutChoice(c, &s.direction, "Direction", []string{"Vertical", "Horizontal", "Both"}, log)
		if input.Button(c, "Top", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.offset.X, s.offset.Y = 0, 0
			log("scroll to start")
		}
		if input.Button(c, "Middle", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.offset.X, s.offset.Y = s.offset.MaxX/2, s.offset.MaxY/2
			log("scroll to middle")
		}
		if input.Button(c, "End", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.offset.X, s.offset.Y = math.MaxFloat32, math.MaxFloat32
			log("scroll to end")
		}
	})
	k := core.Tokens(c)
	ui.Column(c).Margin(16, 0, 0, 0).Width(420).Children(func() {
		layout.ScrollArea(c, &s.offset, layout.ScrollAreaOptions{Direction: dirs[s.direction], Label: "Ledger"}, func() {
			for i := range 40 {
				ui.Text(c, fmt.Sprintf("%02d  Ledger line of the treasury, folio %d, carried over from the last quarter", i+1, i+7)).
					SingleLine().Padding(4, 12)
			}
		}).Height(220).Background(k.Surface).Border(1, k.Border).Radius(4)
	})
	if s.offset.Y != s.lastY || s.offset.X != s.lastX {
		s.lastX, s.lastY = s.offset.X, s.offset.Y
		log(fmt.Sprintf("offset %.0f, %.0f", s.offset.X, s.offset.Y))
	}
	ui.Text(c, fmt.Sprintf("Offset %.0f, %.0f of %.0f, %.0f", s.offset.X, s.offset.Y, s.offset.MaxX, s.offset.MaxY)).
		TextColor(k.TextMuted).Margin(8, 0, 0, 0)
}

const layoutSplitPaneCode = `split := layout.SplitPaneState{Size: 160}
layout.SplitPane(c, &split, layout.SplitPaneOptions{MinFirst: 120, Collapsible: layout.SplitPaneFirst},
	func() { ui.Text(c, "Files") },
	func() { ui.Text(c, "Editor") },
).Fill()`

func layoutSplitPaneDemo(c *ui.Context, log func(string)) {
	type state struct {
		split       layout.SplitPaneState
		vertical    bool
		collapsible string
		minFirst    float64
	}
	s := ui.Local(c.Root(), "008", func() state {
		return state{split: layout.SplitPaneState{Size: 180}, collapsible: "First", minFirst: 120}
	})
	sides := map[string]layout.SplitPaneSide{"Neither": layout.SplitPaneNeither, "First": layout.SplitPaneFirst, "Second": layout.SplitPaneSecond}
	layoutControls(c, func() {
		layoutCheck(c, &s.vertical, "Vertical", log)
		layoutChoice(c, &s.collapsible, "Collapsible", []string{"Neither", "First", "Second"}, log)
		layoutSlider(c, &s.minFirst, "Min first", 40, 200, log)
	})
	s.split.Collapsed = s.split.Collapsed && s.collapsible != "Neither"
	k := core.Tokens(c)
	pane := func(title, body string, bg ui.Color) func() {
		return func() {
			ui.Column(c).Fill().Padding(12).Gap(6).Background(bg).Children(func() {
				ui.Text(c, title).FontSize(16)
				ui.Text(c, body).TextColor(k.TextMuted)
			})
		}
	}
	ui.Column(c).Margin(16, 0, 0, 0).Width(560).Children(func() {
		layoutFrame(c, 280, func() {
			r := layout.SplitPane(c, &s.split, layout.SplitPaneOptions{Vertical: s.vertical, MinFirst: float32(s.minFirst),
				Collapsible: sides[s.collapsible]},
				pane("Files", "Drag the divider, or focus it and use the arrows.", k.Surface),
				pane("Editor", "Enter on the divider folds the collapsible pane.", k.Background))
			r.Element.Fill()
			if r.Resized {
				log(fmt.Sprintf("resized to %.0f", s.split.Size))
			}
			if r.Toggled {
				log(fmt.Sprintf("collapsed: %v", s.split.Collapsed))
			}
		})
	})
}

const layoutAspectRatioCode = `ui.Box(c).Size(320, 240).Children(func() {
	layout.AspectRatio(c, layout.AspectRatioOptions{Ratio: 16.0 / 9}, func() {
		ui.Text(c, "16 : 9")
	})
})`

func layoutAspectRatioDemo(c *ui.Context, log func(string)) {
	type state struct {
		ratio         string
		width, height float64
	}
	s := ui.Local(c.Root(), "009", func() state { return state{ratio: "16 : 9", width: 360, height: 220} })
	ratios := map[string]float32{"1 : 1": 1, "4 : 3": 4.0 / 3, "16 : 9": 16.0 / 9, "3 : 4": 3.0 / 4, "21 : 9": 21.0 / 9}
	layoutControls(c, func() {
		layoutChoice(c, &s.ratio, "Ratio", []string{"1 : 1", "4 : 3", "16 : 9", "3 : 4", "21 : 9"}, log)
		layoutSlider(c, &s.width, "Parent width", 120, 560, log)
		layoutSlider(c, &s.height, "Parent height", 80, 320, log)
	})
	k := core.Tokens(c)
	ui.Box(c).Margin(16, 0, 0, 0).Size(float32(s.width), float32(s.height)).Border(1, k.Border).Children(func() {
		layout.AspectRatio(c, layout.AspectRatioOptions{Ratio: ratios[s.ratio]}, func() {
			ui.Box(c).Fill().Center().Background(k.Surface).Border(1, k.Ornament).Children(func() {
				ui.Text(c, s.ratio).TextColor(k.TextMuted)
			})
		})
	})
}

const layoutDividerCode = `ui.Column(c).Padding(16).Gap(12).Children(func() {
	layout.Divider(c, layout.DividerOptions{Ornament: true})
	layout.Divider(c, layout.DividerOptions{Label: "Chapter II"})
})`

func layoutDividerDemo(c *ui.Context, log func(string)) {
	type state struct {
		vertical, ornament bool
		label              string
	}
	s := ui.Local(c.Root(), "010", func() state { return state{ornament: true, label: "None"} })
	layoutControls(c, func() {
		layoutCheck(c, &s.vertical, "Vertical", log)
		layoutCheck(c, &s.ornament, "Ornament", log)
		if !s.vertical {
			layoutChoice(c, &s.label, "Label", []string{"None", "Chapter II", "第二章"}, log)
		}
	})
	label := s.label
	if label == "None" || s.vertical {
		label = ""
	}
	opts := layout.DividerOptions{Vertical: s.vertical, Ornament: s.ornament, Label: label}
	ui.Column(c).Margin(16, 0, 0, 0).Width(480).Children(func() {
		if s.vertical {
			ui.Row(c).Height(80).Gap(16).Children(func() {
				ui.Text(c, "Overture")
				layout.Divider(c, opts)
				ui.Text(c, "Aria")
			})
			return
		}
		ui.Column(c).Gap(16).Children(func() {
			ui.Text(c, "Overture")
			layout.Divider(c, opts)
			ui.Text(c, "Aria")
		})
	})
}
