package core

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Severity is the meaning of a message. The zero value is SeverityInfo.
type Severity int

const (
	SeverityInfo Severity = iota
	SeveritySuccess
	SeverityWarning
	SeverityError
)

func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeveritySuccess:
		return "success"
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	}
	panic(fmt.Sprintf("mujica: unknown Severity %d", int(s)))
}

// SeverityLook returns the color and icon that, with text, mark s.
func SeverityLook(k theme.Tokens, s Severity) (ui.Color, *ui.SVG) {
	switch s {
	case SeveritySuccess:
		return k.Success, icons.Must("circle-check")
	case SeverityWarning:
		return k.Warning, icons.Must("triangle-alert")
	case SeverityError:
		return k.Danger, icons.Must("circle-x")
	case SeverityInfo:
		return k.Info, icons.Must("info")
	}
	panic(fmt.Sprintf("mujica: unknown Severity %d", int(s)))
}

// AlertResult is what an Alert built and reported this frame.
type AlertResult struct {
	// Element is nil while the alert is closed.
	Element *ui.Element
	// Action reports a click on the action button.
	Action bool
	// Closed reports the user dismissing the alert.
	Closed bool
}
