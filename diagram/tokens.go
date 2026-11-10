package diagram

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// tokens maps Pele's tokens to Mujica's.
func tokens(k theme.Tokens) palette {
	s := internal.Palette(k)
	return palette{bg: k.Background, fg: k.Text, muted: k.TextMuted, line: k.Text, surface: k.Surface,
		alt: k.SurfaceHover, border: k.ControlBorder, accent: k.AccentText, critical: k.Danger, radius: theme.CardRadius,
		series: [8]ui.Color{s[0], s[1], s[2], s[3], s[4], s[5], k.Danger, s[0].Mix(s[1], 0.5)}}
}
