package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

var displayPages = []page{
	{ID: "011", Name: "Text", ZhName: "文本", Category: "Display", Demo: textDemo, Code: textCode,
		Keys: "Selectable 文本：点击后 Cmd/Ctrl+A 全选、Cmd/Ctrl+C 复制。可编辑：Tab 聚焦，Enter 或双击开始编辑，Enter 确认，Escape 取消，移走焦点也确认。"},
	{ID: "012", Name: "Heading", ZhName: "标题", Category: "Display", Demo: headingDemo, Code: headingCode,
		Keys: "纯展示，不获取焦点。"},
	{ID: "013", Name: "Icon", ZhName: "图标", Category: "Display", Demo: iconDemo, Code: iconCode,
		Keys: "纯展示，不获取焦点。"},
	{ID: "014", Name: "Image", ZhName: "图片", Category: "Display", Demo: imageDemo, Code: imageCode,
		Keys: "错误状态的「重试」按钮可 Tab 聚焦，Enter 或 Space 触发。"},
	{ID: "015", Name: "Avatar", ZhName: "头像", Category: "Display", Demo: avatarDemo, Code: avatarCode,
		Keys: "纯展示，不获取焦点。"},
	{ID: "016", Name: "AvatarGroup", ZhName: "头像组", Category: "Display", Demo: avatarGroupDemo, Code: avatarGroupCode,
		Keys: "Tab 聚焦「+N」按钮，Enter 或 Space 展开完整名单，Escape 或点击外部关闭。"},
	{ID: "017", Name: "Badge", ZhName: "徽标", Category: "Display", Demo: badgeDemo, Code: badgeCode,
		Keys: "纯展示，不获取焦点；辅助技术朗读完整数值。"},
	{ID: "018", Name: "Tag", ZhName: "标签", Category: "Display", Demo: tagDemo, Code: tagCode,
		Keys: "可点击标签：Tab 聚焦，Enter 或 Space 点击，Delete 或 Backspace 关闭。关闭按钮是独立的 Tab 停靠点。"},
	{ID: "019", Name: "Link", ZhName: "链接", Category: "Display", Demo: linkDemo, Code: linkCode,
		Keys: "Tab 聚焦链接，Enter 激活。"},
	{ID: "020", Name: "Kbd", ZhName: "快捷键", Category: "Display", Demo: kbdDemo, Code: kbdCode,
		Keys: "纯展示，不注册快捷键，不获取焦点。"},
}

type textState struct {
	maxLines   float64
	selectable bool
	name       string
	saving     bool
	failNext   bool
	saveErr    string
}

func textDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "011", func() textState { return textState{maxLines: 0, selectable: true, name: "Coronation.md"} })
	mujica.Heading(c, "Text", mujica.HeadingOptions{Subtitle: "换行、省略、选择复制与可编辑模式", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			ui.Text(c, fmt.Sprintf("MaxLines %d", int(s.maxLines)))
			displayLog(log, ui.StepSlider(c, &s.maxLines, 0, 4, 1).Width(160), "maxLines", &s.maxLines)
			displayLog(log, ui.Checkbox(c, &s.selectable, "Selectable"), "selectable", &s.selectable)
			displayLog(log, ui.Checkbox(c, &s.saving, "Saving"), "saving", &s.saving)
			displayLog(log, ui.Checkbox(c, &s.failNext, "Next save fails"), "failNext", &s.failNext)
		})
	})
	displaySection(c, "Display", func() {
		ui.Box(c).Width(360).Children(func() {
			mujica.Text(c, "Velvet curtains fall over the long gallery. Every candle leans toward the open door, and the quartet plays on while the guests take their places beneath the chandeliers.",
				mujica.TextOptions{MaxLines: int(s.maxLines), Selectable: s.selectable})
		})
	})
	displaySection(c, "Editable", func() {
		e := mujica.EditableText(c, &s.name, mujica.EditableTextOptions{Saving: s.saving, Error: s.saveErr, Label: "File name"})
		if e.Changed {
			if s.failNext {
				s.saveErr = "Save failed: the name is already taken."
				log("save failed: " + s.name)
			} else {
				s.saveErr = ""
				log("saved: " + s.name)
			}
		}
		if e.Canceled {
			log("edit canceled")
		}
	})
}

const textCode = `mujica.Text(c, "Order ID: 4821-AX", mujica.TextOptions{Selectable: true})
mujica.Text(c, "A long note that stops after two lines.", mujica.TextOptions{MaxLines: 2})
if mujica.EditableText(c, &name, mujica.EditableTextOptions{Error: saveErr}).Changed {
	saveErr = "" // save name here; set saveErr on failure
}`

