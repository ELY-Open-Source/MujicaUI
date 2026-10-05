# 374 LabelManager · 标签管理

## 用途

标签管理用于项目管理界面。

## 最小示例

```go
labels := []mujica.TaskLabel{{ID: "ui", Name: "UI", Color: mujica.EventInfo}}
selected := []string{}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.LabelManager(c, &labels, &selected, mujica.LabelManagerOptions{}).Changed() {
		log.Println("labels changed", len(labels), selected)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`labels *[]TaskLabel`、`selected *[]string`；ID 非空唯一；`Color` 为六种 EventColor；`Disabled`、`Label` 控制交互与名称；`Version` 用于外部原地修改。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示选择、创建、重命名、颜色或删除变更。删除同步移除已选 ID。

## 键盘操作

Tab 移动；Space 勾选；新增输入 Enter 创建；编辑颜色按钮打开编辑；Enter 保存名称；Escape 结束编辑。

## 限制

名称去除首尾空白并按大小写折叠检测重复；颜色选中立即提交。列表虚拟化，最多显示八行高度。
