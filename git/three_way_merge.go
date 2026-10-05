package git

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MergeChoice is how a conflict is resolved.
type MergeChoice int

const (
	MergeUnresolved MergeChoice = iota
	// MergeOurs takes the local (current) side.
	MergeOurs
	// MergeTheirs takes the remote (incoming) side.
	MergeTheirs
	// MergeBoth takes ours, then theirs.
	MergeBoth
	// MergeBase keeps the common ancestor.
	MergeBase
)

// mergeChunk is a merge run; stable when all three agree.
type mergeChunk struct {
	base, ours, theirs []string
	conflict           bool
	// auto is the merged text of a non-conflicting change.
	auto   []string
	stable bool
}

// ThreeWayMergeState holds a merge's chunks and conflict choices.
type ThreeWayMergeState struct {
	list               ui.ListState
	base, ours, theirs string
	chunks             []mergeChunk
	choices            []MergeChoice
	built              bool
}

// Resolved returns chosen conflicts and all conflicts.
func (s *ThreeWayMergeState) Resolved() (done, total int) {
	for i, ch := range s.chunks {
		if ch.conflict {
			total++
			if s.choices[i] != MergeUnresolved {
				done++
			}
		}
	}
	return done, total
}

// Result is the merged text; unresolved conflicts keep conflict markers.
func (s *ThreeWayMergeState) Result() string {
	var out []string
	for i, ch := range s.chunks {
		out = append(out, mergeChunkText(ch, s.choices[i])...)
	}
	return strings.Join(out, "\n")
}

func mergeChunkText(ch mergeChunk, choice MergeChoice) []string {
	switch {
	case ch.stable:
		return ch.base
	case !ch.conflict:
		return ch.auto
	}
	switch choice {
	case MergeOurs:
		return ch.ours
	case MergeTheirs:
		return ch.theirs
	case MergeBoth:
		return slices.Concat(ch.ours, ch.theirs)
	case MergeBase:
		return ch.base
	}
	return slices.Concat([]string{"<<<<<<< ours"}, ch.ours, []string{"======="}, ch.theirs, []string{">>>>>>> theirs"})
}

// ThreeWayMergeOptions configure a ThreeWayMerge.
type ThreeWayMergeOptions struct {
	Disabled bool
}

// ThreeWayMerge shows four columns per chunk; give it a height.
func ThreeWayMerge(c *ui.Context, s *ThreeWayMergeState, base, ours, theirs string, opts ThreeWayMergeOptions) core.ChoiceResult {
	if !s.built || s.base != base || s.ours != ours || s.theirs != theirs {
		s.chunks = mergeChunks(internal.SplitLines(base), internal.SplitLines(ours), internal.SplitLines(theirs))
		s.choices = make([]MergeChoice, len(s.chunks))
		s.base, s.ours, s.theirs, s.built = base, ours, theirs, true
		_, total := s.Resolved()
		slog.Debug("mujica: three way merge chunks", "chunks", len(s.chunks), "conflicts", total)
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	col := ui.Column(c).AlignItems(ui.Stretch).Disabled(opts.Disabled).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius).Clip()
	res := core.ChoiceResult{Element: col}
	col.Children(func() {
		done, total := s.Resolved()
		ui.Row(c).AlignItems(ui.Center).Gap(u*3).Padding(u*2, u*3).Background(k.Background).Children(func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "three_way_merge.resolved"), done, total)).FontWeight(600).Role(ui.RoleStatus)
		})
		ui.Grid(c).ColumnTracks(ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1)).Background(k.Background).Children(func() {
			for _, key := range []string{"ours", "base", "theirs", "result"} {
				ui.Text(c, core.Msg(c, "three_way_merge."+key)).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize)).
					TextColor(k.Ornament.Mix(k.Text, 0.35)).Padding(u, u*3)
			}
		})
		ui.List(c, &s.list, len(s.chunks), func(i int) {
			if threeWayMergeChunk(c, s, i) && !col.IsDisabled() {
				res = core.NewChoiceResult(res.Element, true)
			}
		}).Grow(1)
	})
	return res
}

// threeWayMergeChunk builds chunk i and reports a new choice.
func threeWayMergeChunk(c *ui.Context, s *ThreeWayMergeState, i int) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	ch := s.chunks[i]
	changed := false
	ui.Column(c).AlignItems(ui.Stretch).Children(func() {
		grid := ui.Grid(c).ColumnTracks(ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1))
		grid.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: 1}, k.Border, 0) })
		grid.Children(func() {
			tint := func(lines []string, col ui.Color, on bool) {
				cell := ui.Column(c).Padding(u, u*3).MinWidth(0).Clip()
				if on {
					cell.Background(col.Mix(k.Surface, 0.88))
				}
				cell.Children(func() { mergeLines(c, lines, ch.stable) })
			}
			if ch.stable {
				tint(ch.base, k.Text, false)
				tint(ch.base, k.Text, false)
				tint(ch.base, k.Text, false)
				tint(ch.base, k.Text, false)
				return
			}
			choice := s.choices[i]
			tint(ch.ours, k.Info, ch.conflict && (choice == MergeOurs || choice == MergeBoth) || !ch.conflict && slices.Equal(ch.auto, ch.ours))
			tint(ch.base, k.TextMuted, choice == MergeBase)
			tint(ch.theirs, k.Warning, ch.conflict && (choice == MergeTheirs || choice == MergeBoth) || !ch.conflict && slices.Equal(ch.auto, ch.theirs))
			result := ui.Column(c).Padding(u, u*3).MinWidth(0).Clip()
			result.Children(func() {
				if ch.conflict && choice == MergeUnresolved {
					ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
						ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14)
						ui.Text(c, core.Msg(c, "three_way_merge.unresolved")).FontSize(core.FontSize(c, theme.CaptionSize))
					})
					return
				}
				mergeLines(c, mergeChunkText(ch, choice), false)
			})
		})
		if !ch.conflict {
			return
		}
		ui.Row(c).Gap(u*2).Padding(u, u*3, u*2, u*3).Wrap().Children(func() {
			for _, opt := range []struct {
				choice MergeChoice
				key    string
			}{{MergeOurs, "take_ours"}, {MergeTheirs, "take_theirs"}, {MergeBoth, "take_both"}, {MergeBase, "take_base"}} {
				if mergeChoiceButton(c, core.Msg(c, "three_way_merge."+opt.key), s.choices[i] == opt.choice).Clicked() && s.choices[i] != opt.choice {
					s.choices[i], changed = opt.choice, true
					slog.Debug("mujica: three way merge choose", "chunk", i, "choice", opt.key)
				}
			}
		})
	})
	return changed
}

