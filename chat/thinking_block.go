package chat

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ThinkingBlockOptions configure a ThinkingBlock.
type ThinkingBlockOptions struct {
	// Thinking shows the spinner and a running timer from Started.
	Thinking bool
	Started  time.Time
	// Elapsed is the final thinking time, shown once Thinking is false.
	Elapsed time.Duration
}

// ThinkingBlock folds a model's reasoning under a header with its time; *open's first value sets the default.
func ThinkingBlock(c *ui.Context, open *bool, opts ThinkingBlockOptions, content func()) core.ChoiceResult {
	if opts.Elapsed < 0 || opts.Thinking && opts.Started.IsZero() {
		panic("mujica: ThinkingBlock needs Started while Thinking and a non-negative Elapsed")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	secs := int(opts.Elapsed.Round(time.Second).Seconds())
	title := fmt.Sprintf(core.Msg(c, "thinking_block.done"), secs)
	if opts.Thinking {
		run := c.Now().Sub(opts.Started)
		secs = int(max(run, 0) / time.Second)
		title = fmt.Sprintf(core.Msg(c, "thinking_block.thinking"), secs)
		c.After(time.Second - run%time.Second)
	}
	return data.Collapsible(c, open, data.CollapsibleOptions{Title: title, Trigger: func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.TextMuted).Children(func() {
			if opts.Thinking {
				core.Spinner(c, core.SpinnerOptions{Size: 14, Label: title})
			} else {
				ui.Icon(c, icons.Must("brain")).Size(16, 16)
			}
			ui.Text(c, title).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
		})
	}}, func() {
		ui.Column(c).Gap(u*2).Padding(0, 0, 0, u*3).BorderWidth(0, 0, 0, theme.BorderWidth).BorderColor(k.Border).
			TextColor(k.TextMuted).Children(content)
	})
}
