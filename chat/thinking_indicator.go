package chat

import (
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ThinkingIndicatorOptions configure a ThinkingIndicator.
type ThinkingIndicatorOptions struct {
	// Phase is the caller's stage text; empty uses "Thinking".
	Phase string
}

// ThinkingIndicator shows rising diamonds and a phase, still when reduced.
func ThinkingIndicator(c *ui.Context, opts ThinkingIndicatorOptions) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	phase := opts.Phase
	if phase == "" {
		phase = core.Msg(c, "thinking_indicator.thinking")
	}
	still := c.Preferences().ReduceMotion
	row := ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).MinHeight(d.ControlHeight()).Role(ui.RoleStatus).Label(phase)
	row.Children(func() {
		ui.Box(c).Size(30, 14).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			t := 0.0
			if !still {
				t = float64(p.Now().UnixMilli()%1200) / 1200
				p.AnimationFrame()
			}
			for i := range 3 {
				lift := float32(0)
				alpha := float32(1)
				if !still {
					phase := math.Mod(t-float64(i)*0.18+1, 1)
					wave := float32(math.Max(0, math.Sin(phase*2*math.Pi)))
					lift, alpha = wave*3, 0.45+0.55*wave
				}
				x := internal.Snap(p, r.X+5+float32(i)*10)
				y := internal.Snap(p, r.Y+r.H/2-lift)
				p.FillPath(internal.Diamond(x, y, 7), k.AccentText.Alpha(alpha))
			}
		})
		text := phase
		if still {
			text += "…"
		}
		ui.Text(c, text).FontSize(core.FontSize(c, theme.BodySize)).TextColor(k.TextMuted).Italic()
	})
	return row
}
