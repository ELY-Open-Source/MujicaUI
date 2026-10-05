# 332 OutputPanel · 输出面板

## 用途

显示多个输出通道之一，跟随尾部，支持 ANSI 颜色；清空请求通过结果交给调用方。

## 最小示例

```go
channels := []mujica.OutputChannel{{Name: "Build", Lines: []string{"compiling court", "\x1b[32mok\x1b[0m"}}}
var state mujica.OutputPanelState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.OutputPanel(c, channels, &state, mujica.OutputPanelOptions{})
	r.Element.Height(240)
	if r.Cleared != "" {
		channels[state.Channel].Lines = nil
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func OutputPanel(c *ui.Context, channels []OutputChannel, s *OutputPanelState, opts OutputPanelOptions) OutputPanelResult`。通道为空或 `Channel` 越界时 panic。行列表虚拟化。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OutputChannel.Name`, `Lines` | `string`, `[]string` | — | 通道名与内容，行内可含 ANSI 样式。 |
| `OutputPanelState.Channel` | `int` | 0 | 当前通道。 |
| `OutputPanelState.Follow` | `bool` | 首次为 true | 跟随尾部。 |
| `Disabled` | `bool` | false | 不能切换、清空与改跟随。 |

## 状态

通道下拉、跟随按钮按下、清空按钮、无输出、禁用。

## 事件

`Changed`：切换了通道；`Cleared`：要清空的通道名，调用方清空其 `Lines`。

## 键盘操作

- Tab 到达通道下拉、跟随按钮、清空按钮与输出列表；下拉中方向键选择。

## 限制

- 只渲染 SGR 样式，其余控制序列原样保留（与 `internal.ParseAnsi` 一致）。
