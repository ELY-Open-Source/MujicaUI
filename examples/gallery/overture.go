package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

// overtureState is the continuity proof: values that must survive theme changes.
type overtureState struct {
	note     string
	tempo    float64
	sealed   bool
	switches int
}

func (g *gallery) overturePage(c *ui.Context) {
	k := mujica.Tokens(c)
	w, _ := c.Size()
	mujica.Container(c, mujica.ContainerOptions{MaxWidth: 1080, Center: true}, func() {
		ui.Column(c).FillWidth().Gap(32).Children(func() {
			mujica.Heading(c, "MujicaUI", mujica.HeadingOptions{
				Subtitle: g.t("One hundred court-gothic components for MyGo native UI.", "为 MyGo 原生界面打造的一百件宫廷哥特组件。"),
				Divider:  true,
			})
			ui.Text(c, g.t(
				"Every pixel here is drawn by Go. No HTML, no WebView. Ivory and charcoal, a wine accent, old gold on the edges.",
				"这里每一个像素都由 Go 绘制。没有 HTML，没有 WebView。象牙白与炭黑，酒红为主色，旧金只落在边缘。",
			)).TextColor(k.TextMuted).FillWidth()
			mujica.Grid(c, mujica.GridOptions{Columns: []ui.Track{ui.Fr(1), ui.Fr(1), ui.Fr(1)}}, func() {
				mujica.Statistic(c, float64(len(g.pages)), mujica.StatisticOptions{Title: g.t("Components", "组件")})
				mujica.Statistic(c, 2, mujica.StatisticOptions{Title: g.t("Themes", "主题")})
				mujica.Statistic(c, 3, mujica.StatisticOptions{Title: g.t("Densities", "密度")})
			})
			g.continuity(c)
			g.categoryGrid(c, w)
		})
	})
}

// continuity holds live state while the user flips theme and density.
func (g *gallery) continuity(c *ui.Context) {
	s := &g.overture
	k := mujica.Tokens(c)
	mujica.Card(c, mujica.CardOptions{Title: g.t("Proof of continuity", "连续性证明"), Decorated: true}, func() {
		ui.Column(c).FillWidth().Gap(16).Children(func() {
			ui.Text(c, g.t(
				"Type a line, move the slider, then switch theme or density above, or press ⌘K. Nothing resets: not the text, not the caret, not the focus.",
				"写下一行字，拖动滑块，再切换上方的主题或密度，或按 ⌘K。什么都不会重置：文字、光标、焦点都留在原处。",
			)).FillWidth()
			if mujica.TextInput(c, &s.note, mujica.TextInputOptions{Placeholder: g.t("A line for the court ledger", "写给宫廷账簿的一行"), Label: g.t("Ledger line", "账簿")}).Changed() {
				g.logFor(overtureID)("ledger line edited")
			}
			if mujica.Slider(c, &s.tempo, mujica.SliderOptions{Max: 100, Step: 1, ShowValue: true, Label: g.t("Tempo", "速度")}).Changed() {
				g.logFor(overtureID)(fmt.Sprintf("tempo %.0f", s.tempo))
			}
			if mujica.Checkbox(c, &s.sealed, g.t("Seal with wax", "以火漆封缄"), mujica.CheckboxOptions{}).Changed() {
				g.logFor(overtureID)(fmt.Sprintf("sealed %v", s.sealed))
			}
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				if mujica.Button(c, g.t("Cycle theme", "切换主题"), mujica.ButtonOptions{Icon: icons.Must("moon")}).Clicked() {
					g.mode = (g.mode + 1) % 3
					s.switches++
					g.logFor(overtureID)("theme " + mujica.Mode(g.mode).String())
				}
				if mujica.Button(c, g.t("Cycle density", "切换密度"), mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
					g.density = (g.density + 1) % len(densities)
					s.switches++
					g.logFor(overtureID)("density " + densities[g.density].String())
				}
				ui.Text(c, fmt.Sprintf(g.t("%d switches · %d characters kept", "已切换 %d 次 · 保留 %d 个字符"), s.switches, len([]rune(s.note)))).
					TextColor(k.TextMuted)
			})
		})
	})
}

// categoryGrid lists the ten PRD groups as doors into the gallery.
func (g *gallery) categoryGrid(c *ui.Context, width float32) {
	cols := 3
	if width < 1000 {
		cols = 2
	}
	tracks := make([]ui.Track, cols)
	for i := range tracks {
		tracks[i] = ui.Fr(1)
	}
	var cats []string
	first := map[string]page{}
	count := map[string]int{}
	for _, p := range g.pages {
		if _, ok := first[p.Category]; !ok {
			cats = append(cats, p.Category)
			first[p.Category] = p
		}
		count[p.Category]++
	}
	mujica.Grid(c, mujica.GridOptions{Columns: tracks}, func() {
		for _, cat := range cats {
			ui.Box(c).Key(cat).Children(func() {
				mujica.Card(c, mujica.CardOptions{Title: g.categoryName(cat), Footer: func() {
					if mujica.Button(c, g.t("Enter", "进入"), mujica.ButtonOptions{Variant: mujica.Outline}).Clicked() {
						g.selected = first[cat].ID
						g.logFor(overtureID)("enter " + cat)
					}
				}}, func() {
					ui.Text(c, fmt.Sprintf(g.t("%d components, from %s", "%d 个组件，始于 %s"), count[cat], g.pageName(first[cat]))).
						TextColor(mujica.Tokens(c).TextMuted)
				})
			})
		}
	})
}
