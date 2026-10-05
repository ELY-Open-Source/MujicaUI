package input

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// copyButtonShown is how long the button says Copied.
const copyButtonShown = 2 * time.Second

// CopyButtonOptions configure a CopyButton.
type CopyButtonOptions struct {
	// Label shows beside the icon; empty makes a square icon button.
	Label    string
	Disabled bool
}

// CopyButton puts text on the clipboard, then shows and announces Copied for two seconds.
func CopyButton(c *ui.Context, text string, opts CopyButtonOptions) core.ChoiceResult {
	copied := core.Msg(c, "copy_button.copied")
	name := core.Msg(c, "copy_button.copy")
	label, icon := opts.Label, icons.Must("copy")
	var b *ui.Element
	box := ui.Box(c)
	at := ui.Local(box, "copy_button.at", func() time.Time { return time.Time{} })
	left := copyButtonShown - c.Now().Sub(*at)
	done := !at.IsZero() && left > 0
	if done {
		icon, name = icons.Must("check"), copied
		if label != "" {
			label = copied
		}
		c.After(left)
	}
	box.Children(func() {
		bo := ButtonOptions{Variant: Secondary, Icon: icon, Disabled: opts.Disabled}
		if opts.Label == "" {
			bo.Label = name
		}
		b = Button(c, label, bo)
		if opts.Label == "" {
			b.Padding(0).Size(core.Density(c).ControlHeight(), core.Density(c).ControlHeight()).Tooltip(name)
		} else {
			// Room for the longer word, so the width holds.
			w1, _ := c.MeasureText(0, ui.Span{Text: copied})
			w2, _ := c.MeasureText(0, ui.Span{Text: opts.Label})
			b.MinWidth(max(w1, w2) + buttonSpinnerSize + core.Density(c).Unit()*10 + 2*theme.BorderWidth)
		}
	})
	if b.Clicked() && !b.IsDisabled() {
		c.WriteClipboard(text)
		*at = c.Now()
		c.Announce(copied)
		c.Invalidate()
		slog.Debug("mujica: copy_button copied", "bytes", len(text))
		return core.NewChoiceResult(box, true)
	}
	return core.NewChoiceResult(box, false)
}
