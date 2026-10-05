# 372 AssigneePicker · 负责人选择

## 用途

负责人选择用于项目管理界面。

## 最小示例

```go
members := []mujica.ProjectMember{{ID: "ann", Name: "Ann"}, {ID: "ben", Name: "Ben"}}
ids := []string{"ann"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.AssigneePicker(c, &ids, members, mujica.AssigneePickerOptions{Max: 2}).Changed() {
		log.Println("assignees", ids)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`ids *[]string`、`members []ProjectMember`；`Max=0` 表示自由数量；`Disabled`、`Label` 控制交互与名称；`Version` 用于外部原地修改。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示受控成员 ID 列表更改。

## 键盘操作

Enter/Space 打开；输入搜索；下键进入选项；Space/Enter 勾选；Backspace 移除最后一项；Escape 关闭。

## 限制

成员 ID 唯一且非空；所有已选 ID 必须存在。头像使用姓名回退；网络成员查询由调用方提供。搜索结果按可见行虚拟化。
