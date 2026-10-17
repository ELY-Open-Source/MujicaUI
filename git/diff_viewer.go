package git

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DiffViewerState keeps a DiffViewer's diff, opened folds and scroll place.
type DiffViewerState struct {
	list     ui.ListState
	from, to string
	diff     []internal.DiffLine
	rows     []diffViewerRow
	expanded map[int]bool
	split    bool
	context  int
	built    bool
	dirty    bool
}

// DiffViewerOptions configure a DiffViewer.
type DiffViewerOptions struct {
	// SideBySide shows old and new in two columns.
	SideBySide bool
	// Context is unchanged lines kept around changes; zero means 3.
	Context int
	// Language picks syntax colors, as "go"; empty is plain.
	Language string
}

// diffViewerRow is a fold or one line per side.
type diffViewerRow struct {
	left, right  int
	fold, hidden int
}

// DiffViewer compares two texts by line; give it a height.
func DiffViewer(c *ui.Context, s *DiffViewerState, from, to string, opts DiffViewerOptions) ui.Element {
	ctx := opts.Context
	if ctx == 0 {
		ctx = 3
	}
	if ctx < 0 {
		panic("mujica: DiffViewer Context is negative")
	}
	if !s.built || s.from != from || s.to != to {
		s.diff = internal.LineDiff(internal.SplitLines(from), internal.SplitLines(to))
		s.from, s.to, s.built, s.dirty, s.expanded = from, to, true, true, map[int]bool{}
		add, del := internal.DiffStats(s.diff)
		slog.Debug("mujica: diff viewer diff", "lines", len(s.diff), "added", add, "deleted", del)
	}
	if s.dirty || s.split != opts.SideBySide || s.context != ctx {
		s.rows = diffViewerRows(s.diff, ctx, s.expanded, opts.SideBySide)
		s.split, s.context, s.dirty = opts.SideBySide, ctx, false
	}
	k := core.Tokens(c)
	digits := float32(len(strconv.Itoa(max(len(s.diff), 1))))
	numW := digits*8 + core.Density(c).Unit()*3
	e := ui.List(c, &s.list, len(s.rows), func(i int) {
		r := s.rows[i]
		if r.hidden > 0 {
			diffViewerFold(c, s, r)
			return
		}
		row := ui.Row(c).AlignItems(ui.Stretch).MinHeight(theme.MinHitSize)
		row.Children(func() {
			if !opts.SideBySide {
				diffViewerCell(c, s.diff[r.left], opts.Language, numW, true)
				return
			}
			for _, side := range []int{r.left, r.right} {
				half := ui.Row(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Stretch).Clip()
				half.Children(func() {
					if side < 0 {
						half.Background(k.Background)
						return
					}
					diffViewerCell(c, s.diff[side], opts.Language, numW, false)
				})
			}
		})
	})
	e.Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).Font(theme.MonoFont)
	return e
}

// diffViewerCell shows one diff line with numbers, sign and tint.
func diffViewerCell(c *ui.Context, d internal.DiffLine, language string, numW float32, both bool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	small := core.FontSize(c, theme.CaptionSize)
	sign, col, name := " ", k.TextMuted, ""
	switch d.Op {
	case internal.DiffDelete:
		sign, col, name = "−", k.Danger, core.Msg(c, "diff_viewer.deleted")
	case internal.DiffInsert:
		sign, col, name = "+", k.Success, core.Msg(c, "diff_viewer.added")
	}
	cell := ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center)
	if d.Op != internal.DiffEqual {
		cell.Background(col.Mix(k.Surface, 0.88))
	}
	num := func(n int) {
		t := ""
		if n > 0 {
			t = strconv.Itoa(n)
		}
		ui.Text(c, t).FontSize(small).FontFeatures("tnum").TextColor(k.TextMuted).Width(numW).TextAlign(ui.End).Shrink(0).Padding(0, u, 0, 0)
	}
	cell.Children(func() {
		if both || d.Op != internal.DiffInsert {
			num(d.Old)
		}
		if both || d.Op == internal.DiffInsert {
			num(d.New)
		}
		ui.Text(c, sign).FontWeight(700).TextColor(col).Width(u * 4).TextAlign(ui.Center).Shrink(0)
		spans, _ := internal.CodeSpans(k, language, d.Text, false)
		if len(spans) == 0 || d.Text == "" {
			spans = []ui.Span{{Text: " "}}
		}
		t := ui.RichText(c, spans...).NoWrap().MinWidth(0).Padding(u*0.5, u*2, u*0.5, 0)
		if name != "" {
			t.Label(name + ": " + d.Text)
		}
	})
}

// diffViewerFold is the button standing for hidden unchanged lines.
func diffViewerFold(c *ui.Context, s *DiffViewerState, r diffViewerRow) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	b := ui.ButtonBase(c).Justify(ui.Start).Gap(u*2).MinHeight(theme.MinHitSize).Padding(0, u*3).Background(k.Background).
		TextColor(k.TextMuted).Cursor(ui.CursorPointer)
	label := fmt.Sprintf(core.Msg(c, "diff_viewer.fold"), r.hidden)
	b.Label(label)
	internal.Ring(b, k.Focus, 0)
	if b.Hovered() {
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	b.Transition(core.ColorFade)
	b.Children(func() {
		ui.Text(c, "⋯").FontWeight(700)
		ui.Text(c, label).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize))
	})
	if b.Clicked() && !b.IsDisabled() {
		s.expanded[r.fold] = true
		s.dirty = true
		slog.Debug("mujica: diff viewer unfold", "at", r.fold, "lines", r.hidden)
	}
}

// diffViewerRows folds long unchanged runs and pairs split sides.
func diffViewerRows(diff []internal.DiffLine, ctx int, expanded map[int]bool, split bool) []diffViewerRow {
	var rows []diffViewerRow
	var dels, ins []int
	flush := func() {
		for j := 0; j < max(len(dels), len(ins)); j++ {
			r := diffViewerRow{left: -1, right: -1}
			if j < len(dels) {
				r.left = dels[j]
			}
			if j < len(ins) {
				r.right = ins[j]
			}
			rows = append(rows, r)
		}
		dels, ins = dels[:0], ins[:0]
	}
	line := func(i int) {
		switch {
		case !split:
			rows = append(rows, diffViewerRow{left: i, right: -1})
		case diff[i].Op == internal.DiffDelete:
			dels = append(dels, i)
		case diff[i].Op == internal.DiffInsert:
			ins = append(ins, i)
		default:
			flush()
			rows = append(rows, diffViewerRow{left: i, right: i})
		}
	}
	for i := 0; i < len(diff); {
		if diff[i].Op != internal.DiffEqual {
			line(i)
			i++
			continue
		}
		end := i
		for end < len(diff) && diff[end].Op == internal.DiffEqual {
			end++
		}
		keepHead, keepTail := ctx, ctx
		if i == 0 {
			keepHead = 0
		}
		if end == len(diff) {
			keepTail = 0
		}
		hidden := end - i - keepHead - keepTail
		if hidden <= 1 || expanded[i+keepHead] {
			for j := i; j < end; j++ {
				line(j)
			}
		} else {
			for j := i; j < i+keepHead; j++ {
				line(j)
			}
			flush()
			rows = append(rows, diffViewerRow{left: -1, right: -1, fold: i + keepHead, hidden: hidden})
			for j := end - keepTail; j < end; j++ {
				line(j)
			}
		}
		i = end
	}
	flush()
	return rows
}
