package messaging

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// PresenceStatus is a person's availability. The zero value is offline.
type PresenceStatus int

const (
	PresenceOffline PresenceStatus = iota
	PresenceOnline
	PresenceAway
	PresenceBusy
)

func (p PresenceStatus) key() string {
	switch p {
	case PresenceOffline:
		return "online_status.offline"
	case PresenceOnline:
		return "online_status.online"
	case PresenceAway:
		return "online_status.away"
	case PresenceBusy:
		return "online_status.busy"
	}
	panic(fmt.Sprintf("mujica: unknown PresenceStatus %d", int(p)))
}

// OnlineStatusOptions configure an OnlineStatus.
type OnlineStatusOptions struct {
	// HideText keeps the shape, accessible name and tooltip.
	HideText bool
	// Size is the shape's diameter; zero means 10.
	Size float32
}

// OnlineStatus shows presence as a distinct shape plus its name.
func OnlineStatus(c *ui.Context, p PresenceStatus, opts OnlineStatusOptions) ui.Element {
	name := core.Msg(c, p.key())
	if opts.Size < 0 {
		panic("mujica: OnlineStatus Size must not be negative")
	}
	size := opts.Size
	if size == 0 {
		size = 10
	}
	k := core.Tokens(c)
	row := ui.Row(c).Gap(core.Density(c).Unit() * 1.5).AlignItems(ui.Center).Shrink(0).Role(ui.RoleStatus).Label(name)
	if opts.HideText {
		row.Tooltip(name)
	}
	row.Children(func() {
		ui.Box(c).Size(size, size).Shrink(0).Draw(func(p2 *ui.Painter, r ui.Rect) { presencePaint(p2, k, p, r) })
		if !opts.HideText {
			ui.Text(c, name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		}
	})
	return row
}

func presencePaint(p *ui.Painter, k theme.Tokens, s PresenceStatus, r ui.Rect) {
	cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
	rad := r.W / 2
	var disc ui.Path
	disc.Circle(cx, cy, rad)
	switch s {
	case PresenceOnline:
		p.FillPath(&disc, k.Success)
	case PresenceAway:
		p.FillPath(&disc, k.Warning)
		var bite ui.Path
		bite.Circle(cx-rad*0.45, cy-rad*0.45, rad*0.62)
		p.FillPath(&bite, k.Surface)
	case PresenceBusy:
		p.FillPath(&disc, k.Danger)
		p.Fill(ui.Rect{X: cx - rad*0.6, Y: cy - max(1, rad*0.18), W: rad * 1.2, H: max(2, rad*0.36)}, k.Surface, 0)
	default:
		p.StrokePath(&disc, max(1.5, rad*0.3), k.TextMuted)
	}
}
