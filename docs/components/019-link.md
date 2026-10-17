# 019 Link 链接

## 用途

带下划线的链接文本：打开外部地址、在 MyGo Router 中跳转路径，或作为应用内动作。可以显示打开失败的反馈。

## 最小示例

```go
details := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	display.Link(c, "Read the charter", display.LinkOptions{URL: "https://example.com/charter"})
	if display.Link(c, "Show details", display.LinkOptions{}).Clicked() {
		details = !details
	}
}
```

## 参数

`func Link(c *ui.Context, label string, opts LinkOptions) ui.Element`

`label` 必填，空字符串 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `URL` | `string` | `""` | 有协议的地址用系统默认应用打开。没有协议的路径（如 `/notes/42`、`?tab=info`）在 MyGo Router 中跳转。为空时是应用内动作。 |
| `Error` | `string` | `""` | 调用方提供的打开失败说明，显示在链接下方，也报告给辅助技术。 |

## 状态

- 默认：`AccentText` 颜色，始终带 1 DIP 下划线。外部地址后面带“外部链接”小图标。
- 悬停 / 按下：文字与下划线变为 `Text` 颜色，布局不变。
- 聚焦：2 DIP 焦点环。
- 错误：格式不正确的地址（无法解析，或 http/https 缺少主机名）不会打开，下方显示“无法打开此链接：地址”；`Error` 非空时优先显示 `Error`。

## 事件

- `Clicked()`：每次点击或按 Enter 都会报告，无论是否打开了地址。同时记录 `slog.Debug("mujica: link activate")`。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦链接 |
| Enter | 激活 |

## 限制

- MyGo v0.2.10 不报告打开失败。`ui.Link` 和 `Context.OpenURL` 不返回结果（`ui/context.go:270`、`ui/window.go:583`）。窗口宿主在 goroutine 中调用 `Shell.OpenExternal` 并丢弃它的错误（`content.go:55`）。因此组件只能自行检查地址格式，系统打开失败无法自动显示。需要这种反馈时，传空 `URL`，在 `Clicked()` 后由调用方在 goroutine 中调用 `mygo.Shell.OpenExternal(url)`。它会返回 error，把结果通过 `Window.Update` 写回 `Error`。
- 路径是否交给 Router 由 MyGo 决定。没有 Router 时，路径会被当作 URL 交给系统打开。
- 链接不能嵌入 `ui.RichText` 句子中作为行内元素，因为外部图标和错误说明需要额外的行布局。行内链接请直接使用 `ui.Link`。
