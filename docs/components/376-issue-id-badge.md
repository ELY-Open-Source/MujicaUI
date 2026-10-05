# 376 IssueIdBadge · 任务编号

## 用途

任务编号用于项目管理界面。

## 最小示例

```go
id := "ENG-123"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.IssueIdBadge(c, id, project.IssueIdBadgeOptions{Copyable: true}).Changed() {
		log.Println("copied", id)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`id string`；`Copyable` 默认 false；`Disabled` 禁用复制。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示复制，剪贴板由组件写入；同时朗读复制结果。

## 键盘操作

Copyable 开启时 Tab 聚焦，Enter/Space 复制。

## 限制

编号匹配大写项目名、连字符、正整数，例如 ENG-123。复制提示保持 1500 ms。
