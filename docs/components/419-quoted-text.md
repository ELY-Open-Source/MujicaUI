# 419 QuotedText 引用文字

## 用途

折叠或展开引用文字，并去掉一层引用前缀。

## 最小示例

```go
open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.QuotedText(c, "> Can you send the score?", &open)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`text string` 提供引用文字，`open *bool` 控制展开。`QuoteStrip` 去掉每行一个 `>` 及其后一个空格。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

返回 `ChoiceResult`；`Changed()` 表示折叠状态变化。

## 键盘操作

Tab 聚焦开关；Enter / Space 展开或折叠；展开内容支持选择复制。

## 限制

文本按一层引用规则处理，嵌套引用继续保留标记。
