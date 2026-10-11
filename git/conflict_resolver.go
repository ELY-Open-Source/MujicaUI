package git

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ConflictResolverState holds conflict blocks and their choices.
type ConflictResolverState struct {
	text    string
	parts   []conflictPart
	choices []MergeChoice
	bad     int
	built   bool
}

// conflictPart is plain text or a conflict block.
type conflictPart struct {
	lines              []string
	ours, theirs       []string
	oursRef, theirsRef string
	line               int
	conflict           bool
}

// Resolved returns chosen blocks and all blocks.
func (s *ConflictResolverState) Resolved() (done, total int) {
	for i, p := range s.parts {
		if p.conflict {
			total++
			if s.choices[i] != MergeUnresolved {
				done++
			}
		}
	}
	return done, total
}

// Result applies the choices; unresolved blocks keep markers.
func (s *ConflictResolverState) Result() string {
	var out []string
	for i, p := range s.parts {
		if !p.conflict {
			out = append(out, p.lines...)
			continue
		}
		switch s.choices[i] {
		case MergeOurs:
			out = append(out, p.ours...)
		case MergeTheirs:
			out = append(out, p.theirs...)
		case MergeBoth:
			out = append(append(out, p.ours...), p.theirs...)
		default:
			out = append(out, p.lines...)
		}
	}
	return strings.Join(out, "\n")
}

// ConflictResolverOptions configure a ConflictResolver.
type ConflictResolverOptions struct {
	Disabled bool
}

// ConflictResolver resolves the git conflict blocks of text.
func ConflictResolver(c *ui.Context, s *ConflictResolverState, text string, opts ConflictResolverOptions) core.ChoiceResult {
	if !s.built || s.text != text {
		s.parts, s.bad = conflictParse(internal.SplitLines(text))
		s.choices = make([]MergeChoice, len(s.parts))
		s.text, s.built = text, true
		_, total := s.Resolved()
		slog.Debug("mujica: conflict resolver parse", "conflicts", total, "badLine", s.bad)
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 3).Disabled(opts.Disabled)
	res := core.ChoiceResult{Element: col}
	col.Children(func() {
		done, total := s.Resolved()
		ui.Row(c).AlignItems(ui.Center).Gap(u * 3).Children(func() {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "conflict_resolver.resolved"), done, total)).FontWeight(600).Role(ui.RoleStatus)
			if total > 0 {
				ui.Box(c).Height(4).Grow(1).Radius(2).Background(k.Border).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * float32(done) / float32(total), H: r.H}, k.Success, 2)
				})
			}
		})
		if s.bad > 0 {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
				ui.Text(c, fmt.Sprintf(core.Msg(c, "conflict_resolver.broken"), s.bad)).Role(ui.RoleStatus)
			})
		}
		if total == 0 && s.bad == 0 {
			ui.Text(c, core.Msg(c, "conflict_resolver.none")).TextColor(k.TextMuted)
		}
		n := 0
		for i, p := range s.parts {
			if !p.conflict {
				continue
			}
			n++
			if conflictBlock(c, s, i, n) && !col.IsDisabled() {
				res = core.NewChoiceResult(res.Element, true)
			}
		}
	})
	return res
}

// conflictBlock shows one conflict and its choice buttons.
func conflictBlock(c *ui.Context, s *ConflictResolverState, i, n int) bool {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	p := s.parts[i]
	choice := s.choices[i]
	changed := false
	card := ui.Column(c).AlignItems(ui.Stretch).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).Clip()
	card.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(u*2).Padding(u*2, u*3).Background(k.Background).Children(func() {
			icon, colr, key := "triangle-alert", k.Warning, "conflict_resolver.open"
			if choice != MergeUnresolved {
				icon, colr, key = "circle-check", k.Success, "conflict_resolver.done"
			}
			ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(colr)
			ui.Text(c, fmt.Sprintf(core.Msg(c, "conflict_resolver.title"), n, p.line)).Font(theme.SerifFont).FontWeight(600)
			ui.Text(c, core.Msg(c, key)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(colr)
		})
		ui.Grid(c).ColumnTracks(ui.Fr(1), ui.Fr(1)).Children(func() {
			side := func(title, ref string, lines []string, col ui.Color, on bool) {
				cell := ui.Column(c).Gap(u).Padding(u*2, u*3).MinWidth(0).Clip()
				if on {
					cell.Background(col.Mix(k.Surface, 0.88))
				}
				cell.Children(func() {
					ui.Row(c).Gap(u).AlignItems(ui.Center).Children(func() {
						if on {
							ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
						}
						ui.Text(c, title).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).TextColor(col)
						ui.Text(c, ref).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
					})
					mergeLines(c, lines, false)
				})
			}
			side(core.Msg(c, "conflict_resolver.current"), p.oursRef, p.ours, k.Info, choice == MergeOurs || choice == MergeBoth)
			side(core.Msg(c, "conflict_resolver.incoming"), p.theirsRef, p.theirs, k.Warning, choice == MergeTheirs || choice == MergeBoth)
		})
		ui.Row(c).Gap(u*2).Padding(u*2, u*3).Wrap().Children(func() {
			for _, opt := range []struct {
				choice MergeChoice
				key    string
			}{{MergeOurs, "take_current"}, {MergeTheirs, "take_incoming"}, {MergeBoth, "take_both"}, {MergeUnresolved, "reset"}} {
				if opt.choice == MergeUnresolved && choice == MergeUnresolved {
					continue
				}
				if mergeChoiceButton(c, core.Msg(c, "conflict_resolver."+opt.key), choice == opt.choice && opt.choice != MergeUnresolved).Clicked() && choice != opt.choice {
					s.choices[i], changed = opt.choice, true
					slog.Debug("mujica: conflict resolver choose", "block", n, "choice", opt.key)
				}
			}
		})
	})
	return changed
}

// conflictParse splits at markers; bad is an unclosed block's line.
func conflictParse(lines []string) (parts []conflictPart, bad int) {
	var plain []string
	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "<<<<<<<") {
			plain = append(plain, lines[i])
			continue
		}
		start := i
		p := conflictPart{conflict: true, line: i + 1, oursRef: strings.TrimSpace(lines[i][7:])}
		stage, j := 0, i+1
		for ; j < len(lines); j++ {
			l := lines[j]
			switch {
			case stage == 0 && strings.HasPrefix(l, "|||||||"):
				stage = 1
			case stage < 2 && l == "=======":
				stage = 2
			case stage == 2 && strings.HasPrefix(l, ">>>>>>>"):
				p.theirsRef = strings.TrimSpace(l[7:])
				stage = 3
			case stage == 0:
				p.ours = append(p.ours, l)
			case stage == 2:
				p.theirs = append(p.theirs, l)
			}
			if stage == 3 {
				break
			}
		}
		if stage != 3 {
			return []conflictPart{{lines: lines}}, start + 1
		}
		if len(plain) > 0 {
			parts = append(parts, conflictPart{lines: plain})
			plain = nil
		}
		p.lines = lines[start : j+1]
		parts = append(parts, p)
		i = j
	}
	if len(plain) > 0 {
		parts = append(parts, conflictPart{lines: plain})
	}
	return parts, 0
}
