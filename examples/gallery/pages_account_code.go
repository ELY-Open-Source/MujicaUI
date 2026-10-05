package main

// accountCode contains the compiled Example bodies.
var accountCode = map[string]string{
	"SettingsLayout": `page := ""
sound := true
pages := []mujica.SettingsPage{{ID: "general", Label: "General"}, {ID: "privacy", Label: "Privacy"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.SettingsLayout(c, &page, pages, mujica.SettingsLayoutOptions{}, func(id string) {
		mujica.SettingsSection(c, "Notifications", func() {
			mujica.SettingsRow(c, "Sound", mujica.SettingsRowOptions{ID: "sound"}, func() *ui.Element {
				return mujica.Switch(c, &sound, "", mujica.SwitchOptions{})
			})
		})
	})
}
ui.Render(view, 400, 300, 1)`,
	"SettingsSearch": `query, page := "", "general"
entries := []mujica.SettingEntry{{ID: "sound", Page: "general", Label: "Sound", Keywords: []string{"audio"}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if e, ok := mujica.SettingsSearch(c, &query, entries, mujica.SettingsSearchOptions{}).Chosen(); ok {
		page = e.Page
	}
	fmt.Println("page", page)
}
ui.Render(view, 400, 300, 1)`,
	"ThemeSelector": `mode := mujica.System
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{Mode: mode})
	mujica.ThemeSelector(c, &mode, mujica.ThemeSelectorOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"AccentColorPicker": `accent := ui.Hex("#7D2034")
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AccentColorPicker(c, &accent, mujica.AccentColorPickerOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"ShortcutRecorder": `palette := mujica.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ShortcutRecorder(c, &palette, mujica.ShortcutRecorderOptions{
		Label: "Command palette",
		Taken: []mujica.ShortcutBinding{{Action: "Save", Shortcut: mujica.Shortcut{Mods: ui.Cmd, Key: ui.KeyS}}},
	})
}
ui.Render(view, 400, 300, 1)`,
	"KeyboardShortcutsList": `query := ""
groups := []mujica.ShortcutGroup{{Title: "General", Bindings: []mujica.ShortcutBinding{
	{Action: "Command palette", Shortcut: mujica.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}},
}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.KeyboardShortcutsList(c, &query, groups, mujica.KeyboardShortcutsListOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"LoginForm": `var state mujica.LoginFormState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.LoginForm(c, &state, mujica.LoginFormOptions{OnSubmit: func(account, password string, remember bool) {
		// Submit credentials through the caller-provided auth service.
		state.EndSubmit("")
	}})
}
ui.Render(view, 400, 300, 1)`,
	"OAuthButtons": `providers := []mujica.OAuthProvider{{ID: "sso", Name: "Company SSO"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if id, ok := mujica.OAuthButtons(c, providers, mujica.OAuthButtonsOptions{}).Chosen(); ok {
		fmt.Println("start sign-in with", id)
	}
}
ui.Render(view, 400, 300, 1)`,
	"TwoFactorInput": `code := ""
resendAt := time.Now().Add(30 * time.Second)
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.TwoFactorInput(c, &code, mujica.TwoFactorInputOptions{ResendAt: resendAt})
	if r.Completed() {
		fmt.Println("verify the code")
	}
}
ui.Render(view, 400, 300, 1)`,
	"ProfileEditor": `saved := mujica.Profile{Name: "Sakiko", Bio: "Keyboard."}
profile := saved
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ProfileEditor(c, &profile, mujica.ProfileEditorOptions{Saved: saved}).Save() {
		saved = profile
	}
}
ui.Render(view, 400, 300, 1)`,
	"AccountSwitcher": `current := "a"
accounts := []mujica.Account{{ID: "a", Name: "Sakiko Togawa", Email: "sakiko@court"}, {ID: "b", Name: "Oblivionis"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.AccountSwitcher(c, &current, accounts, mujica.AccountSwitcherOptions{}).Add() {
		fmt.Println("sign in another account")
	}
}
ui.Render(view, 400, 300, 1)`,
	"WorkspaceSwitcher": `current := "ave"
spaces := []mujica.Workspace{{ID: "ave", Name: "Ave Mujica"}, {ID: "mygo", Name: "MyGO!!!!!"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WorkspaceSwitcher(c, &current, spaces, mujica.WorkspaceSwitcherOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"UserMenu": `status := mujica.UserStatusOnline
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.UserMenu(c, &status, mujica.UserMenuOptions{Name: "Sakiko Togawa"}).SignOut() {
		fmt.Println("sign out")
	}
}
ui.Render(view, 400, 300, 1)`,
	"SessionList": `sessions := []mujica.Session{{ID: "s1", Device: "MacBook Pro", Current: true}, {ID: "s2", Device: "iPhone", Location: "Osaka"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if id, ok := mujica.SessionList(c, sessions, mujica.SessionListOptions{}).Revoked(); ok {
		fmt.Println("sign out", id)
	}
}
ui.Render(view, 400, 300, 1)`,
	"ApiKeyManager": `var state mujica.ApiKeyManagerState
keys := []mujica.APIKey{{ID: "k1", Name: "Deploy bot", Prefix: "mk_live_9f2c"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if name, ok := mujica.ApiKeyManager(c, &state, keys, mujica.ApiKeyManagerOptions{}).Create(); ok {
		fmt.Println("create", name)
		state.Reveal("secret from your key service")
	}
}
ui.Render(view, 400, 300, 1)`,
	"UsageQuota": `items := []mujica.QuotaItem{{Label: "API calls", Used: 4200, Limit: 10000}, {Label: "Storage", Used: 9, Limit: 10, Unit: "GB"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.UsageQuota(c, items, mujica.UsageQuotaOptions{ResetAt: time.Now().Add(72 * time.Hour)})
}
ui.Render(view, 400, 300, 1)`,
	"OnboardingWizard": `step := 0
steps := []mujica.WizardStep{{Title: "Welcome"}, {Title: "Pick a theme"}, {Title: "Done"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.OnboardingWizard(c, &step, steps, mujica.OnboardingWizardOptions{}).Finished() {
		fmt.Println("onboarded")
	}
}
ui.Render(view, 400, 300, 1)`,
	"Coachmark": `open, step := true, 0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	compose := mujica.Button(c, "Compose", mujica.ButtonOptions{})
	mujica.Coachmark(c, &open, &step, []mujica.CoachmarkStep{{Anchor: compose, Title: "Write here", Body: "Start a new letter."}})
}
ui.Render(view, 400, 300, 1)`,
	"WhatsNewDialog": `open, hide := true, false
releases := []mujica.Release{{Version: "1.1", Changes: []mujica.ReleaseChange{{Kind: mujica.ChangeNew, Text: "500 components."}}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WhatsNewDialog(c, &open, &hide, releases, mujica.WhatsNewDialogOptions{})
}
ui.Render(view, 400, 300, 1)`,
	"FeedbackWidget": `var state mujica.FeedbackWidgetState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if f, ok := mujica.FeedbackWidget(c, &state, mujica.FeedbackWidgetOptions{}).Submitted(); ok {
		fmt.Println("send", f.Kind)
		state.EndSubmit("")
	}
}
ui.Render(view, 400, 300, 1)`,
}
