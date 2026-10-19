package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/chat"
	"github.com/ELY-Open-Source/MujicaUI/display"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/messaging"
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
	messaging.ChatMessage(c, messaging.ChatEntry{Author: "Uika", Time: at, Text: "Rehearsal at nine.", Own: s.own}, messaging.ChatMessageOptions{})
	opts := messaging.ChatMessageOptions{Continued: s.continued}
	if s.reply {
		opts.Reply = func() {
			if chat.QuoteReply(c, "Sakiko", "When do we start?").Clicked() {
				log("jump to quoted message")
			}
		}
	}
	if s.reactions {
		opts.Reactions = func() {
			if input.Toggle(c, &s.reacted, "Like", input.ToggleOptions{}).Changed() {
				log(fmt.Sprintf("reaction %v", s.reacted))
			}
		}
	}
	messaging.ChatMessage(c, messaging.ChatEntry{Author: "Uika", Time: at.Add(time.Minute), Text: "Bring the masks.", Own: s.own, Edited: s.edited}, opts)
}

func unreadDividerDemo(c *ui.Context, log func(string)) {
	count := ui.Local(c.Root(), "403", func() float64 { return 3 })
	if input.NumberInput(c, count, input.NumberInputOptions{Min: 0, Max: 999, Step: 1, Label: "Unread count"}).Changed() {
		log(fmt.Sprintf("unread %.0f", *count))
	}
	messaging.UnreadDivider(c, int(*count))
	inputSection(c, "Zero count", func() { messaging.UnreadDivider(c, 0) })
}

func readReceiptDemo(c *ui.Context, log func(string)) {
	type state struct {
		receipt messaging.ReceiptState
		readers bool
		max     float64
	}
	s := ui.Local(c.Root(), "404", func() state { return state{receipt: messaging.ReceiptRead, readers: true, max: 3} })
	items := []input.SelectOption[messaging.ReceiptState]{{Value: messaging.ReceiptSent, Label: "Sent"}, {Value: messaging.ReceiptDelivered, Label: "Delivered"}, {Value: messaging.ReceiptRead, Label: "Read"}}
	if input.Select(c, &s.receipt, items, input.SelectOptions{Label: "Receipt"}).Changed() {
		log(fmt.Sprintf("receipt %d", s.receipt))
	}
	mediaCheck(c, log, &s.readers, "Reader avatars")
	if input.NumberInput(c, &s.max, input.NumberInputOptions{Min: 1, Max: 4, Step: 1, Label: "Max avatars"}).Changed() {
		log(fmt.Sprintf("max %.0f", s.max))
	}
	opts := messaging.ReadReceiptOptions{Max: int(s.max)}
	if s.readers {
		opts.Readers = []display.AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}, {Name: "Umiri"}, {Name: "Nyamu"}}
	}
	messaging.ReadReceipt(c, s.receipt, opts)
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
	if r := messaging.EmojiPicker(c, messaging.EmojiPickerOptions{Recent: &s.recent, Tone: &s.tone, Disabled: s.disabled}); r.Picked != "" {
		s.picked = r.Picked
		log("emoji " + r.Picked)
	}
	ui.Text(c, "Picked: "+s.picked)
}

func pinnedMessagesDemo(c *ui.Context, log func(string)) {
	type state struct {
		pins     []messaging.ChatEntry
		disabled bool
		jump     string
	}
	newPins := func() []messaging.ChatEntry {
		return []messaging.ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}, {ID: "b", Author: "Uika", Text: "Masks are in the hall."}}
	}
	s := ui.Local(c.Root(), "413", func() state { return state{pins: newPins()} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.Button(c, "Restore pins", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
		s.pins = newPins()
		log("pins restored")
	}
	r := messaging.PinnedMessages(c, s.pins, messaging.PinnedMessagesOptions{Disabled: s.disabled, Height: 180})
	if r.Jump != "" {
		s.jump = r.Jump
		log("jump " + r.Jump)
	}
	if r.Unpinned != "" {
		s.pins = slices.DeleteFunc(slices.Clone(s.pins), func(m messaging.ChatEntry) bool { return m.ID == r.Unpinned })
		log("unpin " + r.Unpinned)
	}
	ui.Text(c, "Jump target: "+s.jump)
}

func threadPanelDemo(c *ui.Context, log func(string)) {
	type state struct {
		open, disabled bool
		replies        []messaging.ChatEntry
	}
	s := ui.Local(c.Root(), "414", func() state {
		return state{open: true, replies: []messaging.ChatEntry{{ID: "1", Author: "Uika", Text: "Overture first."}}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.Button(c, "Open thread", input.ButtonOptions{Disabled: s.disabled}).Clicked() {
		s.open = true
		log("thread opened")
	}
	before := s.open
	ui.Row(c).Height(420).Children(func() {
		r := messaging.ThreadPanel(c, &s.open, messaging.ChatEntry{ID: "root", Author: "Sakiko", Text: "Set list for Friday?"}, s.replies, messaging.ThreadPanelOptions{Disabled: s.disabled})
		if r.Sent != "" {
			s.replies = append(s.replies, messaging.ChatEntry{ID: fmt.Sprint(len(s.replies) + 1), Author: "Me", Text: r.Sent, Own: true})
			log("reply sent")
		}
	})
	if before && !s.open {
		log("thread closed")
	}
}
