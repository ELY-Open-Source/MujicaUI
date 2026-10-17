# 499 QRCode 二维码

## 用途

在本地用 Go 编码文字为二维码（字节模式，版本 1–40，纠错 L/M/Q/H），不访问网络。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.QRCode(c, "https://example.com/mujica", input.QRCodeOptions{Level: input.QRQuartile})
}
```

签名：`func QRCode(c *ui.Context, text string, opts QRCodeOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `text` | `string` | 必填 | 按 UTF-8 字节编码。 |
| `Level` | `QRLevel` | `QRMedium` | `QRLow`、`QRMedium`、`QRQuartile`、`QRHigh`；未知值 panic。 |
| `Size` | `float32` | 160 | 含静区的边长；负数或非有限 panic。 |
| `Label` | `string` | “QR code” | 无障碍名称。 |

## 状态

始终深色模块印在浅色底上，深色主题也一样，以便扫码；模块对齐设备像素。文字过长时显示错误。

## 事件

无。另有 `EncodeQR(data, level) (QRMatrix, error)` 供调用方导出。

## 键盘操作

不可聚焦。

## 限制

只用字节模式；只在文字或纠错级别变化时重新编码。
