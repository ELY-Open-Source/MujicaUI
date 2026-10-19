package agent

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// FileChangeCardOptions configure a FileChangeCard.
type FileChangeCardOptions struct {
	// Preview is the diff lines shown; zero means 8.
	Preview  int
	Disabled bool
}

// FileChangeResult reports this frame's events.
type FileChangeResult struct {
	Element ui.Element
	// Decided reports that *decision changed.
	Decided bool
	// OpenDiff asks the caller to show the full diff.
	OpenDiff bool
}

// FileChangeCard previews a file change to accept or reject.
func FileChangeCard(c *ui.Context, decision *FileDecision, change FileChange, opts FileChangeCardOptions) FileChangeResult {
	if decision == nil || change.Path == "" {
		panic("mujica: FileChangeCard needs a decision and a Path")
	}
	if opts.Preview < 0 {
		panic("mujica: FileChangeCard Preview is negative")
	}
	preview := opts.Preview
	if preview == 0 {
		preview = 8
	}
	u := core.Density(c).Unit()
	var r FileChangeResult
	r.Element = AgentCard(c).Disabled(opts.Disabled).Label(change.Path)
	dc := agentDiffOf(r.Element, change)
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Children(func() {
			ui.Icon(c, icons.Must("file-text")).Size(16, 16).TextColor(core.Tokens(c).TextMuted).Shrink(0)
			ui.Text(c, change.Path).Font(theme.MonoFont).FontWeight(600).SingleLine().Grow(1).Shrink(1)
			switch {
			case change.Old == "" && change.New != "":
				AgentCaption(c, core.Msg(c, "file_change_card.new"))
			case change.New == "" && change.Old != "":
				AgentCaption(c, core.Msg(c, "file_change_card.deleted"))
			}
			agentDiffStats(c, dc.add, dc.del)
			agentDecisionMark(c, *decision)
		})
		from := max(0, dc.firstEdit-2)
		shown := dc.lines[from:min(len(dc.lines), from+preview)]
		if len(shown) > 0 {
			agentDiffView(c, shown, 320).Font(theme.MonoFont)
		}
		ui.Row(c).Gap(u * 2).Justify(ui.End).Wrap().Children(func() {
			if input.Button(c, core.Msg(c, "file_change_card.open"), input.ButtonOptions{Variant: input.Ghost, Disabled: locked}).Clicked() && !locked {
				r.OpenDiff = true
			}
			if *decision != FilePending {
				if input.Button(c, core.Msg(c, "file_change_card.undo"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked() && !locked {
					*decision, r.Decided = FilePending, true
				}
				return
			}
			if input.Button(c, core.Msg(c, "file_change_card.reject"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked() && !locked {
				*decision, r.Decided = FileRejected, true
			}
			if input.Button(c, core.Msg(c, "file_change_card.accept"), input.ButtonOptions{Disabled: locked}).Clicked() && !locked {
				*decision, r.Decided = FileAccepted, true
			}
		})
	})
	if r.Decided || r.OpenDiff {
		slog.Debug("mujica: file_change_card event", "path", change.Path, "decision", int(*decision), "open", r.OpenDiff)
	}
	return r
}
