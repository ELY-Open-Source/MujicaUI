package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// AudioMessageOptions describe a voice message the caller plays.
type AudioMessageOptions struct {
	// Duration is the length; required.
	Duration time.Duration
	// Position is how far playback has gone, 0 to Duration.
	Position time.Duration
	Playing  bool
	// Waveform holds loudness samples from 0 to 1; empty draws a flat line.
	Waveform []float32
}

// AudioMessageResult is an AudioMessage's element and the user's requests this frame.
type AudioMessageResult struct {
	Element *ui.Element
	toggled bool
	seek    time.Duration
	seeked  bool
}

// Toggled reports a press of play or pause.
func (r AudioMessageResult) Toggled() bool { return r.toggled }

// Seek returns the position the user moved to by click, drag or keys.
func (r AudioMessageResult) Seek() (time.Duration, bool) { return r.seek, r.seeked }

// audioBarWidth and audioBarGap space the waveform's bars.
const audioBarWidth, audioBarGap float32 = 2, 2

// AudioMessage shows a voice bar: play button, waveform with progress, and times.
func AudioMessage(c *ui.Context, opts AudioMessageOptions) AudioMessageResult {
	if opts.Duration <= 0 || opts.Position < 0 || opts.Position > opts.Duration {
		panic("mujica: AudioMessage needs 0 <= Position <= Duration and Duration > 0")
	}
	for _, v := range opts.Waveform {
		if !(v >= 0 && v <= 1) {
			panic("mujica: AudioMessage Waveform samples must be from 0 to 1")
		}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r AudioMessageResult
	r.Element = ui.Row(c).Gap(u*2).AlignItems(ui.Center).Padding(u, u*3, u, u).Width(300).
		Radius(theme.CardRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border)
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		icon, name := "play", msg(c, "audio_message.play")
		if opts.Playing {
			icon, name = "pause", msg(c, "audio_message.pause")
		}
		if iconAction(c, icons.Must(icon), name).Clicked() && !disabled {
			r.toggled = true
			slog.Debug("mujica: audio_message toggle", "playing", !opts.Playing)
		}
		pos := opts.Position.Seconds()
		s := ui.SliderBase(c, &pos, 0, opts.Duration.Seconds()).Grow(1).Height(theme.MinHitSize).
			Label(msg(c, "audio_message.seek")).Description(audioClock(opts.Position) + " / " + audioClock(opts.Duration))
		internal.Ring(s, k.Focus, theme.ControlRadius)
		if s.Changed() && !disabled {
			r.seek, r.seeked = time.Duration(pos*float64(time.Second)), true
			slog.Debug("mujica: audio_message seek", "at", r.seek)
		}
		frac := float32(opts.Position.Seconds() / opts.Duration.Seconds())
		s.Draw(func(p *ui.Painter, b ui.Rect) { audioWave(p, b, opts.Waveform, frac, k) })
		ui.Text(c, audioClock(opts.Position)+" / "+audioClock(opts.Duration)).FontFeatures("tnum").
			FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine().Shrink(0)
	})
	return r
}

// audioWave draws bars, played ones in AccentText, with a position line so progress is not color alone.
func audioWave(p *ui.Painter, b ui.Rect, wave []float32, frac float32, k theme.Tokens) {
	n := int((b.W + audioBarGap) / (audioBarWidth + audioBarGap))
	if n <= 0 {
		return
	}
	cy := b.Y + b.H/2
	for i := range n {
		v := float32(0.12)
		if len(wave) > 0 {
			v = max(v, wave[i*len(wave)/n])
		}
		h := max(2, v*(b.H-8))
		x := b.X + float32(i)*(audioBarWidth+audioBarGap)
		col := k.TextMuted.Alpha(0.55)
		if (x-b.X)/b.W < frac {
			col = k.AccentText
		}
		p.Fill(ui.Rect{X: x, Y: cy - h/2, W: audioBarWidth, H: h}, col, 1)
	}
	at := internal.Snap(p, b.X+b.W*frac)
	p.Fill(ui.Rect{X: at - 1, Y: b.Y + 2, W: 2, H: b.H - 4}, k.Text, 1)
}

// audioClock formats d as m:ss.
func audioClock(d time.Duration) string {
	s := int(d.Round(time.Second).Seconds())
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
