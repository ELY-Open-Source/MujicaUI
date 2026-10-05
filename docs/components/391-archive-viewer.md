# ArchiveViewer

## 用途

展示调用方条目树、大小和压缩率。

## 最小示例

```go
entries := []files.ArchiveEntry{{Path: "court/seal.go", Size: 4000, Packed: 1000}}
var state files.ArchiveViewerState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.ArchiveViewer(c, entries, &state, files.ArchiveViewerOptions{}).Element.Height(240)
}
ui.Render(view, 600, 300, 1)
```

## 参数

Path 唯一；Size、Packed 非负；Status、Error 为加载状态；Revision 原地更新时递增。

## 状态

条目树、加载、失败、空列表。

## 事件

Opened 返回文件打开请求；Retried 返回重试请求。

## 键盘操作

上下选择；左右展开折叠；Enter 打开；Tab 到列操作。

## 限制

压缩包解码由调用方完成；祖先目录自动补齐，文件与目录身份冲突会 panic。
