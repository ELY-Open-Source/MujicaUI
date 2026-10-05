# 380 PomodoroTimer · 番茄钟

## 用途

番茄钟用于项目管理界面。

## 最小示例

```go
var state project.PomodoroState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.PomodoroTimer(c, &state, project.PomodoroTimerOptions{Focus: 25 * time.Minute})
	if r.Changed() {
		log.Println("timer", state.Phase(), state.Running())
	}
	if phase, ok := r.Finished(); ok {
		log.Println("finished", phase)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`state *PomodoroState`；`Focus/Break/LongBreak` 默认 25/5/15 分钟；`LongEvery` 默认 4；`Now` 固定时钟；`Disabled`、`Label` 控制交互与名称。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 返回操作事件；`Finished()` 返回自然完成的阶段，调用方可播放提示。

## 键盘操作

Tab 移动；Enter/Space 开始、暂停、跳过与重置。

## 限制

时长为整秒；自然完成后等待用户启动下一阶段。隐藏组件停止请求帧，恢复显示时读取经过时间。
