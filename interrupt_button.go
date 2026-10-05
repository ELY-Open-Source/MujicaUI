package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
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

// InterruptButton stops a running agent; Changed reports the interrupt, confirmed when asked.
func InterruptButton(c *ui.Context, running bool, opts InterruptButtonOptions) ChoiceResult {
	if opts.Interrupting && !running {
		panic("mujica: InterruptButton cannot be Interrupting while not running")
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "interrupt_button.stop")
	}
	shown := label
	if opts.Interrupting {
		shown = msg(c, "interrupt_button.stopping")
	}
	live := running && !opts.Interrupting
	var res ChoiceResult
	wrap := ui.Row(c).Shrink(0)
	res.Element = wrap
	wrap.Children(func() {
		b := Button(c, shown, ButtonOptions{Variant: Danger, Icon: icons.Must("circle-x"), Loading: opts.Interrupting, Disabled: !running})
		open := ui.Local(wrap, "confirm", func() bool { return false })
		if !live || b.IsDisabled() {
			*open = false
		}
		clicked := b.Clicked() && live
		switch {
		case clicked && opts.Confirm:
			*open = true
		case clicked:
			res.changed = true
		}
		r := Popconfirm(c, b, open, PopconfirmOptions{Title: msg(c, "interrupt_button.confirm"),
			Description: msg(c, "interrupt_button.detail"), ConfirmLabel: label, CancelLabel: msg(c, "interrupt_button.keep"), Danger: true})
		if r.Confirmed && live {
			res.changed = true
		}
	})
	if res.changed {
		slog.Debug("mujica: interrupt button interrupt")
	}
	return res
}
