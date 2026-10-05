# 452 DashboardFilterBar · 仪表盘过滤栏

## 用途

编辑时间范围、环境和标签，并支持清除。

## 最小示例

```go
var filters devtools.DashboardFilters
opts := devtools.DashboardFilterBarOptions{Envs: []string{"production", "staging"}, Tags: []string{"api", "db"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if devtools.DashboardFilterBar(c, &filters, opts).Changed() {
		fmt.Println("filters changed")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Ranges 默认 15 分钟、1 小时、1 天、7 天；Env 与 Tags 来自选项；Range 为 0 使用调用方默认；Disabled 禁用。

## 状态

默认过滤、已选择、已清除、禁用。

## 事件

Changed() 报告过滤或清除。

## 键盘操作

Tab 遍历控件；方向键选择；Space 切换标签；Escape 关闭弹层。

## 限制

过滤请求由调用方执行；选中值须存在于提供的选项中。
