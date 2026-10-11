package chat

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
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

// VoiceInputButton requests recording changes and displays elapsed time.
// Clicked requests recording start or stop for the current state.
func VoiceInputButton(c *ui.Context, opts VoiceInputOptions) core.ChoiceResult {
	if opts.Recording && (opts.Started.IsZero() || opts.Denied) {
		panic("mujica: VoiceInputButton recording needs Started and no Denied")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Shrink(0).Disabled(opts.Disabled)
	r.Element.Children(func() {
		elapsed := ""
		label := core.Msg(c, "voice_input_button.start")
		icon := "mic"
		if opts.Recording {
			d := max(c.Now().Sub(opts.Started), 0)
			elapsed = audioClock(d.Truncate(time.Second))
			label, icon = fmt.Sprintf(core.Msg(c, "voice_input_button.stop"), elapsed), "square"
		}
		h := core.Density(c).ControlHeight()
		b := ui.ButtonBase(c).Size(h, h).Radius(h / 2).Label(label).Tooltip(label).Transition(core.ColorFade)
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
			r = core.NewChoiceResult(r.Element, true)
			slog.Debug("mujica: voice_input_button press", "recording", opts.Recording)
		}
		switch {
		case opts.Recording:
			ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(8, 8).Radius(4).Background(k.Danger)
				ui.Text(c, elapsed).FontFeatures("tnum").TextColor(k.Danger).SingleLine().DrawOver(func(_ *ui.Painter, _ ui.Rect) {
					d := max(c.Now().Sub(opts.Started), 0)
					c.After(time.Second - d%time.Second)
				})
			})
		case opts.Denied:
			ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(k.Danger).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("mic-off")).Size(14, 14)
				ui.Text(c, core.Msg(c, "voice_input_button.denied")).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
			})
		}
	})
	return r
}
