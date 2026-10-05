package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMailReaderPlainQuotesAndAttachments(t *testing.T) {
	m := Mail{From: "Sakiko Togawa", To: []string{"Uika"}, Subject: "Second act", Body: "See the score.\n> Can you send the score?\n> Thanks",
		Attachments: []MailAttachment{{Name: "score.pdf", Size: 2048}}}
	opened := -1
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(500).Children(func() {
			if r := MailReader(c, m); r.Opened >= 0 {
				opened = r.Opened
			}
		})
	}, 600, 600)
	if !tt.HasText("See the score.") || tt.HasText("Can you send the score?") || !tt.HasText("Show quoted text (2 lines)") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if err := tt.Click("Show quoted text (2 lines)"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Can you send the score?\nThanks") {
		t.Fatalf("quote not unfolded: %v", tt.Texts())
	}
	if err := tt.Click("score.pdf"); err != nil {
		t.Fatal(err)
	}
	if opened != 0 {
		t.Fatalf("opened %d", opened)
	}
}

func TestMailReaderMarkdownIsSafe(t *testing.T) {
	m := Mail{From: "A", Subject: "B", Format: MailMarkdown, Body: "**Bold** and [the hall](https://example.com)\n\n<script>alert(1)</script>"}
	link := ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(500).Children(func() {
			if r := MailReader(c, m); r.Link != "" {
				link = r.Link
			}
		})
	}, 600, 400)
	if err := tt.Click("the hall"); err != nil {
		t.Fatal(err)
	}
	if link != "https://example.com" || len(tt.OpenedURLs()) != 0 {
		t.Fatalf("link %q opened %v", link, tt.OpenedURLs())
	}
	if !tt.HasText("<script>alert(1)</script>") {
		t.Fatal("HTML was not shown as text")
	}
	mediaPanics(t, "no subject", func(c *ui.Context) { MailReader(c, Mail{From: "A"}) })
}

func TestMailReaderSnapshot(t *testing.T) {
	m := Mail{From: "Sakiko Togawa", To: []string{"Uika Misumi", "Mutsumi Wakaba"}, Cc: []string{"Umiri Yahata"}, Subject: "The second act",
		Date: time.Date(2026, 10, 5, 21, 7, 0, 0, time.UTC), Body: "The score is attached.\nRehearsal is at nine.\n> Can you send it tonight?",
		Attachments: []MailAttachment{{Name: "score.pdf", Size: 2 << 20}}}
	snapshot(t, "mail_reader", 560, 420, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(520).Children(func() { MailReader(c, m) })
	})
}
