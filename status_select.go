package mujica

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// StatusSelectOptions configure a StatusSelect.
type StatusSelectOptions struct {
	Disabled bool
	Label    string
}

var statusGroupKeys = [...]string{"status_select.backlog", "status_select.unstarted", "status_select.started",
	"status_select.completed", "status_select.canceled"}

// StatusSelect chooses a status grouped by workflow stage.
func StatusSelect(c *ui.Context, value *string, statuses []WorkflowStatus, opts StatusSelectOptions) ChoiceResult {
	if len(statuses) == 0 {
		panic("mujica: StatusSelect needs statuses")
	}
	for i, s := range statuses {
		if s.Group < StatusBacklog || s.Group > StatusCanceled || s.ID == "" {
			panic(fmt.Sprintf("mujica: StatusSelect status %+v is not valid", s))
		}
		for _, o := range statuses[:i] {
			if o.ID == s.ID {
				panic(fmt.Sprintf("mujica: StatusSelect status %q repeats", s.ID))
			}
		}
	}
	cur := projectStatus(statuses, *value)
	k := Tokens(c)
	old := *value
	pick := old
	var s *ui.SelectParts[string]
	root := ui.Row(c).Shrink(0).Disabled(opts.Disabled)
	epoch := ui.Local(root, "status-epoch", func() struct {
		disabled   bool
		generation int
	} {
		return struct {
			disabled   bool
			generation int
		}{}
	})
	if disabled := root.IsDisabled(); disabled != epoch.disabled {
		epoch.disabled = disabled
		epoch.generation++
	}
	root.Children(func() {
		ui.Row(c).Key(epoch.generation).Children(func() { s = ui.SelectBase(c, &pick) })
	})
	tr := s.Trigger
	label := opts.Label
	if label == "" {
		label = msg(c, "status_select.label")
	}
	tr.Label(label).Description(cur.Name)
	selectField(c, tr, s.Open(), func() {
		statusGlyphBox(c, cur.Group)
		ui.Text(c, cur.Name).SingleLine().Grow(1).Shrink(1).TextColor(k.Text)
	})
	if !tr.IsDisabled() {
		s.Popup(func(panel *ui.Element) {
			selectPanel(c, panel)
			_, wh := c.Size()
			ui.Scroll(c).MaxHeight(min(320, wh/2)).Children(func() {
				for g := StatusBacklog; g <= StatusCanceled; g++ {
					first := true
					for _, st := range statuses {
						if st.Group != g {
							continue
						}
						if first {
							ui.Text(c, msg(c, statusGroupKeys[g])).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).
								Padding(Density(c).Unit(), Density(c).Unit()*2).Role(ui.RoleColumnHeader)
							first = false
						}
						it := s.Item(st.ID)
						selectRow(c, it, it.Highlighted(), st.ID == *value, "")
						it.Label(st.Name).Children(func() {
							statusGlyphBox(c, st.Group)
							ui.Text(c, st.Name).SingleLine()
						})
					}
				}
			})
		})
	}
	if pick != *value && !tr.IsDisabled() {
		*value = pick
	}
	changed := *value != old
	if changed {
		slog.Debug("mujica: status_select change", "from", old, "to", *value)
	}
	return ChoiceResult{Element: tr, changed: changed}
}

// statusGlyphBox reserves 14 DIPs for a status glyph.
func statusGlyphBox(c *ui.Context, g StatusGroup) {
	k := Tokens(c)
	ui.Box(c).Size(14, 14).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) { statusGlyph(p, r, g, k) })
}

// statusGlyph draws a distinct shape for each workflow stage.
func statusGlyph(p *ui.Painter, r ui.Rect, g StatusGroup, k theme.Tokens) {
	cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-1
	ring := func(col ui.Color) {
		var path ui.Path
		internal.Arc(&path, cx, cy, rad, 0, 2*math.Pi, true)
		p.StrokePath(&path, 1.5, col)
	}
	switch g {
	case StatusBacklog:
		for i := 0; i < 8; i += 2 {
			var path ui.Path
			internal.Arc(&path, cx, cy, rad, float64(i)*math.Pi/4, float64(i+1)*math.Pi/4, true)
			p.StrokePath(&path, 1.5, k.TextMuted)
		}
	case StatusUnstarted:
		ring(k.TextMuted)
	case StatusStarted:
		ring(k.Warning)
		p.FillPath(internal.Sector(cx, cy, 0, rad-2, 0, math.Pi), k.Warning)
	case StatusCompleted:
		p.FillPath(internal.Sector(cx, cy, 0, rad+1, 0, 2*math.Pi), k.Success)
		internal.Mark(p, ui.Rect{X: r.X + 3, Y: r.Y + 3, W: r.W - 6, H: r.H - 6}, "check", k.Surface)
	case StatusCanceled:
		p.FillPath(internal.Sector(cx, cy, 0, rad+1, 0, 2*math.Pi), k.TextMuted)
		internal.Mark(p, ui.Rect{X: r.X + 3, Y: r.Y + 3, W: r.W - 6, H: r.H - 6}, "x", k.Surface)
	}
}
