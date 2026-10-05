# 166 Level2Quotes

多档行情。

## 用途

并排显示买卖各 N 档的价格、数量与委托笔数；档位变化时高亮一次，并以箭头标出数量增减。

## 最小示例

```go
data := finance.Level2Data{
	Bids: []finance.L2Level{{Price: 100.00, Size: 500, Orders: 12}},
	Asks: []finance.L2Level{{Price: 100.01, Size: 200, Orders: 3}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.Level2Quotes(c, data, finance.Level2QuotesOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `data` | `Level2Data` | 必填 | Bids、Asks，最优价在前。 |
| `Decimals / SizeDecimals` | `int` | 2 / 0 | 价格与数量小数位。 |
| `Levels` | `int` | 5 | 每侧档数；不足的档位显示“—”。 |

## 状态

默认、档位变化高亮（减少动画时只显示箭头）、空档。

## 事件

无。

## 键盘操作

不获取焦点；每行有包含档位、价格、数量和笔数的无障碍名称。

## 限制

- 价格变化视为新档位，不标增减箭头。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
