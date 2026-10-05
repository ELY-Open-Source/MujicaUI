package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

const settingsID = "c-settings"

// settingsValues is one full set of preferences; drafts compare to saved ones with ==.
type settingsValues struct {
	workspace string
	language  string
	startPage string
	launch    bool
	theme     string
	fontSize  float64
	sidebar   float64
	reduce    bool
	telemetry bool
	crashes   bool
	history   float64
	lock      string
}

var settingsInitial = settingsValues{
	workspace: "Mujica", language: "en", startPage: "last", launch: true,
	theme: "system", fontSize: 14, sidebar: 240,
	crashes: true, history: 30, lock: "15",
}

type settingsDemo struct {
	tab          int
	saved, draft settingsValues
}

func newSettingsDemo() settingsDemo {
	return settingsDemo{saved: settingsInitial, draft: settingsInitial}
}

func (g *gallery) settingsView(c *ui.Context, log func(string)) {
	s := &g.cSettings
	mujica.Card(c, mujica.CardOptions{Title: g.t("Preferences", "偏好设置"), Footer: func() { g.settingsFooter(c, log) }}, func() {
		tabs := []mujica.Tab{
			{Label: g.t("General", "通用"), Icon: icons.Must("settings"), Panel: func() { g.settingsGeneral(c, log) }},
			{Label: g.t("Appearance", "外观"), Icon: icons.Must("palette"), Panel: func() { g.settingsAppearance(c, log) }},
			{Label: g.t("Privacy", "隐私"), Icon: icons.Must("lock"), Panel: func() { g.settingsPrivacy(c, log) }},
		}
		if mujica.Tabs(c, &s.tab, tabs, mujica.TabsOptions{Label: g.t("Preference groups", "设置分组")}).Changed() {
			log("tab " + tabs[s.tab].Label)
		}
	})
}

// settingsPanel lays one tab's fields in a readable column.
func settingsPanel(c *ui.Context, fields func()) {
	ui.Column(c).FillWidth().MaxWidth(560).Gap(20).Padding(16, 0, 4, 0).Children(fields)
}

// workspaceError checks the one free-text preference.
func (g *gallery) workspaceError(name string) string {
	if strings.TrimSpace(name) == "" {
		return g.t("Name the workspace.", "请为工作区命名。")
	}
	return ""
}

func (g *gallery) settingsGeneral(c *ui.Context, log func(string)) {
	d := &g.cSettings.draft
	settingsPanel(c, func() {
		bad := g.workspaceError(d.workspace)
		mujica.FormField(c, g.t("Workspace name", "工作区名称"), mujica.FormFieldOptions{Required: true, Error: bad,
			Description: g.t("Shown in the title bar.", "显示在标题栏中。")}, func() *ui.Element {
			r := mujica.TextInput(c, &d.workspace, mujica.TextInputOptions{Error: bad})
			if r.Changed() {
				log(fmt.Sprintf("workspace %q", d.workspace))
			}
			return r.Input
		})
		mujica.FormField(c, g.t("Language", "语言"), mujica.FormFieldOptions{}, func() *ui.Element {
			r := mujica.Select(c, &d.language, []mujica.SelectOption[string]{
				{Value: "en", Label: "English"}, {Value: "zh", Label: "简体中文"}, {Value: "ja", Label: "日本語"},
			}, mujica.SelectOptions{})
			if r.Changed() {
				log("language " + d.language)
			}
			return r.Element
		})
		mujica.FormField(c, g.t("On startup", "启动时"), mujica.FormFieldOptions{}, func() *ui.Element {
			r := mujica.RadioGroup(c, &d.startPage, []mujica.RadioItem[string]{
				{Value: "blank", Label: g.t("Open a blank page", "打开空白页")},
				{Value: "last", Label: g.t("Reopen the last session", "恢复上次会话")},
				{Value: "home", Label: g.t("Open the home page", "打开主页")},
			}, mujica.RadioGroupOptions{})
			if r.Changed() {
				log("startup " + d.startPage)
			}
			return r.Element
		})
		mujica.FormField(c, g.t("Login item", "登录项"), mujica.FormFieldOptions{}, func() *ui.Element {
			e := mujica.Switch(c, &d.launch, g.t("Launch at login", "登录时启动"), mujica.SwitchOptions{})
			if e.Changed() {
				log(fmt.Sprintf("launch at login %v", d.launch))
			}
			return e
		})
	})
}

