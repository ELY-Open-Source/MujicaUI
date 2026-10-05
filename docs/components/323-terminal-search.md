# 323 TerminalSearch · 终端搜索

## 用途

在终端的滚动缓冲与屏幕中查找文字，显示匹配数，上下跳转并把当前匹配滚动到可见处。

## 最小示例

```go
var term mujica.TerminalState
term.Write("build ok\r\nerror: missing seal\r\n")
var search mujica.TerminalSearchState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	ui.Column(c).Gap(8).Children(func() {
		mujica.TerminalSearch(c, &term, &search, mujica.TerminalSearchOptions{})
		mujica.Terminal(c, &term, mujica.TerminalOptions{}).Height(160)
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TerminalSearch(c *ui.Context, term *TerminalState, s *TerminalSearchState, opts TerminalSearchOptions) ChoiceResult`。匹配在查询或终端内容变化时重算一次，不逐帧计算；匹配结果写入 `term`，由同一 `TerminalState` 的 `Terminal` 绘制下划线，当前匹配加框。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `TerminalSearchState.Query` | `string` | "" | 查找文字，不区分大小写。 |
| `TerminalSearchState.Current` | `int` | 0 | 当前匹配序号，从 0 开始。 |
| `Disabled` | `bool` | false | 不能输入与跳转。 |

## 状态

无查询、有匹配（“第 n 项，共 m 项”）、无匹配（“无匹配”）、禁用；终端中其余匹配为 AccentText 下划线，当前匹配加 AccentText 框。

## 事件

`Changed()`：查询变化或跳到另一匹配。

## 键盘操作

- Enter / Shift+Enter：下一项 / 上一项；macOS 上 Cmd+G / Shift+Cmd+G，其他平台 F3 / Shift+F3。
- 上一项、下一项按钮可用 Tab 到达。

## 限制

- 只做纯文字匹配，不支持正则与整词；跨行文字不会匹配。
