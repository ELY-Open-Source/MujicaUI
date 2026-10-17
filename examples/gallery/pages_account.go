package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/account"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

const accountCategory = "Settings & Account"

var accountPages = []page{
	{ID: "461", Name: "SettingsLayout", ZhName: "设置布局", Category: accountCategory, Code: accountCode["SettingsLayout"], Demo: accountSettingsDemo,
		Keys: "Arrows move in the navigation; Tab enters the page; clicking a row's text focuses its control."},
	{ID: "462", Name: "SettingsSearch", ZhName: "设置搜索", Category: accountCategory, Code: accountCode["SettingsSearch"], Demo: accountSearchDemo,
		Keys: "Type to search; Up and Down move; Enter jumps to the setting and focuses it."},
	{ID: "463", Name: "ThemeSelector", ZhName: "主题选择", Category: accountCategory, Code: accountCode["ThemeSelector"], Demo: accountThemeDemo,
		Keys: "Left and Right move between cards; Space picks."},
	{ID: "464", Name: "AccentColorPicker", ZhName: "主色选择", Category: accountCategory, Code: accountCode["AccentColorPicker"], Demo: accountAccentDemo,
		Keys: "Left and Right move between swatches; Space picks."},
	{ID: "465", Name: "ShortcutRecorder", ZhName: "快捷键录制", Category: accountCategory, Code: accountCode["ShortcutRecorder"], Demo: accountRecorderDemo,
		Keys: "Enter or Space records; press the combination; Escape cancels."},
	{ID: "466", Name: "KeyboardShortcutsList", ZhName: "快捷键列表", Category: accountCategory, Code: accountCode["KeyboardShortcutsList"], Demo: accountShortcutsDemo,
		Keys: "Type in the search field to filter by action or key."},
	{ID: "467", Name: "LoginForm", ZhName: "登录表单", Category: accountCategory, Code: accountCode["LoginForm"], Demo: accountLoginDemo,
		Keys: "Enter in any field signs in; an empty field takes the focus."},
	{ID: "468", Name: "OAuthButtons", ZhName: "第三方登录按钮", Category: accountCategory, Code: accountCode["OAuthButtons"], Demo: accountOAuthDemo,
		Keys: "Tab between providers; Enter or Space continues."},
	{ID: "469", Name: "TwoFactorInput", ZhName: "两步验证码", Category: accountCategory, Code: accountCode["TwoFactorInput"], Demo: accountTwoFactorDemo,
		Keys: "Type or paste digits; Backspace clears; Tab reaches Resend."},
	{ID: "470", Name: "ProfileEditor", ZhName: "个人资料编辑", Category: accountCategory, Code: accountCode["ProfileEditor"], Demo: accountProfileDemo,
		Keys: "Tab through picture, name, bio, Discard and Save."},
	{ID: "471", Name: "AccountSwitcher", ZhName: "账户切换", Category: accountCategory, Code: accountCode["AccountSwitcher"], Demo: accountSwitcherDemo,
		Keys: "Enter, Space or Down opens; arrows move; Enter chooses; Escape closes."},
	{ID: "472", Name: "WorkspaceSwitcher", ZhName: "工作区切换", Category: accountCategory, Code: accountCode["WorkspaceSwitcher"], Demo: accountWorkspaceDemo,
		Keys: "Open, type to search, Up and Down move, Enter switches, Escape closes."},
	{ID: "473", Name: "UserMenu", ZhName: "用户菜单", Category: accountCategory, Code: accountCode["UserMenu"], Demo: accountUserMenuDemo,
		Keys: "Enter, Space or Down opens; arrows move; Enter chooses."},
	{ID: "474", Name: "SessionList", ZhName: "登录会话", Category: accountCategory, Code: accountCode["SessionList"], Demo: accountSessionsDemo,
		Keys: "Tab to each Sign out; the confirmation focuses Cancel."},
	{ID: "475", Name: "ApiKeyManager", ZhName: "API 密钥管理", Category: accountCategory, Code: accountCode["ApiKeyManager"], Demo: accountKeysDemo,
		Keys: "Enter in the name field creates; Copy key copies; Revoke asks first."},
	{ID: "476", Name: "UsageQuota", ZhName: "用量配额", Category: accountCategory, Code: accountCode["UsageQuota"], Demo: accountQuotaDemo,
		Keys: "Display only; each bar reads as a meter."},
	{ID: "477", Name: "OnboardingWizard", ZhName: "引导向导", Category: accountCategory, Code: accountCode["OnboardingWizard"], Demo: accountWizardDemo,
		Keys: "Next has the focus; Tab reaches Back and Skip."},
	{ID: "478", Name: "Coachmark", ZhName: "功能提示", Category: accountCategory, Code: accountCode["Coachmark"], Demo: accountCoachDemo,
		Keys: "Next has the focus; Escape closes the tour."},
	{ID: "479", Name: "WhatsNewDialog", ZhName: "更新内容", Category: accountCategory, Code: accountCode["WhatsNewDialog"], Demo: accountWhatsNewDemo,
		Keys: "Got it has the focus; Escape closes."},
	{ID: "480", Name: "FeedbackWidget", ZhName: "反馈组件", Category: accountCategory, Code: accountCode["FeedbackWidget"], Demo: accountFeedbackDemo,
		Keys: "Arrows switch the type; Tab reaches the message and buttons."},
}

