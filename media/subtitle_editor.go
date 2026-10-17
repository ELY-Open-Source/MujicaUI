package media

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SubtitleCue is one subtitle: shown from Start until End.
type SubtitleCue struct {
	Start, End time.Duration
	Text       string
}

// SubtitleIssueKind is what is wrong with a cue.
type SubtitleIssueKind int

const (
	// SubtitleEndBeforeStart: End is not after Start.
	SubtitleEndBeforeStart SubtitleIssueKind = iota
	// SubtitleOverlap: Start is before the previous cue's End.
	SubtitleOverlap
)

// SubtitleIssue is a problem with cue Index.
type SubtitleIssue struct {
	Index int
	Kind  SubtitleIssueKind
}

// CheckSubtitles lists misordered cues and cues overlapping the previous.
func CheckSubtitles(cues []SubtitleCue) []SubtitleIssue {
	var out []SubtitleIssue
	for i, q := range cues {
		if q.End <= q.Start {
			out = append(out, SubtitleIssue{i, SubtitleEndBeforeStart})
		}
		if i > 0 && q.Start < cues[i-1].End {
			out = append(out, SubtitleIssue{i, SubtitleOverlap})
		}
	}
	return out
}

// SubtitleEditorOptions configure a SubtitleEditor.
type SubtitleEditorOptions struct {
	// Height is the visible list's height; zero means 360.
	Height   float32
	Disabled bool
}

type subtitleState struct {
	scroll ui.ScrollState
	drafts map[subtitleField]string
}

type subtitleField struct{ row, part int }

type subtitleKey struct{}

