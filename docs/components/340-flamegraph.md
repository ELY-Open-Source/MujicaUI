# 340 Flamegraph · 火焰图

## 用途

以层叠框绘制性能剖析：点击放大子树，搜索高亮匹配帧，悬停或键盘移动时显示并朗读耗时读数。

## 最小示例

```go
root := &mujica.FlameNode{Name: "main", Value: 100, Children: []mujica.FlameNode{
	{Name: "court.Seal", Value: 60}, {Name: "ledger.Write", Value: 30},
}}
var state mujica.FlamegraphState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Flamegraph(c, root, &state, mujica.FlamegraphOptions{Unit: "ms"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Flamegraph(c *ui.Context, root *FlameNode, s *FlamegraphState, opts FlamegraphOptions) ChoiceResult`。帧宽按耗时占比；放大的帧占满宽度，祖先以淡色整宽显示。子帧耗时之和大于父帧或出现负值时 panic。帧颜色按名称取自主题调色板（AccentText、Info、Success、Warning、加深的 Ornament、TextMuted，均对 Surface 至少 3:1）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `FlameNode.Name`, `Value`, `Children` | `string`, `float64`, `[]FlameNode` | — | 帧名、含子帧的总耗时与子帧。 |
| `FlamegraphState.Zoom`, `Focus` | `[]int` | nil（根） | 放大与键盘焦点所在帧的子序号路径。 |
| `FlamegraphState.Query` | `string` | "" | 搜索高亮。 |
| `Unit` | `string` | "" | 读数中的单位。 |
| `Status`, `Error` | `DataStatus`, `string` | `DataReady` | 加载中或失败，与图同高显示。 |
| `Disabled` | `bool` | false | 不能放大与移动。 |

## 状态

帧为浅色底加同色边框与名称；搜索匹配加 AccentText 粗框，其余帧变淡；焦点帧带焦点环；读数行显示名称、耗时与占比；加载、失败、无数据。

## 事件

`Changed()`：放大或焦点变化（失败状态下点击“重试”也为 true）。

## 键盘操作

- Tab 聚焦图表；↑ 到父帧，↓ 到第一个子帧，←/→ 在兄弟帧间移动；Enter 放大当前帧，Escape 缩小一级。
- 读数随焦点更新并通过屏幕阅读器朗读。

## 限制

- 窄于 1 DIP 的帧不绘制，窄于 24 DIP 的帧不写名称；用键盘可到达所有帧。
- 帧行高 28 DIP，窄帧的点击区域可能小于 28 DIP，请用键盘或先放大。
