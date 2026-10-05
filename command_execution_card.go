package mujica

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
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
func CommandExecutionCard(c *ui.Context, run CommandRun, opts CommandExecutionCardOptions) ChoiceResult {
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
	k := Tokens(c)
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = agentCard(c).Disabled(opts.Disabled).Label(run.Command)
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
				agentCaption(c, msg(c, "agent_status.running"))
				locked := r.Element.IsDisabled()
				if Button(c, msg(c, "command_execution_card.abort"), ButtonOptions{Variant: Outline, Disabled: locked}).Clicked() && !locked {
					r.changed = true
				}
				return
			}
			commandExit(c, run.ExitCode)
			agentCaption(c, agentDuration(run.Duration))
		})
		if run.Dir != "" {
			ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Must("folder")).Size(14, 14).TextColor(k.TextMuted)
				agentCaption(c, run.Dir).Font(theme.MonoFont)
			})
		}
		well := ui.List(c, &o.list, len(o.lines), func(i int) {
			ui.RichText(c, agentAnsiSpans(k, o.lines[i])...)
		})
		agentWell(c, well, h).Height(h).Label(msg(c, "command_execution_card.output")).Children(func() {
			if len(o.lines) == 0 {
				agentCaption(c, msg(c, "command_execution_card.no_output"))
			}
		})
	})
	if r.changed {
		slog.Debug("mujica: command_execution_card abort", "command", run.Command)
	}
	return r
}

// commandExit shows the exit code by icon and word.
func commandExit(c *ui.Context, code int) {
	k := Tokens(c)
	tone, icon := k.Success, "circle-check"
	if code != 0 {
		tone, icon = k.Danger, "circle-x"
	}
	ui.Row(c).Gap(Density(c).Unit()).AlignItems(ui.Center).Shrink(0).Children(func() {
		ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(tone)
		ui.Text(c, fmt.Sprintf(msg(c, "command_execution_card.exit"), code)).TextColor(tone).
			FontSize(FontSize(c, theme.CaptionSize+1)).SingleLine()
	})
}
