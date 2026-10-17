package data

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TimelineStatus marks an event's state with a shape and a word, not only color.
type TimelineStatus int

const (
	TimelineDefault TimelineStatus = iota
	TimelineDone
	TimelineActive
	TimelineWarning
	TimelineError
)

// TimelineItem is one event of a Timeline.
type TimelineItem struct {
	Time    string
	Title   string
	Content string
	Status  TimelineStatus
}

// Timeline lists events top to bottom with markers joined by a line; long text wraps.
func Timeline(c *ui.Context, items []TimelineItem) ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	return ui.Column(c).Role(ui.RoleList).Children(func() {
		for i, it := range items {
			last := i == len(items)-1
			ui.Row(c).Gap(u * 3).AlignItems(ui.Stretch).Role(ui.RoleListItem).Children(func() {
				ui.Column(c).Width(20).Shrink(0).AlignItems(ui.Center).Children(func() {
					timelineMarker(c, it.Status)
					if !last {
						ui.Box(c).Width(1).Grow(1).MinHeight(u * 3).Background(k.Border)
					}
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(u).Padding(0, 0, u*5, 0).Children(func() {
					ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Wrap().Children(func() {
						if it.Time != "" {
							ui.Text(c, it.Time).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
						}
						if w := timelineWord(c, it.Status); w != "" {
							ui.Text(c, w).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(timelineColor(k, it.Status)).FontWeight(600)
						}
					})
					if it.Title != "" {
						ui.Text(c, it.Title).FontWeight(600)
					}
					if it.Content != "" {
						ui.Text(c, it.Content).TextColor(k.TextMuted)
					}
				})
			})
		}
	})
}

func timelineColor(k theme.Tokens, s TimelineStatus) ui.Color {
	switch s {
	case TimelineDone:
		return k.Success
	case TimelineActive:
		return k.AccentText
	case TimelineWarning:
		return k.Warning
	case TimelineError:
		return k.Danger
	}
	return k.TextMuted
}

func timelineWord(c *ui.Context, s TimelineStatus) string {
	switch s {
	case TimelineDone:
		return core.Msg(c, "timeline.done")
	case TimelineActive:
		return core.Msg(c, "timeline.active")
	case TimelineWarning:
		return core.Msg(c, "timeline.warning")
	case TimelineError:
		return core.Msg(c, "timeline.error")
	}
	return ""
}

// timelineMarker draws a hollow diamond, or a status icon.
func timelineMarker(c *ui.Context, s TimelineStatus) {
	k := core.Tokens(c)
	col := timelineColor(k, s)
	box := ui.Box(c).Size(20, 20).Shrink(0).Center().TextColor(col)
	switch s {
	case TimelineDone:
		box.Children(func() { ui.Icon(c, icons.Must("circle-check")).Size(18, 18) })
	case TimelineWarning:
		box.Children(func() { ui.Icon(c, icons.Must("triangle-alert")).Size(18, 18) })
	case TimelineError:
		box.Children(func() { ui.Icon(c, icons.Must("circle-x")).Size(18, 18) })
	default:
		box.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
			if s == TimelineActive {
				p.FillPath(internal.Diamond(cx, cy, 12), col)
			} else {
				p.StrokePath(internal.Diamond(cx, cy, 11), 1.5, col)
			}
		})
	}
}
