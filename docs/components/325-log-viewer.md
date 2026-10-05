# 325 LogViewer · 日志查看器

## 用途

显示日志：按级别过滤、跟随尾部、关键字高亮，十万行也只构建可见行。

## 最小示例

```go
lines := []code.LogLine{{Level: code.LogInfo, Text: "court opened"}, {Level: code.LogError, Text: "seal missing"}}
var state code.LogViewerState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	code.LogViewer(c, lines, &state, code.LogViewerOptions{}).Element.Height(240)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func LogViewer(c *ui.Context, lines []LogLine, s *LogViewerState, opts LogViewerOptions) ChoiceResult`。新追加的行增量过滤；级别变化时重新过滤一次，不逐帧计算。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `LogLine.Level` | `LogLevel` | `LogDebug` | `LogDebug`、`LogInfo`、`LogWarn`、`LogError`；越界 panic。 |
| `LogLine.Time`, `Text` | `time.Time`, `string` | — | 时间（零值不显示）与内容。 |
| `LogViewerState.Hidden` | `[4]bool` | 全显示 | 被隐藏的级别。 |
| `LogViewerState.Query` | `string` | "" | 高亮关键字，不区分大小写。 |
| `LogViewerState.Follow` | `bool` | 首次为 true | 跟随尾部；向上滚动后 MyGo 暂停跟随，滚回底部恢复。 |
| `Disabled` | `bool` | false | 不能改过滤与跟随。 |

## 状态

级别按钮按下表示显示；级别同时用文字（DEBUG/INFO/WARN/ERROR）与颜色区分；高亮为 Selection 背景、粗体与下划线；无日志空状态；禁用。

## 事件

`Changed()`：级别过滤或跟随开关变化。

## 键盘操作

- Tab 依次到达级别按钮、高亮输入框、跟随按钮与日志列表；Space 切换按钮。
- 列表中方向键、PageUp/PageDown、Home/End 滚动。

## 限制

- 只识别向末尾追加；替换为更短的切片会重新过滤。
- 高亮不过滤行，只标出匹配。
