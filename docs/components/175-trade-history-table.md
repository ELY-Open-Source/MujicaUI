# 175 TradeHistoryTable

成交记录表。

## 用途

成交记录：时间、代码、方向、价格、数量与手续费；按日期范围筛选；导出当前范围的记录交给调用方。

## 最小示例

```go
var state finance.TradeHistoryTableState
rows := []finance.FillRow{{ID: "f1", Time: time.Now(), Symbol: "MUJC", Side: finance.SideBuy, Price: 400, Quantity: 10, Fee: 1.2}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if fills, ok := finance.TradeHistoryTable(c, &state, rows, finance.TradeHistoryTableOptions{}).Exported(); ok {
		fmt.Println("export", len(fills))
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*TradeHistoryTableState` | 必填 | Range 为日期范围，零值不筛选。 |
| `rows` | `[]FillRow` | 必填 | 成交记录。 |
| `Decimals / QuantityDecimals` | `int` | 2 / 0 | 小数位。 |
| `Version` | `int` | 0 | 原地修改行时递增。 |
| `Status / Error` | `DataStatus / string` | Ready | 加载与失败状态。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

全部日期、指定范围、空范围、加载、失败、禁用。

## 事件

`TradeHistoryResult.Exported()` 返回导出请求时当前范围内的记录。

## 键盘操作

日期范围选择器与按钮按 Tab 顺序；表格沿用 DataTable 键盘操作。

## 限制

- 只回传数据，不生成文件：导出格式与写入由调用方负责。
- 日期按每条记录自身时区的日期比较。
