package account

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

var switcherAccounts = []Account{{ID: "a", Name: "Sakiko Togawa", Email: "sakiko@court"}, {ID: "b", Name: "Oblivionis", Email: "ob@mask"}}

func TestAccountSwitcherSwitchesAndAdds(t *testing.T) {
	cur, n, adds := "a", 0, 0
	disabled := false
	tt := accountTester(500, 400, func(c *ui.Context) {
		r := AccountSwitcher(c, &cur, switcherAccounts, AccountSwitcherOptions{Disabled: disabled})
		if r.Changed() {
			n++
		}
		if r.Add() {
			adds++
		}
	})
	if err := tt.Click("Sakiko Togawa"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Oblivionis · ob@mask"); err != nil {
		t.Fatal(err)
	}
	if cur != "b" || n != 1 {
		t.Fatalf("switch: %q %d", cur, n)
	}
	tt.Key(0, ui.KeyDown)
	if err := tt.Click("Add account…"); err != nil {
		t.Fatal(err)
	}
	if adds != 1 || cur != "b" {
		t.Fatal("add not reported")
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Oblivionis")
	if tt.HasText("Add account…") {
		t.Fatal("disabled switcher opened")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown current accepted")
		}
	}()
	bad := "z"
	accountTester(100, 100, func(c *ui.Context) { AccountSwitcher(c, &bad, switcherAccounts, AccountSwitcherOptions{}) })
}

func TestAccountSwitcherSnapshot(t *testing.T) {
	snapshot(t, "account_switcher", 300, 90, func(c *ui.Context) {
		cur := ui.Local(c.Root(), "cur", func() string { return "a" })
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { AccountSwitcher(c, cur, switcherAccounts, AccountSwitcherOptions{}) })
	})
}

func TestAccountSwitcherLargeDataBuildsVisibleRows(t *testing.T) {
	items := make([]Account, 100000)
	for i := range items {
		items[i] = Account{ID: fmt.Sprint(i), Name: fmt.Sprintf("Account %d", i)}
	}
	cur := "0"
	tt := accountTester(500, 500, func(c *ui.Context) { AccountSwitcher(c, &cur, items, AccountSwitcherOptions{}) })
	if err := tt.Click("Account 0"); err != nil {
		t.Fatal(err)
	}
	if n := len(tt.Texts()); n > 100 {
		t.Fatalf("built %d text elements for a bounded viewport", n)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if cur != "99999" {
		t.Fatalf("large-list keyboard choice %q", cur)
	}
}

func TestAccountSwitcherTabChoosesFocusedLargeRow(t *testing.T) {
	items := make([]Account, 9)
	for i := range items {
		items[i] = Account{ID: fmt.Sprint(i), Name: fmt.Sprintf("Account %d", i)}
	}
	cur := "0"
	tt := accountTester(500, 500, func(c *ui.Context) { AccountSwitcher(c, &cur, items, AccountSwitcherOptions{}) })
	if err := tt.Click("Account 0"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Account 1") {
		t.Fatal("Tab did not reach the second account")
	}
	tt.Key(0, ui.KeyEnter)
	if cur != "1" {
		t.Fatal("Enter did not choose the focused account")
	}
}

func TestAccountSwitcherFocusFollowsReorderedID(t *testing.T) {
	items := []Account{{ID: "a", Name: "A"}, {ID: "b", Name: "B"}, {ID: "c", Name: "C"}}
	cur := "a"
	tt := accountTester(500, 500, func(c *ui.Context) { AccountSwitcher(c, &cur, items, AccountSwitcherOptions{}) })
	if err := tt.Click("A"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("B") {
		t.Fatal("Down did not reach B")
	}
	items[1], items[2] = items[2], items[1]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if cur != "b" {
		t.Fatalf("reordered focus chose %q", cur)
	}
}

func TestAccountSwitcherArrowsContinueFromTabFocus(t *testing.T) {
	items := make([]Account, 9)
	for i := range items {
		items[i] = Account{ID: fmt.Sprint(i), Name: fmt.Sprintf("Account %d", i)}
	}
	cur := "0"
	tt := accountTester(500, 500, func(c *ui.Context) { AccountSwitcher(c, &cur, items, AccountSwitcherOptions{}) })
	if err := tt.Click("Account 0"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if cur != "2" {
		t.Fatalf("arrow choice after Tab: %q", cur)
	}
}
