# 497 ResizablePanelGroup 可调面板组

## 用途

多个面板以分隔条相连，拖动或键盘调整尺寸；可折叠面板拖到最小值一半以下即折叠。尺寸以比例保存在调用方的切片中。

## 最小示例

```go
var sizes []float32
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.ResizablePanelGroup(c, &sizes, []mujica.ResizablePanel{
		{Label: "Files", Collapsible: true, Build: func() { ui.Text(c, "Files") }},
		{Label: "Editor", Build: func() { ui.Text(c, "Editor") }},
	}, mujica.ResizablePanelGroupOptions{})
	if r.Committed {
		// save sizes
	}
}
```

签名：`func ResizablePanelGroup(c *ui.Context, sizes *[]float32, panels []ResizablePanel, opts ResizablePanelGroupOptions) ResizablePanelGroupResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `sizes` | `*[]float32` | 空则平分 | 各面板比例，和为 1，否则 panic；0 只允许可折叠面板。 |
| `panels` | `[]ResizablePanel` | 必填 | Label、Min（默认 10 个间距单位）、Collapsible、Build；至少两个。 |
| `Vertical` | `bool` | `false` | 上下排列。 |
| `Disabled` | `bool` | `false` | 不可调整。 |

## 状态

分隔条悬停或拖动时变为强调色；折叠后分隔条报告“已折叠”。

## 事件

`Resized`：拖动中尺寸变化；`Committed`：一次调整完成（松开、按键、折叠），用于保存。

## 键盘操作

Tab 聚焦分隔条；方向键每次 4 个间距单位；Home/End 到最小与最大；Enter 折叠或恢复相邻的可折叠面板。

## 限制

无。
