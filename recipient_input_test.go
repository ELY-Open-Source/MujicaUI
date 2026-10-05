package mujica

import (
	"net/mail"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestRecipientInputRefreshesContacts(t *testing.T) {
	var recipients []string
	contacts := []Contact{{Name: "Alice", Email: "alice@example.com"}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		RecipientInput(c, &recipients, RecipientInputOptions{Contacts: contacts, Label: "To"})
	}, 500, 300)
	if err := tt.Click("To"); err != nil {
		t.Fatal(err)
	}
	tt.Type("ali")
	contacts = []Contact{{Name: "Alina", Email: "alina@example.com"}}
	tt.Frame()
	if !tt.HasText("Alina") || tt.HasText("Alice") {
		t.Fatalf("stale contacts: %v", tt.Texts())
	}
}

func TestRecipientInputPopupFitsWindow(t *testing.T) {
	var recipients []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Absolute().Left(190).Top(220).Width(110).Children(func() {
			RecipientInput(c, &recipients, RecipientInputOptions{Contacts: recipientContacts, Label: "To"})
		})
	}, 320, 280)
	if err := tt.Click("To"); err != nil {
		t.Fatal(err)
	}
	tt.Type("ui")
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	r, ok := tt.Find("Contacts")
	if !ok || r.X < 0 || r.Y < 0 || r.X+r.W > 320 || r.Y+r.H > 280 {
		t.Fatalf("popup bounds: %+v", r)
	}
}

func TestRecipientInputDisabledPreservesAddresses(t *testing.T) {
	recipients := []string{"a@example.com"}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			RecipientInput(c, &recipients, RecipientInputOptions{Contacts: recipientContacts, Label: "To"})
		})
	}, 500, 300)
	_ = tt.Click("To")
	tt.Type("b@example.com,")
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Remove a@example.com")
	if !slices.Equal(recipients, []string{"a@example.com"}) {
		t.Fatalf("disabled recipients: %v", recipients)
	}
}

func TestRecipientInputValidates(t *testing.T) {
	var list []string
	tt := recipientTester(&list)
	if err := tt.Click("To"); err != nil {
		t.Fatal(err)
	}
	tt.Type("nobody")
	tt.Key(0, ui.KeyEnter)
	if len(list) != 0 || !tt.HasText(`"nobody" is not an email address`) {
		t.Fatalf("invalid address kept: %v %v", list, tt.Texts())
	}
	tt.Command("selectAll")
	tt.Type("guest@example.com,")
	if !slices.Equal(list, []string{"guest@example.com"}) {
		t.Fatalf("list %v", list)
	}
}

func TestRecipientInputSuggests(t *testing.T) {
	var list []string
	tt := recipientTester(&list)
	if err := tt.Click("To"); err != nil {
		t.Fatal(err)
	}
	tt.Type("yah")
	if !tt.HasText("Umiri Yahata") || tt.HasText("Sakiko Togawa") {
		t.Fatalf("suggestions %v", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(list, []string{"Umiri Yahata <umiri@example.com>"}) {
		t.Fatalf("list %v", list)
	}
	tt.Type("ui")
	if err := tt.Click("Uika Misumi <uika@example.com>"); err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[1] != "Uika Misumi <uika@example.com>" {
		t.Fatalf("click pick %v", list)
	}
	tt.Type("sa")
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("sakiko@example.com") {
		t.Fatal("Escape did not hide suggestions")
	}
}

func TestRecipientInputSnapshot(t *testing.T) {
	list := []string{"Uika Misumi <uika@example.com>", "guest@example.com"}
	snapshot(t, "recipient_input", 460, 120, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(420).Children(func() {
			RecipientInput(c, &list, RecipientInputOptions{Contacts: recipientContacts, Label: "To"})
		})
	})
}

func TestRecipientInputQuotesDisplayNames(t *testing.T) {
	var recipients []string
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		RecipientInput(c, &recipients, RecipientInputOptions{Label: "To", Contacts: []Contact{{Name: "Doe, Jane", Email: "jane@example.com"}}})
	}, 500, 300)
	if err := tt.Click("To"); err != nil {
		t.Fatal(err)
	}
	tt.Type("jane")
	if err := tt.Click(`"Doe, Jane" <jane@example.com>`); err != nil {
		t.Fatal(err)
	}
	if len(recipients) != 1 {
		t.Fatalf("recipients=%v texts=%v", recipients, tt.Texts())
	}
	addr, err := mail.ParseAddress(recipients[0])
	if err != nil || addr.Name != "Doe, Jane" || addr.Address != "jane@example.com" {
		t.Fatalf("quoted address=%v error=%v", recipients, err)
	}
}

func TestRecipientInputTabShowsFocusRing(t *testing.T) {
	var recipients []string
	var focus ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		focus = Tokens(c).Focus
		ui.Column(c).Padding(20).Children(func() {
			RecipientInput(c, &recipients, RecipientInputOptions{Contacts: recipientContacts, Label: "To"})
		})
	}, 500, 300)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("To"); err != nil {
		t.Fatal(err)
	}
	tt.Type("ui")
	tt.Key(0, ui.KeyTab)
	label := "Uika Misumi <uika@example.com>"
	if !tt.Focused(label) {
		t.Fatal("Tab did not focus suggestion")
	}
	r, ok := tt.Find(label)
	if !ok {
		t.Fatal("suggestion missing")
	}
	pixels := 0
	img := tt.Image()
	for y := int(r.Y) - 4; y < int(r.Y+r.H)+4; y++ {
		for x := int(r.X) - 4; x < int(r.X+r.W)+4; x++ {
			if x >= int(r.X) && x < int(r.X+r.W) && y >= int(r.Y) && y < int(r.Y+r.H) {
				continue
			}
			pixel := img.RGBAAt(x, y)
			if pixel.R == focus.R && pixel.G == focus.G && pixel.B == focus.B {
				pixels++
			}
		}
	}
	if pixels == 0 {
		t.Fatal("keyboard suggestion has no focus ring")
	}
	tt.Key(0, ui.KeyEnter)
	if len(recipients) != 1 {
		t.Fatal("focused suggestion did not commit")
	}
}
