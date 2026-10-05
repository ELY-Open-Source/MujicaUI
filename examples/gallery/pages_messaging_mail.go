package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/messaging"
)

func mailListDemo(c *ui.Context, log func(string)) {
	type state struct {
		list            data.ListState[string]
		mails           []messaging.MailItem
		disabled, empty bool
		opened          string
	}
	s := ui.Local(c.Root(), "415", func() state {
		return state{mails: []messaging.MailItem{
			{ID: "1", From: "Sakiko", Subject: "Second act", Snippet: "The score is attached.", Unread: true, Attachments: 2},
			{ID: "2", From: "Uika", Subject: "Costumes", Snippet: "Fittings on Thursday.", Starred: true},
			{ID: "3", From: "Umiri", Subject: "Venue", Snippet: "The hall is booked."},
		}}
	})
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.empty, "Empty")
	if input.Button(c, "Load 100,000 mails", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
		s.mails = make([]messaging.MailItem, 100000)
		for i := range s.mails {
			s.mails[i] = messaging.MailItem{ID: fmt.Sprint(i), From: "Sender", Subject: fmt.Sprintf("Mail %d", i), Unread: i%2 == 0}
		}
		log("loaded 100,000 mails")
	}
	mails := s.mails
	if s.empty {
		mails = nil
	}
	r := messaging.MailList(c, &s.list, mails, messaging.MailListOptions{Disabled: s.disabled, Height: 240})
	if r.Changed() {
		log(fmt.Sprintf("selected %d mails", s.list.Selection.Len()))
	}
	if r.Opened != "" {
		s.opened = r.Opened
		for i := range s.mails {
			if s.mails[i].ID == r.Opened {
				s.mails[i].Unread = false
				break
			}
		}
		log("open " + r.Opened)
	}
	if r.Starred != "" {
		for i := range s.mails {
			if s.mails[i].ID == r.Starred {
				s.mails[i].Starred = !s.mails[i].Starred
				break
			}
		}
		log("star " + r.Starred)
	}
	ui.Text(c, "Opened: "+s.opened)
}

func mailReaderDemo(c *ui.Context, log func(string)) {
	type state struct{ markdown, attachments, disabled bool }
	s := ui.Local(c.Root(), "416", func() state { return state{attachments: true} })
	mediaCheck(c, log, &s.markdown, "Safe Markdown")
	mediaCheck(c, log, &s.attachments, "Attachments")
	mediaCheck(c, log, &s.disabled, "Disabled actions")
	m := messaging.Mail{From: "Sakiko Togawa", To: []string{"Uika Misumi"}, Cc: []string{"Umiri Yahata"}, Subject: "Second act", Date: time.Date(2026, 10, 5, 21, 7, 0, 0, time.Local), Body: "The score is attached.\n> Can you send it tonight?"}
	if s.markdown {
		m.Format = messaging.MailMarkdown
		m.Body = "**Rehearsal** at nine. [Venue](https://example.com)\n\n<script>alert(1)</script>"
	}
	if s.attachments {
		m.Attachments = []messaging.MailAttachment{{Name: "score.pdf", Size: 2048}}
	}
	ui.Column(c).Disabled(s.disabled).Children(func() {
		r := messaging.MailReader(c, m)
		if r.Link != "" {
			log("link " + r.Link)
		}
		if r.Opened >= 0 {
			log(fmt.Sprintf("open attachment %d", r.Opened))
		}
		if r.Downloaded >= 0 {
			log(fmt.Sprintf("download attachment %d", r.Downloaded))
		}
	})
}

func mailComposerDemo(c *ui.Context, log func(string)) {
	type state struct {
		draft    messaging.MailDraft
		disabled bool
		saved    string
	}
	s := ui.Local(c.Root(), "417", func() state { return state{draft: messaging.MailDraft{To: []string{"uika@example.com"}}} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	r := messaging.MailComposer(c, &s.draft, messaging.MailComposerOptions{
		Disabled: s.disabled, Contacts: messagingContacts,
		OnSaveDraft: func(d messaging.MailDraft) { s.saved = d.Subject; log("draft saved") },
		OnAttach: func() {
			s.draft.Attachments = append(s.draft.Attachments, messaging.MailAttachment{Name: fmt.Sprintf("score-%d.pdf", len(s.draft.Attachments)+1), Size: 2048})
			log("attachment added")
		},
	})
	if r.Sent {
		log("mail sent")
		s.draft = messaging.MailDraft{}
	}
	if r.Discarded {
		log("draft discarded")
		s.draft = messaging.MailDraft{}
	}
	ui.Text(c, "Saved subject: "+s.saved)
}

var messagingContacts = []messaging.Contact{{Name: "Uika Misumi", Email: "uika@example.com"}, {Name: "Umiri Yahata", Email: "umiri@example.com"}, {Name: "Sakiko Togawa", Email: "sakiko@example.com"}}

func recipientInputDemo(c *ui.Context, log func(string)) {
	type state struct {
		to                 []string
		disabled, contacts bool
	}
	s := ui.Local(c.Root(), "418", func() state { return state{contacts: true} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	mediaCheck(c, log, &s.contacts, "Contact suggestions")
	contacts := messagingContacts
	if !s.contacts {
		contacts = nil
	}
	if messaging.RecipientInput(c, &s.to, messaging.RecipientInputOptions{Contacts: contacts, Label: "Recipients", Disabled: s.disabled}).Changed() {
		log(fmt.Sprintf("recipients %d", len(s.to)))
	}
	if input.Button(c, "Clear recipients", input.ButtonOptions{Variant: input.Secondary, Disabled: s.disabled}).Clicked() {
		s.to = nil
		log("recipients cleared")
	}
	ui.Text(c, "Type a name or email; Enter or comma commits valid addresses.")
}

func quotedTextDemo(c *ui.Context, log func(string)) {
	type state struct {
		open, disabled bool
		text           string
	}
	s := ui.Local(c.Root(), "419", func() state { return state{text: "> Can you send the score?\n> Thanks."} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.TextArea(c, &s.text, input.TextAreaOptions{Label: "Quoted source", MinHeight: 80}).Changed() {
		log("quote source edited")
	}
	ui.Column(c).Disabled(s.disabled).Children(func() {
		if messaging.QuotedText(c, s.text, &s.open).Changed() {
			log(fmt.Sprintf("quote expanded %v", s.open))
		}
	})
}

func snoozePickerDemo(c *ui.Context, log func(string)) {
	type state struct {
		until    time.Time
		disabled bool
		morning  float64
	}
	s := ui.Local(c.Root(), "420", func() state { return state{morning: 9} })
	mediaCheck(c, log, &s.disabled, "Disabled")
	if input.NumberInput(c, &s.morning, input.NumberInputOptions{Min: 1, Max: 23, Step: 1, Label: "Morning hour"}).Changed() {
		log(fmt.Sprintf("morning %.0f", s.morning))
	}
	if messaging.SnoozePicker(c, &s.until, messaging.SnoozePickerOptions{Disabled: s.disabled, MorningHour: int(s.morning)}).Changed() {
		log("snooze " + s.until.Format(time.RFC3339))
	}
	if input.Button(c, "Clear snooze", input.ButtonOptions{Variant: input.Secondary, Disabled: s.disabled}).Clicked() {
		s.until = time.Time{}
		log("snooze cleared")
	}
	ui.Text(c, fmt.Sprint("Until: ", s.until))
}
