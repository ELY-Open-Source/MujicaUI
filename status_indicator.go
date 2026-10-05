package mujica

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// Status is a system or item state. The zero value is StatusSuccess.
type Status int

const (
	StatusSuccess Status = iota
	StatusWarning
	StatusError
	StatusOffline
	StatusBusy
)

func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusWarning:
		return "warning"
	case StatusError:
		return "error"
	case StatusOffline:
		return "offline"
	case StatusBusy:
		return "busy"
	}
	panic(fmt.Sprintf("mujica: unknown Status %d", int(s)))
}

// StatusIndicatorOptions configure a StatusIndicator.
type StatusIndicatorOptions struct {
	// Label replaces the status's built-in word, as "Database online".
	Label string
}

// StatusIndicator shows a status by its own shape and a word, never by color alone.
func StatusIndicator(c *ui.Context, status Status, opts StatusIndicatorOptions) *ui.Element {
	key := "status_indicator." + status.String() // panics on an unknown Status
	k := Tokens(c)
	u := Density(c).Unit()
	label := opts.Label
	if label == "" {
		label = msg(c, key)
	}
	var tone ui.Color
	var icon *ui.SVG
	switch status {
	case StatusSuccess:
		tone, icon = k.Success, icons.Must("circle-check")
	case StatusWarning:
		tone, icon = k.Warning, icons.Must("triangle-alert")
	case StatusError:
		tone, icon = k.Danger, icons.Must("circle-x")
	case StatusOffline:
		tone, icon = k.TextMuted, icons.Must("wifi-off")
	case StatusBusy:
		tone = k.Info
	}
	row := ui.Row(c).AlignItems(ui.Center).Gap(u * 1.5).Role(ui.RoleStatus).Label(label)
	row.Children(func() {
		if status == StatusBusy {
			Spinner(c, SpinnerOptions{Size: 16, Color: tone, Label: label})
		} else {
			ui.Icon(c, icon).Size(16, 16).TextColor(tone)
		}
		ui.Text(c, label).TextColor(k.Text).FontSize(FontSize(c, theme.BodySize)).SingleLine()
	})
	return row
}
