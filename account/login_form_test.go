package account

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// accountLogs records every debug log until the test ends.
func accountLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var b bytes.Buffer
	saved := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&b, &slog.HandlerOptions{Level: slog.LevelDebug})))
	t.Cleanup(func() { slog.SetDefault(saved) })
	return &b
}

type loginCall struct {
	account, password string
	remember          bool
}

func loginTester(s *LoginFormState, calls *[]loginCall, disabled *bool) *ui.Tester {
	return accountTester(420, 400, func(c *ui.Context) {
		LoginForm(c, s, LoginFormOptions{Disabled: *disabled, OnSubmit: func(a, p string, r bool) {
			*calls = append(*calls, loginCall{a, p, r})
		}})
	})
}

func TestLoginFormSubmitsOnceThroughCallback(t *testing.T) {
	logs := accountLogs(t)
	const secret = "s3cret-Sakiko"
	var s LoginFormState
	var calls []loginCall
	disabled := false
	tt := loginTester(&s, &calls, &disabled)
	_ = tt.Click("Email or username")
	tt.Type(" sakiko@court ")
	_ = tt.Click("Password")
	tt.Type(secret)
	_ = tt.Click("Remember me")
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Sign in")
	if len(calls) != 1 || calls[0].account != "sakiko@court" || calls[0].password != secret || !calls[0].remember {
		t.Fatalf("%d calls", len(calls))
	}
	if strings.Contains(logs.String(), secret) {
		t.Fatal("the password reached the log")
	}
	s.EndSubmit("Too many attempts. Try again in 5 minutes.")
	tt.Frame()
	if !tt.HasText("Too many attempts. Try again in 5 minutes.") || s.password != "" || s.Account != " sakiko@court " {
		t.Fatalf("shown %v wiped %v account %q", tt.HasText("Too many attempts. Try again in 5 minutes."), s.password == "", s.Account)
	}
}

func TestLoginFormRequiresBothAndBlocksDisabled(t *testing.T) {
	var s LoginFormState
	var calls []loginCall
	disabled := false
	tt := loginTester(&s, &calls, &disabled)
	_ = tt.Click("Sign in")
	if len(calls) != 0 || !tt.HasText("Required") || !tt.Focused("Email or username") {
		t.Fatalf("empty submit: %d calls", len(calls))
	}
	tt.Type("uika")
	tt.Key(0, ui.KeyEnter)
	if len(calls) != 0 || !tt.Focused("Password") {
		t.Fatal("submitted without a password")
	}
	disabled = true
	tt.Type("pw")
	_ = tt.Click("Sign in")
	if len(calls) != 0 || s.password != "" {
		t.Fatal("disabled form took input or submitted")
	}
}

func TestLoginFormSnapshot(t *testing.T) {
	snapshot(t, "login_form", 400, 330, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() LoginFormState {
			return LoginFormState{Account: "sakiko@court", err: "Wrong account or password."}
		})
		LoginForm(c, s, LoginFormOptions{OnSubmit: func(string, string, bool) {}})
	})
}
