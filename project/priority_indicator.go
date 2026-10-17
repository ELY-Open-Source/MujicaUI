package project

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PriorityIndicatorOptions configure a PriorityIndicator.
type PriorityIndicatorOptions struct {
	// IconOnly hides the text; the name stays in the label.
	IconOnly bool
}

var priorityKeys = [...]string{"priority_indicator.none", "priority_indicator.low", "priority_indicator.medium",
	"priority_indicator.high", "priority_indicator.urgent"}

// PriorityName is a priority's localized name.
func PriorityName(c *ui.Context, p Priority) string {
	projectCheckPriority(p)
	return core.Msg(c, priorityKeys[p])
}

// PriorityIndicator displays urgency with shapes and localized text.
func PriorityIndicator(c *ui.Context, p Priority, opts PriorityIndicatorOptions) ui.Element {
	name := PriorityName(c, p)
	k := core.Tokens(c)
	d := core.Density(c)
	col := k.TextMuted
	switch p {
	case PriorityUrgent:
		col = k.Danger
	case PriorityHigh:
		col = k.Warning
	}
	row := ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Shrink(0).Role(ui.RoleImage).Label(name)
	row.Children(func() {
		ui.Box(c).Size(14, 14).Shrink(0).Draw(func(pp *ui.Painter, r ui.Rect) { priorityGlyph(pp, r, p, col, k.Surface, k.Border) })
		if !opts.IconOnly {
			ui.Text(c, name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(col).SingleLine()
		}
	})
	return row
}

// priorityGlyph draws level bars, urgency marks or empty dashes.
func priorityGlyph(p *ui.Painter, r ui.Rect, pr Priority, col, bg, faint ui.Color) {
	switch pr {
	case PriorityUrgent:
		p.Fill(r, col, 2)
		p.Fill(ui.Rect{X: r.X + r.W/2 - 1, Y: r.Y + 3, W: 2, H: r.H - 8}, bg, 0)
		p.Fill(ui.Rect{X: r.X + r.W/2 - 1, Y: r.Y + r.H - 4, W: 2, H: 2}, bg, 0)
	case PriorityNone:
		for i := range 3 {
			p.Fill(ui.Rect{X: r.X + 1 + float32(i)*4.5, Y: r.Y + r.H/2 - 1, W: 3, H: 2}, col, 0)
		}
	default:
		for i := range 3 {
			h := r.H * float32(i+1) / 3
			c := faint
			if i < int(pr) {
				c = col
			}
			p.Fill(ui.Rect{X: r.X + 1 + float32(i)*4.5, Y: r.Y + r.H - h, W: 3, H: h}, c, 1)
		}
	}
}
