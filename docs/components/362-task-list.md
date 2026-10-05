# 362 TaskList · 任务列表

## 用途

任务列表用于项目管理界面。

## 最小示例

```go
statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
tasks := []mujica.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
filter := mujica.TaskShowAll
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.TaskList(c, tasks, &filter, mujica.TaskListOptions{GroupBy: mujica.TaskGroupStatus, Statuses: statuses})
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

`tasks []Task`、`filter *TaskFilter`；`GroupBy` 为无分组、状态或优先级。按状态分组必传 `Statuses`。`Members` 解析负责人；`Today` 为逾期基准；`Version`、`Disabled`、`Label` 控制数据刷新与交互。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 返回编辑任务 ID；`Opened()` 返回打开任务 ID。筛选值直接写入 `filter`。

## 键盘操作

Tab 在筛选器和任务控件间移动；空格切换完成；Enter 打开；F2 编辑标题。

## 限制

可见行虚拟化；外部原地数据修改递增 `Version`。列表需父容器提供有限高度。
