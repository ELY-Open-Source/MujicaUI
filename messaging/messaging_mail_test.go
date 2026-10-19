package messaging

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func mailTestItems(n int) []MailItem {
	at := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	out := make([]MailItem, n)
	for i := range out {
		out[i] = MailItem{ID: fmt.Sprint("m", i), From: fmt.Sprint("Sender ", i), Subject: fmt.Sprint("Subject ", i),
			Snippet: "A short preview", Time: at, Unread: i%2 == 0, Starred: i == 1, Attachments: i % 3}
	}
	return out
}

var recipientContacts = []Contact{{Name: "Uika Misumi", Email: "uika@example.com"}, {Name: "Umiri Yahata", Email: "umiri@example.com"},
	{Name: "Sakiko Togawa", Email: "sakiko@example.com"}}

func recipientTester(list *[]string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(400).Children(func() {
			RecipientInput(c, list, RecipientInputOptions{Contacts: recipientContacts, Label: "To"})
		})
	}, 520, 400)
}
