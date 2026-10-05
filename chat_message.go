package mujica

import (
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ChatEntry is one message of a person-to-person chat.
type ChatEntry struct {
	// ID is unique and stable; Author is required.
	ID, Author string
	Image      *ui.Bitmap
	Time       time.Time
	Text       string
	Edited     bool
	// Own marks the user's own messages, shown on the right.
	Own bool
}

// ChatMessageOptions configure a ChatMessage.
type ChatMessageOptions struct {
	// Continued hides repeated author details.
	Continued bool
	// Reply builds a quoted message above the text, as QuoteReply.
	Reply func()
	// Reactions builds the reaction row below the text.
	Reactions func()
}

// chatMergeGap is how close two messages must be to merge.
const chatMergeGap = 5 * time.Minute

// ChatContinues merges matching authors within five minutes.
func ChatContinues(prev, cur ChatEntry) bool {
	d := cur.Time.Sub(prev.Time)
	return prev.Author == cur.Author && prev.Own == cur.Own && d >= 0 && d < chatMergeGap
}

// ChatMessage displays text, timestamps, edits and interaction slots.
func ChatMessage(c *ui.Context, m ChatEntry, opts ChatMessageOptions) *ui.Element {
	if m.Author == "" {
		panic("mujica: ChatMessage needs an Author")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	role := MessageAssistant
	if m.Own {
		role = MessageUser
	}
	bo := MessageBubbleOptions{}
	const avatar = 32
	if !opts.Continued {
		bo.Name = m.Author
		bo.Avatar = func() { Avatar(c, m.Author, AvatarOptions{Image: m.Image, Size: avatar}) }
	} else {
		bo.Avatar = func() { ui.Box(c).Width(avatar).Shrink(0) }
	}
	e := MessageBubble(c, role, bo, func() {
		if opts.Reply != nil {
			opts.Reply()
		}
		ui.Text(c, m.Text).Selectable().FocusRing(false).LineHeight(theme.LineHeight)
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			if !m.Time.IsZero() {
				ui.Text(c, m.Time.Format(msg(c, "chat_message.time"))).FontSize(FontSize(c, theme.CaptionSize)).
					TextColor(k.TextMuted).FontFeatures("tnum").SingleLine()
			}
			if m.Edited {
				ui.Text(c, msg(c, "chat_message.edited")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).
					Italic().SingleLine()
			}
		})
		if opts.Reactions != nil {
			ui.Row(c).Gap(u).Wrap().Children(opts.Reactions)
		}
	})
	if opts.Continued {
		e.Label(m.Author)
		e.Margin(-u*2, 0, 0, 0)
	}
	return e
}
