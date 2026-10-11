package chat

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// VoiceWaveformOptions configure a VoiceWaveform.
type VoiceWaveformOptions struct {
	// Silence sets the quiet threshold; zero means 0.05.
	Silence float32
	// SilentFor sets consecutive quiet samples; zero means 20.
	SilentFor int
	// Width sets waveform width; zero means 200 DIPs.
	Width float32
}

// VoiceWaveform draws recent input levels and a silence hint.
func VoiceWaveform(c *ui.Context, levels []float32, opts VoiceWaveformOptions) *ui.Element {
	if opts.Silence < 0 || opts.Silence >= 1 || opts.SilentFor < 0 || opts.Width < 0 {
		panic("mujica: VoiceWaveform needs 0 <= Silence < 1 and non-negative SilentFor and Width")
	}
	for _, v := range levels {
		if !(v >= 0 && v <= 1) {
			panic("mujica: VoiceWaveform levels must be from 0 to 1")
		}
	}
	silence, need, w := opts.Silence, opts.SilentFor, opts.Width
	if silence == 0 {
		silence = 0.05
	}
	if need == 0 {
		need = 20
	}
	if w == 0 {
		w = 200
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	silent := len(levels) >= need
	for _, v := range levels[max(0, len(levels)-need):] {
		silent = silent && v < silence
	}
	row := ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Shrink(0).Role(ui.RoleMeter).Label(core.Msg(c, "voice_waveform.label"))
	row.Children(func() {
		ui.Box(c).Size(w, 28).Shrink(0).Draw(func(p *ui.Painter, b ui.Rect) {
			n := int((b.W + audioBarGap) / (audioBarWidth + audioBarGap))
			cy := b.Y + b.H/2
			for i := range n {
				// The newest sample sits at the right edge.
				j := len(levels) - n + i
				v := float32(0)
				if j >= 0 {
					v = levels[j]
				}
				h := max(2, v*(b.H-4))
				col := k.AccentText
				if v < silence {
					col = k.TextMuted.Alpha(0.55)
				}
				x := b.X + float32(i)*(audioBarWidth+audioBarGap)
				p.Fill(ui.Rect{X: x, Y: cy - h/2, W: audioBarWidth, H: h}, col, 1)
			}
		})
		if silent {
			ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(k.Warning).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("mic-off")).Size(14, 14)
				ui.Text(c, core.Msg(c, "voice_waveform.silent")).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
			})
		}
	})
	return row
}
