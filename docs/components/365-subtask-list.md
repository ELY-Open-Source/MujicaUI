# 365 SubtaskList · 子任务列表

## 用途

子任务列表用于项目管理界面。

## 最小示例

```go
subtasks := []project.Subtask{{ID: "one", Title: "Review", Done: true}, {ID: "two", Title: "Ship"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.SubtaskList(c, &subtasks, project.SubtaskListOptions{}).Changed() {
		log.Println("subtasks changed", len(subtasks))
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`subs *[]Subtask` 保存子任务；`Disabled` 禁用所有操作；`Label` 为进度与列表名称；`Version` 用于外部原地修改。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示勾选、新增、删除或排序；变更已写入切片。

## 键盘操作

Tab 移动；Space 切换完成；新增输入框 Enter 提交；Alt+上下排序；拖动握柄排序。

## 限制

新增 ID 按 `sub-N` 生成；原有 ID 必须唯一且非空。空白新增内容保留输入。列表虚拟化，最多显示八行高度。
