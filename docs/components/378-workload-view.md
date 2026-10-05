# 378 WorkloadView · 成员负载

## 用途

成员负载用于项目管理界面。

## 最小示例

```go
loads := []mujica.Workload{{Member: mujica.ProjectMember{ID: "ann", Name: "Ann"}, Assigned: 46, Capacity: 40}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WorkloadView(c, loads, mujica.WorkloadViewOptions{})
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`loads []Workload`；`Unit` 默认小时；`Label` 为列表名称；`Version` 刷新原地数据更改。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

静态展示。

## 键盘操作

静态行通过无障碍名称读取负载与超出量；滚动查看更多成员。

## 限制

成员 ID 唯一；Assigned 有限非负，Capacity 有限正数。可见行虚拟化；父容器需提供有限高度；原地修改递增 Version。
