# 180 OrderConfirmDialog

下单确认对话框。

## 用途

以模态对话框展示下单摘要与风险提示，默认聚焦取消；可要求勾选确认后才能下单。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.OrderConfirmDialog(c, &open, mujica.OrderConfirmDialogOptions{
		Lines:       []mujica.OrderSummaryLine{{Label: "Symbol", Value: "MUJC"}, {Label: "Quantity", Value: "100"}},
		Risk:        "Leveraged products can lose more than the deposit.",
		Acknowledge: "I understand the risk",
	})
	if r.Confirmed {
		fmt.Println("send order")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `open` | `*bool` | 必填 | 是否显示。 |
| `Lines` | `[]OrderSummaryLine` | 必填 | 摘要行，空值 panic。 |
| `Title` | `string` | 内置“Confirm order” | 标题。 |
| `Risk` | `string` | "" | 风险提示。 |
| `Acknowledge` | `string` | "" | 二次确认复选框文字；设置后须勾选才能下单。 |
| `ConfirmLabel` | `string` | 内置“Place order” | 确认按钮文字。 |

## 状态

打开、未勾选（下单按钮禁用）、已勾选。

## 事件

`AlertDialogResult.Confirmed` / `Canceled`；Escape 视为取消。

## 键盘操作

打开时焦点在取消上；Tab 在复选框与按钮间移动；Escape 取消。点击遮罩不关闭。

## 限制

- 每次打开都需重新勾选。
