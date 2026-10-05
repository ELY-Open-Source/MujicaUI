# 319 CommandBlock 命令块

## 用途

把一次命令执行显示为一块：提示符与高亮的命令、退出码、耗时与可折叠的 ANSI 输出，可复制命令或输出。

## 最小示例

```go
open := true
run := agent.CommandRun{Command: "go test ./...", Output: "\x1b[32mok\x1b[0m  example  0.4s", Duration: 400 * time.Millisecond}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.CommandBlock(c, run, &open, git.CommandBlockOptions{})
	if what, ok := r.Copied(); ok {
		_ = what // "command" or "output" is on the clipboard
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func CommandBlock(c *ui.Context, run CommandRun, open *bool, opts CommandBlockOptions) CommandBlockResult`

## 参数

`CommandRun`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Command | string | 必填 | 命令；为空时 panic。 |
| Output | string | "" | 输出，可含 ANSI 转义。 |
| ExitCode | int | 0 | 退出码。 |
| Duration | time.Duration | 0 | 耗时；为负 panic。 |
| Running | bool | false | 运行中：显示加载标记，不显示退出码与耗时。 |

`CommandBlockOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Prompt | string | "$" | 命令前的提示符。 |
| MaxOutputHeight | float32 | 0（即 320） | 输出区最大高度，超出滚动；为负 panic。 |

`open` 控制输出是否展开。

## 状态

- 成功：Success 色对勾与“退出码 0”；失败：Danger 色叉号与退出码。
- 展开箭头 120 ms 转动，减少动画时直接切换。
- 输出区为 Background 底色，等宽字体，横纵滚动。

## 事件

- `Toggled()`：展开或折叠。
- `Copied()`：`"command"` 或 `"output"`；复制的输出已去除 ANSI 转义。

## 键盘操作

- Tab 到标题行，Enter 或 Space 展开或折叠；复制按钮用 Tab 到达。

## 限制

- 不运行命令：输出与状态由调用方提供。
