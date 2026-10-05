package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ThreadPanelOptions configure a ThreadPanel.
type ThreadPanelOptions struct {
	// Width is the panel's width; zero means 360.
	Width    float32
	Disabled bool
}

// ThreadPanelResult is the panel's element and this frame's events.
type ThreadPanelResult struct {
	Element *ui.Element
	// Sent is a reply the user sent; empty is none.
	Sent string
}

type threadState struct {
	list  MessageListState
	draft string
}

type threadKey struct{}

// ThreadPanel displays a root message, replies and an editor.
func ThreadPanel(c *ui.Context, open *bool, root ChatEntry, replies []ChatEntry, opts ThreadPanelOptions) ThreadPanelResult {
	var r ThreadPanelResult
	if opts.Width < 0 {
		panic("mujica: ThreadPanel Width must not be negative")
	}
	w := opts.Width
	if w == 0 {
		w = 360
	}
	ww, _ := c.Size()
	w = min(w, ww)
	k := Tokens(c)
	u := Density(c).Unit()
	panel := ui.Column(c).Width(w).FillHeight().Shrink(0).Background(k.Surface).BorderWidth(0, 0, 0, 1).
		BorderColor(k.Border).Role(ui.RoleGroup).Label(msg(c, "thread_panel.title")).Disabled(opts.Disabled)
	if panel.IsDisabled() {
		*open = false
	}
	if !*open {
		panel.Invisible()
		return r
	}
	overlayStyle(c, panel, overlayEnter(c, panel))
	panel.BorderWidth(0, 0, 0, 1).BorderColor(k.Border)
	r.Element = panel
	st := ui.Local(panel, threadKey{}, func() *threadState { return &threadState{} })
	s := *st
	panel.Children(func() {
		ui.Row(c).Padding(u*3, u*3, u*2, u*4).AlignItems(ui.Center).Children(func() {
			ui.Text(c, msg(c, "thread_panel.title")).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
			if iconAction(c, msgClose, msg(c, "thread_panel.close")).Clicked() {
				*open = false
			}
		})
		ui.Scroll(c).MaxHeight(160).Shrink(1).Padding(0, u*2).Children(func() { ChatMessage(c, root, ChatMessageOptions{}) })
		ui.Row(c).Gap(u*2).Padding(u*2, u*4).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf(msg(c, "thread_panel.replies"), len(replies))).FontSize(FontSize(c, theme.CaptionSize)).
				TextColor(k.TextMuted).SingleLine()
			ui.Box(c).Height(1).Grow(1).Background(k.Border)
		})
		MessageList(c, &s.list, len(replies), MessageListOptions{
			ID:    func(i int) string { return replies[i].ID },
			Label: func(i int) string { return replies[i].Author + ": " + replies[i].Text },
		}, func(i int) {
			ChatMessage(c, replies[i], ChatMessageOptions{Continued: i > 0 && ChatContinues(replies[i-1], replies[i])})
		})
		ui.Row(c).Gap(u*2).Padding(u*3).AlignItems(ui.Center).BorderWidth(1, 0, 0, 0).BorderColor(k.Border).Children(func() {
			in := TextInput(c, &s.draft, TextInputOptions{Placeholder: msg(c, "thread_panel.reply"), Label: msg(c, "thread_panel.reply"), Disabled: opts.Disabled})
			in.Element.Grow(1)
			send := IconButton(c, msgSend, msg(c, "thread_panel.send"), IconButtonOptions{Disabled: strings.TrimSpace(s.draft) == ""})
			if (in.Submitted() || send.Clicked()) && strings.TrimSpace(s.draft) != "" && !panel.IsDisabled() {
				r.Sent, s.draft = strings.TrimSpace(s.draft), ""
				s.list.ScrollToEnd()
			}
		})
	})
	if !panel.IsDisabled() && panel.Shortcut(0, ui.KeyEscape) {
		*open = false
	}
	switch {
	case r.Sent != "":
		slog.Debug("mujica: thread_panel send", "len", len(r.Sent))
	case !*open:
		slog.Debug("mujica: thread_panel close")
	}
	return r
}
