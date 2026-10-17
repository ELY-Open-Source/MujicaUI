# 320 AnsiText ANSI 文本

## 用途

渲染终端输出：16 色、256 色、真彩，以及粗体、斜体、下划线、删除线、暗淡与反显。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.AnsiText(c, "\x1b[1;32mPASS\x1b[0m 12 tests\n\x1b[31mFAIL\x1b[0m 1 test", git.AnsiTextOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func AnsiText(c *ui.Context, text string, opts AnsiTextOptions) ui.Element`

## 参数

`AnsiTextOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Wrap | bool | false | 按宽度换行；否则每行保持完整。 |

样式跨行延续。

## 状态

- 16 色映射为在两种主题 Surface 上都达到 4.5:1 的颜色；256 色与真彩按原值显示。
- 默认前景为 Text；反显以 Surface 为前景。

## 事件

- 无。

## 键盘操作

- 无：纯展示。

## 限制

- 只解释 SGR（m 结尾）序列；光标移动、清屏等其他 CSI 序列原样保留为文字，并以 slog.Debug 记录。
- 256 色与真彩不做对比度调整。
