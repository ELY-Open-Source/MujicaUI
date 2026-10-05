# 272 CommandExecutionCard 命令执行卡片

## 用途

显示一条命令的执行：命令、工作目录、带 ANSI 颜色的输出、退出码与耗时，执行中可请求中止。

## 最小示例

```go
run := mujica.CommandRun{Command: "go test ./...", Dir: "~/court", Running: true, Output: "\x1b[32mok\x1b[0m court\n"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.CommandExecutionCard(c, run, mujica.CommandExecutionCardOptions{}).Changed() {
		fmt.Println("abort the process")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func CommandExecutionCard(c *ui.Context, run CommandRun, opts CommandExecutionCardOptions) ChoiceResult`。`Command` 为空或 `Height` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Height` | `float32` | 240 | 输出区高度。 |
| `Disabled` | `bool` | `false` | 禁止中止。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- `CommandRun` 字段：`Command`、`Dir`、`Output`、`Running`、`ExitCode`、`Duration`。
- 输出由 `internal.ParseAnsi` 解析，颜色映射到主题；样式可跨行延续；追加输出只解析新增行。
- 输出虚拟化并跟随末尾，用户向上滚动后停止跟随。
- 退出码以图标和“退出码 N”文字显示。

## 事件

- `Changed()`：用户点击了“中止”；终止进程由调用方完成。

## 键盘操作

- Tab 到“中止”，Enter / Space 触发；输出区可滚动。

## 限制

- 只支持 SGR 颜色与样式，其它转义序列原样显示。
- 不运行命令，也不提供 PTY。
