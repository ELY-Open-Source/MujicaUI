# 222 AttachmentChip 附件标签

## 用途

附件名、类型、大小、上传进度、移除；失败重试。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := chat.AttachmentChip(c, chat.Attachment{Name: "guests.csv", Size: 4096, Uploading: true, Progress: 0.4})
	if r.Removed() {
		fmt.Println("cancel the upload")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`Attachment.Name` 必填；`Type` 为 MIME 类型，空值按扩展名选择图标；`Size` 为非负字节数；`Uploading` 默认 false，`Progress` 为 0–1；`Err` 表示上传失败。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Removed() 表示移除；Retried() 表示重试。

## 键盘操作

Tab 到移除与重试；Enter、Space 激活。

## 限制

上传由调用方执行，组件只展示进度与结果。
