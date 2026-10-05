package media

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MediaDevice is one device; IDs are unique within a kind.
type MediaDevice struct {
	ID, Name string
}

// MediaDevices are the caller's lists; nil hides that kind.
type MediaDevices struct {
	Microphones, Cameras, Speakers []MediaDevice
}

// DeviceSelection holds each kind's chosen ID; empty is none.
type DeviceSelection struct {
	Microphone, Camera, Speaker string
}

// DeviceSelectorOptions configure a DeviceSelector.
type DeviceSelectorOptions struct {
	Disabled bool
}

// DeviceSelector chooses a microphone, camera and speaker.
func DeviceSelector(c *ui.Context, sel *DeviceSelection, devices MediaDevices, opts DeviceSelectorOptions) core.ChoiceResult {
	old := *sel
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	col := ui.Column(c).Gap(u * 3).Role(ui.RoleGroup).Label(core.Msg(c, "device_selector.label")).Disabled(opts.Disabled)
	col.Children(func() {
		deviceRow(c, k, core.MediaMic, "device_selector.microphone", &sel.Microphone, devices.Microphones, opts.Disabled)
		deviceRow(c, k, core.MediaCamera, "device_selector.camera", &sel.Camera, devices.Cameras, opts.Disabled)
		deviceRow(c, k, core.MediaSpeaker, "device_selector.speaker", &sel.Speaker, devices.Speakers, opts.Disabled)
	})
	changed := *sel != old
	if changed {
		slog.Debug("mujica: device_selector changed", "selection", *sel)
	}
	return core.NewChoiceResult(col, changed)
}

func deviceRow(c *ui.Context, k theme.Tokens, icon *ui.SVG, key string, id *string, list []MediaDevice, disabled bool) {
	if list == nil {
		return
	}
	options := make([]input.SelectOption[string], len(list))
	for i, d := range list {
		if d.ID == "" {
			panic("mujica: DeviceSelector device without an ID")
		}
		for _, o := range options[:i] {
			if o.Value == d.ID {
				panic("mujica: DeviceSelector duplicate device ID " + d.ID)
			}
		}
		options[i] = input.SelectOption[string]{Value: d.ID, Label: d.Name}
	}
	if _, ok := input.SelectLabel(options, *id); *id != "" && !ok {
		panic("mujica: DeviceSelector selection " + *id + " is not in its list")
	}
	ui.Row(c).Gap(core.Density(c).Unit() * 2).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icon).Size(16, 16).TextColor(k.TextMuted)
		ui.Text(c, core.Msg(c, key)).Width(96).SingleLine()
		hint := core.Msg(c, "device_selector.choose")
		if len(list) == 0 {
			hint = core.Msg(c, "device_selector.none")
		}
		input.Select(c, id, options, input.SelectOptions{Disabled: disabled || len(list) == 0, Label: core.Msg(c, key),
			Placeholder: hint, Empty: core.Msg(c, "device_selector.none")}).Element.Grow(1)
	})
}
