package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
	"mujicaui/theme"
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
}

// densities maps the density control's segments to theme densities.
var densities = []theme.Density{theme.Compact, theme.Default, theme.Comfortable}

func newGallery() *gallery {
	g := &gallery{pages: catalog(), byID: map[string]int{}, selected: overtureID, density: 1, scrolls: map[string]*ui.ScrollState{}}
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

func (g *gallery) settings() mujica.Settings {
	s := mujica.Settings{Mode: mujica.Mode(g.mode), Density: densities[g.density], Locale: mujica.EnUS}
	if g.zh {
		s.Locale = mujica.ZhCN
	}
	return s
}

func (g *gallery) view(c *ui.Context) {
	mujica.Use(c, g.settings())
	w, _ := c.Size()
	if narrow := w < 900; narrow != g.autoFold {
		g.autoFold, g.folded = narrow, narrow
	}
	scroll := g.scrolls[g.selected]
	if scroll == nil {
		scroll = &ui.ScrollState{}
		g.scrolls[g.selected] = scroll
	}
	mujica.AppShell(c, &g.folded, mujica.AppShellOptions{
		TitleBar:      func() { g.titleBar(c) },
		Sidebar:       func() { g.sidebar(c) },
		Content:       func() { g.content(c) },
		StatusBar:     func() { g.statusBar(c) },
		ContentScroll: scroll,
	})
	g.commandPalette(c)
	mujica.Toaster(c)
}

func (g *gallery) titleBar(c *ui.Context) {
	mujica.TitleBar(c, "MujicaUI", mujica.TitleBarOptions{Trailing: func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			mode := g.mode
			if mujica.SegmentedControl(c, &g.mode, []mujica.Segment{
				{Label: g.t("System", "跟随系统"), Icon: icons.Must("settings"), HideLabel: true},
				{Label: g.t("Light", "浅色"), Icon: icons.Must("sun"), HideLabel: true},
				{Label: g.t("Dark", "深色"), Icon: icons.Must("moon"), HideLabel: true},
			}, mujica.SegmentedControlOptions{Label: g.t("Theme", "主题")}).Changed() && mode != g.mode {
				g.logFor("shell")("theme " + mujica.Mode(g.mode).String())
				g.overture.switches++
			}
			density := g.density
			if mujica.SegmentedControl(c, &g.density, []mujica.Segment{
				{Label: g.t("Compact", "紧凑")}, {Label: g.t("Default", "默认")}, {Label: g.t("Comfortable", "宽松")},
			}, mujica.SegmentedControlOptions{Label: g.t("Density", "密度")}).Changed() && density != g.density {
				g.logFor("shell")("density " + densities[g.density].String())
				g.overture.switches++
			}
			if mujica.Button(c, g.t("中文", "English"), mujica.ButtonOptions{Variant: mujica.Ghost}).Clicked() {
				g.zh = !g.zh
				g.logFor("shell")(fmt.Sprintf("locale %s", g.settings().Locale))
			}
			if mujica.IconButton(c, icons.Must("command"), g.t("Command palette", "命令面板"), mujica.IconButtonOptions{Variant: mujica.Ghost}).Clicked() {
				g.palette = true
			}
		})
	}})
}

func (g *gallery) statusBar(c *ui.Context) {
	left := []mujica.StatusItem{
		{ID: "count", Text: fmt.Sprintf(g.t("%d components", "%d 个组件"), len(g.pages)), Icon: icons.Must("layout-grid")},
		{ID: "look", Text: fmt.Sprintf("%s · %s", mujica.Mode(g.mode), densities[g.density]), Secondary: true},
	}
	var right []mujica.StatusItem
	if n := len(g.events); n > 0 {
		e := g.events[n-1]
		right = append(right, mujica.StatusItem{ID: "last", Text: e.page + " · " + e.text, Icon: icons.Must("bell"), Secondary: true})
	}
	mujica.StatusBar(c, mujica.StatusBarOptions{Left: left, Right: right})
}

func (g *gallery) content(c *ui.Context) {
	ui.Box(c).Key(g.selected).FillWidth().Children(func() {
		switch g.selected {
		case overtureID:
			g.overturePage(c)
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
