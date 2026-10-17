package media

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MicLevelMeterOptions configure a MicLevelMeter.
type MicLevelMeterOptions struct {
	// Hold keeps the peak mark this long; zero is 1.5s.
	Hold time.Duration
	// Label names the meter for assistive technology.
	Label string
}

const micSegments = 24

type micPeak struct {
	level float32
	at    time.Time
}

// MicLevelMeter shows level as segments, held peak and dB.
func MicLevelMeter(c *ui.Context, level float32, opts MicLevelMeterOptions) ui.Element {
	if level < 0 || level > 1 || math.IsNaN(float64(level)) || opts.Hold < 0 {
		panic("mujica: MicLevelMeter needs level in 0..1 and Hold >= 0")
	}
	hold := opts.Hold
	if hold == 0 {
		hold = 1500 * time.Millisecond
	}
	k := core.Tokens(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "mic_level_meter.label")
	}
	row := ui.Row(c).Gap(core.Density(c).Unit() * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize)
	pk := ui.Local(row, micPeak{}, func() micPeak { return micPeak{} })
	now := c.Now()
	switch {
	case level >= pk.level:
		*pk = micPeak{level, now}
	case now.Sub(pk.at) >= hold:
		*pk = micPeak{level, now}
	default:
		// Repaint once when the hold ends; nothing runs while idle.
		c.After(hold - now.Sub(pk.at))
	}
	peak := pk.level
	row.Role(ui.RoleMeter).Label(label).Description(micDB(level) + ", " + core.Msg(c, "mic_level_meter.peak") + " " + micDB(peak))
	row.Children(func() {
		ui.Box(c).Height(12).MinWidth(120).Grow(1).Draw(func(p *ui.Painter, r ui.Rect) {
			gap := float32(2)
			w := (r.W - gap*(micSegments-1)) / micSegments
			lit := int(math.Round(float64(level * micSegments)))
			mark := int(math.Round(float64(peak*micSegments))) - 1
			for i := range micSegments {
				x := r.X + float32(i)*(w+gap)
				seg := ui.Rect{X: x, Y: r.Y, W: w, H: r.H}
				col := micZone(k, i)
				switch {
				case i < lit:
					p.Fill(seg, col, 1)
				case i == mark:
					p.Fill(seg, col, 1)
					p.Fill(ui.Rect{X: x + 1, Y: r.Y + 2, W: w - 2, H: r.H - 4}, k.Surface, 0)
				default:
					p.Stroke(seg, k.Border, 1, 1)
				}
			}
		})
		ui.Text(c, micDB(level)).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").
			TextColor(k.TextMuted).MinWidth(56).TextAlign(ui.End).SingleLine()
	})
	return row
}

// micZone colors the safe, loud and clipping zones.
func micZone(k theme.Tokens, i int) ui.Color {
	switch {
	case i >= micSegments*9/10:
		return k.Danger
	case i >= micSegments*7/10:
		return k.Warning
	}
	return k.Success
}

// micDB writes a full-scale level in decibels.
func micDB(v float32) string {
	if v <= 0 {
		return "-∞ dB"
	}
	// Adding 0 turns a rounded -0 into 0.
	return fmt.Sprintf("%.0f dB", math.Round(20*math.Log10(float64(v)))+0)
}
