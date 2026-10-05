# FilePreview

## 用途

展示图片、文本、代码和不支持类型说明。

## 最小示例

```go
file := files.PreviewFile{Name: "seal.go", Text: "package court\nfunc Seal() {}"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FilePreview(c, file, files.FilePreviewOptions{}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)
```

## 参数

Name 必填；Image、Text、Language、Size 为调用方数据；Status、Error 描述加载状态。

## 状态

内容就绪、加载、失败、不支持类型。

## 事件

Closed、Opened、Retried 返回关闭、打开和重试请求。

## 键盘操作

Escape 关闭；Tab 到操作；文本支持选择与复制。

## 限制

读取和解码由调用方提供；预览作为内嵌面板布局。
