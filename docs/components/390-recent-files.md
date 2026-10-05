# RecentFiles

## 用途

按今天、昨天、本周与更早分组展示最近文件。

## 最小示例

```go
now := time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)
recent := []files.RecentFile{{Path: "/court/seal.go", Opened: now}}
var state files.RecentFilesState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.RecentFiles(c, recent, &state, files.RecentFilesOptions{Now: now}).Element.Height(240)
}
ui.Render(view, 500, 300, 1)
```

## 参数

Now 默认帧时间；Revision 原地更新时递增；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Opened、Removed、Cleared 为打开、移除、清空请求。

## 键盘操作

上下选择文件；Enter 打开；Tab 到清除和移除按钮。

## 限制

时间按 Now 时区分组；切片替换或 Revision 更新会重算排序。