// mergeChoiceButton is a toggle marking the chosen side.
func mergeChoiceButton(c *ui.Context, label string, on bool) *ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).Padding(0, u*3).Gap(u*1.5).Radius(theme.ControlRadius).
		Border(theme.BorderWidth, k.ControlBorder).Label(label).TextColor(k.Text)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case on:
		b.Background(k.Selection).Border(theme.BorderWidth, k.AccentText).FontWeight(600)
	case b.Pressed():
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	b.Transition(core.ColorFade)
	b.Children(func() {
		if on {
			ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
		}
		ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
	})
	return b
}

// mergeLines shows code lines, folding long stable runs.
func mergeLines(c *ui.Context, lines []string, fold bool) {
	k := core.Tokens(c)
	show := func(ls []string) {
		for _, l := range ls {
			if l == "" {
				l = " "
			}
			t := ui.Text(c, l).Font(theme.MonoFont).NoWrap().FontSize(core.FontSize(c, theme.CaptionSize))
			if fold {
				t.TextColor(k.TextMuted)
			}
		}
	}
	if len(lines) == 0 {
		ui.Text(c, "∅").TextColor(k.TextMuted).Label(core.Msg(c, "three_way_merge.empty"))
		return
	}
	if !fold || len(lines) <= 5 {
		show(lines)
		return
	}
	show(lines[:2])
	ui.Text(c, fmt.Sprintf(core.Msg(c, "diff_viewer.fold"), len(lines)-4)).Font(theme.SerifFont).
		FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
	show(lines[len(lines)-2:])
}

// mergeHunk replaces base lines [from, to) with lines.
type mergeHunk struct {
	from, to int
	lines    []string
}

// mergeHunks lists the changes a diff makes to base.
func mergeHunks(d []internal.DiffLine) []mergeHunk {
	var out []mergeHunk
	at, open := 0, -1
	for _, l := range d {
		if l.Op == internal.DiffEqual {
			at = l.Old
			open = -1
			continue
		}
		if open < 0 {
			out = append(out, mergeHunk{from: at, to: at})
			open = len(out) - 1
		}
		h := &out[open]
		if l.Op == internal.DiffDelete {
			h.to = l.Old
		} else {
			h.lines = append(h.lines, l.Text)
		}
	}
	return out
}

// mergeChunks splits a merge into stable, clean and conflicting runs.
func mergeChunks(base, ours, theirs []string) []mergeChunk {
	type tagged struct {
		mergeHunk
		side int
	}
	var all []tagged
	for side, other := range [][]string{ours, theirs} {
		for _, h := range mergeHunks(internal.LineDiff(base, other)) {
			all = append(all, tagged{h, side})
		}
	}
	slices.SortStableFunc(all, func(a, b tagged) int { return a.from - b.from })
	var chunks []mergeChunk
	at := 0
	for i := 0; i < len(all); {
		from, to := all[i].from, all[i].to
		j := i + 1
		for j < len(all) && (all[j].from < to || all[j].from == to && (all[j].from == all[j].to || from == to)) {
			to = max(to, all[j].to)
			j++
		}
		if from > at {
			chunks = append(chunks, mergeChunk{base: base[at:from], stable: true})
		}
		sides := [2][]string{}
		touched := [2]bool{}
		for s := range 2 {
			cur := from
			for _, h := range all[i:j] {
				if h.side != s {
					continue
				}
				touched[s] = true
				sides[s] = append(sides[s], base[cur:h.from]...)
				sides[s] = append(sides[s], h.lines...)
				cur = h.to
			}
			sides[s] = append(sides[s], base[cur:to]...)
		}
		ch := mergeChunk{base: base[from:to], ours: sides[0], theirs: sides[1]}
		switch {
		case !touched[1]:
			ch.auto = ch.ours
		case !touched[0], slices.Equal(ch.ours, ch.theirs):
			ch.auto = ch.theirs
		default:
			ch.conflict = true
		}
		chunks = append(chunks, ch)
		at, i = to, j
	}
	if at < len(base) {
		chunks = append(chunks, mergeChunk{base: base[at:], stable: true})
	}
	return chunks
}
