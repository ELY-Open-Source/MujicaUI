package chat

import (
	"cmp"
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// SessionRole is who wrote a message; zero is SessionUser.
type SessionRole int

const (
	SessionUser SessionRole = iota
	SessionAssistant
	SessionSystem
)

// SessionMessage is one saved message, as plain text.
type SessionMessage struct {
	Role SessionRole
	// Author names the writer; empty uses the role's name.
	Author string
	Text   string
	Time   time.Time
}

// SharedConversationViewOptions configure a SharedConversationView.
type SharedConversationViewOptions struct {
	Title    string
	SharedBy string
	SharedAt time.Time
	Disabled bool
}

// SharedConversationView shows a shared conversation read-only with Copy link.
func SharedConversationView(c *ui.Context, messages []SessionMessage, opts SharedConversationViewOptions) core.ChoiceResult {
	if opts.Title == "" {
		panic("mujica: SharedConversationView needs a Title")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	col := ui.Column(c).Grow(1).MinHeight(0).Gap(d.Unit() * 3).Disabled(opts.Disabled)
	rows := *ui.Local(col, sharedConversationKey{}, func() *ui.ListState { return &ui.ListState{} })
	copied := false
	col.Children(func() {
		ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				ui.Text(c, opts.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H2Size)).TextColor(k.Text)
				by := fmt.Sprintf(core.Msg(c, "shared_conversation_view.shared_by"), cmp.Or(opts.SharedBy, core.Msg(c, "shared_conversation_view.someone")))
				if !opts.SharedAt.IsZero() {
					by += " · " + dateSeparatorLabel(c, core.DateOf(opts.SharedAt), core.DateOf(c.Now()))
				}
				ui.Text(c, by).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			})
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Must("lock")).Size(14, 14).TextColor(k.TextMuted)
				ui.Text(c, core.Msg(c, "shared_conversation_view.read_only")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			})
			copied = input.Button(c, core.Msg(c, "shared_conversation_view.copy_link"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("copy"),
				Disabled: opts.Disabled}).Clicked()
		})
		ui.Box(c).FillWidth().Height(9).Draw(func(p *ui.Painter, r ui.Rect) { internal.RuleDiamond(p, r, 6, 5, k.Ornament) })
		rows.Key = func(i int) any { return i }
		ui.List(c, rows, len(messages), func(i int) { sessionMessageView(c, k, messages[i]) }).Grow(1).MinHeight(0).Gap(d.Unit() * 3)
	})
	if copied {
		slog.Debug("mujica: shared_conversation_view copy_link")
	}
	return core.NewChoiceResult(col, copied)
}

type sharedConversationKey struct{}

// sessionMessageView shows one read-only message with author and time.
func sessionMessageView(c *ui.Context, k theme.Tokens, m SessionMessage) {
	d := core.Density(c)
	if m.Role < SessionUser || m.Role > SessionSystem {
		panic("mujica: SessionMessage Role out of range")
	}
	author := cmp.Or(m.Author, core.Msg(c, sessionRoleKeys[m.Role]))
	box := ui.Column(c).Gap(d.Unit()).Padding(d.Unit()*2, d.Unit()*3).Radius(theme.CardRadius)
	if m.Role == SessionUser {
		box.Background(k.Background).Border(theme.BorderWidth, k.Border)
	}
	box.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, author).Bold().TextColor(k.Text)
			if !m.Time.IsZero() {
				ui.Text(c, m.Time.Format("15:04")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			}
		})
		// Plain text: Markdown is not rendered here.
		display.Text(c, m.Text, display.TextOptions{Selectable: true})
	})
}

var sessionRoleKeys = map[SessionRole]string{
	SessionUser: "shared_conversation_view.user", SessionAssistant: "shared_conversation_view.assistant",
	SessionSystem: "shared_conversation_view.system",
}
