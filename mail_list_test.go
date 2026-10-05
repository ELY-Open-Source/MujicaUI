package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestMailListSelectStarOpen(t *testing.T) {
	var s ListState[string]
	mails := mailTestItems(6)
	var opened, starred []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := MailList(c, &s, mails, MailListOptions{})
		if r.Opened != "" {
			opened = append(opened, r.Opened)
		}
		if r.Starred != "" {
			starred = append(starred, r.Starred)
		}
	}, 600, 500)
	if err := tt.Click("Subject 0"); err != nil {
		t.Fatal(err)
	}
	r2, _ := tt.Find("Subject 2")
	tt.ClickAtWith(ui.Cmd, r2.X+4, r2.Y+4)
	if s.Selection.Len() != 2 || !s.Selection.Has("m0") || !s.Selection.Has("m2") {
		t.Fatalf("selection %d", s.Selection.Len())
	}
	if err := tt.Click("Star Subject 3"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyS)
	tt.Key(0, ui.KeyEnter)
	if len(starred) != 2 || starred[0] != "m3" || len(opened) != 1 {
		t.Fatalf("starred %v opened %v", starred, opened)
	}
	if _, ok := tt.Find("Unread, Sender 0, Subject 0"); !ok {
		t.Fatal("unread not named")
	}
}

func TestMailListBuildsOnlyVisible(t *testing.T) {
	var s ListState[string]
	mails := mailTestItems(100000)
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		MailList(c, &s, mails, MailListOptions{})
	}, 600, 500)
	tt.Frame()
	built := 0
	for _, x := range tt.Texts() {
		if strings.HasPrefix(x, "Subject ") {
			built++
		}
	}
	if built == 0 || built > 40 {
		t.Fatalf("built %d rows of 100000", built)
	}
}

func TestMailListDisabled(t *testing.T) {
	var s ListState[string]
	starred := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if MailList(c, &s, mailTestItems(3), MailListOptions{Disabled: true}).Starred != "" {
			starred++
		}
	}, 600, 400)
	_ = tt.Click("Star Subject 0")
	_ = tt.Click("Subject 1")
	if starred != 0 || s.Selection.Len() != 0 {
		t.Fatal("disabled list reacted")
	}
}

func TestMailListSnapshot(t *testing.T) {
	mails := []MailItem{
		{ID: "1", From: "Sakiko Togawa", Subject: "Second act", Snippet: "The score is attached.", Unread: true, Attachments: 2},
		{ID: "2", From: "Uika Misumi", Subject: "Costumes", Snippet: "Fittings on Thursday.", Starred: true},
		{ID: "3", From: "Umiri Yahata", Subject: "Venue", Snippet: "The hall is booked."},
	}
	snapshot(t, "mail_list", 480, 260, func(c *ui.Context) {
		s := ui.Local(c.Root(), "mail", func() *ListState[string] { return &ListState[string]{} })
		ui.Column(c).Padding(20).Width(440).Children(func() { MailList(c, *s, mails, MailListOptions{Height: 220}) })
	})
}
