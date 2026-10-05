package mujica

import (
	"log/slog"
	"runtime"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// StopGeneratingButtonOptions configure a StopGeneratingButton.
type StopGeneratingButtonOptions struct {
	// Generating shows the button; it hides once generation ends.
	Generating bool
	Disabled   bool
}

// StopGeneratingButton asks to stop a reply while one is generating; Cmd+. (Ctrl+. elsewhere) does too.
func StopGeneratingButton(c *ui.Context, opts StopGeneratingButtonOptions) ChoiceResult {
	if !opts.Generating {
		return ChoiceResult{}
	}
	k := Tokens(c)
	d := Density(c)
	mods := ui.Ctrl
	if runtime.GOOS == "darwin" {
		mods = ui.Cmd
	}
	var b *ui.Element
	wrap := ui.Row(c).Shrink(0).Disabled(opts.Disabled)
	wrap.Children(func() {
		b = ui.ButtonBase(c).MinHeight(max(theme.MinHitSize, d.ControlHeight()-4)).Padding(0, d.Unit()*3).Gap(d.Unit()*2).
			Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Label(msg(c, "stop_generating_button.stop"))
	})
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Transition(colorFade)
	b.Children(func() {
		ui.Box(c).Size(10, 10).Shrink(0).Background(k.Text).Radius(1)
		ui.Text(c, msg(c, "stop_generating_button.stop")).SingleLine()
		Kbd(c, "Mod", ".")
	})
	stop := b.Clicked()
	if c.Shortcut(mods, ui.KeyPeriod) {
		stop = true
	}
	stop = stop && !b.IsDisabled()
	if stop {
		slog.Debug("mujica: stop_generating_button stop")
	}
	return ChoiceResult{Element: b, changed: stop}
}
