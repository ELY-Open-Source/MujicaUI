package git

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/agent"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CommandBlockOptions configure a CommandBlock.
type CommandBlockOptions struct {
	// Prompt precedes the command; empty means "$".
	Prompt string
	// MaxOutputHeight scrolls longer output; zero means 320 DIPs.
	MaxOutputHeight float32
}

// CommandBlockResult is a built CommandBlock and its events.
type CommandBlockResult struct {
	Element *ui.Element
	toggled bool
	copied  string
}

// Toggled reports the output was shown or hidden this frame.
func (r CommandBlockResult) Toggled() bool { return r.toggled }

// Copied returns "command" or "output" when copied.
func (r CommandBlockResult) Copied() (string, bool) { return r.copied, r.copied != "" }

// CommandBlock shows a command run; *open shows its output.
func CommandBlock(c *ui.Context, run agent.CommandRun, open *bool, opts CommandBlockOptions) CommandBlockResult {
	if run.Command == "" {
		panic("mujica: CommandBlock needs a Command")
	}
	if run.Duration < 0 || opts.MaxOutputHeight < 0 {
		panic("mujica: CommandBlock Duration and MaxOutputHeight must not be negative")
	}
	prompt := opts.Prompt
	if prompt == "" {
		prompt = "$"
	}
	maxH := opts.MaxOutputHeight
	if maxH == 0 {
		maxH = 320
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	box := ui.Column(c).AlignItems(ui.Stretch).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).Clip()
	res := CommandBlockResult{Element: box}
	box.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(u).Padding(0, u*2, 0, 0).MinHeight(d.ControlHeight()).Children(func() {
			was := *open
			tr := ui.CollapsibleBase(c, open).Trigger
			if box.IsDisabled() {
				*open = was
			}
			if *open != was {
				res.toggled = true
				slog.Debug("mujica: command block toggle", "open", *open)
			}
			tr.Grow(1).MinWidth(0).MinHeight(d.ControlHeight()).AlignItems(ui.Center).Gap(u*2).Padding(0, u*3).
				Radius(theme.ControlRadius).Label(prompt + " " + run.Command).Cursor(ui.CursorPointer)
			internal.Ring(tr, k.Focus, theme.ControlRadius)
			if tr.Hovered() {
				tr.Background(k.SurfaceHover)
			}
			tr.Transition(core.ColorFade)
			turn := tr.Animate("mujica.turn", commandBlockTurn(*open), core.Motion(c, theme.StateDuration))
			tr.Children(func() {
				ui.Icon(c, icons.Must("chevron-right")).Size(14, 14).Shrink(0).Rotate(turn).TextColor(k.TextMuted)
				ui.Text(c, prompt).Font(theme.MonoFont).TextColor(k.Ornament.Mix(k.Text, 0.35)).Shrink(0)
				spans, _ := internal.CodeSpans(k, "shell", run.Command, false)
				ui.RichText(c, spans...).Font(theme.MonoFont).SingleLine().MinWidth(0)
			})
			commandBlockStatus(c, run)
			if core.IconAction(c, icons.Must("copy"), core.Msg(c, "command_block.copy_command")).Clicked() && !box.IsDisabled() {
				c.WriteClipboard(run.Command)
				res.copied = "command"
			}
			if core.IconAction(c, icons.Must("file-text"), core.Msg(c, "command_block.copy_output")).Disabled(run.Output == "").Clicked() && !box.IsDisabled() {
				c.WriteClipboard(internal.StripAnsi(run.Output))
				res.copied = "output"
			}
		})
		if *open {
			ui.Box(c).Height(1).Background(k.Border).Role(ui.RoleNone)
			out := ui.ScrollBoth(c).MaxHeight(maxH).Padding(u*3, u*4).Background(k.Background).Label(core.Msg(c, "command_block.output"))
			out.Children(func() {
				if run.Output == "" {
					ui.Text(c, core.Msg(c, "command_block.no_output")).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
					return
				}
				AnsiText(c, run.Output, AnsiTextOptions{})
			})
		}
	})
	if res.copied != "" {
		slog.Debug("mujica: command block copy", "what", res.copied)
	}
	return res
}

// commandBlockStatus shows running, success or failure with icon and words.
func commandBlockStatus(c *ui.Context, run agent.CommandRun) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ui.Row(c).AlignItems(ui.Center).Gap(u * 1.5).Shrink(0).FontSize(core.FontSize(c, theme.CaptionSize)).Children(func() {
		if run.Running {
			core.Spinner(c, core.SpinnerOptions{Label: core.Msg(c, "command_block.running"), ShowLabel: true, Size: 12})
			return
		}
		icon, col := "circle-check", k.Success
		if run.ExitCode != 0 {
			icon, col = "circle-x", k.Danger
		}
		ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(col).Children(func() {
			ui.Icon(c, icons.Must(icon)).Size(14, 14)
			ui.Text(c, fmt.Sprintf(core.Msg(c, "command_block.exit"), run.ExitCode)).FontFeatures("tnum").SingleLine()
		})
		ui.Text(c, commandBlockDuration(run.Duration)).TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
	})
}

func commandBlockTurn(open bool) float32 {
	if open {
		return 90
	}
	return 0
}

// commandBlockDuration prints milliseconds, seconds or minutes.
func commandBlockDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1f s", d.Seconds())
	}
	return fmt.Sprintf("%dm %02ds", int(d.Minutes()), int(d.Seconds())%60)
}
