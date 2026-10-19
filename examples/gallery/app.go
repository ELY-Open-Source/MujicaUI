package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/feedback"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/layout"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// event is one line of the gallery's event log.
type event struct {
	at   time.Time
	page string
	text string
}

type gallery struct {
	pages    []page
	byID     map[string]int
	selected string
	query    string
	mode     int
	density  int
	zh       bool
	folded   bool
	autoFold bool
	palette  bool
	events   []event
	scrolls  map[string]*ui.ScrollState
	overture overtureState
	// Composition state lives here so it survives page switches.
	cSettings settingsDemo
	cFiles    filesDemo
	cForm     formSubmitDemo
}

// densities maps the density control's segments to theme densities.
var densities = []theme.Density{theme.Compact, theme.Default, theme.Comfortable}

func newGallery() *gallery {
	g := &gallery{pages: catalog(), byID: map[string]int{}, selected: overtureID, density: 1, scrolls: map[string]*ui.ScrollState{},
		cSettings: newSettingsDemo(), cFiles: newFilesDemo(), cForm: newFormSubmitDemo()}
	for i, p := range g.pages {
		if _, dup := g.byID[p.ID]; dup {
			panic("gallery: duplicate page " + p.ID)
		}
		g.byID[p.ID] = i
	}
	return g
}

// logFor returns the event logger of one page.
func (g *gallery) logFor(id string) func(string) {
	return func(s string) {
		g.events = append(g.events, event{time.Now(), id, s})
		if len(g.events) > 200 {
			g.events = g.events[len(g.events)-200:]
		}
		slog.Info("gallery: event", "page", id, "event", s)
	}
}

func (g *gallery) settings() core.Settings {
	s := core.Settings{Mode: core.Mode(g.mode), Density: densities[g.density], Locale: core.EnUS}
	if g.zh {
		s.Locale = core.ZhCN
	}
	return s
}

func (g *gallery) view(c *ui.Context) {
	core.Use(c, g.settings())
	w, _ := c.Size()
	if narrow := w < 900; narrow != g.autoFold {
		g.autoFold, g.folded = narrow, narrow
	}
	scroll := g.scrolls[g.selected]
	if scroll == nil {
		scroll = &ui.ScrollState{}
		g.scrolls[g.selected] = scroll
	}
	layout.AppShell(c, &g.folded, layout.AppShellOptions{
		TitleBar:      func() { g.titleBar(c) },
		Sidebar:       func() { g.sidebar(c) },
		Content:       func() { g.content(c) },
		StatusBar:     func() { g.statusBar(c) },
		ContentScroll: scroll,
	})
	g.commandPalette(c)
	feedback.Toaster(c)
}

func (g *gallery) titleBar(c *ui.Context) {
	layout.TitleBar(c, "MujicaUI", layout.TitleBarOptions{Trailing: func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			mode := g.mode
			if navigation.SegmentedControl(c, &g.mode, []navigation.Segment{
				{Label: g.t("System", "跟随系统"), Icon: icons.Must("settings"), HideLabel: true},
				{Label: g.t("Light", "浅色"), Icon: icons.Must("sun"), HideLabel: true},
				{Label: g.t("Dark", "深色"), Icon: icons.Must("moon"), HideLabel: true},
			}, navigation.SegmentedControlOptions{Label: g.t("Theme", "主题")}).Changed() && mode != g.mode {
				g.logFor("shell")("theme " + core.Mode(g.mode).String())
				g.overture.switches++
			}
			density := g.density
			if navigation.SegmentedControl(c, &g.density, []navigation.Segment{
				{Label: g.t("Compact", "紧凑")}, {Label: g.t("Default", "默认")}, {Label: g.t("Comfortable", "宽松")},
			}, navigation.SegmentedControlOptions{Label: g.t("Density", "密度")}).Changed() && density != g.density {
				g.logFor("shell")("density " + densities[g.density].String())
				g.overture.switches++
			}
			if input.Button(c, g.t("中文", "English"), input.ButtonOptions{Variant: input.Ghost}).Clicked() {
				g.zh = !g.zh
				g.logFor("shell")(fmt.Sprintf("locale %s", g.settings().Locale))
			}
			if input.IconButton(c, icons.Must("command"), g.t("Command palette", "命令面板"), input.IconButtonOptions{Variant: input.Ghost}).Clicked() {
				g.palette = true
			}
		})
	}})
}

func (g *gallery) statusBar(c *ui.Context) {
	left := []layout.StatusItem{
		{ID: "count", Text: fmt.Sprintf(g.t("%d components", "%d 个组件"), len(g.pages)), Icon: icons.Must("layout-grid")},
		{ID: "look", Text: fmt.Sprintf("%s · %s", core.Mode(g.mode), densities[g.density]), Secondary: true},
	}
	var right []layout.StatusItem
	if n := len(g.events); n > 0 {
		e := g.events[n-1]
		right = append(right, layout.StatusItem{ID: "last", Text: e.page + " · " + e.text, Icon: icons.Must("bell"), Secondary: true})
	}
	layout.StatusBar(c, layout.StatusBarOptions{Left: left, Right: right})
}

func (g *gallery) content(c *ui.Context) {
	ui.Box(c.Key(g.selected)).FillWidth().Children(func() {
		switch g.selected {
		case overtureID:
			g.overturePage(c)
		case settingsID:
			g.compositePage(c, settingsID, g.settingsView)
		case filesID:
			g.compositePage(c, filesID, g.filesView)
		case formID:
			g.compositePage(c, formID, g.formView)
		default:
			i, ok := g.byID[g.selected]
			if !ok {
				panic("gallery: no page " + g.selected)
			}
			g.componentPage(c, g.pages[i])
		}
	})
}

// t picks the English or Chinese copy for the gallery's own text.
func (g *gallery) t(en, zh string) string {
	if g.zh {
		return zh
	}
	return en
}
