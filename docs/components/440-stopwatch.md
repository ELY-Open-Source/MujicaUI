# 440 Stopwatch · 秒表

## 用途

秒表：开始、暂停、继续、计次与复位，计次列表显示单圈与累计时间。

## 最小示例

```go
var sw mujica.StopwatchState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Stopwatch(c, &sw, mujica.StopwatchOptions{}).Changed() {
		fmt.Println("laps", len(sw.Laps))
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Stopwatch(c *ui.Context, s *StopwatchState, opts StopwatchOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Now` | `time.Time` | 零值（系统时钟） | 固定当前时间并停止刷新。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

未开始、运行中、暂停（显示“继续”）、有计次、禁用。

## 事件

`Changed()` 报告开始、暂停、计次或复位；`StopwatchState.Elapsed` 与 `Laps` 读取结果。

## 键盘操作

- Tab 到达开始/暂停、计次与复位，Enter 或 Space 按下。

## 限制

- 运行且可见时每 100 ms 用 `c.After` 预约重绘；暂停或隐藏后停止。
