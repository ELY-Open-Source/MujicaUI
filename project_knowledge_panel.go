package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// KnowledgeFile is a file in a project's knowledge, sized in tokens by the caller.
type KnowledgeFile struct {
	ID, Name string
	Tokens   int
}

// ProjectKnowledgePanelOptions configure a ProjectKnowledgePanel. Capacity is required.
type ProjectKnowledgePanelOptions struct {
	// Capacity is how many tokens the project's knowledge holds.
	Capacity int
	// Title replaces the built-in "Project knowledge".
	Title    string
	Disabled bool
}

// ProjectKnowledgePanelView is a built panel and what the user asked for.
type ProjectKnowledgePanelView struct {
	Element *ui.Element
	// Add asks the caller to pick files to add.
	Add bool
	// Removed is the ID of a file to remove, empty for none.
	Removed string
}

// ProjectKnowledgePanel lists a project's knowledge files with add, remove and how full it is.
func ProjectKnowledgePanel(c *ui.Context, files []KnowledgeFile, opts ProjectKnowledgePanelOptions) ProjectKnowledgePanelView {
	if opts.Capacity <= 0 {
		panic("mujica: ProjectKnowledgePanel needs a positive Capacity")
	}
	used := 0
	for _, f := range files {
		if f.Tokens < 0 {
			panic("mujica: KnowledgeFile with negative Tokens")
		}
		used += f.Tokens
	}
	k := Tokens(c)
	d := Density(c)
	col := ui.Column(c).Gap(d.Unit()*2).Padding(d.Unit()*3).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.CardRadius).Disabled(opts.Disabled)
	v := ProjectKnowledgePanelView{Element: col}
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			title := opts.Title
			if title == "" {
				title = msg(c, "project_knowledge_panel.title")
			}
			ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).TextColor(k.Text).Grow(1)
			v.Add = Button(c, msg(c, "project_knowledge_panel.add"), ButtonOptions{Variant: Secondary, Icon: icons.Must("plus"), Disabled: opts.Disabled}).Clicked()
		})
		ui.Box(c).FillWidth().Height(9).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.RuleDiamond(p, r, 6, 5, k.Ornament)
		})
		projectKnowledgeMeter(c, k, used, opts.Capacity)
		if len(files) == 0 {
			ui.Text(c, msg(c, "project_knowledge_panel.empty")).TextColor(k.TextMuted).Padding(d.Unit()*2, 0)
			return
		}
		ui.Column(c).Gap(2).FocusGroup(ui.Vertical).Children(func() {
			for _, f := range files {
				ui.Row(c).Key(f.ID).MinHeight(d.ControlHeight()).Gap(d.Unit()*2).AlignItems(ui.Center).Padding(0, d.Unit()).Children(func() {
					ui.Icon(c, icons.Must("file-text")).Size(16, 16).TextColor(k.TextMuted)
					ui.Text(c, f.Name).SingleLine().Grow(1).MinWidth(0)
					ui.Text(c, fmt.Sprintf(msg(c, "project_knowledge_panel.tokens"), f.Tokens)).
						FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
					x := iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "project_knowledge_panel.remove"), f.Name))
					if x.Clicked() && !x.IsDisabled() {
						v.Removed = f.ID
					}
				})
			}
		})
	})
	if v.Add || v.Removed != "" {
		slog.Debug("mujica: project_knowledge_panel action", "add", v.Add, "remove", v.Removed)
	}
	return v
}

// projectKnowledgeMeter shows used of capacity; near and over the limit add a word and an icon, not only color.
func projectKnowledgeMeter(c *ui.Context, k theme.Tokens, used, capacity int) {
	d := Density(c)
	frac := float64(used) / float64(capacity)
	ink, note, icon := k.AccentText, "", ""
	switch {
	case used > capacity:
		ink, note, icon = k.Danger, msg(c, "project_knowledge_panel.over"), "circle-alert"
	case frac >= 0.8:
		ink, note, icon = k.Warning, msg(c, "project_knowledge_panel.near"), "triangle-alert"
	}
	pct := int(frac * 100)
	label := fmt.Sprintf(msg(c, "project_knowledge_panel.usage"), used, capacity, pct)
	ui.Column(c).Gap(d.Unit()).Role(ui.RoleMeter).Label(label).Children(func() {
		ui.Box(c).FillWidth().Height(6).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, k.Border, 3)
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(min(frac, 1)), H: r.H}, ink, 3)
		})
		ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Children(func() {
			if icon != "" {
				ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(ink)
			}
			t := ui.Text(c, label).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			if note != "" {
				t.TextColor(ink)
				ui.Text(c, note).FontSize(FontSize(c, theme.CaptionSize)).TextColor(ink).Role(ui.RoleStatus)
			}
		})
	})
}
