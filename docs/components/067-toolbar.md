# 067 Toolbar 工具栏

## 用途

一行操作按钮：分组、分隔线；宽度不足时，放不下的按钮从末尾起移入溢出菜单。基于 MyGo `ui.Toolbar`，保留其焦点组与溢出行为。

## 最小示例

```go
navigation.Toolbar(c, navigation.ToolbarOptions{Label: "Format"}, func() {
	navigation.ToolbarGroup(c, func() {
		navigation.ToolbarButton(c, "Bold", navigation.ToolbarButtonOptions{Icon: icons.Must("bold")})
		navigation.ToolbarButton(c, "Italic", navigation.ToolbarButtonOptions{Icon: icons.Must("italic")})
	})
	navigation.ToolbarSeparator(c)
	if navigation.ToolbarButton(c, "Copy", navigation.ToolbarButtonOptions{Icon: icons.Must("copy"), ShowLabel: true}).Clicked() {
		fmt.Println("copy")
	}
})
```

签名：

- `func Toolbar(c *ui.Context, opts ToolbarOptions, fn func()) ui.Element`
- `func ToolbarButton(c *ui.Context, label string, opts ToolbarButtonOptions) ui.Element`
- `func ToolbarGroup(c *ui.Context, fn func()) ui.Element`
- `func ToolbarSeparator(c *ui.Context)`

## 参数

`ToolbarOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | "" | 工具栏的无障碍名称。 |

`ToolbarButton` 的 `label` 必填（为空时 panic）：它是按钮名称、仅图标时的提示文字和溢出菜单中的文字。

`ToolbarButtonOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Icon` | `*ui.SVG` | nil | 16 DIP 图标。 |
| `ShowLabel` | `bool` | false | 图标旁显示文字；没有图标时总是显示。 |
| `Disabled` | `bool` | false | 禁用；溢出菜单中同样禁用。 |

## 状态

- 工具栏：Surface 背景、Border 边框、4 DIP 圆角。
- 按钮：无底色，悬停 SurfaceHover，按下 Selection，焦点环 2 DIP。
- 分隔线：旧金细竖线，不进入无障碍树。
- 溢出：末尾显示 “»” 按钮，菜单中按原顺序列出被收起的操作，组之间以分隔线隔开。

## 事件

- `ToolbarButton(...).Clicked()`：点击、Enter / Space，或在溢出菜单中选择。

## 键盘操作

- Tab：进入工具栏（一个停靠点），再按 Tab 离开。
- ←/→：在操作之间移动（跳过禁用项并循环）；Home / End 到首末。
- Enter / Space：执行。溢出按钮用 Enter、Space 或 ↓ 打开菜单。

## 限制

- 溢出菜单是 MyGo 打开的系统菜单（`ui/toolbar.go:45`），外观为系统样式，不是 MujicaUI 自绘。
- 溢出按钮的无障碍名称固定为英文 “More”（`ui/toolbar.go:44`），按钮样式来自 MyGo 主题。
