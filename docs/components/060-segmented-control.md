# 060 SegmentedControl 分段控件

## 用途

在少量互斥选项之间切换视图或范围，如“日 / 周 / 月”或“列表 / 网格”。基于 MyGo `SegmentedBase`，辅助技术看到单选组。

## 最小示例

```go
if navigation.SegmentedControl(c, &view, []navigation.Segment{
	{Label: "List", Icon: icons.Must("list"), HideLabel: true},
	{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: true},
}, navigation.SegmentedControlOptions{Label: "View"}).Changed() {
	fmt.Println("view", view)
}
```

签名：`func SegmentedControl(c *ui.Context, selected *int, segments []Segment, opts SegmentedControlOptions) ui.Element`

## 参数

`selected`：当前分段下标；越界时 panic。`segments` 为空时 panic。

`Segment`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 必填 | 分段文字与无障碍名称；为空时 panic。 |
| `Icon` | `*ui.SVG` | nil | 文字前的图标。 |
| `HideLabel` | `bool` | false | 只显示图标，`Label` 作为名称与提示文字；没有 `Icon` 时 panic。 |

`SegmentedControlOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Equal` | `bool` | false | 所有分段等宽（取最宽者）；否则按内容自适应。 |
| `Label` | `string` | "" | 单选组的无障碍名称。 |
| `Disabled` | `bool` | false | 整组禁用。 |

## 状态

- 选中：Accent 填充、OnAccent 文字、600 字重。
- 悬停 SurfaceHover，按下 Selection；焦点环 2 DIP。
- 禁用：整组降低不透明度，点击与方向键不改变值。

## 事件

- 返回元素的 `Changed()`：本帧用户选择了新分段。

## 键盘操作

- Tab：进入控件，落在当前分段（一个 Tab 停靠点）。
- 方向键：选择上一个/下一个分段并循环；Home / End 到两端。

## 限制

- 不支持单个分段禁用：MyGo `SegmentedBase` 的焦点组在方向键移动时直接点击目标分段。
