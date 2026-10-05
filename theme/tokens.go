// Package theme holds MujicaUI's design tokens: colors, type, metrics and motion.
package theme

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Tokens are the colors of one appearance.
type Tokens struct {
	Dark          bool
	Background    ui.Color
	Surface       ui.Color
	SurfaceHover  ui.Color
	Text          ui.Color
	TextMuted     ui.Color
	Border        ui.Color
	ControlBorder ui.Color
	Accent        ui.Color
	AccentHover   ui.Color
	AccentPressed ui.Color
	OnAccent      ui.Color
	AccentText    ui.Color
	Selection     ui.Color
	Focus         ui.Color
	Ornament      ui.Color
	Success       ui.Color
	Warning       ui.Color
	Danger        ui.Color
	Info          ui.Color
	DangerFill    ui.Color
	OnDanger      ui.Color
}

// Light returns the ivory appearance.
func Light() Tokens {
	return Tokens{
		Background:    ui.Hex("#F2EBDD"),
		Surface:       ui.Hex("#FBF7EF"),
		SurfaceHover:  ui.Hex("#E8DFD0"),
		Text:          ui.Hex("#211C1D"),
		TextMuted:     ui.Hex("#6A5B55"),
		Border:        ui.Hex("#D2C3AF"),
		ControlBorder: ui.Hex("#877662"),
		Accent:        ui.Hex("#7D2034"),
		AccentHover:   ui.Hex("#651729"),
		AccentPressed: ui.Hex("#511020"),
		OnAccent:      ui.Hex("#FFF8EE"),
		AccentText:    ui.Hex("#7D2034"),
		Selection:     ui.Hex("#EAD7D9"),
		Focus:         ui.Hex("#7D2034"),
		Ornament:      ui.Hex("#AD9164"),
		Success:       ui.Hex("#326247"),
		Warning:       ui.Hex("#805714"),
		Danger:        ui.Hex("#A12635"),
		Info:          ui.Hex("#355E80"),
		DangerFill:    ui.Hex("#A12635"),
		OnDanger:      ui.Hex("#FFF8EE"),
	}
}

// Dark returns the charcoal appearance.
func Dark() Tokens {
	return Tokens{
		Dark:          true,
		Background:    ui.Hex("#191618"),
		Surface:       ui.Hex("#242024"),
		SurfaceHover:  ui.Hex("#342C32"),
		Text:          ui.Hex("#F2EBDD"),
		TextMuted:     ui.Hex("#BCAEAA"),
		Border:        ui.Hex("#51424A"),
		ControlBorder: ui.Hex("#968278"),
		Accent:        ui.Hex("#9C3C53"),
		AccentHover:   ui.Hex("#AA425C"),
		AccentPressed: ui.Hex("#813047"),
		OnAccent:      ui.Hex("#FFF8EE"),
		AccentText:    ui.Hex("#E4A0AE"),
		Selection:     ui.Hex("#492A36"),
		Focus:         ui.Hex("#E4A0AE"),
		Ornament:      ui.Hex("#BDA172"),
		Success:       ui.Hex("#94C6A4"),
		Warning:       ui.Hex("#DFC080"),
		Danger:        ui.Hex("#F0A1AA"),
		Info:          ui.Hex("#9DBDDA"),
		DangerFill:    ui.Hex("#A12635"),
		OnDanger:      ui.Hex("#FFF8EE"),
	}
}

// HighContrast returns t with secondary text and borders raised for readability.
func (t Tokens) HighContrast() Tokens {
	t.TextMuted = t.Text
	t.Border = t.ControlBorder
	return t
}

// Density sets control heights and the spacing unit.
type Density int

const (
	Default Density = iota
	Compact
	Comfortable
)

// Unit is the base spacing unit in DIPs.
func (d Density) Unit() float32 {
	switch d {
	case Compact:
		return 3
	case Comfortable:
		return 5
	}
	return 4
}

// ControlHeight is the minimum height of single-line controls and data rows.
func (d Density) ControlHeight() float32 {
	switch d {
	case Compact:
		return 28
	case Comfortable:
		return 44
	}
	return 36
}

func (d Density) String() string {
	switch d {
	case Compact:
		return "Compact"
	case Comfortable:
		return "Comfortable"
	}
	return "Default"
}

// Type sizes and families.
const (
	H1Size      float32 = 28
	H2Size      float32 = 22
	H3Size      float32 = 18
	BodySize    float32 = 14
	CaptionSize float32 = 12
	LineHeight  float32 = 1.5

	SerifFont = "Georgia, Songti SC, SimSun, Noto Serif CJK SC, serif"
	MonoFont  = "monospace"
)

// Shape metrics in DIPs.
const (
	ControlRadius float32 = 2
	CardRadius    float32 = 4
	PopupRadius   float32 = 6
	BorderWidth   float32 = 1
	FocusWidth    float32 = 2
	FocusGap      float32 = 2
	MinHitSize    float32 = 28
)

// Motion durations.
const (
	StateDuration = 120 * time.Millisecond
	PopupDuration = 160 * time.Millisecond
)
