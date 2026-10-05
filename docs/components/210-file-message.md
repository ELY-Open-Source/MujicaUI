# 210 FileMessage 文件消息

## 用途

显示消息中的文件：按类型的图标、文件名、大小，以及打开与下载回调。

## 最小示例

```go
r := chat.FileMessage(c, chat.FileMessageOptions{Name: "treaty.pdf", Size: 2621440})
if r.Downloaded() {
	fmt.Println("save treaty.pdf")
}
```

签名：`func FileMessage(c *ui.Context, opts FileMessageOptions) FileMessageResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Name` | `string` | 必填 | 文件名；为空 panic。 |
| `Size` | `int64` | 0 | 字节数，按 B、KB、MB、GB 显示（1024 进制）；负数 panic。 |
| `Type` | `string` | "" | MIME 类型；为空时按扩展名选择图标。 |

## 状态

悬停时文件区域加深；禁用时两个操作都不响应。

## 事件

`Opened()`：点击或在文件区按 Enter。`Downloaded()`：点击下载按钮。

## 键盘操作

Tab 依次到文件区与下载按钮，Enter 或 Space 触发。

## 限制

- 组件不读取文件、不下载：打开与下载由调用方完成。
