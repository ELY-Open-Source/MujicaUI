package feedback

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/overlay"
)

// InterruptButtonOptions configure an InterruptButton.
type InterruptButtonOptions struct {
	// Confirm asks before interrupting.
	Confirm bool
	// Interrupting shows the stop in progress until the run ends.
	Interrupting bool
	// Label replaces the built-in "Stop".
	Label string
}

// InterruptButton stops a running agent; Changed reports the interrupt.
func InterruptButton(c *ui.Context, running bool, opts InterruptButtonOptions) core.ChoiceResult {
	if opts.Interrupting && !running {
		panic("mujica: InterruptButton cannot be Interrupting while not running")
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "interrupt_button.stop")
	}
	shown := label
	if opts.Interrupting {
		shown = core.Msg(c, "interrupt_button.stopping")
	}
	live := running && !opts.Interrupting
	var res core.ChoiceResult
	wrap := ui.Row(c).Shrink(0)
	res.Element = wrap
	wrap.Children(func() {
		b := input.Button(c, shown, input.ButtonOptions{Variant: input.Danger, Icon: icons.Must("circle-x"), Loading: opts.Interrupting, Disabled: !running})
		open := ui.Local(wrap, "confirm", func() bool { return false })
		if !live || b.IsDisabled() {
			*open = false
		}
		clicked := b.Clicked() && live && !b.IsDisabled()
		switch {
		case clicked && opts.Confirm:
			*open = true
		case clicked:
			res = core.NewChoiceResult(res.Element, true)
		}
		r := overlay.Popconfirm(c, b, open, overlay.PopconfirmOptions{Title: core.Msg(c, "interrupt_button.confirm"),
			Description: core.Msg(c, "interrupt_button.detail"), ConfirmLabel: label, CancelLabel: core.Msg(c, "interrupt_button.keep"), Danger: true})
		if r.Confirmed && live && !b.IsDisabled() {
			res = core.NewChoiceResult(res.Element, true)
		}
	})
	if res.Changed() {
		slog.Debug("mujica: interrupt button interrupt")
	}
	return res
}
