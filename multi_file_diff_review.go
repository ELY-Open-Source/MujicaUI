package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// MultiFileDiffReviewOptions configure a MultiFileDiffReview.
type MultiFileDiffReviewOptions struct {
	// Height is 420 DIPs when zero.
	Height   float32
	Disabled bool
}

type multiFileDiffState struct{ selected int }

// MultiFileDiffReview reviews many file changes, one diff at a time.
func MultiFileDiffReview(c *ui.Context, decisions *[]FileDecision, files []FileChange, opts MultiFileDiffReviewOptions) ChoiceResult {
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
	k := Tokens(c)
	u := Density(c).Unit()
	var r ChoiceResult
	r.Element = agentCard(c).Height(h).Disabled(opts.Disabled).Label(msg(c, "multi_file_diff_review.title"))
	st := ui.Local(r.Element, "mujica.mfdr", func() multiFileDiffState { return multiFileDiffState{} })
	st.selected = max(0, min(st.selected, len(files)-1))
	locked := r.Element.IsDisabled()
	r.Element.Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			ui.Text(c, msg(c, "multi_file_diff_review.title")).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600)
			agentCaption(c, fmt.Sprintf(msg(c, "multi_file_diff_review.progress"), reviewed, len(files))).Grow(1)
			all := reviewed < len(files)
			if Button(c, msg(c, "multi_file_diff_review.accept_all"), ButtonOptions{Variant: Secondary, Disabled: locked || !all}).Clicked() && !locked && all {
				for i := range *decisions {
					if (*decisions)[i] == FilePending {
						(*decisions)[i] = FileAccepted
					}
				}
				r.changed = true
			}
		})
		frac := 0.0
		if len(files) > 0 {
			frac = float64(reviewed) / float64(len(files))
		}
		agentBar(c, frac, fmt.Sprintf(msg(c, "multi_file_diff_review.progress"), reviewed, len(files)))
		if len(files) == 0 {
			agentCaption(c, msg(c, "multi_file_diff_review.empty"))
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
					r.changed = true
				}
			})
		})
	})
	if r.changed {
		slog.Debug("mujica: multi_file_diff_review decide", "reviewed", reviewed)
	}
	return r
}

// multiFileRow is one listed file; it reports a choice.
func multiFileRow(c *ui.Context, f FileChange, d FileDecision, selected bool) bool {
	k := Tokens(c)
	u := Density(c).Unit()
	row := ui.ButtonBase(c).Key(f.Path).Justify(ui.Start).Gap(u*2).MinHeight(Density(c).ControlHeight()).
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
		ui.Text(c, f.Path).Font(theme.MonoFont).FontSize(FontSize(c, theme.CaptionSize+1)).SingleLine().Grow(1).Shrink(1)
		agentDiffStats(c, dc.add, dc.del)
	})
	return row.Clicked()
}

// fileChangeDecide shows a full diff and its decision buttons.
func fileChangeDecide(c *ui.Context, d *FileDecision, f FileChange, locked bool) bool {
	u := Density(c).Unit()
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
				if Button(c, msg(c, "file_change_card.undo"), ButtonOptions{Variant: Secondary, Disabled: locked}).Clicked() && !locked {
					*d, changed = FilePending, true
				}
				return
			}
			if Button(c, msg(c, "file_change_card.reject"), ButtonOptions{Variant: Secondary, Disabled: locked}).Clicked() && !locked {
				*d, changed = FileRejected, true
			}
			if Button(c, msg(c, "file_change_card.accept"), ButtonOptions{Disabled: locked}).Clicked() && !locked {
				*d, changed = FileAccepted, true
			}
		})
	})
	return changed
}
