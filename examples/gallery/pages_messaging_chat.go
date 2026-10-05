package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	mujica "mujicaui"
)

func chatMessageDemo(c *ui.Context, log func(string)) {
	type state struct {
		continued, edited, own, reply, reactions bool
		reacted                                  bool
	}
	s := ui.Local(c.Root(), "402", func() state { return state{continued: true, edited: true, reply: true, reactions: true} })
	mediaCheck(c, log, &s.continued, "Merge consecutive messages")
	mediaCheck(c, log, &s.edited, "Edited")
	mediaCheck(c, log, &s.own, "Own message")
	mediaCheck(c, log, &s.reply, "Reply slot")
	mediaCheck(c, log, &s.reactions, "Reaction slot")
	at := time.Date(2026, 10, 5, 21, 7, 0, 0, time.Local)
	mujica.ChatMessage(c, mujica.ChatEntry{Author: "Uika", Time: at, Text: "Rehearsal at nine.", Own: s.own}, mujica.ChatMessageOptions{})
	opts := mujica.ChatMessageOptions{Continued: s.continued}
	if s.reply {
		opts.Reply = func() {
			if mujica.QuoteReply(c, "Sakiko", "When do we start?").Clicked() {
				log("jump to quoted message")
			}
		}
	}
	if s.reactions {
		opts.Reactions = func() {
			if mujica.Toggle(c, &s.reacted, "Like", mujica.ToggleOptions{}).Changed() {
				log(fmt.Sprintf("reaction %v", s.reacted))
			}
		}
	}
	mujica.ChatMessage(c, mujica.ChatEntry{Author: "Uika", Time: at.Add(time.Minute), Text: "Bring the masks.", Own: s.own, Edited: s.edited}, opts)
}

func unreadDividerDemo(c *ui.Context, log func(string)) {
	count := ui.Local(c.Root(), "403", func() float64 { return 3 })
	if mujica.NumberInput(c, count, mujica.NumberInputOptions{Min: 0, Max: 999, Step: 1, Label: "Unread count"}).Changed() {
		log(fmt.Sprintf("unread %.0f", *count))
	}
	mujica.UnreadDivider(c, int(*count))
	inputSection(c, "Zero count", func() { mujica.UnreadDivider(c, 0) })
}

func readReceiptDemo(c *ui.Context, log func(string)) {
	type state struct {
		receipt mujica.ReceiptState
		readers bool
		max     float64
	}
	s := ui.Local(c.Root(), "404", func() state { return state{receipt: mujica.ReceiptRead, readers: true, max: 3} })
	items := []mujica.SelectOption[mujica.ReceiptState]{{Value: mujica.ReceiptSent, Label: "Sent"}, {Value: mujica.ReceiptDelivered, Label: "Delivered"}, {Value: mujica.ReceiptRead, Label: "Read"}}
	if mujica.Select(c, &s.receipt, items, mujica.SelectOptions{Label: "Receipt"}).Changed() {
		log(fmt.Sprintf("receipt %d", s.receipt))
	}
	mediaCheck(c, log, &s.readers, "Reader avatars")
	if mujica.NumberInput(c, &s.max, mujica.NumberInputOptions{Min: 1, Max: 4, Step: 1, Label: "Max avatars"}).Changed() {
		log(fmt.Sprintf("max %.0f", s.max))
	}
	opts := mujica.ReadReceiptOptions{Max: int(s.max)}
	if s.readers {
		opts.Readers = []mujica.AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}, {Name: "Umiri"}, {Name: "Nyamu"}}
	}
	mujica.ReadReceipt(c, s.receipt, opts)
}

func emojiPickerDemo(c *ui.Context, log func(string)) {
	type state struct {
		recent   []string
		tone     int
		picked   string
		disabled bool
	}
	s := ui.Local(c.Root(), "412", func() state { return state{} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if r := mujica.EmojiPicker(c, mujica.EmojiPickerOptions{Recent: &s.recent, Tone: &s.tone, Disabled: s.disabled}); r.Picked != "" {
		s.picked = r.Picked
		log("emoji " + r.Picked)
	}
	ui.Text(c, "Picked: "+s.picked)
}

func pinnedMessagesDemo(c *ui.Context, log func(string)) {
	type state struct {
		pins     []mujica.ChatEntry
		disabled bool
		jump     string
	}
	newPins := func() []mujica.ChatEntry {
		return []mujica.ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}, {ID: "b", Author: "Uika", Text: "Masks are in the hall."}}
	}
	s := ui.Local(c.Root(), "413", func() state { return state{pins: newPins()} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if mujica.Button(c, "Restore pins", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
		s.pins = newPins()
		log("pins restored")
	}
	r := mujica.PinnedMessages(c, s.pins, mujica.PinnedMessagesOptions{Disabled: s.disabled, Height: 180})
	if r.Jump != "" {
		s.jump = r.Jump
		log("jump " + r.Jump)
	}
	if r.Unpinned != "" {
		s.pins = slices.DeleteFunc(slices.Clone(s.pins), func(m mujica.ChatEntry) bool { return m.ID == r.Unpinned })
		log("unpin " + r.Unpinned)
	}
	ui.Text(c, "Jump target: "+s.jump)
}

func threadPanelDemo(c *ui.Context, log func(string)) {
	type state struct {
		open, disabled bool
		replies        []mujica.ChatEntry
	}
	s := ui.Local(c.Root(), "414", func() state {
		return state{open: true, replies: []mujica.ChatEntry{{ID: "1", Author: "Uika", Text: "Overture first."}}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	if mujica.Button(c, "Open thread", mujica.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.open = true
		log("thread opened")
	}
	before := s.open
	ui.Row(c).Height(420).Children(func() {
		r := mujica.ThreadPanel(c, &s.open, mujica.ChatEntry{ID: "root", Author: "Sakiko", Text: "Set list for Friday?"}, s.replies, mujica.ThreadPanelOptions{Disabled: s.disabled})
		if r.Sent != "" {
			s.replies = append(s.replies, mujica.ChatEntry{ID: fmt.Sprint(len(s.replies) + 1), Author: "Me", Text: r.Sent, Own: true})
			log("reply sent")
		}
	})
	if before && !s.open {
		log("thread closed")
	}
}
