package messaging

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func TestMailComposerDisabledDefersSave(t *testing.T) {
	draft := MailDraft{To: []string{"a@example.com"}}
	disabled, saves := false, 0
	var state *composerState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			r := MailComposer(c, &draft, MailComposerOptions{OnSaveDraft: func(MailDraft) { saves++ }})
			state = *ui.Local(r.Element, composerKey{}, func() *composerState { panic("missing state") })
		})
	}, 600, 600)
	if err := tt.Click("Subject"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Draft")
	disabled = true
	state.edited = time.Now().Add(-10 * time.Second)
	tt.Frame()
	if saves != 0 {
		t.Fatalf("disabled save callbacks: %d", saves)
	}
	disabled = false
	tt.Frame()
	if saves != 1 {
		t.Fatalf("resumed save callbacks: %d", saves)
	}
}

func TestMailComposerManualSaveAndFocus(t *testing.T) {
	draft := MailDraft{To: []string{"a@example.com"}}
	saves, sent := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if MailComposer(c, &draft, MailComposerOptions{OnSaveDraft: func(MailDraft) { saves++ }}).Sent {
			sent++
		}
		input.TextInput(c, new(string), input.TextInputOptions{Label: "Other"})
	}, 600, 700)
	if err := tt.Click("Subject"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyS)
	if saves != 1 {
		t.Fatalf("manual save callbacks: %d", saves)
	}
	if err := tt.Click("Other"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyEnter)
	if sent != 0 {
		t.Fatalf("composer submitted from another field: %d", sent)
	}
}

func TestMailComposerAutosaveAndSend(t *testing.T) {
	d := MailDraft{To: []string{"uika@example.com"}}
	var saves []MailDraft
	sent := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(560).Children(func() {
			if MailComposer(c, &d, MailComposerOptions{OnSaveDraft: func(m MailDraft) { saves = append(saves, m) }}).Sent {
				sent++
			}
		})
	}, 700, 700)
	if err := tt.Click("Subject"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Second act")
	if len(saves) != 0 || !tt.HasText("Unsaved changes") {
		t.Fatal("saved before the draft rested")
	}
	tt.Key(ui.Cmd, ui.KeyS)
	if len(saves) != 1 || saves[0].Subject != "Second act" {
		t.Fatalf("Cmd+S saves %+v", saves)
	}
	if err := tt.Click("Cc"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Message"); err != nil {
		t.Fatal(err)
	}
	tt.Type("The score is attached.")
	tt.Key(ui.Cmd, ui.KeyEnter)
	if sent != 1 || d.Body != "The score is attached." {
		t.Fatalf("sent %d body %q", sent, d.Body)
	}
}

func TestMailComposerNeedsRecipientsAndDisabled(t *testing.T) {
	d := MailDraft{Subject: "x"}
	sent := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if MailComposer(c, &d, MailComposerOptions{OnSaveDraft: func(MailDraft) {}}).Sent {
			sent++
		}
	}, 700, 700)
	_ = tt.Click("Send")
	tt.Key(ui.Cmd, ui.KeyEnter)
	if sent != 0 {
		t.Fatal("sent without recipients")
	}
	mediaPanics(t, "no save callback", func(c *ui.Context) { MailComposer(c, &d, MailComposerOptions{}) })
}

func TestMailComposerSnapshot(t *testing.T) {
	d := MailDraft{To: []string{"Uika Misumi <uika@example.com>"}, Subject: "The second act", Body: "The score is attached.",
		Attachments: []MailAttachment{{Name: "score.pdf", Size: 2048}}}
	snapshot(t, "mail_composer", 600, 460, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(560).Children(func() {
			MailComposer(c, &d, MailComposerOptions{OnSaveDraft: func(MailDraft) {}, OnAttach: func() {}})
		})
	})
}

func TestMailComposerAttachmentsAndAddressFields(t *testing.T) {
	d := MailDraft{To: []string{"uika@example.com"}}
	attachments, discards := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := MailComposer(c, &d, MailComposerOptions{OnSaveDraft: func(MailDraft) {}, OnAttach: func() {
			d.Attachments = append(d.Attachments, MailAttachment{Name: "score.pdf", Size: 2048})
			attachments++
		}})
		if r.Discarded {
			discards++
		}
	}, 700, 700)
	for _, label := range []string{"Cc", "Bcc", "Attach"} {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	if attachments != 1 || len(d.Attachments) != 1 {
		t.Fatal("attachment callback did not update draft")
	}
	if err := tt.Click("Cc"); err != nil {
		t.Fatal(err)
	}
	tt.Type("cc@example.com")
	tt.Key(0, ui.KeyEnter)
	if err := tt.Click("Bcc"); err != nil {
		t.Fatal(err)
	}
	tt.Type("bcc@example.com")
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(d.Cc, []string{"cc@example.com"}) || !slices.Equal(d.Bcc, []string{"bcc@example.com"}) {
		t.Fatal("address fields did not commit")
	}
	if err := tt.Click("Remove score.pdf"); err != nil {
		t.Fatal(err)
	}
	if len(d.Attachments) != 0 {
		t.Fatal("attachment removal failed")
	}
	if err := tt.Click("Discard"); err != nil {
		t.Fatal(err)
	}
	if discards != 1 {
		t.Fatalf("discard events=%d", discards)
	}
}

func TestMailComposerLoadedAddressFields(t *testing.T) {
	d := MailDraft{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		MailComposer(c, &d, MailComposerOptions{OnSaveDraft: func(MailDraft) {}})
	}, 700, 700)
	d = MailDraft{Cc: []string{"cc@example.com"}, Bcc: []string{"bcc@example.com"}}
	tt.Frame()
	if !tt.HasText("cc@example.com") || !tt.HasText("bcc@example.com") {
		t.Fatal("loaded recipients remain hidden")
	}
}
