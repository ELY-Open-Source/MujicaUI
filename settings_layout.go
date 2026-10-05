package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// SettingsPage is one page of a SettingsLayout's navigation.
type SettingsPage struct {
	ID, Label string
	Icon      *ui.SVG
}

// SettingsLayoutOptions configure a SettingsLayout.
type SettingsLayoutOptions struct {
	// Title heads the navigation; empty uses the built-in "Settings".
	Title string
}

// SettingsLayout builds section navigation and independently scrolling content.
func SettingsLayout(c *ui.Context, page *string, pages []SettingsPage, opts SettingsLayoutOptions, content func(id string)) ChoiceResult {
	if len(pages) == 0 || content == nil {
		panic("mujica: SettingsLayout needs pages and content")
	}
	found := false
	items := make([]SidebarItem, len(pages))
	for i, p := range pages {
		if p.ID == "" || p.Label == "" {
			panic("mujica: SettingsLayout page needs an ID and a Label")
		}
		found = found || p.ID == *page
		items[i] = SidebarItem{ID: p.ID, Label: p.Label, Icon: p.Icon}
	}
	if *page == "" {
		*page, found = pages[0].ID, true
	}
	if !found {
		panic("mujica: SettingsLayout page " + *page + " is not one of pages")
	}
	title := opts.Title
	if title == "" {
		title = msg(c, "settings_layout.title")
	}
	k := Tokens(c)
	d := Density(c)
	var r ChoiceResult
	r.Element = ui.Row(c).Fill().AlignItems(ui.Stretch).Background(k.Background)
	r.Element.Children(func() {
		nav := Sidebar(c, page, []SidebarSection{{Title: title, Items: items}}, SidebarOptions{Width: 220})
		if nav.Changed() {
			r.changed = true
			slog.Debug("mujica: settings_layout page", "page", *page)
		}
		ui.Scroll(c).Grow(1).MinWidth(0).Padding(d.Unit()*6, d.Unit()*8).Children(func() {
			ui.Column(c).Key(*page).MaxWidth(720).Gap(d.Unit() * 8).Children(func() { content(*page) })
		})
	})
	return r
}

// SettingsSection groups rows under a serif heading.
func SettingsSection(c *ui.Context, title string, rows func()) *ui.Element {
	k := Tokens(c)
	d := Density(c)
	e := ui.Column(c).Gap(d.Unit() * 2).Role(ui.RoleGroup).Label(title)
	e.Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).Role(ui.RoleNone)
		ui.Column(c).Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).Padding(0, d.Unit()*4).
			Dividers(theme.BorderWidth, k.Border).Children(rows)
	})
	return e
}

// SettingsRowOptions configure a SettingsRow.
type SettingsRowOptions struct {
	Description string
	// ID identifies this row for SettingsSearch navigation.
	ID string
}

type settingsRevealKey struct{}

// settingsReveal queues focus and scrolling for the matching row.
func settingsReveal(c *ui.Context, id string) {
	*ui.Local(c.Root(), settingsRevealKey{}, func() string { return "" }) = id
}

type settingsRowState struct{ at time.Time }

// SettingsRow associates a label and description with its control.
func SettingsRow(c *ui.Context, label string, opts SettingsRowOptions, control func() *ui.Element) *ui.Element {
	if label == "" || control == nil {
		panic("mujica: SettingsRow needs a label and a control")
	}
	k := Tokens(c)
	d := Density(c)
	row := ui.Row(c).MinHeight(d.ControlHeight()+d.Unit()*4).Gap(d.Unit()*4).AlignItems(ui.Center).Padding(d.Unit()*3, 0)
	if opts.ID != "" {
		row.Key(opts.ID)
	}
	st := ui.Local(row, "reveal", func() settingsRowState { return settingsRowState{} })
	reveal := ui.Local(c.Root(), settingsRevealKey{}, func() string { return "" })
	var ctrl, text *ui.Element
	row.Children(func() {
		text = ui.Column(c).Grow(1).MinWidth(0).Gap(d.Unit())
		text.Children(func() {
			ui.Text(c, label).FontWeight(500).Role(ui.RoleNone)
			if opts.Description != "" {
				ui.Text(c, opts.Description).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize)).Role(ui.RoleNone)
			}
		})
		ui.Row(c).Shrink(0).Children(func() { ctrl = control() })
	})
	if ctrl == nil {
		panic("mujica: SettingsRow " + label + ": control returned nil")
	}
	ctrl.Label(label).Description(opts.Description)
	if text.Clicked() && !ctrl.IsDisabled() {
		ctrl.Focus()
	}
	if opts.ID != "" && *reveal == opts.ID {
		*reveal = ""
		st.at = c.Now()
		ctrl.Focus().ScrollIntoView()
		slog.Debug("mujica: settings_layout reveal", "row", opts.ID)
	}
	// A revealed row's face fades from the selection color.
	if dur := Motion(c, theme.StateDuration); !st.at.IsZero() && dur > 0 {
		if t := c.Now().Sub(st.at); t < dur {
			row.Background(k.Selection.Alpha(1 - float32(t)/float32(dur)))
			c.AnimationFrame()
		}
	}
	return row
}
