# 255 CostEstimator · 费用估算

## 用途

显示输入、输出 token 数及按调用方单价估算的费用。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.CostEstimator(c, chat.CostEstimatorOptions{InputTokens: 12000, OutputTokens: 3400, InputPrice: 3, OutputPrice: 15})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CostEstimator(c *ui.Context, opts CostEstimatorOptions) ui.Element`。另有 `EstimateCost(inputTokens, outputTokens, inputPrice, outputPrice) CostEstimate`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `InputTokens`, `OutputTokens` | `int` | 0 | token 数，不能为负。 |
| `InputPrice`, `OutputPrice` | `float64` | 0 | 每百万 tokens 单价，须有限且不为负。 |
| `Currency` | `string` | "$" | 金额前缀。 |
| `Digits` | `int` | 0 → 4 | 金额小数位；负数 panic。 |

## 状态

只读表格：输入、输出、合计（衬线标签、酒红加粗金额）与单价说明。

## 事件

无。

## 键盘操作

不获取焦点；整体读作“预计费用 金额”。

## 限制

- 估算按 tokens × 单价 / 1,000,000 的浮点运算；显示按 `Digits` 位小数四舍五入，不用于结算。
