# 335 DebugToolbar · 调试工具栏

## 用途

调试会话的继续/暂停、单步跳过、单步进入、单步跳出、重启与停止，按会话状态禁用不适用的按钮，并注册常用快捷键。

## 最小示例

```go
state := code.DebugPaused
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	switch code.DebugToolbar(c, state, code.DebugToolbarOptions{}).Action {
	case code.DebugContinue:
		state = code.DebugRunning
	case code.DebugStop:
		state = code.DebugStopped
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DebugToolbar(c *ui.Context, state DebugState, opts DebugToolbarOptions) DebugToolbarResult`。`DebugEnabled(state, action)` 给出禁用规则：运行中只能暂停、重启、停止；暂停时可继续、单步、重启、停止；停止时只能“开始”（`DebugContinue`）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `state` | `DebugState` | — | `DebugStopped`、`DebugRunning`、`DebugPaused`；越界 panic。 |
| `Disabled` | `bool` | false | 全部按钮与快捷键无效。 |

## 状态

运行中第一个按钮显示暂停，否则显示继续（停止时名为“开始”）；不适用的按钮禁用；按钮悬停、按下、焦点环与带快捷键的提示。

## 事件

`Action`：本帧请求的 `DebugAction`，无则为 `DebugNone`。

## 键盘操作

- F5 继续/开始，F6 暂停，F10 单步跳过，F11 单步进入，Shift+F11 单步跳出，Cmd/Ctrl+Shift+F5 重启，Shift+F5 停止（窗口内全局，焦点控件优先处理）。
- 工具栏是一个 Tab 停靠点，←/→ 在按钮间移动。

## 限制

- 快捷键固定，不可配置；同一窗口只放一个调试工具栏。
- 图标为自绘 SVG。
