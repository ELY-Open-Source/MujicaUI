# 032 Toggle · 切换按钮

## 用途

保持按下状态的按钮，例如编辑器工具栏中的“加粗”。可显示文字、图标或二者兼有。

## 最小示例

```go
bold := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.Toggle(c, &bold, "", input.ToggleOptions{Icon: icons.Must("bold"), Label: "Bold"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Toggle(c *ui.Context, on *bool, label string, opts ToggleOptions) ui.Element`。`label` 与 `opts.Label` 至少提供一个，否则 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Icon` | `*ui.SVG` | nil | 显示在文字前的图标。 |
| `Disabled` | `bool` | false | 禁用：不响应输入，退出 Tab 序列。 |
| `Label` | `string` | "" | 无文字按钮的无障碍名称。 |

## 状态

未按下（象牙底、描边）、按下（玫瑰底、酒红边框与文字，底部 2 DIP 横条标记——选中不只依赖颜色）、悬停、按住、聚焦（2 DIP 焦点环）、禁用。

## 事件

`Changed()` 在用户切换的那一帧返回 true。

## 键盘操作

- Tab：聚焦。
- Space：按下或松开。

## 限制

- 基于 `ui.ToggleBase`，无障碍角色为 ToggleButton，按下状态由 MyGo 报告。
