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

// VoiceInputOptions describe the caller's recording.
type VoiceInputOptions struct {
	// Recording shows the stop state and the time since Started.
	Recording bool
	Started   time.Time
	// Denied shows that microphone access failed; the button then retries.
	Denied   bool
	Disabled bool
}

// VoiceInputButton asks the caller to start or stop recording; while recording it shows the time.
// Clicked reports a press: start when idle or denied, stop when recording.
func VoiceInputButton(c *ui.Context, opts VoiceInputOptions) ChoiceResult {
	if opts.Recording && (opts.Started.IsZero() || opts.Denied) {
		panic("mujica: VoiceInputButton recording needs Started and no Denied")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Shrink(0).Disabled(opts.Disabled)
	r.Element.Children(func() {
		elapsed := ""
		label := msg(c, "voice_input_button.start")
		icon := "mic"
		if opts.Recording {
			d := max(c.Now().Sub(opts.Started), 0)
			elapsed = audioClock(d.Truncate(time.Second))
			label, icon = fmt.Sprintf(msg(c, "voice_input_button.stop"), elapsed), "square"
			c.After(time.Second - d%time.Second)
		}
		h := Density(c).ControlHeight()
		b := ui.ButtonBase(c).Size(h, h).Radius(h / 2).Label(label).Tooltip(label).Transition(colorFade)
		internal.Ring(b, k.Focus, h/2)
		switch {
		case opts.Recording:
			b.Background(k.DangerFill).TextColor(k.OnDanger)
		case b.Pressed():
			b.Background(k.Selection).TextColor(k.Text)
		case b.Hovered():
			b.Background(k.SurfaceHover).TextColor(k.Text)
		default:
			b.TextColor(k.TextMuted)
		}
		b.Children(func() { ui.Icon(c, icons.Must(icon)).Size(16, 16) })
		if b.Clicked() && !b.IsDisabled() {
			r.changed = true
			slog.Debug("mujica: voice_input_button press", "recording", opts.Recording)
		}
		switch {
		case opts.Recording:
			ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(8, 8).Radius(4).Background(k.Danger)
				ui.Text(c, elapsed).FontFeatures("tnum").TextColor(k.Danger).SingleLine()
			})
		case opts.Denied:
			ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(k.Danger).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("mic-off")).Size(14, 14)
				ui.Text(c, msg(c, "voice_input_button.denied")).FontSize(FontSize(c, theme.CaptionSize)).SingleLine()
			})
		}
	})
	return r
}
