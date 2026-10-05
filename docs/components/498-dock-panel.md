# 498 DockPanel 停靠面板

## 用途

左、中、右、下四个停靠区，每区以标签页显示面板；把标签拖到区域或窗口边缘即停靠到那里。布局由调用方保存。

## 最小示例

```go
var layout input.DockLayout
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.DockPanel(c, &layout, []input.DockItem{
		{ID: "files", Title: "Files", Build: func() { ui.Text(c, "Files") }},
		{ID: "terminal", Title: "Terminal", Build: func() { ui.Text(c, "$") }},
	}, input.DockPanelOptions{})
	if r.Changed() {
		// save layout
	}
}
```

签名：`func DockPanel(c *ui.Context, layout *DockLayout, items []DockItem, opts DockPanelOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `layout` | `*DockLayout` | 全部在中央 | 各区面板 ID 与当前标签；未知或重复 ID panic。 |
| `items` | `[]DockItem` | 必填 | ID、Title、Build；ID 重复 panic。 |
| `Disabled` | `bool` | `false` | 不可切换与停靠。 |

## 状态

拖动时四边出现停靠条，悬停的目标显示选中底色与强调边框；当前标签有下划线与选中底色。

## 事件

`Changed()`：布局变化（停靠或切换标签），用于保存。

## 键盘操作

每区标签是一个 Tab 停靠点，左右键切换；在标签上按 Shift+Alt+←/→/↓ 停靠到左、右、下，Shift+Alt+↑ 回到中央，焦点跟随标签。

## 限制

区域尺寸固定（侧栏 220 DIP，底部 160 DIP）；需要调整时与 ResizablePanelGroup 组合。
