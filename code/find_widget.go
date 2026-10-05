package code

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// FindState is a FindWidget's query, options and current match.
type FindState struct {
	Query, Replacement string
	// MatchCase, WholeWord and Regex refine the query.
	MatchCase, WholeWord, Regex bool
	// ShowReplace opens the replace row.
	ShowReplace bool
	// Current is the chosen match, from 0.
	Current int
	matches [][]int
	err     string
	key     findKey
}

type findKey struct {
	text, query            string
	matchCase, word, regex bool
	ran                    bool
}

// Matches returns the byte ranges of matches.
func (s *FindState) Matches() [][2]int {
	out := make([][2]int, len(s.matches))
	for i, m := range s.matches {
		out[i] = [2]int{m[0], m[1]}
	}
	return out
}

// FindWidgetOptions configure a FindWidget.
type FindWidgetOptions struct {
	// ReadOnly hides replacing.
	ReadOnly bool
	Disabled bool
}

// FindWidgetResult is a FindWidget and what happened in this frame.
type FindWidgetResult struct {
	Element *ui.Element
	// Moved reports navigation; Replaced counts replacements.
	Moved    bool
	Replaced int
}

// FindWidget finds and replaces in *text.
func FindWidget(c *ui.Context, text *string, s *FindState, opts FindWidgetOptions) FindWidgetResult {
	if text == nil || s == nil {
		panic("mujica: FindWidget needs text and a state")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var res FindWidgetResult
	res.Element = CodePanel(c).Padding(u * 2).Gap(u * 2).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	findRun(*text, s)
	res.Element.Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			if !opts.ReadOnly {
				chev := "chevron-right"
				if s.ShowReplace {
					chev = "chevron-down"
				}
				if core.IconAction(c, icons.Must(chev), core.Msg(c, "find_widget.toggleReplace")).Clicked() && !disabled {
					s.ShowReplace = !s.ShowReplace
				}
			}
			r := input.SearchInput(c, &s.Query, input.SearchInputOptions{Find: true, Matches: len(s.matches), Current: &s.Current, Disabled: opts.Disabled})
			r.Element.Grow(1).MinWidth(0)
			changed := r.Changed() || r.Moved()
			for _, o := range []struct {
				on    *bool
				short string
				key   string
			}{{&s.MatchCase, "Aa", "find_widget.case"}, {&s.WholeWord, "ab", "find_widget.word"}, {&s.Regex, ".*", "find_widget.regex"}} {
				if input.Toggle(c, o.on, o.short, input.ToggleOptions{Label: core.Msg(c, o.key), Disabled: opts.Disabled}).Tooltip(core.Msg(c, o.key)).Changed() {
					changed = true
				}
			}
			if changed {
				findRun(*text, s)
				res.Moved = true
				slog.Debug("mujica: find_widget move", "matches", len(s.matches), "current", s.Current)
			}
		})
		if s.err != "" {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "find_widget.invalid"), s.err)).TextColor(k.Danger).Role(ui.RoleStatus).
				FontSize(core.FontSize(c, 12))
		}
		if !s.ShowReplace || opts.ReadOnly {
			return
		}
		ui.Row(c).Gap(u*2).AlignItems(ui.Center).Padding(0, 0, 0, 28+u*2).Children(func() {
			in := input.TextInput(c, &s.Replacement, input.TextInputOptions{Placeholder: core.Msg(c, "find_widget.replace"), Label: core.Msg(c, "find_widget.replace"), Disabled: opts.Disabled})
			in.Element.Grow(1).MinWidth(0)
			none := len(s.matches) == 0
			if (input.Button(c, core.Msg(c, "find_widget.replaceOne"), input.ButtonOptions{Variant: input.Secondary, Disabled: none || opts.Disabled}).Clicked() || in.Submitted()) && !none && !disabled {
				res.Replaced = findReplace(text, s, false)
			}
			if input.Button(c, core.Msg(c, "find_widget.replaceAll"), input.ButtonOptions{Variant: input.Secondary, Disabled: none || opts.Disabled}).Clicked() && !disabled {
				res.Replaced = findReplace(text, s, true)
			}
		})
	})
	return res
}

// findPattern compiles the query with its options.
func findPattern(s *FindState) (*regexp.Regexp, error) {
	q := s.Query
	if !s.Regex {
		q = regexp.QuoteMeta(q)
	}
	if s.WholeWord {
		q = `\b(?:` + q + `)\b`
	}
	if !s.MatchCase {
		q = "(?i)" + q
	}
	return regexp.Compile(q)
}

// findRun refinds matches when inputs change.
func findRun(text string, s *FindState) {
	key := findKey{text, s.Query, s.MatchCase, s.WholeWord, s.Regex, true}
	if key == s.key {
		return
	}
	s.key, s.err, s.matches = key, "", nil
	if s.Query != "" {
		re, err := findPattern(s)
		if err != nil {
			s.err = err.Error()
		} else {
			for _, m := range re.FindAllStringSubmatchIndex(text, -1) {
				if m[1] > m[0] {
					s.matches = append(s.matches, m)
				}
			}
		}
	}
	s.Current = max(0, min(s.Current, len(s.matches)-1))
}

// findReplace replaces one or all matches.
func findReplace(text *string, s *FindState, all bool) int {
	re, err := findPattern(s)
	if err != nil {
		panic("mujica: FindWidget replaced with an invalid pattern")
	}
	targets := s.matches
	if !all {
		targets = s.matches[s.Current : s.Current+1]
	}
	out := []byte{}
	last := 0
	for _, m := range targets {
		out = append(out, (*text)[last:m[0]]...)
		if s.Regex {
			out = re.ExpandString(out, s.Replacement, *text, m)
		} else {
			out = append(out, s.Replacement...)
		}
		last = m[1]
	}
	*text = string(append(out, (*text)[last:]...))
	findRun(*text, s)
	slog.Debug("mujica: find_widget replace", "count", len(targets))
	return len(targets)
}
