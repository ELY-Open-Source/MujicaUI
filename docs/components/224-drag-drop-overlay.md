# 224 DragDropOverlay 拖放遮罩

## 用途

拖入文件时的全区域遮罩与提示；类型检查回调。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	images := func(path string) bool { return strings.HasSuffix(path, ".png") }
	r := chat.DragDropOverlay(c, chat.DragDropOverlayOptions{Accept: images}, func() {
		ui.Text(c, "Drag images into the chat")
	})
	for _, path := range r.Accepted() {
		fmt.Println("attach", path)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`content` 必填；`Hint` 默认“松开以添加文件”；`Accept` 接收路径并返回类型检查结果，nil 接受全部。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Accepted() 与 Rejected() 返回本帧接受和拒绝的文件路径。

## 键盘操作

内容中的操作按其自身键盘规则访问；文件拖入由操作系统驱动。

## 限制

读取文件与上传由调用方完成。MyGo ui/headless.go:93–163 的 Tester 通过平台事件驱动，公开接口未提供外部文件拖放注入；类型检查与遮罩分别测试。