func (g *gallery) settingsAppearance(c *ui.Context, log func(string)) {
	d := &g.cSettings.draft
	settingsPanel(c, func() {
		mujica.FormField(c, g.t("Theme", "主题"), mujica.FormFieldOptions{}, func() *ui.Element {
			r := mujica.RadioGroup(c, &d.theme, []mujica.RadioItem[string]{
				{Value: "system", Label: g.t("System", "跟随系统")}, {Value: "light", Label: g.t("Light", "浅色")}, {Value: "dark", Label: g.t("Dark", "深色")},
			}, mujica.RadioGroupOptions{Horizontal: true})
			if r.Changed() {
				log("theme " + d.theme)
			}
			return r.Element
		})
		mujica.FormField(c, g.t("Text size", "文字大小"), mujica.FormFieldOptions{}, func() *ui.Element {
			r := mujica.Slider(c, &d.fontSize, mujica.SliderOptions{Min: 12, Max: 20, Step: 1, ShowValue: true,
				Marks:  []mujica.SliderMark{{Value: 12, Label: "12"}, {Value: 16, Label: "16"}, {Value: 20, Label: "20"}},
				Format: func(v float64) string { return fmt.Sprintf("%.0f pt", v) }})
			if r.Changed() {
				log(fmt.Sprintf("text size %.0f", d.fontSize))
			}
			return r.Element
		})
		mujica.FormField(c, g.t("Sidebar width", "侧栏宽度"), mujica.FormFieldOptions{Description: g.t("In DIPs, from 180 to 360.", "单位 DIP，180 至 360。")}, func() *ui.Element {
			r := mujica.NumberInput(c, &d.sidebar, mujica.NumberInputOptions{Min: 180, Max: 360, Step: 10})
			if r.Changed() {
				log(fmt.Sprintf("sidebar width %.0f", d.sidebar))
			}
			return r.Input
		})
		mujica.FormField(c, g.t("Motion", "动效"), mujica.FormFieldOptions{}, func() *ui.Element {
			e := mujica.Switch(c, &d.reduce, g.t("Reduce motion", "减少动效"), mujica.SwitchOptions{})
			if e.Changed() {
				log(fmt.Sprintf("reduce motion %v", d.reduce))
			}
			return e
		})
	})
}

func (g *gallery) settingsPrivacy(c *ui.Context, log func(string)) {
	d := &g.cSettings.draft
	settingsPanel(c, func() {
		mujica.FormField(c, g.t("Usage data", "使用数据"), mujica.FormFieldOptions{Description: g.t("Nothing leaves this demo.", "本示例不会发送任何数据。")}, func() *ui.Element {
			e := mujica.Switch(c, &d.telemetry, g.t("Share anonymous usage data", "共享匿名使用数据"), mujica.SwitchOptions{})
			if e.Changed() {
				log(fmt.Sprintf("usage data %v", d.telemetry))
			}
			return e
		})
		mujica.FormField(c, g.t("Crash reports", "崩溃报告"), mujica.FormFieldOptions{}, func() *ui.Element {
			e := mujica.Switch(c, &d.crashes, g.t("Send crash reports", "发送崩溃报告"), mujica.SwitchOptions{})
			if e.Changed() {
				log(fmt.Sprintf("crash reports %v", d.crashes))
			}
			return e
		})
		mujica.FormField(c, g.t("Keep history", "保留历史"), mujica.FormFieldOptions{Description: g.t("Days, from 0 to 365; 0 keeps none.", "天数，0 至 365；0 表示不保留。")}, func() *ui.Element {
			r := mujica.NumberInput(c, &d.history, mujica.NumberInputOptions{Min: 0, Max: 365, Step: 1})
			if r.Changed() {
				log(fmt.Sprintf("history %.0f days", d.history))
			}
			return r.Input
		})
		mujica.FormField(c, g.t("Lock after", "自动锁定"), mujica.FormFieldOptions{}, func() *ui.Element {
			r := mujica.Select(c, &d.lock, []mujica.SelectOption[string]{
				{Value: "never", Label: g.t("Never", "从不")}, {Value: "5", Label: g.t("5 minutes", "5 分钟")},
				{Value: "15", Label: g.t("15 minutes", "15 分钟")}, {Value: "60", Label: g.t("1 hour", "1 小时")},
			}, mujica.SelectOptions{})
			if r.Changed() {
				log("lock after " + d.lock)
			}
			return r.Element
		})
	})
}

// settingsFooter flags unsaved changes and holds Reset and Save.
func (g *gallery) settingsFooter(c *ui.Context, log func(string)) {
	s := &g.cSettings
	dirty := s.draft != s.saved
	invalid := g.workspaceError(s.draft.workspace) != ""
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(12).Children(func() {
		if dirty {
			mujica.StatusIndicator(c, mujica.StatusWarning, mujica.StatusIndicatorOptions{Label: g.t("Unsaved changes", "有未保存的更改")})
		} else {
			mujica.StatusIndicator(c, mujica.StatusSuccess, mujica.StatusIndicatorOptions{Label: g.t("All changes saved", "更改均已保存")})
		}
		ui.Box(c).Grow(1)
		if mujica.Button(c, g.t("Reset", "重置"), mujica.ButtonOptions{Variant: mujica.Secondary, Disabled: s.draft == settingsInitial}).Clicked() {
			s.draft = settingsInitial
			log("reset to initial values")
		}
		if mujica.Button(c, g.t("Save", "保存"), mujica.ButtonOptions{Disabled: !dirty || invalid}).Clicked() {
			s.saved = s.draft
			log("saved")
			mujica.ShowToast(c, mujica.ToastMessage{ID: "settings-saved", Severity: mujica.SeveritySuccess, Message: g.t("Settings saved", "设置已保存")})
		}
	})
}
