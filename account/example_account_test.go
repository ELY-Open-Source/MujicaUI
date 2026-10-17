package account_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func ExampleSettingsLayout() {
	page := ""
	sound := true
	pages := []account.SettingsPage{{ID: "general", Label: "General"}, {ID: "privacy", Label: "Privacy"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.SettingsLayout(c, &page, pages, account.SettingsLayoutOptions{}, func(id string) {
			account.SettingsSection(c, "Notifications", func() {
				account.SettingsRow(c, "Sound", account.SettingsRowOptions{ID: "sound"}, func() ui.Element {
					return input.Switch(c, &sound, "", input.SwitchOptions{})
				})
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSettingsSearch() {
	query, page := "", "general"
	entries := []account.SettingEntry{{ID: "sound", Page: "general", Label: "Sound", Keywords: []string{"audio"}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if e, ok := account.SettingsSearch(c, &query, entries, account.SettingsSearchOptions{}).Chosen(); ok {
			page = e.Page
		}
		fmt.Println("page", page)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleThemeSelector() {
	mode := core.System
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: mode})
		account.ThemeSelector(c, &mode, account.ThemeSelectorOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAccentColorPicker() {
	accent := ui.Hex("#7D2034")
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.AccentColorPicker(c, &accent, account.AccentColorPickerOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleShortcutRecorder() {
	palette := account.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.ShortcutRecorder(c, &palette, account.ShortcutRecorderOptions{
			Label: "Command palette",
			Taken: []account.ShortcutBinding{{Action: "Save", Shortcut: account.Shortcut{Mods: ui.Cmd, Key: ui.KeyS}}},
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleKeyboardShortcutsList() {
	query := ""
	groups := []account.ShortcutGroup{{Title: "General", Bindings: []account.ShortcutBinding{
		{Action: "Command palette", Shortcut: account.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}},
	}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.KeyboardShortcutsList(c, &query, groups, account.KeyboardShortcutsListOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleLoginForm() {
	var state account.LoginFormState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.LoginForm(c, &state, account.LoginFormOptions{OnSubmit: func(account, password string, remember bool) {
			// Submit credentials through the caller-provided auth service.
			state.EndSubmit("")
		}})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOAuthButtons() {
	providers := []account.OAuthProvider{{ID: "sso", Name: "Company SSO"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if id, ok := account.OAuthButtons(c, providers, account.OAuthButtonsOptions{}).Chosen(); ok {
			fmt.Println("start sign-in with", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTwoFactorInput() {
	code := ""
	resendAt := time.Now().Add(30 * time.Second)
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := account.TwoFactorInput(c, &code, account.TwoFactorInputOptions{ResendAt: resendAt})
		if r.Completed() {
			fmt.Println("verify the code")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleProfileEditor() {
	saved := account.Profile{Name: "Sakiko", Bio: "Keyboard."}
	profile := saved
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if account.ProfileEditor(c, &profile, account.ProfileEditorOptions{Saved: saved}).Save() {
			saved = profile
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAccountSwitcher() {
	current := "a"
	accounts := []account.Account{{ID: "a", Name: "Sakiko Togawa", Email: "sakiko@court"}, {ID: "b", Name: "Oblivionis"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if account.AccountSwitcher(c, &current, accounts, account.AccountSwitcherOptions{}).Add() {
			fmt.Println("sign in another account")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWorkspaceSwitcher() {
	current := "ave"
	spaces := []account.Workspace{{ID: "ave", Name: "Ave Mujica"}, {ID: "mygo", Name: "MyGO!!!!!"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.WorkspaceSwitcher(c, &current, spaces, account.WorkspaceSwitcherOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUserMenu() {
	status := account.UserStatusOnline
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if account.UserMenu(c, &status, account.UserMenuOptions{Name: "Sakiko Togawa"}).SignOut() {
			fmt.Println("sign out")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSessionList() {
	sessions := []account.Session{{ID: "s1", Device: "MacBook Pro", Current: true}, {ID: "s2", Device: "iPhone", Location: "Osaka"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if id, ok := account.SessionList(c, sessions, account.SessionListOptions{}).Revoked(); ok {
			fmt.Println("sign out", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleApiKeyManager() {
	var state account.ApiKeyManagerState
	keys := []account.APIKey{{ID: "k1", Name: "Deploy bot", Prefix: "mk_live_9f2c"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if name, ok := account.ApiKeyManager(c, &state, keys, account.ApiKeyManagerOptions{}).Create(); ok {
			fmt.Println("create", name)
			state.Reveal("secret from your key service")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleUsageQuota() {
	items := []account.QuotaItem{{Label: "API calls", Used: 4200, Limit: 10000}, {Label: "Storage", Used: 9, Limit: 10, Unit: "GB"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.UsageQuota(c, items, account.UsageQuotaOptions{ResetAt: time.Now().Add(72 * time.Hour)})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOnboardingWizard() {
	step := 0
	steps := []account.WizardStep{{Title: "Welcome"}, {Title: "Pick a theme"}, {Title: "Done"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if account.OnboardingWizard(c, &step, steps, account.OnboardingWizardOptions{}).Finished() {
			fmt.Println("onboarded")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCoachmark() {
	open, step := true, 0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		compose := input.Button(c, "Compose", input.ButtonOptions{})
		account.Coachmark(c, &open, &step, []account.CoachmarkStep{{Anchor: compose, Title: "Write here", Body: "Start a new letter."}})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWhatsNewDialog() {
	open, hide := true, false
	releases := []account.Release{{Version: "1.1", Changes: []account.ReleaseChange{{Kind: account.ChangeNew, Text: "500 components."}}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		account.WhatsNewDialog(c, &open, &hide, releases, account.WhatsNewDialogOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleFeedbackWidget() {
	var state account.FeedbackWidgetState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if f, ok := account.FeedbackWidget(c, &state, account.FeedbackWidgetOptions{}).Submitted(); ok {
			fmt.Println("send", f.Kind)
			state.EndSubmit("")
		}
	}
	ui.Render(view, 400, 300, 1)
}
