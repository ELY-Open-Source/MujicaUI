package git

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// CommitMessage is a commit's title and body.
type CommitMessage struct {
	Title string
	Body  string
}

// CommitInputOptions configure a CommitInput.
type CommitInputOptions struct {
	// TitleLimit warns past this many title characters; zero means 50.
	TitleLimit int
	// Busy shows a commit in progress and blocks another.
	Busy     bool
	Disabled bool
}

// CommitInputResult is a built CommitInput and its events.
type CommitInputResult struct {
	Element          *ui.Element
	committed, amend bool
}

// Committed reports Commit, or Cmd+Enter, this frame.
func (r CommitInputResult) Committed() bool { return r.committed }

// Amended reports Amend this frame.
func (r CommitInputResult) Amended() bool { return r.amend }

type commitInputState struct {
	// pending blocks commits until Busy comes and goes.
	pending, sawBusy bool
}

// CommitInput edits a commit message and commits or amends once.
func CommitInput(c *ui.Context, m *CommitMessage, opts CommitInputOptions) CommitInputResult {
	limit := opts.TitleLimit
	if limit == 0 {
		limit = 50
	}
	if limit < 0 {
		panic("mujica: CommitInput TitleLimit is negative")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	col := ui.Column(c).AlignItems(ui.Stretch).Gap(u * 2).Disabled(opts.Disabled)
	res := CommitInputResult{Element: col}
	st := ui.Local(col, "state", func() *commitInputState { return &commitInputState{} })
	s := *st
	if opts.Busy {
		s.sawBusy = true
	} else if s.sawBusy {
		s.pending, s.sawBusy = false, false
	}
	blocked := opts.Busy || s.pending
	n := utf8.RuneCountInString(m.Title)
	col.Children(func() {
		title := input.TextInput(c, &m.Title, input.TextInputOptions{Placeholder: core.Msg(c, "commit_input.title"), Label: core.Msg(c, "commit_input.title"),
			Disabled: opts.Disabled, ReadOnly: blocked,
			Trailing: func() {
				col := k.TextMuted
				if n > limit {
					col = k.Warning
				}
				name := fmt.Sprintf(core.Msg(c, "commit_input.count"), n, limit)
				if n > limit {
					name = fmt.Sprintf(core.Msg(c, "commit_input.long"), n-limit)
				}
				ui.Text(c, fmt.Sprintf("%d/%d", n, limit)).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").TextColor(col).Shrink(0).
					Label(name).Tooltip(name)
			}})
		body := input.TextArea(c, &m.Body, input.TextAreaOptions{Placeholder: core.Msg(c, "commit_input.body"), Label: core.Msg(c, "commit_input.body"),
			MinHeight: 72, MaxHeight: 200, Disabled: opts.Disabled, ReadOnly: blocked})
		empty := strings.TrimSpace(m.Title) == ""
		send := (title.Submitted() || body.Submitted()) && !empty
		ui.Row(c).Gap(u * 2).Justify(ui.End).Children(func() {
			if input.Button(c, core.Msg(c, "commit_input.amend"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("rotate-ccw"), Disabled: blocked || empty}).Clicked() {
				res.amend = true
			}
			if input.Button(c, core.Msg(c, "commit_input.commit"), input.ButtonOptions{Icon: icons.Must("check"), Loading: opts.Busy, Disabled: s.pending || empty}).Clicked() {
				send = true
			}
		})
		res.committed = send && !res.amend
	})
	if blocked || col.IsDisabled() {
		res.committed, res.amend = false, false
	}
	if res.committed || res.amend {
		s.pending = true
		slog.Debug("mujica: commit input submit", "amend", res.amend)
	}
	return res
}
