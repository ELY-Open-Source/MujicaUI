package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// BulkActionBarOptions configure a BulkActionBar.
type BulkActionBarOptions struct {
	Disabled bool
}

// BulkActionBar shows while selected > 0 with the count, actions and a clear button; Changed reports a clear.
func BulkActionBar(c *ui.Context, selected int, opts BulkActionBarOptions, actions func()) ChoiceResult {
	if selected < 0 {
		panic(fmt.Sprintf("mujica: BulkActionBar selected %d is negative", selected))
	}
	if selected == 0 {
		return ChoiceResult{}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	count := fmt.Sprintf(msg(c, "bulk_action_bar.selected"), selected)
	bar := ui.Row(c).AlignItems(ui.Center).Gap(u*3).Padding(u*2, u*2, u*2, u*4).MinHeight(Density(c).ControlHeight() + 4*u).
		Role(ui.RoleToolbar).Label(count).Disabled(opts.Disabled)
	overlayStyle(c, bar, overlayEnter(c, bar))
	said := ui.Local(bar, "bulk_action_bar.said", func() int { return 0 })
	if *said != selected {
		*said = selected
		c.Announce(count)
	}
	var clear bool
	bar.Children(func() {
		ui.Text(c, count).Bold().TextColor(k.Text).SingleLine()
		ui.Box(c).Size(theme.BorderWidth, u*5).Background(k.Border).Role(ui.RoleNone)
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Grow(1).Children(func() {
			if actions != nil {
				actions()
			}
		})
		b := Button(c, msg(c, "bulk_action_bar.clear"), ButtonOptions{Variant: Ghost, Icon: icons.Must("x")})
		clear = b.Clicked() || bar.Shortcut(0, ui.KeyEscape) && !bar.IsDisabled()
	})
	if clear {
		slog.Debug("mujica: bulk_action_bar clear", "selected", selected)
	}
	return ChoiceResult{bar, clear}
}
