# FileOperationProgress

## 用途

显示复制、移动、删除进度、取消和冲突回答。

## 最小示例

```go
op := files.FileOperation{Kind: files.FileCopy, Done: 2, Total: 10, Current: "seal.go"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FileOperationProgress(c, op, files.FileOperationOptions{})
}
ui.Render(view, 500, 240, 1)
```

## 参数

Done、Total、BytesDone、BytesTotal 为计数；Conflict 为冲突名称；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Canceled、Resolved、ApplyAll 返回操作意图。

## 键盘操作

Tab 到取消、全选规则和冲突选项；Enter / Space 激活。

## 限制

文件操作由调用方执行；非法计数或 Kind 会 panic。