type accountSettingsState struct {
	page, query           string
	sound, digest, report bool
}

var accountSettingPages = []account.SettingsPage{{ID: "general", Label: "General"}, {ID: "privacy", Label: "Privacy"}}

var accountSettingEntries = []account.SettingEntry{
	{ID: "sound", Page: "general", Label: "Sound", Description: "Chime on new messages.", Keywords: []string{"audio", "notification"}},
	{ID: "digest", Page: "general", Label: "Weekly digest", Description: "A summary every Monday.", Keywords: []string{"email"}},
	{ID: "report", Page: "privacy", Label: "Usage reports", Description: "Send anonymous usage."},
}

// accountSettingsPage builds the demo pages shared by 461 and 462.
func accountSettingsPage(c *ui.Context, d *accountSettingsState, log func(string)) {
	ui.Box(c).Width(900).Height(320).Children(func() {
		if account.SettingsLayout(c, &d.page, accountSettingPages, account.SettingsLayoutOptions{}, func(id string) {
			switch id {
			case "general":
				account.SettingsSection(c, "Notifications", func() {
					accountSwitchRow(c, "sound", "Sound", "Chime on new messages.", &d.sound, log)
					accountSwitchRow(c, "digest", "Weekly digest", "A summary every Monday.", &d.digest, log)
				})
			case "privacy":
				account.SettingsSection(c, "Data", func() {
					accountSwitchRow(c, "report", "Usage reports", "Send anonymous usage.", &d.report, log)
				})
			}
		}).Changed() {
			log("page " + d.page)
		}
	})
}

// accountSwitchRow is a settings row with a switch.
func accountSwitchRow(c *ui.Context, id, label, desc string, v *bool, log func(string)) {
	account.SettingsRow(c, label, account.SettingsRowOptions{ID: id, Description: desc}, func() ui.Element {
		sw := input.Switch(c, v, "", input.SwitchOptions{})
		if sw.Changed() {
			log(fmt.Sprintf("%s %v", label, *v))
		}
		return sw
	})
}

func accountSettingsDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "461", func() accountSettingsState { return accountSettingsState{sound: true} })
	accountSettingsPage(c, d, log)
}

func accountSearchDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "462", func() accountSettingsState { return accountSettingsState{sound: true} })
	ui.Box(c).Width(320).Children(func() {
		if e, ok := account.SettingsSearch(c, &d.query, accountSettingEntries, account.SettingsSearchOptions{}).Chosen(); ok {
			d.page = e.Page
			log("jump to " + e.ID)
		}
	})
	accountSettingsPage(c, d, log)
}

type accountThemeState struct {
	mode     core.Mode
	disabled bool
}

func accountThemeDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "463", func() accountThemeState { return accountThemeState{} })
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	if account.ThemeSelector(c, &d.mode, account.ThemeSelectorOptions{Disabled: d.disabled}).Changed() {
		log("theme " + d.mode.String())
	}
}

type accountAccentState struct {
	color    ui.Color
	disabled bool
}

func accountAccentDemo(c *ui.Context, log func(string)) {
	d := ui.Local(c.Root(), "464", func() accountAccentState { return accountAccentState{color: ui.Hex("#7D2034")} })
	accountToggles(c, log, []string{"Disabled"}, &d.disabled)
	if account.AccentColorPicker(c, &d.color, account.AccentColorPickerOptions{Disabled: d.disabled}).Changed() {
		log(fmt.Sprintf("accent %v", d.color))
	}
	input.Button(c, "Preview button", input.ButtonOptions{}).Background(d.color)
}

func accountToggles(c *ui.Context, log func(string), labels []string, values ...*bool) {
	before := make([]bool, len(values))
	for i, v := range values {
		before[i] = *v
	}
	formToggles(c, labels, values...)
	for i, v := range values {
		if before[i] != *v {
			log(fmt.Sprintf("%s %v", labels[i], *v))
		}
	}
}
