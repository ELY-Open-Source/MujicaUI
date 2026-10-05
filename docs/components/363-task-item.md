# 363 TaskItem · 任务行

## 用途

任务行用于项目管理界面。

## 最小示例

```go
task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.TaskItem(c, &task, project.TaskItemOptions{})
	if id, ok := r.Changed(); ok {
		log.Println("edit", id)
	}
	if id, ok := r.Opened(); ok {
		log.Println("open", id)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`task *Task` 保存标题、完成、优先级、日期、负责人；`Members` 解析负责人；`Today` 默认当前日期；`Disabled` 禁用输入。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 返回完成或标题编辑的任务 ID；`Opened()` 返回打开任务 ID。

## 键盘操作

Tab 在复选框和标题间移动；空格切换完成；Enter 打开；F2 或双击开始编辑；Enter 保存；Escape 取消。

## 限制

任务标题保存时去除首尾空白；空标题保留编辑态。业务字段来自调用方。
