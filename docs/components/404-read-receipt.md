# 404 ReadReceipt 已读回执

## 用途

通过勾选标记与文字展示发送进度，已读时显示多人头像。

## 最小示例

```go
readers := []display.AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.ReadReceipt(c, messaging.ReceiptRead, messaging.ReadReceiptOptions{Readers: readers})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`state ReceiptState` 为 `ReceiptSent`、`ReceiptDelivered`、`ReceiptRead`，零值表示已发送。`Readers []AvatarItem` 提供已读头像；`Max` 默认 3，必须 ≥ 0。

## 状态

状态由调用方数据控制；浅深主题、密度与文字缩放随窗口设置变化。适用插槽操作沿用控件焦点和禁用状态。

## 事件

纯展示，回执进度与读者由调用方更新。

## 键盘操作

展示组件保持自然焦点序列。

## 限制

已读数据由调用方提供；头像超过 Max 时显示剩余数量。
