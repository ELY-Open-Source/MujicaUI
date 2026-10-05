package media

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// VolumeControlOptions configure a VolumeControl.
type VolumeControlOptions struct {
	Disabled bool
}

// VolumeControl sets *volume in 5% steps and toggles *muted.
func VolumeControl(c *ui.Context, volume *float64, muted *bool, opts VolumeControlOptions) core.ChoiceResult {
	if *volume < 0 || *volume > 1 || math.IsNaN(*volume) {
		panic("mujica: VolumeControl volume outside 0..1")
	}
	oldV, oldM := *volume, *muted
	u := core.Density(c).Unit()
	row := ui.Row(c).Gap(u).AlignItems(ui.Center).Role(ui.RoleGroup).Label(core.Msg(c, "volume_control.label")).Disabled(opts.Disabled)
	row.Children(func() {
		icon, key := core.MediaVolume, "volume_control.mute"
		if *muted || *volume == 0 {
			icon, key = core.MediaMute, "volume_control.unmute"
		}
		b := MediaToggleAction(c, icon, core.Msg(c, key), *muted)
		if b.Clicked() && !b.IsDisabled() {
			*muted = !*muted
		}
		s := input.Slider(c, volume, input.SliderOptions{Min: 0, Max: 1, Step: 0.05, Disabled: opts.Disabled,
			Label: core.Msg(c, "volume_control.label"), ShowValue: true, Format: func(v float64) string { return fmt.Sprintf("%d%%", int(math.Round(v*100))) }})
		s.Element.Width(120)
		// Raising the volume while muted unmutes, as players do.
		if s.Changed() && *muted && *volume > oldV {
			*muted = false
		}
	})
	if !row.IsDisabled() && row.Shortcut(0, ui.KeyM) {
		*muted = !*muted
	}
	changed := *volume != oldV || *muted != oldM
	if changed {
		slog.Debug("mujica: volume changed", "volume", *volume, "muted", *muted)
	}
	return core.NewChoiceResult(row, changed)
}
