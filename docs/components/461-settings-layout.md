# 461 SettingsLayout

设置布局。

## 用途

设置页框架：左侧分区导航（Sidebar），右侧显示当前分区；配套 `SettingsSection` 分节与 `SettingsRow` 设置行（标签、说明、控件）。

## 最小示例

```go
page := ""
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
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `page` | `*string` | 必填 | 当前分区 ID；空值取第一个，未知 ID panic。 |
| `pages` | `[]SettingsPage` | 必填 | ID、Label、Icon。 |
| `Title` | `string` | 内置“Settings” | 导航标题。 |
| `content` | `func(id string)` | 必填 | 构建分区内容。 |
| `SettingsRow control` | `func() *ui.Element` | 必填 | 返回控件；行标签与说明关联到它。 |
| `SettingsRowOptions.ID` | `string` | "" | 供 SettingsSearch 跳转定位。 |

## 状态

导航选中、行悬停、被搜索定位的行短暂高亮（减少动画时不高亮）。

## 事件

`ChoiceResult.Changed()` 在切换分区时为 true。

## 键盘操作

导航沿用 Sidebar 键盘操作；Tab 进入内容区各控件；点击行文字聚焦其控件。

## 限制

- 内容区独立滚动，最大宽度 720 DIP。
