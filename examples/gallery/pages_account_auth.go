package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

type accountRecorderState struct {
	palette, zen mujica.Shortcut
	disabled     bool
}

var accountTaken = []mujica.ShortcutBinding{{Action: "Save", Shortcut: mujica.Shortcut{Mods: ui.Cmd, Key: ui.KeyS}}}

func accountRecorderDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "465", func() accountRecorderState {
		return accountRecorderState{palette: mujica.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}}
	})
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	ui.Text(c, "Cmd/Ctrl+S is taken by Save.")
	for _, r := range []struct {
		name string
		v    *mujica.Shortcut
	}{{"Command palette", &d.palette}, {"Zen mode", &d.zen}} {
		if mujica.ShortcutRecorder(c, r.v, mujica.ShortcutRecorderOptions{Label: r.name, Taken: accountTaken, Disabled: d.disabled}).Changed() {
			log(fmt.Sprintf("%s = %s", r.name, r.v.String()))
		}
	}
}

type accountShortcutsState struct{ query string }

func accountShortcutsDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "466", func() accountShortcutsState { return accountShortcutsState{} })
	groups := []mujica.ShortcutGroup{
		{Title: "General", Bindings: []mujica.ShortcutBinding{
			{Action: "Command palette", Shortcut: mujica.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}},
			{Action: "Settings", Shortcut: mujica.Shortcut{Mods: ui.Cmd, Key: ui.KeyComma}},
			{Action: "Save", Shortcut: mujica.Shortcut{Mods: ui.Cmd, Key: ui.KeyS}},
		}},
		{Title: "Editing", Bindings: []mujica.ShortcutBinding{
			{Action: "Duplicate line", Shortcut: mujica.Shortcut{Mods: ui.Alt | ui.Shift, Key: ui.KeyDown}},
			{Action: "Toggle comment", Shortcut: mujica.Shortcut{Mods: ui.Cmd, Key: ui.KeySlash}},
			{Action: "Zen mode"},
		}},
	}
	before := d.query
	ui.Box(c).Width(480).Children(func() { mujica.KeyboardShortcutsList(c, &d.query, groups, mujica.KeyboardShortcutsListOptions{}) })
	if d.query != before {
		log(fmt.Sprintf("search %q", d.query))
	}
}

type accountLoginState struct {
	form             mujica.LoginFormState
	disabled, reject bool
	attempts         int
}

func accountLoginDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "467", func() accountLoginState { return accountLoginState{reject: true} })
	done := *ui.Local(c.Root(), "467.done", func() chan string { return make(chan string, 1) })
	select {
	case e := <-done:
		d.form.EndSubmit(e)
		log(fmt.Sprintf("sign-in finished, error %q", e))
	default:
	}
	accountToggles(c, log, []string{"Disabled", "Server rejects"}, &d.disabled, &d.reject)
	mujica.LoginForm(c, &d.form, mujica.LoginFormOptions{Disabled: d.disabled, OnSubmit: func(account, password string, remember bool) {
		// Credentials stay in the caller's sign-in request.
		d.attempts++
		log(fmt.Sprintf("submit %q, remember %v", account, remember))
		fail := ""
		switch {
		case d.attempts >= 3:
			fail = "Too many attempts. Try again in 5 minutes."
		case d.reject:
			fail = "Wrong account or password."
		}
		go func() {
			time.Sleep(700 * time.Millisecond)
			done <- fail
			c.Invalidate()
		}()
	}})
}

type accountOAuthState struct {
	pending  string
	disabled bool
}

func accountOAuthDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "468", func() accountOAuthState { return accountOAuthState{} })
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	providers := []mujica.OAuthProvider{{ID: "git", Name: "Git host"}, {ID: "sso", Name: "Company SSO"}, {ID: "mail", Name: "Email link"}}
	if id, ok := mujica.OAuthButtons(c, providers, mujica.OAuthButtonsOptions{Pending: d.pending, Disabled: d.disabled}).Chosen(); ok {
		d.pending = id
		log("continue with " + id)
	}
	if d.pending != "" && tradingAction(c, "Finish sign-in") {
		log("sign-in with " + d.pending + " finished")
		d.pending = ""
	}
}

type accountTwoFactorState struct {
	code     string
	resendAt time.Time
	disabled bool
	err      string
}

func accountTwoFactorDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "469", func() accountTwoFactorState { return accountTwoFactorState{} })
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	r := mujica.TwoFactorInput(c, &d.code, mujica.TwoFactorInputOptions{ResendAt: d.resendAt, Error: d.err, Disabled: d.disabled})
	if r.Completed() {
		// The code is checked here and never logged.
		d.err = ""
		if d.code != "123456" {
			d.err = "That code is not right."
		}
		log(fmt.Sprintf("code complete, accepted %v", d.err == ""))
	}
	if r.Resend() {
		d.resendAt, d.code, d.err = c.Now().Add(30*time.Second), "", ""
		log("resend requested")
	}
}

type accountProfileState struct {
	saved, edit      mujica.Profile
	saving, disabled bool
}

func accountProfileDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "470", func() accountProfileState {
		p := mujica.Profile{Name: "Sakiko Togawa", Bio: "Keyboard, Ave Mujica."}
		return accountProfileState{saved: p, edit: p}
	})
	accountToggles(c, log, []string{"Saving", "Disabled"}, &d.saving, &d.disabled)
	r := mujica.ProfileEditor(c, &d.edit, mujica.ProfileEditorOptions{Saved: d.saved, Saving: d.saving, Disabled: d.disabled})
	if r.Save() {
		d.saved = d.edit
		log("saved")
	}
	if r.Discarded() {
		log("discarded")
	}
	if r.ChangeAvatar() {
		log("change picture")
	}
}