type headingState struct {
	level            float64
	subtitle, divide bool
}

func headingDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "012", func() headingState { return headingState{level: 1, subtitle: true, divide: true} })
	mujica.Heading(c, "Heading", mujica.HeadingOptions{Subtitle: "三级衬线标题、副标题与装饰分隔线", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			ui.Text(c, fmt.Sprintf("Level %d", int(s.level)))
			if ui.StepSlider(c, &s.level, 1, 3, 1).Width(120).Changed() {
				log(fmt.Sprintf("level %d", int(s.level)))
			}
			displayLog(log, ui.Checkbox(c, &s.subtitle, "Subtitle"), "subtitle", &s.subtitle)
			displayLog(log, ui.Checkbox(c, &s.divide, "Divider"), "divide", &s.divide)
		})
	})
	displaySection(c, "Preview", func() {
		sub := ""
		if s.subtitle {
			sub = "Ceremonies of the winter season"
		}
		mujica.Heading(c, "The Grand Hall", mujica.HeadingOptions{Level: int(s.level), Subtitle: sub, Divider: s.divide})
	})
	displaySection(c, "Levels", func() {
		for l := 1; l <= 3; l++ {
			mujica.Heading(c, fmt.Sprintf("Level %d · 宫廷", l), mujica.HeadingOptions{Level: l})
		}
	})
}

const headingCode = `mujica.Heading(c, "The Grand Hall", mujica.HeadingOptions{Subtitle: "Winter ceremonies", Divider: true})
mujica.Heading(c, "Guest Ledger", mujica.HeadingOptions{Level: 2})`

type iconState struct {
	name  string
	size  string
	tone  string
	label bool
}

var iconTones = map[string]mujica.IconTone{"Inherit": mujica.IconInherit, "Muted": mujica.IconMuted, "Accent": mujica.IconAccent,
	"Success": mujica.IconSuccess, "Warning": mujica.IconWarning, "Danger": mujica.IconDanger, "Info": mujica.IconInfo}

func iconDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "013", func() iconState { return iconState{name: "star", size: "24", tone: "Accent"} })
	mujica.Heading(c, "Icon", mujica.HeadingOptions{Subtitle: "按名称取图标；尺寸与语义颜色；缺失时可见回退", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			if ui.Select(c, &s.name, []string{"star", "bell", "settings", "moon", "sun", "no-such-icon"}).Changed() {
				log("name " + s.name)
			}
			if ui.Select(c, &s.size, []string{"16", "20", "24"}).Changed() {
				log("size " + s.size)
			}
			if ui.Select(c, &s.tone, []string{"Inherit", "Muted", "Accent", "Success", "Warning", "Danger", "Info"}).Changed() {
				log("tone " + s.tone)
			}
			displayLog(log, ui.Checkbox(c, &s.label, "Label"), "label", &s.label)
		})
	})
	displaySection(c, "Preview", func() {
		var size float32
		fmt.Sscan(s.size, &size)
		label := ""
		if s.label {
			label = "Preview icon"
		}
		mujica.Icon(c, s.name, mujica.IconOptions{Size: size, Tone: iconTones[s.tone], Label: label})
	})
	displaySection(c, "Tones and fallback", func() {
		ui.Row(c).Gap(14).Children(func() {
			for _, t := range []string{"Inherit", "Muted", "Accent", "Success", "Warning", "Danger", "Info"} {
				mujica.Icon(c, "circle-check", mujica.IconOptions{Size: 20, Tone: iconTones[t]}).Key(t)
			}
			mujica.Icon(c, "missing-name", mujica.IconOptions{Size: 20})
		})
	})
}

const iconCode = `ui.Row(c).Gap(8).Children(func() {
	mujica.Icon(c, "star", mujica.IconOptions{})
	mujica.Icon(c, "circle-check", mujica.IconOptions{Size: 20, Tone: mujica.IconSuccess, Label: "Done"})
})`

var displayPeople = []mujica.AvatarItem{
	{Name: "Sakiko Togawa"}, {Name: "Uika Misumi"}, {Name: "Mutsumi Wakaba"}, {Name: "Umiri Yahata"},
	{Name: "Nyamu Yutenji"}, {Name: "Tomori Takamatsu"}, {Name: "Anon Chihaya"}, {Name: "Rana Kaname"},
}

type avatarGroupState struct{ max, count float64 }

func avatarGroupDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "016", func() avatarGroupState { return avatarGroupState{max: 3, count: 6} })
	mujica.Heading(c, "AvatarGroup", mujica.HeadingOptions{Subtitle: "最大显示数量、剩余数量与完整名单", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			ui.Text(c, fmt.Sprintf("Max %d", int(s.max)))
			if ui.StepSlider(c, &s.max, 1, 8, 1).Width(160).Changed() {
				log(fmt.Sprintf("max %d", int(s.max)))
			}
			ui.Text(c, fmt.Sprintf("People %d", int(s.count)))
			if ui.StepSlider(c, &s.count, 1, 8, 1).Width(160).Changed() {
				log(fmt.Sprintf("people %d", int(s.count)))
			}
		})
	})
	displaySection(c, "Preview", func() {
		for i, g := range []mujica.AvatarGroupView{
			mujica.AvatarGroup(c, displayPeople[:int(s.count)], mujica.AvatarGroupOptions{Max: int(s.max)}),
			mujica.AvatarGroup(c, displayPeople, mujica.AvatarGroupOptions{Max: 4, Size: 24}),
		} {
			if g.Toggled {
				log(fmt.Sprintf("group %d open=%v", i+1, g.Open))
			}
		}
	})
}

const avatarGroupCode = `mujica.AvatarGroup(c, people, mujica.AvatarGroupOptions{Max: 3})`

type badgeState struct {
	count, max float64
	text       string
	tone, pos  string
}

var badgeTones = map[string]mujica.BadgeTone{"Accent": mujica.BadgeAccent, "Neutral": mujica.BadgeNeutral, "Success": mujica.BadgeSuccess,
	"Warning": mujica.BadgeWarning, "Danger": mujica.BadgeDanger, "Info": mujica.BadgeInfo}
var badgePositions = map[string]mujica.BadgePosition{"TopRight": mujica.BadgeTopRight, "TopLeft": mujica.BadgeTopLeft,
	"BottomRight": mujica.BadgeBottomRight, "BottomLeft": mujica.BadgeBottomLeft}

func badgeDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "017", func() badgeState { return badgeState{count: 120, max: 99, tone: "Accent", pos: "TopRight"} })
	mujica.Heading(c, "Badge", mujica.HeadingOptions{Subtitle: "文字或计数、上限、位置与完整朗读", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			ui.Text(c, fmt.Sprintf("Count %d", int(s.count)))
			displayLog(log, ui.StepSlider(c, &s.count, 0, 200, 1).Width(160), "count", &s.count)
			ui.Text(c, fmt.Sprintf("Max %d", int(s.max)))
			displayLog(log, ui.StepSlider(c, &s.max, 9, 999, 10).Width(120), "max", &s.max)
			displayLog(log, ui.TextInput(c, &s.text).Placeholder("Text overrides count").Width(180), "text", &s.text)
			if ui.Select(c, &s.tone, []string{"Accent", "Neutral", "Success", "Warning", "Danger", "Info"}).Changed() {
				log("tone " + s.tone)
			}
			if ui.Select(c, &s.pos, []string{"TopRight", "TopLeft", "BottomRight", "BottomLeft"}).Changed() {
				log("position " + s.pos)
			}
		})
	})
	opts := mujica.BadgeOptions{Text: s.text, Count: int(s.count), Max: int(s.max), Tone: badgeTones[s.tone], Label: "Unread"}
	displaySection(c, "Inline", func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Inbox")
			mujica.Badge(c, opts, nil)
		})
	})
	displaySection(c, "Corner", func() {
		opts.Position = badgePositions[s.pos]
		ui.Box(c).Padding(10).Children(func() {
			mujica.Badge(c, opts, func() { mujica.Avatar(c, "Sakiko Togawa", mujica.AvatarOptions{Size: 40}) })
		})
	})
}

const badgeCode = `mujica.Badge(c, mujica.BadgeOptions{Text: "New", Tone: mujica.BadgeNeutral}, nil)
mujica.Badge(c, mujica.BadgeOptions{Count: unread, Label: "Unread", Position: mujica.BadgeTopRight}, func() {
	mujica.Avatar(c, "Sakiko Togawa", mujica.AvatarOptions{})
})`

type tagState struct {
	tags                           []string
	icon, clickable, closable, dis bool
}

func tagDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "018", func() tagState {
		return tagState{tags: []string{"Baroque", "Gothic", "Rococo", "Court"}, icon: true, clickable: true, closable: true}
	})
	mujica.Heading(c, "Tag", mujica.HeadingOptions{Subtitle: "标签文字、图标；点击与关闭分离", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			displayLog(log, ui.Checkbox(c, &s.icon, "Icon"), "icon", &s.icon)
			displayLog(log, ui.Checkbox(c, &s.clickable, "Clickable"), "clickable", &s.clickable)
			displayLog(log, ui.Checkbox(c, &s.closable, "Closable"), "closable", &s.closable)
			displayLog(log, ui.Checkbox(c, &s.dis, "Disabled"), "dis", &s.dis)
			if mujica.Button(c, "Reset tags", mujica.ButtonOptions{Variant: mujica.Ghost}).Clicked() {
				s.tags = []string{"Baroque", "Gothic", "Rococo", "Court"}
				log("reset")
			}
		})
	})
	displaySection(c, "Preview", func() {
		icon := ""
		if s.icon {
			icon = "star"
		}
		ui.Row(c).Gap(8).Wrap().Children(func() {
			for i := 0; i < len(s.tags); i++ {
				t := s.tags[i]
				var v mujica.TagView
				ui.Box(c).Key(t).Children(func() {
					v = mujica.Tag(c, t, mujica.TagOptions{Icon: icon, Clickable: s.clickable, Closable: s.closable, Disabled: s.dis})
				})
				if v.Clicked {
					log("clicked " + t)
				}
				if v.Closed {
					log("closed " + t)
					s.tags = append(s.tags[:i:i], s.tags[i+1:]...)
					i--
				}
			}
		})
	})
}

const tagCode = `for i, t := range tags {
	var v mujica.TagView
	ui.Box(c).Key(t).Children(func() {
		v = mujica.Tag(c, t, mujica.TagOptions{Icon: "star", Clickable: true, Closable: true})
	})
	if v.Closed {
		tags = append(tags[:i:i], tags[i+1:]...)
	}
}`

type linkState struct {
	url      string
	external string
}

func linkDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "019", func() linkState { return linkState{url: "https://example.com/charter"} })
	mujica.Heading(c, "Link", mujica.HeadingOptions{Subtitle: "外部地址或应用内动作；键盘激活；打开失败反馈", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			displayLog(log, ui.TextInput(c, &s.url).Placeholder("URL (empty = in-app action)").Width(320), "url", &s.url)
			if mujica.Button(c, "Simulate open failure", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
				s.external = "No application can open this link."
				log("open failure set")
			}
			if mujica.Button(c, "Clear error", mujica.ButtonOptions{Variant: mujica.Ghost}).Clicked() {
				s.external = ""
				log("open failure cleared")
			}
		})
	})
	displaySection(c, "Preview", func() {
		if mujica.Link(c, "Read the court charter", mujica.LinkOptions{URL: s.url, Error: s.external}).Clicked() {
			log("activated " + s.url)
		}
		if mujica.Link(c, "Show the guest list", mujica.LinkOptions{}).Clicked() {
			log("in-app action")
		}
	})
}

const linkCode = `mujica.Link(c, "Read the charter", mujica.LinkOptions{URL: "https://example.com/charter"})
if mujica.Link(c, "Show details", mujica.LinkOptions{}).Clicked() {
	details = !details
}`

type kbdState struct{ combo string }

func kbdDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "020", func() kbdState { return kbdState{combo: "Mod+Shift+P"} })
	mujica.Heading(c, "Kbd", mujica.HeadingOptions{Subtitle: "单键与组合键；按平台显示 Command 或 Ctrl", Divider: true})
	displaySection(c, "Parameters", func() {
		displayControls(c, func() {
			if ui.TextInput(c, &s.combo).Placeholder("Keys joined by +").Width(220).Changed() {
				log("keys " + s.combo)
			}
		})
	})
	displaySection(c, "Preview", func() {
		var keys []string
		for _, k := range strings.Split(s.combo, "+") {
			if k = strings.TrimSpace(k); k != "" {
				keys = append(keys, k)
			}
		}
		if len(keys) > 0 {
			mujica.Kbd(c, keys...)
		}
	})
	displaySection(c, "Common", func() {
		for _, row := range [][]string{{"Mod", "K"}, {"Mod", "Shift", "P"}, {"Alt", "Up"}, {"Esc"}, {"Ctrl", "Tab"}} {
			ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
				mujica.Kbd(c, row...)
				ui.Text(c, strings.Join(row, " + ")).TextColor(mujica.Tokens(c).TextMuted)
			})
		}
	})
}

const kbdCode = `ui.Row(c).Gap(6).Children(func() {
	mujica.Text(c, "Command palette", mujica.TextOptions{})
	mujica.Kbd(c, "Mod", "Shift", "P")
})`
