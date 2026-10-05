// Package mujica is a court-gothic component library for MyGo native UI.
package mujica

import (
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// Mode chooses the appearance of a window.
type Mode int

const (
	System Mode = iota
	Light
	Dark
)

func (m Mode) String() string {
	switch m {
	case Light:
		return "Light"
	case Dark:
		return "Dark"
	}
	return "System"
}

// Settings are a window's theme. Nil token sets use the library defaults.
type Settings struct {
	Mode    Mode
	Density theme.Density
	Light   *theme.Tokens
	Dark    *theme.Tokens
	Locale  Locale
	// Messages overrides built-in copy by key.
	Messages map[string]string
}

type windowTheme struct {
	tokens   theme.Tokens
	density  theme.Density
	locale   Locale
	messages map[string]string
	set      bool
}

type themeKey struct{}

// Use applies s to the window. Call it first in every frame of the view.
func Use(c *ui.Context, s Settings) {
	dark := c.Theme().Dark
	switch s.Mode {
	case Light:
		dark = false
	case Dark:
		dark = true
	}
	k := theme.Light()
	if s.Light != nil {
		k = *s.Light
	}
	if dark {
		k = theme.Dark()
		if s.Dark != nil {
			k = *s.Dark
		}
	}
	p := c.Preferences()
	if p.HighContrast {
		k = k.HighContrast()
	}
	size := theme.BodySize
	if p.TextScale > 0 {
		size *= p.TextScale
	}
	c.SetTheme(&ui.Theme{
		Dark:           k.Dark,
		Background:     k.Background,
		Surface:        k.Surface,
		SurfaceHover:   k.SurfaceHover,
		SurfacePressed: k.Selection,
		Border:         k.ControlBorder,
		Text:           k.Text,
		TextMuted:      k.TextMuted,
		Accent:         k.Accent,
		AccentHover:    k.AccentHover,
		AccentPressed:  k.AccentPressed,
		AccentText:     k.OnAccent,
		Danger:         k.Danger,
		Warning:        k.Warning,
		Success:        k.Success,
		Selection:      k.Selection,
		Focus:          k.Focus,
		Scrollbar:      k.TextMuted.Alpha(0.5),
		Radius:         theme.ControlRadius,
		Spacing:        s.Density.Unit(),
		FontSize:       size,
	})
	*ui.Local(c.Root(), themeKey{}, func() windowTheme { return windowTheme{} }) = windowTheme{k, s.Density, s.Locale, s.Messages, true}
}

func current(c *ui.Context) windowTheme {
	w := *ui.Local(c.Root(), themeKey{}, func() windowTheme { return windowTheme{} })
	if !w.set {
		panic("mujica: call mujica.Use at the start of the view")
	}
	return w
}

// Tokens returns the window's colors.
func Tokens(c *ui.Context) theme.Tokens { return current(c).tokens }

// Density returns the window's density.
func Density(c *ui.Context) theme.Density { return current(c).density }

// Motion returns d, or zero when the desktop asks for reduced motion.
func Motion(c *ui.Context, d time.Duration) time.Duration {
	if c.Preferences().ReduceMotion {
		return 0
	}
	return d
}
