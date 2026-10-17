# 007 ScrollArea 滚动区域

## 用途

纵向、横向或双向滚动内容。支持滚轮、触控板、滚动条拖动与键盘；可读取、设置并保留滚动位置。基于 MyGo `ui.Scroll`。

## 最小示例

```go
var offset ui.ScrollState
layout.ScrollArea(c, &offset, layout.ScrollAreaOptions{Label: "Log"}, func() {
	for range 50 {
		ui.Text(c, "A line of the log")
	}
}).Height(200)
```

签名：`func ScrollArea(c *ui.Context, state *ui.ScrollState, opts ScrollAreaOptions, children func()) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `state` | `*ui.ScrollState` | nil | 滚动位置。设置 `X`/`Y` 即滚动到该处，`math.MaxFloat32` 表示末尾；组件写回 `MaxX`/`MaxY`。nil 时按区域在视图中的位置保留。 |
| `Direction` | `ScrollDirection` | `ScrollVertical` | `ScrollVertical`、`ScrollHorizontal`、`ScrollBoth`；非法值 panic。 |
| `Label` | `string` | "" | 区域的无障碍名称。 |
| `children` | `func()` | 必填 | 内容；nil 时 panic。 |

需给区域尺寸（`Height`、`Grow` 等）。

## 状态

- 内容溢出时区域可获得焦点，键盘焦点环画在区域内侧：先留 2 DIP 间隔，再画 2 DIP Focus 色环线，避免被父级裁剪。溢出状态在布局后变化时组件立即重建，焦点停靠随之更新。
- 不溢出时不进入 Tab 序列。
- 滚动条为 MyGo 叠加式滑块，颜色取主题 Scrollbar（TextMuted 50%），指针悬停时显示，不占布局空间。
- 切换主题、折叠周围布局都不重置位置。

## 事件

无独立事件；比较 `state.X`/`state.Y` 的变化即可得知滚动。

## 键盘操作

焦点在区域或其内容上时：方向键滚动 40 DIP，Page Up/Down 与 Space/Shift+Space 翻页，Home/End 到顶部/底部（macOS 另有 Cmd+↑/↓）。文本输入等自行处理这些键的控件优先。

## 限制

- 滚动条只在悬停时显示，不提供常驻滚动条（MyGo 滚动条不占空间的设计）。
- 大量行请使用虚拟化列表，ScrollArea 会构建全部子项。
