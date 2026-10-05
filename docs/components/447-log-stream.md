# 447 LogStream · 日志流

## 用途

展示实时 ANSI 日志，支持暂停、级别过滤与跟随末尾。

## 最小示例

```go
lines := []mujica.LogLine{{Level: mujica.LogInfo, Text: "server ready"}, {Level: mujica.LogWarn, Text: "slow query"}}
var state mujica.LogStreamState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.LogStream(c, lines, &state, mujica.LogStreamOptions{}).Changed() {
		fmt.Println("filter changed")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Height 默认 320 DIP；Disabled 禁用；LogStreamState.Paused 控制暂停，MinLevel 指定最低级别。

## 状态

跟随末尾、暂停、级别过滤、空数据、禁用。

## 事件

Changed() 报告暂停或过滤变化。

## 键盘操作

Tab 聚焦暂停和级别；方向键切换级别；Space 激活暂停。

## 限制

追加或替换日志提交新切片；暂停冻结条数。