// SubtitleEditor edits cue times and text, building visible rows.
func SubtitleEditor(c *ui.Context, cues *[]SubtitleCue, opts SubtitleEditorOptions) core.ChoiceResult {
	if opts.Height < 0 {
		panic("mujica: SubtitleEditor Height must not be negative")
	}
	h := opts.Height
	if h == 0 {
		h = 360
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	rowH := core.Density(c).ControlHeight() + core.FontSize(c, theme.CaptionSize)*theme.LineHeight + u*3
	col := ui.Column(c).Gap(u * 2).Disabled(opts.Disabled)
	st := ui.Local(col, subtitleKey{}, func() *subtitleState { return &subtitleState{drafts: map[subtitleField]string{}} })
	s := *st
	issues := map[int]string{}
	for _, is := range CheckSubtitles(*cues) {
		if is.Kind == SubtitleOverlap {
			issues[is.Index] = fmt.Sprintf(core.Msg(c, "subtitle_editor.overlap"), is.Index)
		} else {
			issues[is.Index] = core.Msg(c, "subtitle_editor.order")
		}
	}
	changed := false
	col.Children(func() {
		list := ui.Scroll(c).Height(h).TrackScroll(&s.scroll).Background(k.Surface).
			Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
		list.Children(func() {
			n := len(*cues)
			// The container keeps Y inside the content; match it.
			y := max(0, min(s.scroll.Y, float32(n)*rowH-h))
			first := max(0, int(y/rowH)-1)
			last := min(n, int((y+h)/rowH)+2)
			ui.Box(c).Height(float32(n) * rowH).Children(func() {
				for i := first; i < last; i++ {
					before := len(*cues)
					if subtitleRow(c, k, cues, s, i, rowH, issues[i]) {
						changed = true
					}
					if len(*cues) != before {
						// A removal shifts rows below; the next frame builds them.
						c.AnimationFrame()
						return
					}
				}
			})
			if n == 0 {
				ui.Text(c, core.Msg(c, "subtitle_editor.empty")).TextColor(k.TextMuted).Padding(u * 3)
			}
		})
		if input.Button(c, core.Msg(c, "subtitle_editor.add"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("plus")}).Clicked() && !col.IsDisabled() {
			var at time.Duration
			if n := len(*cues); n > 0 {
				at = (*cues)[n-1].End
			}
			*cues = append(*cues, SubtitleCue{Start: at, End: at + 2*time.Second})
			s.scroll.Y = float32(len(*cues)) * rowH
			changed = true
		}
	})
	if changed {
		slog.Debug("mujica: subtitle_editor changed", "cues", len(*cues))
	}
	return core.NewChoiceResult(col, changed)
}

// subtitleRow edits cue i; true means changed or removed.
func subtitleRow(c *ui.Context, k theme.Tokens, cues *[]SubtitleCue, s *subtitleState, i int, rowH float32, issue string) bool {
	u := core.Density(c).Unit()
	q := &(*cues)[i]
	changed, removed := false, false
	row := ui.Column(c.Key(i)).Absolute().Left(0).Right(0).Top(float32(i)*rowH).Height(rowH).Padding(u*1.5, u*2).Gap(2).
		BorderWidth(0, 0, 1, 0).BorderColor(k.Border)
	row.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Textf(c, "%d", i+1).Width(28).TextColor(k.TextMuted).FontFeatures("tnum")
			for part, t := range []*time.Duration{&q.Start, &q.End} {
				f := subtitleField{i, part}
				text, editing := s.drafts[f]
				if !editing {
					text = subtitleClock(*t)
				}
				name := core.Msg(c, []string{"subtitle_editor.start", "subtitle_editor.end"}[part])
				in := input.TextInput(c, &text, input.TextInputOptions{Label: fmt.Sprintf("%s %d", name, i+1)})
				in.Element.Width(112)
				bad := false
				if in.Changed() {
					if v, ok := subtitleParse(text); ok {
						*t, changed = v, true
						delete(s.drafts, f)
					} else {
						s.drafts[f], bad = text, true
					}
				} else if editing {
					_, ok := subtitleParse(text)
					bad = !ok
				}
				if bad {
					in.Input.Error(core.Msg(c, "subtitle_editor.badTime"))
				} else if issue != "" {
					in.Input.Error(issue)
				}
			}
			txt := input.TextInput(c, &q.Text, input.TextInputOptions{Label: fmt.Sprintf("%s %d", core.Msg(c, "subtitle_editor.text"), i+1)})
			txt.Element.Grow(1)
			if txt.Changed() {
				changed = true
			}
			if b := core.IconAction(c, icons.Must("trash-2"), fmt.Sprintf("%s %d", core.Msg(c, "subtitle_editor.remove"), i+1)); b.Clicked() && !b.IsDisabled() {
				removed = true
			}
		})
		ui.Row(c).Gap(u).PaddingX(28 + u*2).Children(func() {
			note := issue
			for part := range 2 {
				if d, ok := s.drafts[subtitleField{i, part}]; ok {
					if _, valid := subtitleParse(d); !valid {
						note = core.Msg(c, "subtitle_editor.badTime")
					}
				}
			}
			if note != "" {
				ui.Text(c, note).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Danger).SingleLine()
			}
		})
	})
	if removed {
		*cues = slices.Delete(*cues, i, i+1)
		clear(s.drafts)
		return true
	}
	return changed
}

// subtitleClock writes d as m:ss.mmm.
func subtitleClock(d time.Duration) string {
	ms := d.Milliseconds()
	return fmt.Sprintf("%d:%02d.%03d", ms/60000, ms/1000%60, ms%1000)
}

// subtitleParse reads m:ss.mmm or h:mm:ss.mmm; the fraction is optional.
func subtitleParse(s string) (time.Duration, bool) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}
	acc := 0
	for i, p := range parts[:len(parts)-1] {
		n, err := strconv.Atoi(p)
		if err != nil || !subtitleDigits(p) || (i > 0 && n > 59) {
			return 0, false
		}
		acc = acc*60 + n
	}
	sec := parts[len(parts)-1]
	whole, frac, _ := strings.Cut(sec, ".")
	s0, err := strconv.Atoi(whole)
	if err != nil || len(whole) != 2 || s0 < 0 || s0 > 59 || len(frac) > 3 || !subtitleDigits(whole+frac) {
		return 0, false
	}
	ms := 0
	if frac != "" {
		ms, _ = strconv.Atoi(frac + strings.Repeat("0", 3-len(frac)))
	}
	return time.Duration(acc*60+s0)*time.Second + time.Duration(ms)*time.Millisecond, true
}

// subtitleDigits reports whether s is ASCII digits only.
func subtitleDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
