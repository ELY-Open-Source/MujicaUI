package chat

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

const errorMessageText = `429 Too Many Requests: {"error":"rate_limit","retry_after":30}`

func TestErrorMessageKeepsTextAndRetries(t *testing.T) {
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ErrorMessage(c, errorMessageText, ErrorMessageOptions{}).Retried() {
			retries++
		}
	}, 500, 200)
	if !tt.HasText(errorMessageText) || !tt.HasText("Something went wrong") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Retry")
	if retries != 1 {
		t.Fatalf("retries %d", retries)
	}
}

func TestErrorMessageRateLimitCountdown(t *testing.T) {
	retries := 0
	until := time.Now().Add(2500 * time.Millisecond)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ErrorMessage(c, errorMessageText, ErrorMessageOptions{Title: "Rate limited by the API", RetryAt: until}).Retried() {
			retries++
		}
	}, 500, 200)
	if !tt.HasText("Retry in 3 s") {
		t.Fatalf("countdown missing: %q", tt.Texts())
	}
	_ = tt.Click("Retry in 3 s")
	if retries != 0 {
		t.Fatal("retried during the rate limit")
	}
	until = time.Now().Add(-time.Second)
	tt.Frame()
	_ = tt.Click("Retry")
	if retries != 1 {
		t.Fatalf("retries %d after the limit lifted", retries)
	}
}

func TestErrorMessageNeedsText(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without error text")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ErrorMessage(c, "", ErrorMessageOptions{})
	}, 100, 100)
}

func TestErrorMessageSnapshot(t *testing.T) {
	snapshot(t, "error_message", 520, 150, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ErrorMessage(c, errorMessageText, ErrorMessageOptions{}) })
	})
}
