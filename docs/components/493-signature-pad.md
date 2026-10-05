# 493 SignaturePad 签名板

## 用途

用指针书写签名，可清除，并可用 `SignaturePath` 导出 SVG 路径数据。

## 最小示例

```go
var strokes [][]mujica.SignaturePoint
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.SignaturePad(c, &strokes, mujica.SignaturePadOptions{}).Changed() {
		_ = mujica.SignaturePath(strokes) // SVG path data to store
	}
}
```

签名：`func SignaturePad(c *ui.Context, strokes *[][]SignaturePoint, opts SignaturePadOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `strokes` | `*[][]SignaturePoint` | 必填 | 笔画，坐标为相对画板的 DIP。 |
| `Height` | `float32` | 160 | 画板高度；宽度填满。 |
| `Label` | `string` | “Signature” | 名称。 |
| `Disabled` | `bool` | `false` | 不可书写与清除。 |

## 状态

空白时显示“请在此签名”；有签名基线；清除按钮在无笔画时禁用。

## 事件

`Changed()`：完成一笔或清除。

## 键盘操作

Tab 聚焦画板，Delete/Backspace 清除；Tab 到“清除”按钮。

## 限制

书写只能用指针，键盘无法替代手写。
