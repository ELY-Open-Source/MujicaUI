package account

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestTwoFactorInputPasteCompleteAndResend(t *testing.T) {
	logs := accountLogs(t)
	code, done, resent := "", 0, 0
	opts := TwoFactorInputOptions{ResendAt: time.Now().Add(-time.Second)}
	tt := accountTester(420, 220, func(c *ui.Context) {
		r := TwoFactorInput(c, &code, opts)
		if r.Completed() {
			done++
		}
		if r.Resend() {
			resent++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.SetClipboard("381 942")
	tt.Command("paste")
	if code != "381942" || done != 1 {
		t.Fatalf("paste: %d digits, %d completions", len(code), done)
	}
	if strings.Contains(logs.String(), "381942") {
		t.Fatal("the code reached the log")
	}
	if err := tt.Click("Resend code"); err != nil {
		t.Fatal(err)
	}
	if resent != 1 {
		t.Fatal("resend not reported")
	}
	opts.ResendAt = tt.Now().Add(65 * time.Second)
	opts.Error = "That code has expired."
	tt.Frame()
	if !tt.HasText("Resend in 1:05") && !tt.HasText("Resend in 1:04") || !tt.HasText("That code has expired.") {
		t.Fatalf("countdown: %q", tt.Texts())
	}
	_ = tt.Click("Resend in 1:05")
	if resent != 1 {
		t.Fatal("resent during countdown")
	}
}

func TestTwoFactorInputDisabled(t *testing.T) {
	code := "12"
	tt := accountTester(420, 220, func(c *ui.Context) {
		TwoFactorInput(c, &code, TwoFactorInputOptions{Disabled: true})
	})
	tt.Key(0, ui.KeyTab)
	tt.Type("3456")
	if code != "12" {
		t.Fatal("disabled input changed")
	}
}

func TestTwoFactorInputSnapshot(t *testing.T) {
	snapshot(t, "two_factor_input", 420, 160, func(c *ui.Context) {
		code := ui.Local(c.Root(), "code", func() string { return "3819" })
		ui.Column(c).Padding(16).Children(func() {
			TwoFactorInput(c, code, TwoFactorInputOptions{ResendAt: c.Now().Add(42 * time.Second)})
		})
	})
}

func TestTwoFactorInputCountdownRoundsUp(t *testing.T) {
	tt := accountTester(420, 220, func(c *ui.Context) {
		code := ""
		TwoFactorInput(c, &code, TwoFactorInputOptions{ResendAt: c.Now().Add(100 * time.Millisecond)})
	})
	if !tt.HasText("Resend in 0:01") {
		t.Fatal("positive wait must show one remaining second")
	}
}

func TestTwoFactorInputInvalidValueKeepsSecretOutOfPanic(t *testing.T) {
	const secret = "1234567"
	defer func() {
		v := recover()
		if v == nil {
			t.Fatal("oversized code accepted")
		}
		if strings.Contains(fmt.Sprint(v), secret) {
			t.Fatal("code appears in panic")
		}
	}()
	code := secret
	accountTester(420, 220, func(c *ui.Context) { TwoFactorInput(c, &code, TwoFactorInputOptions{}) })
}
