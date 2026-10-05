package mujica

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// Milestone holds a goal, due date and completion counts.
type Milestone struct {
	Title       string
	Due         Date
	Done, Total int
}

// MilestoneProgressOptions configure a MilestoneProgress.
type MilestoneProgressOptions struct {
	// Today anchors the date note; zero uses the system clock.
	Today Date
}

// MilestoneProgress shows completion counts, progress and due dates.
func MilestoneProgress(c *ui.Context, m Milestone, opts MilestoneProgressOptions) *ui.Element {
	datetimeCheckDate("MilestoneProgress Due", m.Due)
	if m.Total < 0 || m.Done < 0 || m.Done > m.Total {
		panic(fmt.Sprintf("mujica: MilestoneProgress %d of %d is not valid", m.Done, m.Total))
	}
	k := Tokens(c)
	d := Density(c)
	today := datetimeToday(c, opts.Today)
	frac := 0.0
	if m.Total > 0 {
		frac = float64(m.Done) / float64(m.Total)
	}
	reached := m.Total > 0 && m.Done == m.Total
	pct := fmt.Sprintf("%d%%", int(frac*100))
	col := ui.Column(c).Gap(d.Unit()*2).Padding(d.Unit()*3).Radius(theme.CardRadius).Border(theme.BorderWidth, k.Border).
		Background(k.Surface).Role(ui.RoleGroup).Label(m.Title).Description(pct)
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(12, 12).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, r.W), k.Ornament)
			})
			ui.Text(c, m.Title).Font(theme.SerifFont).TextColor(k.Text).Grow(1).SingleLine().Ellipsis("…")
			ui.Text(c, pct).TextColor(k.AccentText).FontFeatures("tnum")
		})
		projectBar(c, frac, m.Title)
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf(msg(c, "milestone_progress.count"), m.Done, m.Total)).FontSize(FontSize(c, theme.CaptionSize)).
				TextColor(k.TextMuted).Grow(1).FontFeatures("tnum")
			text, late := projectDueText(c, m.Due, today, reached)
			fg, icon := k.TextMuted, "calendar"
			switch {
			case reached:
				text, fg, icon = msg(c, "milestone_progress.reached"), k.Success, "circle-check"
			case late:
				fg, icon = k.Danger, "triangle-alert"
			}
			if text != "" {
				ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(fg).Role(ui.RoleStatus).Children(func() {
					ui.Icon(c, icons.Must(icon)).Size(14, 14)
					ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize))
				})
			}
		})
	})
	return col
}
