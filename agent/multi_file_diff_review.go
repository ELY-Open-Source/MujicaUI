package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MultiFileDiffReviewOptions configure a MultiFileDiffReview.
type MultiFileDiffReviewOptions struct {
	// Height is 420 DIPs when zero.
	Height   float32
	Disabled bool
}

type multiFileDiffState struct{ selected int }

// MultiFileDiffReview reviews many file changes, one diff at a time.
func MultiFileDiffReview(c *ui.Context, decisions *[]FileDecision, files []FileChange, opts MultiFileDiffReviewOptions) core.ChoiceResult {
	if decisions == nil || len(*decisions) != len(files) {
		panic("mujica: MultiFileDiffReview needs one decision per file")
	}
	if opts.Height < 0 {
		panic("mujica: MultiFileDiffReview Height is negative")
	}
	seen := map[string]bool{}
	reviewed := 0
	for i, f := range files {
		if f.Path == "" || seen[f.Path] {
			panic(fmt.Sprintf("mujica: MultiFileDiffReview file %d needs a unique Path", i))
		}
		seen[f.Path] = true
		reviewed += agentB2I((*decisions)[i] != FilePending)
	}
	h := opts.Height
	if h == 0 {
		h = 420
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r core.ChoiceResult
	r.Element = AgentCard(c).Height(h).Disabled(opts.Disabled).Label(core.Msg(c, "multi_file_diff_review.title"))
	st := ui.Local(r.Element, "mujica.mfdr", func() multiFileDiffState { return multiFileDiffState{} })
	st.selected = max(0, min(st.selected, len(files)-1))
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, "multi_file_diff_review.title")).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600)
			AgentCaption(c, fmt.Sprintf(core.Msg(c, "multi_file_diff_review.progress"), reviewed, len(files))).Grow(1)
			all := reviewed < len(files)
			if input.Button(c, core.Msg(c, "multi_file_diff_review.accept_all"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked || !all}).Clicked() && !locked && all {
				for i := range *decisions {
					if (*decisions)[i] == FilePending {
						(*decisions)[i] = FileAccepted
					}
				}
				r = core.NewChoiceResult(r.Element, true)
			}
		})
		frac := 0.0
		if len(files) > 0 {
			frac = float64(reviewed) / float64(len(files))
		}
		agentBar(c, frac, fmt.Sprintf(core.Msg(c, "multi_file_diff_review.progress"), reviewed, len(files)))
		if len(files) == 0 {
			AgentCaption(c, core.Msg(c, "multi_file_diff_review.empty"))
			return
		}
		ui.Row(c).Grow(1).Shrink(1).Gap(u * 3).AlignItems(ui.Stretch).Children(func() {
			row := func(i int) {
				if multiFileRow(c, files[i], (*decisions)[i], i == st.selected) && !locked {
					st.selected = i
					slog.Debug("mujica: multi_file_diff_review select", "path", files[i].Path)
				}
			}
			var list *ui.Element
			if len(files) > agentVirtualRows {
				list = ui.List(c, nil, len(files), row)
			} else {
				list = ui.Scroll(c).Children(func() {
					for i := range files {
						row(i)
					}
				})
			}
			list.Width(240).Shrink(0).FocusGroup(ui.Vertical).Border(theme.BorderWidth, k.Border).
				Radius(theme.ControlRadius).Padding(u)
			ui.Column(c).Grow(1).Shrink(1).Gap(u * 2).Children(func() {
				i := st.selected
				if fileChangeDecide(c, &(*decisions)[i], files[i], locked) {
					r = core.NewChoiceResult(r.Element, true)
				}
			})
		})
	})
	if r.Changed() {
		slog.Debug("mujica: multi_file_diff_review decide", "reviewed", reviewed)
	}
	return r
}

// multiFileRow is one listed file; it reports a choice.
func multiFileRow(c *ui.Context, f FileChange, d FileDecision, selected bool) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	row := ui.ButtonBase(c).Key(f.Path).Justify(ui.Start).Gap(u*2).MinHeight(core.Density(c).ControlHeight()).
		Padding(0, u*2).Radius(theme.ControlRadius).Label(f.Path)
	internal.Ring(row, k.Focus, theme.ControlRadius)
	dc := agentDiffOf(row, f)
	switch {
	case selected:
		row.Background(k.Selection).DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 2, H: r.H - 8}, k.AccentText, 1)
		})
	case row.Hovered():
		row.Background(k.SurfaceHover)
	}
	row.Children(func() {
		agentDecisionIcon(c, d)
		ui.Text(c, f.Path).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize+1)).SingleLine().Grow(1).Shrink(1)
		agentDiffStats(c, dc.add, dc.del)
	})
	return row.Clicked()
}

// fileChangeDecide shows a full diff and its decision buttons.
func fileChangeDecide(c *ui.Context, d *FileDecision, f FileChange, locked bool) bool {
	u := core.Density(c).Unit()
	changed := false
	holder := ui.Column(c).Key(f.Path).Grow(1).Shrink(1).Gap(u * 2)
	dc := agentDiffOf(holder, f)
	holder.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, f.Path).Font(theme.MonoFont).FontWeight(600).SingleLine().Grow(1).Shrink(1)
			agentDecisionMark(c, *d)
		})
		agentDiffView(c, dc.lines, 1e5).Grow(1).Font(theme.MonoFont)
		ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
			if *d != FilePending {
				if input.Button(c, core.Msg(c, "file_change_card.undo"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked() && !locked {
					*d, changed = FilePending, true
				}
				return
			}
			if input.Button(c, core.Msg(c, "file_change_card.reject"), input.ButtonOptions{Variant: input.Secondary, Disabled: locked}).Clicked() && !locked {
				*d, changed = FileRejected, true
			}
			if input.Button(c, core.Msg(c, "file_change_card.accept"), input.ButtonOptions{Disabled: locked}).Clicked() && !locked {
				*d, changed = FileAccepted, true
			}
		})
	})
	return changed
}
