package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SendButtonOptions configure a SendButton.
type SendButtonOptions struct {
	// Streaming selects Stop while the caller generates a reply.
	Streaming bool
	// Disabled blocks sending; Stop remains available.
	Disabled bool
	// Shortcut supplies tooltip text; empty hides the hint.
	Shortcut string
}

// SendResult carries the button and requested action.
type SendResult struct {
	Element    ui.Element
	send, stop bool
}

// Sent requests sending; Stopped requests generation cancellation.
func (r SendResult) Sent() bool    { return r.send }
func (r SendResult) Stopped() bool { return r.stop }

// SendButton sends the prompt, or stops the reply while streaming.
func SendButton(c *ui.Context, opts SendButtonOptions) SendResult {
	k := core.Tokens(c)
	h := core.Density(c).ControlHeight()
	label, icon := core.Msg(c, "send_button.send"), "arrow-up"
	if opts.Streaming {
		label, icon = core.Msg(c, "send_button.stop"), "square"
	}
	tip := label
	if opts.Shortcut != "" && !opts.Streaming {
		tip += " (" + opts.Shortcut + ")"
	}
	var r SendResult
	var b ui.Element
	r.Element = ui.Row(c).Shrink(0).Disabled(opts.Disabled && !opts.Streaming)
	r.Element.Children(func() {
		b = ui.ButtonBase(c).Size(h, h).Radius(h / 2).Label(label).Tooltip(tip).Transition(core.ColorFade)
		internal.Ring(b, k.Focus, h/2)
		bg, fg := k.Accent, k.OnAccent
		switch {
		case opts.Streaming:
			bg, fg = k.Surface, k.Text
			b.Border(theme.BorderWidth, k.ControlBorder)
			if b.Hovered() {
				bg = k.SurfaceHover
			}
		case b.Pressed():
			bg = k.AccentPressed
		case b.Hovered():
			bg = k.AccentHover
		}
		b.Background(bg).TextColor(fg)
		b.Children(func() { ui.Icon(c, icons.Must(icon)).Size(16, 16) })
	})
	if b.Clicked() && !b.IsDisabled() {
		r.send, r.stop = !opts.Streaming, opts.Streaming
		slog.Debug("mujica: send_button press", "stop", r.stop)
	}
	return r
}
