# 451 RefreshIntervalSelector · 刷新间隔选择

## 用途

选择自动刷新间隔，报告手动与定时刷新请求。

## 最小示例

```go
interval := 30 * time.Second
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.RefreshIntervalSelector(c, &interval, mujica.RefreshIntervalSelectorOptions{})
	if r.Refresh {
		fmt.Println("refresh")
	}
	if r.Changed {
		fmt.Println("interval", interval)
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Intervals 默认提供 5、10、30 秒、1、5 分钟；0 表示关闭；选项须为正整数秒；Refreshing 暂停计时，Disabled 禁用。

## 状态

关闭自动刷新、自动刷新、刷新中、禁用。

## 事件

Refresh 表示请求刷新；Changed 报告间隔变化。

## 键盘操作

Tab 到选择器与按钮；方向键选择间隔；Enter 激活。

## 限制

调用方收到 Refresh 后执行加载并设置 Refreshing。
