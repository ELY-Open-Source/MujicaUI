# 465 ShortcutRecorder

快捷键录制。

## 用途

点击后录制下一次按下的组合键；与已占用的快捷键冲突时提示并保留原值；按平台显示修饰键。

## 最小示例

```go
palette := mujica.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ShortcutRecorder(c, &palette, mujica.ShortcutRecorderOptions{
		Label: "Command palette",
		Taken: []mujica.ShortcutBinding{{Action: "Save", Shortcut: mujica.Shortcut{Mods: ui.Cmd, Key: ui.KeyS}}},
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*Shortcut` | 必填 | Mods 与 Key；零值表示未设置。 |
| `Label` | `string` | 必填 | 动作名称。 |
| `Taken` | `[]ShortcutBinding` | nil | 已占用的快捷键。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

显示、录制中、冲突提示、缺少修饰键提示、禁用。

## 事件

`ChoiceResult.Changed()` 在录到新快捷键时为 true。

## 键盘操作

Enter/Space 开始录制；Escape 取消；Tab 离开。字母、数字等须带 Ctrl、Alt 或 Command，F1–F12 可单独使用。

## 限制

- 系统或应用菜单先拦截的组合键（如 Cmd+Q）无法录制。
