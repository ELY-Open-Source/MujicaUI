package mujica

import (
	"fmt"
	"log/slog"
	"net/mail"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// MailDraft is a mail being written.
type MailDraft struct {
	To, Cc, Bcc   []string
	Subject, Body string
	Attachments   []MailAttachment
}

// MailComposerOptions configure a MailComposer.
type MailComposerOptions struct {
	Contacts []Contact
	// OnSaveDraft saves rested edits or manual requests; required.
	OnSaveDraft func(MailDraft)
	// OnAttach asks the caller to pick files; nil hides Attach.
	OnAttach func()
	Disabled bool
}

// MailComposerResult is the composer's element and this frame's events.
type MailComposerResult struct {
	Element *ui.Element
	// Sent requests sending; Discarded requests discarding the draft.
	Sent, Discarded bool
}

const mailSaveDelay = 2 * time.Second

type composerState struct {
	cc, bcc bool
	dirty   bool
	edited  time.Time
	saved   time.Time
	last    MailDraft
	started bool
}

type composerKey struct{}

// MailComposer edits mail fields and reports draft save requests.
func MailComposer(c *ui.Context, d *MailDraft, opts MailComposerOptions) MailComposerResult {
	if opts.OnSaveDraft == nil {
		panic("mujica: MailComposer needs OnSaveDraft")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	var r MailComposerResult
	col := ui.Column(c).Gap(u * 3).Role(ui.RoleGroup).Label(msg(c, "mail_composer.label")).Disabled(opts.Disabled)
	r.Element = col
	st := ui.Local(col, composerKey{}, func() *composerState { return &composerState{} })
	s := *st
	live := !col.IsDisabled()
	if !live {
		copy := composerCopy(*d)
		d = &copy
	}
	if !s.started {
		s.started, s.last = true, composerCopy(*d)
		s.cc, s.bcc = len(d.Cc) > 0, len(d.Bcc) > 0
	}
	s.cc = s.cc || len(d.Cc) > 0
	s.bcc = s.bcc || len(d.Bcc) > 0
	col.Children(func() {
		field := func(key string, list *[]string, extra func()) {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
				label := ui.Text(c, msg(c, key)).Width(56).TextColor(k.TextMuted).Margin(u*2, 0, 0, 0).SingleLine()
				parts := RecipientInput(c, list, RecipientInputOptions{Contacts: opts.Contacts, Label: msg(c, key), Disabled: opts.Disabled})
				parts.Element.Grow(1).MinWidth(0)
				if label.Clicked() && !parts.Input.IsDisabled() {
					parts.Input.Focus()
				}
				if extra != nil {
					extra()
				}
			})
		}
		field("mail_composer.to", &d.To, func() {
			if !s.cc && Button(c, msg(c, "mail_composer.cc"), ButtonOptions{Variant: Ghost}).Clicked() {
				s.cc = true
			}
			if !s.bcc && Button(c, msg(c, "mail_composer.bcc"), ButtonOptions{Variant: Ghost}).Clicked() {
				s.bcc = true
			}
		})
		if s.cc {
			field("mail_composer.cc", &d.Cc, nil)
		}
		if s.bcc {
			field("mail_composer.bcc", &d.Bcc, nil)
		}
		TextInput(c, &d.Subject, TextInputOptions{Placeholder: msg(c, "mail_composer.subject"), Label: msg(c, "mail_composer.subject"), Disabled: opts.Disabled})
		if TextArea(c, &d.Body, TextAreaOptions{Label: msg(c, "mail_composer.body"), MinHeight: 160, Disabled: opts.Disabled}).Submitted() &&
			composerReady(*d) {
			r.Sent = true
		}
		if len(d.Attachments) > 0 {
			composerAttachments(c, k, d, opts.Disabled)
		}
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ready := composerReady(*d)
			if Button(c, msg(c, "mail_composer.send"), ButtonOptions{Icon: msgSend, Disabled: !ready}).Clicked() {
				r.Sent = true
			}
			if opts.OnAttach != nil && Button(c, msg(c, "mail_composer.attach"), ButtonOptions{Variant: Secondary, Icon: msgPaperclip}).Clicked() {
				slog.Debug("mujica: mail_composer attach")
				opts.OnAttach()
			}
			ui.Box(c).Grow(1)
			composerStatus(c, k, s)
			if Button(c, msg(c, "mail_composer.discard"), ButtonOptions{Variant: Ghost, Icon: icons.Must("trash-2")}).Clicked() {
				r.Discarded = true
			}
		})
	})
	if live && col.Shortcut(ui.Cmd, ui.KeyEnter) && composerReady(*d) {
		r.Sent = true
	}
	if live {
		composerAutosave(c, s, *d, opts.OnSaveDraft, col.Shortcut(ui.Cmd, ui.KeyS))
	}
	if !live {
		r.Sent, r.Discarded = false, false
	}
	if r.Sent || r.Discarded {
		slog.Debug("mujica: mail_composer done", "sent", r.Sent)
	}
	return r
}

func composerAutosave(c *ui.Context, s *composerState, d MailDraft, save func(MailDraft), now bool) {
	if !composerSame(s.last, d) {
		s.last, s.dirty, s.edited = composerCopy(d), true, c.Now()
	}
	if !s.dirty && !now {
		return
	}
	if wait := mailSaveDelay - c.Now().Sub(s.edited); wait > 0 && !now {
		c.After(wait)
		return
	}
	s.dirty, s.saved = false, c.Now()
	slog.Debug("mujica: mail_composer save draft")
	save(composerCopy(d))
}

// composerStatus says whether the draft is saved.
func composerStatus(c *ui.Context, k theme.Tokens, s *composerState) {
	text := ""
	switch {
	case s.dirty:
		text = msg(c, "mail_composer.unsaved")
	case !s.saved.IsZero():
		text = fmt.Sprintf(msg(c, "mail_composer.saved"), s.saved.Format("15:04"))
	}
	if text != "" {
		ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleStatus)
	}
}

// composerAttachments lists attached files with remove buttons.
func composerAttachments(c *ui.Context, k theme.Tokens, d *MailDraft, disabled bool) {
	ui.Row(c).Gap(Density(c).Unit() * 2).Wrap().Children(func() {
		for i, a := range d.Attachments {
			var v TagView
			ui.Box(c).Key(fmt.Sprint(i, a.Name)).Children(func() {
				v = Tag(c, a.Name, TagOptions{Icon: "file", Closable: true, Disabled: disabled})
			})
			if v.Closed {
				d.Attachments = slices.Delete(slices.Clone(d.Attachments), i, i+1)
				return
			}
		}
	})
}

// composerReady reports at least one recipient and only valid addresses.
func composerReady(d MailDraft) bool {
	if len(d.To)+len(d.Cc)+len(d.Bcc) == 0 {
		return false
	}
	for _, list := range [][]string{d.To, d.Cc, d.Bcc} {
		for _, a := range list {
			if _, err := mail.ParseAddress(a); err != nil {
				return false
			}
		}
	}
	return true
}

func composerCopy(d MailDraft) MailDraft {
	d.To, d.Cc, d.Bcc = slices.Clone(d.To), slices.Clone(d.Cc), slices.Clone(d.Bcc)
	d.Attachments = slices.Clone(d.Attachments)
	return d
}

func composerSame(a, b MailDraft) bool {
	return a.Subject == b.Subject && a.Body == b.Body && slices.Equal(a.To, b.To) && slices.Equal(a.Cc, b.Cc) &&
		slices.Equal(a.Bcc, b.Bcc) && slices.Equal(a.Attachments, b.Attachments)
}
