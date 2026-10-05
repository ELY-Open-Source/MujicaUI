package mujica

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var apiKeys = []APIKey{{ID: "k1", Name: "Deploy bot", Prefix: "mk_live_9f2c", Created: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)}}

func TestApiKeyManagerCreateRevealOnceAndRevoke(t *testing.T) {
	logs := accountLogs(t)
	const secret = "mk_live_9f2c_SECRET_DO_NOT_LOG"
	var s ApiKeyManagerState
	var created, revoked []string
	tt := accountTester(700, 600, func(c *ui.Context) {
		r := ApiKeyManager(c, &s, apiKeys, ApiKeyManagerOptions{})
		if name, ok := r.Create(); ok {
			created = append(created, name)
		}
		if id, ok := r.Revoked(); ok {
			revoked = append(revoked, id)
		}
	})
	if !tt.HasText("mk_live_9f2c••••••••••••") || !tt.HasText("Created 2026-09-01 · last used never") {
		t.Fatalf("masked list: %q", tt.Texts())
	}
	_ = tt.Click("Create key")
	tt.Type("CI")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEnter)
	if len(created) != 1 || created[0] != "CI" {
		t.Fatalf("created %v", created)
	}
	s.Reveal(secret)
	tt.Frame()
	if s.secret != secret || !tt.HasText("Copy key") {
		t.Fatal("new key not shown")
	}
	if err := tt.Click("Copy key"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != secret {
		t.Fatal("copy did not take the key")
	}
	_ = tt.Click("I've saved it")
	if tt.HasText("Copy key") || s.secret != "" {
		t.Fatal("key still shown or kept after dismissal")
	}
	if strings.Contains(logs.String(), secret) {
		t.Fatal("the key reached the log")
	}
	_ = tt.Click("Revoke Deploy bot")
	if !tt.Focused("Cancel") {
		t.Fatal("revoke confirmation not focused on Cancel")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if len(revoked) != 1 || revoked[0] != "k1" {
		t.Fatalf("revoked %v", revoked)
	}
}

func TestApiKeyManagerFailAndDisabled(t *testing.T) {
	var s ApiKeyManagerState
	disabled, made := false, 0
	tt := accountTester(700, 600, func(c *ui.Context) {
		if _, ok := ApiKeyManager(c, &s, nil, ApiKeyManagerOptions{Disabled: disabled}).Create(); ok {
			made++
		}
	})
	if !tt.HasText("No API keys") {
		t.Fatal("no empty state")
	}
	_ = tt.Click("Create key")
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEscape)
	if made != 1 || !s.pending {
		t.Fatal("escape left a pending create")
	}
	s.Fail("Key limit reached.")
	tt.Frame()
	if !tt.HasText("Key limit reached.") {
		t.Fatalf("failure: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	disabled = true
	tt.Frame()
	_ = tt.Click("Create key")
	if tt.HasText("New API key") {
		t.Fatal("disabled manager opened the dialog")
	}
}

func TestApiKeyManagerSnapshot(t *testing.T) {
	snapshot(t, "api_key_manager", 560, 220, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() ApiKeyManagerState { return ApiKeyManagerState{} })
		ui.Column(c).Padding(16).Children(func() { ApiKeyManager(c, s, apiKeys, ApiKeyManagerOptions{}) })
	})
}

func TestApiKeyManagerDisablingClosesEveryDialog(t *testing.T) {
	for _, phase := range []string{"name", "pending", "secret", "revoke"} {
		t.Run(phase, func(t *testing.T) {
			var s ApiKeyManagerState
			disabled := false
			tt := accountTester(700, 600, func(c *ui.Context) {
				ui.Column(c).Disabled(disabled).Children(func() { ApiKeyManager(c, &s, apiKeys, ApiKeyManagerOptions{}) })
			})
			if phase == "revoke" {
				if err := tt.Click("Revoke Deploy bot"); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := tt.Click("Create key"); err != nil {
					t.Fatal(err)
				}
				if phase != "name" {
					tt.Type("CI")
					tt.Key(0, ui.KeyEnter)
				}
				if phase == "secret" {
					s.Reveal("one-time-test-secret")
					tt.Frame()
				}
			}
			disabled = true
			tt.Frame()
			if tt.HasText("New API key") || tt.HasText("Copy your new key") || tt.HasText("Revoke Deploy bot?") || s.secret != "" {
				t.Fatal("disabled manager retains an active dialog or secret")
			}
		})
	}
}

func TestApiKeyManagerRemovedKeyCancelsConfirmation(t *testing.T) {
	var s ApiKeyManagerState
	keys := apiKeys
	tt := accountTester(700, 600, func(c *ui.Context) { ApiKeyManager(c, &s, keys, ApiKeyManagerOptions{}) })
	if err := tt.Click("Revoke Deploy bot"); err != nil {
		t.Fatal(err)
	}
	keys = nil
	tt.Frame()
	if s.revoking != "" || tt.HasText("Apps using this key stop working at once.") {
		t.Fatal("removed key keeps confirmation open")
	}
}

func TestApiKeyManagerLargeDataBuildsVisibleRows(t *testing.T) {
	items := make([]APIKey, 100000)
	for i := range items {
		items[i] = APIKey{ID: fmt.Sprint(i), Name: fmt.Sprintf("Key %d", i), Prefix: "public"}
	}
	var s ApiKeyManagerState
	tt := accountTester(600, 500, func(c *ui.Context) { ApiKeyManager(c, &s, items, ApiKeyManagerOptions{}) })
	if n := len(tt.Texts()); n > 100 {
		t.Fatalf("built %d text elements for a bounded viewport", n)
	}
}

func TestApiKeyManagerFailureSurvivesDisabledRequest(t *testing.T) {
	for _, whileDisabled := range []bool{false, true} {
		t.Run(fmt.Sprint(whileDisabled), func(t *testing.T) {
			var s ApiKeyManagerState
			disabled, requests := false, 0
			tt := accountTester(700, 600, func(c *ui.Context) {
				if _, ok := ApiKeyManager(c, &s, nil, ApiKeyManagerOptions{Disabled: disabled}).Create(); ok {
					requests++
				}
			})
			if err := tt.Click("Create key"); err != nil {
				t.Fatal(err)
			}
			tt.Type("CI")
			tt.Key(0, ui.KeyEnter)
			disabled = true
			tt.Frame()
			if whileDisabled {
				s.Fail("Create failed")
				tt.Frame()
			}
			disabled = false
			if !whileDisabled {
				s.Fail("Create failed")
			}
			tt.Frame()
			if !tt.HasText("Create failed") || s.name != "CI" || s.pending {
				t.Fatal("disabled request lost its failure or draft")
			}
			tt.Key(0, ui.KeyEnter)
			if requests != 2 || !s.pending {
				t.Fatal("failed request cannot be retried")
			}
		})
	}
}
