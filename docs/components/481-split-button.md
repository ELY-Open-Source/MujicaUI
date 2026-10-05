# 481 SplitButton 拆分按钮

## 用途

主操作与下拉菜单合为一体：左半执行主操作，右半打开 `PopupMenu` 列出相关操作。两半各自获得焦点。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := overlay.SplitButton(c, "Save", overlay.SplitButtonOptions{}, func(m *overlay.PopupMenu) {
		if m.Item("Save as…", overlay.PopupMenuItemOptions{}) {
			// ask for a file name
		}
	})
	if r.Changed() {
		// save
	}
}
```

签名：`func SplitButton(c *ui.Context, label string, opts SplitButtonOptions, build func(m *PopupMenu)) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `label` | `string` | 必填 | 主操作文字。 |
| `build` | `func(*PopupMenu)` | 必填 | 填充菜单项，`Item` 报告选择。 |
| `Variant` | `ButtonVariant` | `Primary` | 只接受 `Primary` 或 `Secondary`，其余 panic。 |
| `Icon` | `*ui.SVG` | `nil` | 主操作图标。 |
| `MenuLabel` | `string` | “More options” | 菜单半的无障碍名称。 |
| `Disabled` | `bool` | `false` | 两半都不可用，打开的菜单关闭。 |

## 状态

两半分别有悬停、按下、焦点环；菜单打开时右半保持按下色。禁用时 MyGo 统一减淡。

## 事件

`Changed()`：主操作被触发。菜单项通过 `m.Item` 的返回值报告。

## 键盘操作

Tab 依次聚焦主操作与菜单按钮；Enter/Space 触发主操作；在菜单按钮上 Enter、Space 或 Down 打开菜单，方向键选择，Escape 关闭。

## 限制

菜单宽度至少等于整个按钮组；菜单内容由调用方提供。
