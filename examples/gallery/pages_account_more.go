package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type accountSwitchState struct {
	account, space string
	accounts       []account.Account
	disabled       bool
}

func accountSwitcherDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "471", func() accountSwitchState {
		return accountSwitchState{account: "a", accounts: []account.Account{
			{ID: "a", Name: "Sakiko Togawa", Email: "sakiko@court"}, {ID: "b", Name: "Oblivionis", Email: "ob@mask"}}}
	})
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	r := account.AccountSwitcher(c, &d.account, d.accounts, account.AccountSwitcherOptions{Disabled: d.disabled})
	if r.Changed() {
		log("switched to " + d.account)
	}
	if r.Add() {
		n := len(d.accounts) + 1
		d.accounts = append(d.accounts, account.Account{ID: fmt.Sprint(n), Name: fmt.Sprintf("Guest %d", n)})
		log("add account")
	}
}

var accountSpaces = []account.Workspace{
	{ID: "ave", Name: "Ave Mujica", Detail: "5 members"}, {ID: "mygo", Name: "MyGO!!!!!", Detail: "5 members"},
	{ID: "crychic", Name: "CRYCHIC", Detail: "Archived"}, {ID: "sumimi", Name: "sumimi", Detail: "2 members"},
}

func accountWorkspaceDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "472", func() accountSwitchState { return accountSwitchState{space: "ave"} })
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	if account.WorkspaceSwitcher(c, &d.space, accountSpaces, account.WorkspaceSwitcherOptions{Disabled: d.disabled}).Changed() {
		log("workspace " + d.space)
	}
}

type accountUserState struct {
	status   account.UserStatus
	disabled bool
}

func accountUserMenuDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "473", func() accountUserState { return accountUserState{} })
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	r := account.UserMenu(c, &d.status, account.UserMenuOptions{Name: "Sakiko Togawa", Email: "sakiko@court", Disabled: d.disabled})
	if r.StatusChanged() {
		log(fmt.Sprintf("status %d", d.status))
	}
	if r.Settings() {
		log("open settings")
	}
	if r.SignOut() {
		log("sign out")
	}
}

type accountSessionsState struct {
	sessions []account.Session
	disabled bool
}

func accountSessionsDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "474", func() accountSessionsState {
		return accountSessionsState{sessions: []account.Session{
			{ID: "s1", Device: "MacBook Pro · Safari", Location: "Tokyo", LastActive: time.Now(), Current: true},
			{ID: "s2", Device: "iPhone · App", Location: "Osaka", LastActive: time.Now().Add(-26 * time.Hour)},
			{ID: "s3", Device: "Windows · Edge", Location: "Kyoto"},
		}}
	})
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	r := account.SessionList(c, d.sessions, account.SessionListOptions{Disabled: d.disabled})
	if id, ok := r.Revoked(); ok {
		for i, s := range d.sessions {
			if s.ID == id {
				d.sessions = append(d.sessions[:i:i], d.sessions[i+1:]...)
				break
			}
		}
		log("signed out " + id)
	}
	if r.OthersRevoked() {
		d.sessions = d.sessions[:1]
		log("signed out every other session")
	}
}

type accountKeysState struct {
	state            account.ApiKeyManagerState
	keys             []account.APIKey
	disabled, reject bool
}

func accountKeysDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "475", func() accountKeysState {
		return accountKeysState{keys: []account.APIKey{{ID: "k1", Name: "Deploy bot", Prefix: "mk_live_9f2c", Created: time.Now().AddDate(0, -1, 0)}}}
	})
	accountToggles(c, log, []string{"Disabled", "Create fails"}, &d.disabled, &d.reject)
	r := account.ApiKeyManager(c, &d.state, d.keys, account.ApiKeyManagerOptions{Disabled: d.disabled})
	if name, ok := r.Create(); ok {
		if d.reject {
			d.state.Fail("Key limit reached.")
			log("create failed")
			return
		}
		// The caller supplies the generated key.
		n := len(d.keys) + 1
		prefix := fmt.Sprintf("mk_live_%04x", n*7919)
		d.keys = append(d.keys, account.APIKey{ID: fmt.Sprint("k", n), Name: name, Prefix: prefix, Created: time.Now()})
		d.state.Reveal(fmt.Sprintf("%s_%x", prefix, time.Now().UnixNano()))
		log("created key " + name)
	}
	if id, ok := r.Revoked(); ok {
		for i, k := range d.keys {
			if k.ID == id {
				d.keys = append(d.keys[:i:i], d.keys[i+1:]...)
				break
			}
		}
		log("revoked " + id)
	}
}

type accountQuotaState struct{ calls float64 }

func accountQuotaDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "476", func() accountQuotaState { return accountQuotaState{calls: 4200} })
	if input.Slider(c, &d.calls, input.SliderOptions{Min: 0, Max: 12000, Step: 100, Label: "API calls used", ShowValue: true}).Changed() {
		log(fmt.Sprintf("calls %.0f", d.calls))
	}
	ui.Box(c).Width(440).Children(func() {
		account.UsageQuota(c, []account.QuotaItem{{Label: "API calls", Used: d.calls, Limit: 10000}, {Label: "Storage", Used: 9, Limit: 10, Unit: "GB"}},
			account.UsageQuotaOptions{ResetAt: time.Now().Add(52 * time.Hour)})
	})
}

type accountWizardState struct {
	step   int
	theme  core.Mode
	name   string
	noSkip bool
}

func accountWizardDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "477", func() accountWizardState { return accountWizardState{} })
	accountToggles(c, log, []string{"Hide Skip"}, &d.noSkip)
	steps := []account.WizardStep{
		{Title: "Welcome", Description: "Three short steps to set up your court."},
		{Title: "Your name", Ready: func() bool { return d.name != "" }, Content: func() {
			if input.TextInput(c, &d.name, input.TextInputOptions{Label: "Name", Placeholder: "Name"}).Changed() {
				log("name changed")
			}
		}},
		{Title: "Pick a theme", Content: func() {
			if account.ThemeSelector(c, &d.theme, account.ThemeSelectorOptions{}).Changed() {
				log("theme " + d.theme.String())
			}
		}},
	}
	from := d.step
	r := account.OnboardingWizard(c, &d.step, steps, account.OnboardingWizardOptions{NoSkip: d.noSkip})
	if d.step != from {
		log(fmt.Sprintf("step %d", d.step+1))
	}
	if r.Finished() || r.Skipped() {
		log(fmt.Sprintf("finished %v skipped %v", r.Finished(), r.Skipped()))
		d.step = 0
	}
}

type accountCoachState struct {
	open, disabled bool
	step           int
}

func accountCoachDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "478", func() accountCoachState { return accountCoachState{} })
	accountToggles(c, log, []string{"Disabled anchor"}, &d.disabled)
	if tradingAction(c, "Start tour") {
		d.open, d.step = true, 0
		log("tour started")
	}
	var compose, archive *ui.Element
	ui.Row(c).Gap(48).Padding(8).Children(func() {
		compose = input.Button(c, "Compose", input.ButtonOptions{Disabled: d.disabled})
		archive = input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
	})
	if compose.Clicked() {
		log("compose")
	}
	if archive.Clicked() {
		log("archive")
	}
	from := d.step
	r := account.Coachmark(c, &d.open, &d.step, []account.CoachmarkStep{
		{Anchor: compose, Title: "Write here", Body: "Start a new letter."},
		{Anchor: archive, Title: "Keep it tidy", Body: "Archive threads you have finished."},
	})
	if d.step != from {
		log(fmt.Sprintf("tip %d", d.step+1))
	}
	if r.Finished() {
		log("tour finished")
	}
	if r.Dismissed() {
		log("tour dismissed")
	}
}

type accountWhatsNewState struct{ open, hide bool }

func accountWhatsNewDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "479", func() accountWhatsNewState { return accountWhatsNewState{} })
	if tradingAction(c, "Show what's new") {
		d.open = true
		log("open")
	}
	releases := []account.Release{
		{Version: "1.1", Date: "2026-10-05", Changes: []account.ReleaseChange{{Kind: account.ChangeNew, Text: "400 new components."}, {Kind: account.ChangeFixed, Text: "Snapshots no longer drift."}}},
		{Version: "1.0", Date: "2026-09-01", Changes: []account.ReleaseChange{{Kind: account.ChangeImproved, Text: "Court-gothic theme."}}},
	}
	was := d.open
	if account.WhatsNewDialog(c, &d.open, &d.hide, releases, account.WhatsNewDialogOptions{}).Changed() {
		log(fmt.Sprintf("don't show again %v", d.hide))
	}
	if was && !d.open {
		log("closed")
	}
}

type accountFeedbackState struct {
	state    account.FeedbackWidgetState
	disabled bool
}

func accountFeedbackDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "480", func() accountFeedbackState { return accountFeedbackState{} })
	done := *ui.Local(c.Root(), "480.done", func() chan string { return make(chan string, 1) })
	select {
	case e := <-done:
		d.state.EndSubmit(e)
		log("feedback sent")
	default:
	}
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	before := d.state.Draft
	r := account.FeedbackWidget(c, &d.state, account.FeedbackWidgetOptions{Disabled: d.disabled, Attachments: func() {
		if tradingAction(c, "Attach screenshot") {
			log("attach screenshot")
		}
	}})
	if before.Kind != d.state.Draft.Kind {
		log("type changed")
	}
	if before.Text != d.state.Draft.Text {
		log("message changed")
	}
	if f, ok := r.Submitted(); ok {
		log(fmt.Sprintf("send kind %d, %d characters", f.Kind, len([]rune(f.Text))))
		go func() {
			time.Sleep(600 * time.Millisecond)
			done <- ""
			c.Invalidate()
		}()
	}
}
