# 321 Terminal · 终端

## 用途

显示终端网格：调用方把 PTY 输出写入 `TerminalState.Write`，组件解析 ANSI 样式（`internal.ParseAnsi`）、光标移动与清屏，保留滚动缓冲，支持鼠标选择复制，并把键盘输入编码后交给调用方写回 PTY。PTY 本身由调用方提供。

## 最小示例

```go
var term code.TerminalState
term.Write("\x1b[1m~/court\x1b[0m $ ls\r\n\x1b[34mdecrees\x1b[0m  ledger.txt\r\n$ ")
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	code.Terminal(c, &term, code.TerminalOptions{
		Input:   func(data string) { fmt.Printf("to pty: %q\n", data) },
		Resized: func(cols, rows int) { fmt.Println("pty size", cols, rows) },
	}).Size(380, 200)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Terminal(c *ui.Context, s *TerminalState, opts TerminalOptions) ui.Element`。`TerminalState` 零值是 80×24 空屏，由调用方持有；`Write(data)` 写入输出（须在界面线程，例如经 `Window.Update`），`Resize(cols, rows)`、`Lines()`、`LineText(i)`、`Select(a, b)`、`Selection()`、`ScrollBack(n)`、`Offset()`、`Cursor()`、`Size()` 供程序读取与控制。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Input` | `func(data string)` | nil | 用户输入编码后的字节：文字原样，Enter 为 `\r`，退格 `\x7f`，方向键 `ESC [ A–D`，Ctrl+字母为控制码。 |
| `Resized` | `func(cols, rows int)` | nil | 网格尺寸随元素大小变化时回调，用于同步 PTY 窗口大小。 |
| `Scrollback` | `int` | 1000 | 滚动缓冲最多保留的行数；负数 panic。 |
| `Label` | `string` | 内置“终端” | 无障碍名称。 |
| `Disabled` | `bool` | false | 不接收任何输入。 |

## 状态

默认（Surface 底、细边框）、聚焦（外侧 2 DIP 焦点环，光标为 AccentText 实心块并反白字符）、未聚焦（光标只描边）、选择（Selection 背景）、回看滚动缓冲（隐藏光标）、输入法组合（光标处显示带下划线的组合文字）、禁用。

## 事件

`Input`：每次按键、输入文字或粘贴；`Resized`：网格尺寸变化。复制写入系统剪贴板。

## 键盘操作

- 焦点在终端时，文字、Enter、Backspace、Tab、Shift+Tab、Escape、方向键、Home/End、Delete、PageUp/PageDown 与 Ctrl+字母都发送给 PTY。
- 复制：macOS 上 Cmd+C，其他平台 Ctrl+Shift+C（复制选中文字）；粘贴：Cmd+V / Ctrl+Shift+V。
- Shift+PageUp / Shift+PageDown：翻看滚动缓冲；滚轮同样可用；任意输入回到底部。

## 限制

- 状态机覆盖常用序列：SGR 样式、CUU/CUD/CUF/CUB、CHA、CUP、ED（0/1/2/3）、EL、CR/LF/BS/TAB 与 OSC（忽略）。其余 CSI 序列只记录 Debug 日志后忽略；不支持备用屏幕、滚动区域、鼠标上报与宽字符占两格。
- 需要完整终端仿真时请使用 MyGo 的 `plugins/terminal`（基于 libghostty-vt）。
- 字形用 `ui.Shape` 逐格排布并按行缓存，等宽字体由系统回退决定。
