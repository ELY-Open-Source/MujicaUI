package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/feedback"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
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
	data.Card(c, data.CardOptions{Title: g.t("Preferences", "偏好设置"), Footer: func() { g.settingsFooter(c, log) }}, func() {
		tabs := []navigation.Tab{
			{Label: g.t("General", "通用"), Icon: icons.Must("settings"), Panel: func() { g.settingsGeneral(c, log) }},
			{Label: g.t("Appearance", "外观"), Icon: icons.Must("palette"), Panel: func() { g.settingsAppearance(c, log) }},
			{Label: g.t("Privacy", "隐私"), Icon: icons.Must("lock"), Panel: func() { g.settingsPrivacy(c, log) }},
		}
		if navigation.Tabs(c, &s.tab, tabs, navigation.TabsOptions{Label: g.t("Preference groups", "设置分组")}).Changed() {
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
		input.FormField(c, g.t("Workspace name", "工作区名称"), input.FormFieldOptions{Required: true, Error: bad,
			Description: g.t("Shown in the title bar.", "显示在标题栏中。")}, func() *ui.Element {
			r := input.TextInput(c, &d.workspace, input.TextInputOptions{Error: bad})
			if r.Changed() {
				log(fmt.Sprintf("workspace %q", d.workspace))
			}
			return r.Input
		})
		input.FormField(c, g.t("Language", "语言"), input.FormFieldOptions{}, func() *ui.Element {
			r := input.Select(c, &d.language, []input.SelectOption[string]{
				{Value: "en", Label: "English"}, {Value: "zh", Label: "简体中文"}, {Value: "ja", Label: "日本語"},
			}, input.SelectOptions{})
			if r.Changed() {
				log("language " + d.language)
			}
			return r.Element
		})
		input.FormField(c, g.t("On startup", "启动时"), input.FormFieldOptions{}, func() *ui.Element {
			r := input.RadioGroup(c, &d.startPage, []input.RadioItem[string]{
				{Value: "blank", Label: g.t("Open a blank page", "打开空白页")},
				{Value: "last", Label: g.t("Reopen the last session", "恢复上次会话")},
				{Value: "home", Label: g.t("Open the home page", "打开主页")},
			}, input.RadioGroupOptions{})
			if r.Changed() {
				log("startup " + d.startPage)
			}
			return r.Element
		})
		input.FormField(c, g.t("Login item", "登录项"), input.FormFieldOptions{}, func() *ui.Element {
			e := input.Switch(c, &d.launch, g.t("Launch at login", "登录时启动"), input.SwitchOptions{})
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
		input.FormField(c, g.t("Theme", "主题"), input.FormFieldOptions{}, func() *ui.Element {
			r := input.RadioGroup(c, &d.theme, []input.RadioItem[string]{
				{Value: "system", Label: g.t("System", "跟随系统")}, {Value: "light", Label: g.t("Light", "浅色")}, {Value: "dark", Label: g.t("Dark", "深色")},
			}, input.RadioGroupOptions{Horizontal: true})
			if r.Changed() {
				log("theme " + d.theme)
			}
			return r.Element
		})
		input.FormField(c, g.t("Text size", "文字大小"), input.FormFieldOptions{}, func() *ui.Element {
			r := input.Slider(c, &d.fontSize, input.SliderOptions{Min: 12, Max: 20, Step: 1, ShowValue: true,
				Marks:  []input.SliderMark{{Value: 12, Label: "12"}, {Value: 16, Label: "16"}, {Value: 20, Label: "20"}},
				Format: func(v float64) string { return fmt.Sprintf("%.0f pt", v) }})
			if r.Changed() {
				log(fmt.Sprintf("text size %.0f", d.fontSize))
			}
			return r.Element
		})
		input.FormField(c, g.t("Sidebar width", "侧栏宽度"), input.FormFieldOptions{Description: g.t("In DIPs, from 180 to 360.", "单位 DIP，180 至 360。")}, func() *ui.Element {
			r := input.NumberInput(c, &d.sidebar, input.NumberInputOptions{Min: 180, Max: 360, Step: 10})
			if r.Changed() {
				log(fmt.Sprintf("sidebar width %.0f", d.sidebar))
			}
			return r.Input
		})
		input.FormField(c, g.t("Motion", "动效"), input.FormFieldOptions{}, func() *ui.Element {
			e := input.Switch(c, &d.reduce, g.t("Reduce motion", "减少动效"), input.SwitchOptions{})
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
		input.FormField(c, g.t("Usage data", "使用数据"), input.FormFieldOptions{Description: g.t("Nothing leaves this demo.", "本示例不会发送任何数据。")}, func() *ui.Element {
			e := input.Switch(c, &d.telemetry, g.t("Share anonymous usage data", "共享匿名使用数据"), input.SwitchOptions{})
			if e.Changed() {
				log(fmt.Sprintf("usage data %v", d.telemetry))
			}
			return e
		})
		input.FormField(c, g.t("Crash reports", "崩溃报告"), input.FormFieldOptions{}, func() *ui.Element {
			e := input.Switch(c, &d.crashes, g.t("Send crash reports", "发送崩溃报告"), input.SwitchOptions{})
			if e.Changed() {
				log(fmt.Sprintf("crash reports %v", d.crashes))
			}
			return e
		})
		input.FormField(c, g.t("Keep history", "保留历史"), input.FormFieldOptions{Description: g.t("Days, from 0 to 365; 0 keeps none.", "天数，0 至 365；0 表示不保留。")}, func() *ui.Element {
			r := input.NumberInput(c, &d.history, input.NumberInputOptions{Min: 0, Max: 365, Step: 1})
			if r.Changed() {
				log(fmt.Sprintf("history %.0f days", d.history))
			}
			return r.Input
		})
		input.FormField(c, g.t("Lock after", "自动锁定"), input.FormFieldOptions{}, func() *ui.Element {
			r := input.Select(c, &d.lock, []input.SelectOption[string]{
				{Value: "never", Label: g.t("Never", "从不")}, {Value: "5", Label: g.t("5 minutes", "5 分钟")},
				{Value: "15", Label: g.t("15 minutes", "15 分钟")}, {Value: "60", Label: g.t("1 hour", "1 小时")},
			}, input.SelectOptions{})
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
			feedback.StatusIndicator(c, feedback.StatusWarning, feedback.StatusIndicatorOptions{Label: g.t("Unsaved changes", "有未保存的更改")})
		} else {
			feedback.StatusIndicator(c, feedback.StatusSuccess, feedback.StatusIndicatorOptions{Label: g.t("All changes saved", "更改均已保存")})
		}
		ui.Box(c).Grow(1)
		if input.Button(c, g.t("Reset", "重置"), input.ButtonOptions{Variant: input.Secondary, Disabled: s.draft == settingsInitial}).Clicked() {
			s.draft = settingsInitial
			log("reset to initial values")
		}
		if input.Button(c, g.t("Save", "保存"), input.ButtonOptions{Disabled: !dirty || invalid}).Clicked() {
			s.saved = s.draft
			log("saved")
			feedback.ShowToast(c, feedback.ToastMessage{ID: "settings-saved", Severity: core.SeveritySuccess, Message: g.t("Settings saved", "设置已保存")})
		}
	})
}
