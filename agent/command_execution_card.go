package agent

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CommandRun is a command run; the caller appends Output.
type CommandRun struct {
	Command string
	Dir     string
	// Output is the raw terminal output; ANSI colors are kept.
	Output   string
	Running  bool
	ExitCode int
	Duration time.Duration
}

// CommandExecutionCardOptions configure a CommandExecutionCard.
type CommandExecutionCardOptions struct {
	// Height of the output well; 240 DIPs when zero.
	Height   float32
	Disabled bool
	Still    bool
}

// commandOutput keeps parsed lines; appends parse only new lines.
type commandOutput struct {
	src   string
	lines [][]internal.AnsiSpan
	style internal.AnsiStyle
	list  ui.ListState
}

func (o *commandOutput) update(out string) {
	if out == o.src {
		return
	}
	if !strings.HasPrefix(out, o.src) || o.src != "" && !strings.HasSuffix(o.src, "\n") {
		o.lines, o.style = nil, internal.AnsiStyle{}
		o.src = ""
	}
	for _, line := range internal.SplitLines(out[len(o.src):]) {
		var spans []internal.AnsiSpan
		spans, o.style = internal.ParseAnsi(line, o.style)
		o.lines = append(o.lines, spans)
	}
	o.src = out
}

// CommandExecutionCard shows a command run; Changed requests an abort.
func CommandExecutionCard(c *ui.Context, run CommandRun, opts CommandExecutionCardOptions) core.ChoiceResult {
	if run.Command == "" {
		panic("mujica: CommandExecutionCard needs a Command")
	}
	if opts.Height < 0 {
		panic("mujica: CommandExecutionCard Height is negative")
	}
	h := opts.Height
	if h == 0 {
		h = 240
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = AgentCard(c).Disabled(opts.Disabled).Label(run.Command)
	out := ui.Local(r.Element, "mujica.cmd.output", func() *commandOutput { return &commandOutput{} })
	o := *out
	o.list.FollowEnd = true
	o.update(run.Output)
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Children(func() {
			ui.Text(c, "$").Font(theme.MonoFont).TextColor(k.TextMuted).Shrink(0)
			ui.Text(c, run.Command).Font(theme.MonoFont).FontWeight(600).Grow(1).Shrink(1).Selectable()
			if run.Running {
				agentMark(c, AgentRunning, 16, opts.Still)
				AgentCaption(c, core.Msg(c, "agent_status.running"))
				locked := r.Element.IsDisabled()
				if input.Button(c, core.Msg(c, "command_execution_card.abort"), input.ButtonOptions{Variant: input.Outline, Disabled: locked}).Clicked() && !locked {
					r = core.NewChoiceResult(r.Element, true)
				}
				return
			}
			commandExit(c, run.ExitCode)
			AgentCaption(c, AgentDuration(run.Duration))
		})
		if run.Dir != "" {
			ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Must("folder")).Size(14, 14).TextColor(k.TextMuted)
				AgentCaption(c, run.Dir).Font(theme.MonoFont)
			})
		}
		well := ui.List(c, &o.list, len(o.lines), func(i int) {
			ui.RichText(c, agentAnsiSpans(k, o.lines[i])...)
		})
		agentWell(c, well, h).Height(h).Label(core.Msg(c, "command_execution_card.output")).Children(func() {
			if len(o.lines) == 0 {
				AgentCaption(c, core.Msg(c, "command_execution_card.no_output"))
			}
		})
	})
	if r.Changed() {
		slog.Debug("mujica: command_execution_card abort", "command", run.Command)
	}
	return r
}

// commandExit shows the exit code by icon and word.
func commandExit(c *ui.Context, code int) {
	k := core.Tokens(c)
	tone, icon := k.Success, "circle-check"
	if code != 0 {
		tone, icon = k.Danger, "circle-x"
	}
	ui.Row(c).Gap(core.Density(c).Unit()).AlignItems(ui.Center).Shrink(0).Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(tone)
		ui.Text(c, fmt.Sprintf(core.Msg(c, "command_execution_card.exit"), code)).TextColor(tone).
			FontSize(core.FontSize(c, theme.CaptionSize+1)).SingleLine()
	})
}
