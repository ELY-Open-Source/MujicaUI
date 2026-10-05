package agent

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ScreenshotFrame is one caller-supplied frame of the agent's screen.
type ScreenshotFrame struct {
	Image   ui.ImageSource
	Caption string
	// At is the time since the run started.
	At time.Duration
}

// ScreenshotStreamState is the shown frame; unpaused, it follows the newest.
type ScreenshotStreamState struct {
	Index  int
	Paused bool
}

// ScreenshotStreamOptions configure a ScreenshotStream.
type ScreenshotStreamOptions struct {
	// Height of the viewport; 280 DIPs when zero.
	Height   float32
	Disabled bool
}

// ScreenshotStream shows screen frames live, paused or stepped.
func ScreenshotStream(c *ui.Context, st *ScreenshotStreamState, frames []ScreenshotFrame, opts ScreenshotStreamOptions) core.ChoiceResult {
	if st == nil {
		panic("mujica: ScreenshotStream needs a state")
	}
	if opts.Height < 0 {
		panic("mujica: ScreenshotStream Height is negative")
	}
	for i, f := range frames {
		if f.Image == nil {
			panic(fmt.Sprintf("mujica: ScreenshotStream frame %d has no Image", i))
		}
	}
	h := opts.Height
	if h == 0 {
		h = 280
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	n := len(frames)
	if st.Index < 0 {
		panic(fmt.Sprintf("mujica: ScreenshotStream Index %d is negative", st.Index))
	}
	if !st.Paused || st.Index >= n { // frames may be dropped while paused

		st.Index = max(0, n-1)
	}
	var r core.ChoiceResult
	r.Element = AgentCard(c).Disabled(opts.Disabled).Label(core.Msg(c, "screenshot_stream.label"))
	locked := r.Element.IsDisabled()
	step := func(d int) {
		if n == 0 || locked {
			return
		}
		st.Index = max(0, min(n-1, st.Index+d))
		st.Paused, r = true, core.NewChoiceResult(r.Element, true)
	}
	r.Element.Children(func() {
		view := ui.Box(c).FillWidth().Height(h).Center().Clip().Radius(theme.ControlRadius).
			Background(k.Background).Border(theme.BorderWidth, k.Border).Focusable().
			Label(core.Msg(c, "screenshot_stream.label"))
		internal.Ring(view, k.Focus, theme.ControlRadius)
		view.Children(func() {
			if n == 0 {
				AgentCaption(c, core.Msg(c, "screenshot_stream.waiting"))
				return
			}
			ui.Image(c, frames[st.Index].Image).Fill().Fit(ui.Contain).Label(frames[st.Index].Caption)
		})
		if !locked {
			switch {
			case view.Shortcut(0, ui.KeyLeft):
				step(-1)
			case view.Shortcut(0, ui.KeyRight):
				step(1)
			case view.Shortcut(0, ui.KeySpace):
				st.Paused, r = !st.Paused, core.NewChoiceResult(r.Element, true)
			}
		}
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			if core.IconAction(c, icons.Must("chevron-left"), core.Msg(c, "screenshot_stream.prev")).Disabled(locked || st.Index == 0).Clicked() {
				step(-1)
			}
			pause := core.Msg(c, "screenshot_stream.pause")
			if st.Paused {
				pause = core.Msg(c, "screenshot_stream.live")
			}
			if input.Button(c, pause, input.ButtonOptions{Variant: input.Secondary, Disabled: locked || n == 0}).Clicked() && !locked {
				st.Paused, r = !st.Paused, core.NewChoiceResult(r.Element, true)
			}
			if core.IconAction(c, icons.Must("chevron-right"), core.Msg(c, "screenshot_stream.next")).Disabled(locked || st.Index >= n-1).Clicked() {
				step(1)
			}
			counter := fmt.Sprintf(core.Msg(c, "screenshot_stream.counter"), min(st.Index+1, n), n)
			AgentCaption(c, counter).Role(ui.RoleStatus)
			if n > 0 {
				f := frames[st.Index]
				ui.Text(c, f.Caption).SingleLine().Grow(1).Shrink(1)
				AgentCaption(c, AgentDuration(f.At))
			}
			if !st.Paused && n > 0 {
				AgentStatus(c, AgentRunning, AgentStatusOptions{Label: core.Msg(c, "screenshot_stream.following"), Still: true})
			}
		})
	})
	if r.Changed() {
		slog.Debug("mujica: screenshot_stream move", "index", st.Index, "paused", st.Paused)
	}
	return r
}
