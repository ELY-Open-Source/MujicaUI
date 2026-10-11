package chat

import (
	"log/slog"
	"runtime"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// StopGeneratingButtonOptions configure a StopGeneratingButton.
type StopGeneratingButtonOptions struct {
	// Generating shows the button; it hides once generation ends.
	Generating bool
	Disabled   bool
}

// StopGeneratingButton asks to stop generating; Mod+. does too.
func StopGeneratingButton(c *ui.Context, opts StopGeneratingButtonOptions) core.ChoiceResult {
	if !opts.Generating {
		return core.ChoiceResult{}
	}
	k := core.Tokens(c)
	d := core.Density(c)
	mods := ui.Ctrl
	if runtime.GOOS == "darwin" {
		mods = ui.Cmd
	}
	var b *ui.Element
	wrap := ui.Row(c).Shrink(0).Disabled(opts.Disabled)
	wrap.Children(func() {
		b = ui.ButtonBase(c).MinHeight(max(theme.MinHitSize, d.ControlHeight()-4)).Padding(0, d.Unit()*3).Gap(d.Unit()*2).
			Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Label(core.Msg(c, "stop_generating_button.stop"))
	})
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() {
		b.Background(k.SurfaceHover)
	}
	b.Transition(core.ColorFade)
	b.Children(func() {
		ui.Box(c).Size(10, 10).Shrink(0).Background(k.Text).Radius(1)
		ui.Text(c, core.Msg(c, "stop_generating_button.stop")).SingleLine()
		display.Kbd(c, "Mod", ".")
	})
	stop := b.Clicked()
	if c.Shortcut(mods, ui.KeyPeriod) {
		stop = true
	}
	stop = stop && !b.IsDisabled()
	if stop {
		slog.Debug("mujica: stop_generating_button stop")
	}
	return core.NewChoiceResult(b, stop)
}
