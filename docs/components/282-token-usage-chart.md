# 282 TokenUsageChart Token 用量图

## 用途

按会话或步骤显示输入与输出 token 用量：默认堆叠柱状图，也可切换为折线。

## 最小示例

```go
usage := []mujica.TokenUsage{{Label: "Plan", Input: 1200, Output: 300}, {Label: "Write", Input: 2100, Output: 2900}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TokenUsageChart(c, usage, mujica.TokenUsageChartOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func TokenUsageChart(c *ui.Context, usage []TokenUsage, opts TokenUsageChartOptions) ChartResult`

## 参数

`TokenUsage`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Label | string | "" | 会话或步骤名，作为柱的类别。 |
| Input, Output | int | 0 | 输入、输出 token 数；为负 panic。 |

`TokenUsageChartOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Frame | ChartFrame | 零值 | 图表名称、高度、加载/失败/空状态；Label 为空时用内置“Token 用量”。 |
| Line | bool | false | 折线代替堆叠柱。 |
| Version | int | 0 | 原地修改数据时递增；新切片或长度变化自动刷新。 |

返回 `ChartResult`。序列在 usage 切片、长度或 Version 变化时重建。超过 256 项时，相邻项目汇总为最多 256 组，保留输入输出总数；柱类别显示首尾名称。

## 状态

- 输入、输出分色（取自图表调色板，对 Surface 至少 3:1），图例同时用形状区分。
- 加载、失败、无数据状态与图表同尺寸。

## 事件

- `Changed()`：读数移到另一个标记；`Toggled()`：图例开关某个序列。

## 键盘操作

- Tab 聚焦图表，←/→ 在柱或点之间移动，读数显示并朗读。

## 限制

- 折线模式的横轴为序号，类别名只在柱状模式显示。
