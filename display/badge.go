package display

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// BadgeTone is a badge's emphasis. The zero value is BadgeAccent.
type BadgeTone int

const (
	BadgeAccent BadgeTone = iota
	BadgeNeutral
	BadgeSuccess
	BadgeWarning
	BadgeDanger
	BadgeInfo
)

// BadgePosition places a badge; zero is inline.
type BadgePosition int

const (
	BadgeInline BadgePosition = iota
	BadgeTopRight
	BadgeTopLeft
	BadgeBottomRight
	BadgeBottomLeft
)

// BadgeOptions configure a Badge.
type BadgeOptions struct {
	// Text shows as is; when empty, Count shows instead.
	Text string
	// Count must not be negative; zero without Text hides it.
	Count int
	// Max caps the shown count as "Max+"; zero means 99.
	Max      int
	Tone     BadgeTone
	Position BadgePosition
	// Label says what is counted; read before the full value.
	Label string
}

// Badge shows text or a count, inline or cornered.
func Badge(c *ui.Context, opts BadgeOptions, content func()) ui.Element {
	if opts.Count < 0 || opts.Max < 0 {
		panic("mujica: Badge Count and Max must not be negative")
	}
	if (opts.Position == BadgeInline) != (content == nil) {
		panic("mujica: Badge needs content exactly when Position is not BadgeInline")
	}
	shown, full := badgeText(opts)
	if opts.Position == BadgeInline {
		if shown == "" {
			return ui.Box(c).Size(0, 0)
		}
		return badgePill(c, opts, shown, full)
	}
	at := badgeAnchor(opts.Position)
	wrap := ui.Box(c).Shrink(0).AlignSelf(ui.Start).AlignItems(ui.Start)
	wrap.Children(func() {
		content()
		if shown != "" {
			badgePill(c, opts, shown, full).Attach(at, ui.AnchorCenter)
		}
	})
	return wrap
}

// badgeText returns what the badge shows and what it reads.
func badgeText(opts BadgeOptions) (shown, full string) {
	if opts.Text != "" {
		return opts.Text, opts.Text
	}
	if opts.Count == 0 {
		return "", ""
	}
	limit := opts.Max
	if limit == 0 {
		limit = 99
	}
	full = strconv.Itoa(opts.Count)
	if opts.Count > limit {
		return strconv.Itoa(limit) + "+", full
	}
	return full, full
}

func badgePill(c *ui.Context, opts BadgeOptions, shown, full string) ui.Element {
	k := core.Tokens(c)
	bg, fg, border := badgeColors(k, opts.Tone)
	label := full
	if opts.Label != "" {
		label = opts.Label + " " + full
	}
	p := ui.Box(c).Row().Center().MinWidth(core.FontSize(c, 18)).MinHeight(core.FontSize(c, 18)).Padding(0, 5).Radius(theme.CardRadius).
		Background(bg).Shrink(0).Role(ui.RoleStatus).Label(label)
	if border.A > 0 {
		p.Border(theme.BorderWidth, border)
	}
	p.Children(func() {
		ui.Text(c, shown).FontSize(core.FontSize(c, theme.CaptionSize-1)).FontWeight(600).TextColor(fg).
			FontFeatures("tnum").SingleLine()
	})
	return p
}

func badgeColors(k theme.Tokens, t BadgeTone) (bg, fg, border ui.Color) {
	switch t {
	case BadgeAccent:
		return k.Accent, k.OnAccent, ui.Color{}
	case BadgeNeutral:
		return k.SurfaceHover, k.Text, k.Border
	case BadgeSuccess:
		return k.Surface, k.Success, k.Success
	case BadgeWarning:
		return k.Surface, k.Warning, k.Warning
	case BadgeDanger:
		return k.DangerFill, k.OnDanger, ui.Color{}
	case BadgeInfo:
		return k.Surface, k.Info, k.Info
	}
	panic(fmt.Sprintf("mujica: unknown BadgeTone %d", t))
}

func badgeAnchor(p BadgePosition) ui.Anchor {
	switch p {
	case BadgeTopRight:
		return ui.AnchorTopRight
	case BadgeTopLeft:
		return ui.AnchorTopLeft
	case BadgeBottomRight:
		return ui.AnchorBottomRight
	case BadgeBottomLeft:
		return ui.AnchorBottomLeft
	}
	panic(fmt.Sprintf("mujica: unknown BadgePosition %d", p))
}
