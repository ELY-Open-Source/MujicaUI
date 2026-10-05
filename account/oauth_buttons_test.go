package account

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var oauthProviders = []OAuthProvider{{ID: "gh", Name: "GitHub"}, {ID: "sso", Name: "Company SSO"}}

func TestOAuthButtonsChooseAndPending(t *testing.T) {
	var chosen []string
	opts := OAuthButtonsOptions{}
	tt := accountTester(420, 200, func(c *ui.Context) {
		if id, ok := OAuthButtons(c, oauthProviders, opts).Chosen(); ok {
			chosen = append(chosen, id)
		}
	})
	if err := tt.Click("Continue with Company SSO"); err != nil {
		t.Fatal(err)
	}
	opts.Pending = "sso"
	tt.Frame()
	_ = tt.Click("Continue with GitHub")
	_ = tt.Click("Continue with Company SSO")
	if len(chosen) != 1 || chosen[0] != "sso" {
		t.Fatalf("chosen %v", chosen)
	}
	opts = OAuthButtonsOptions{Disabled: true}
	tt.Frame()
	_ = tt.Click("Continue with GitHub")
	if len(chosen) != 1 {
		t.Fatal("disabled buttons chose")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("nameless provider accepted")
		}
	}()
	accountTester(100, 100, func(c *ui.Context) { OAuthButtons(c, []OAuthProvider{{ID: "x"}}, OAuthButtonsOptions{}) })
}

func TestOAuthButtonsSnapshot(t *testing.T) {
	snapshot(t, "oauth_buttons", 380, 130, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { OAuthButtons(c, oauthProviders, OAuthButtonsOptions{Pending: "gh"}) })
	})
}

func TestOAuthButtonsKeyboard(t *testing.T) {
	chosen := ""
	tt := accountTester(600, 300, func(c *ui.Context) {
		if id, ok := OAuthButtons(c, oauthProviders, OAuthButtonsOptions{}).Chosen(); ok {
			chosen = id
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if chosen != "sso" {
		t.Fatalf("keyboard provider %q", chosen)
	}
}
