# 012 Heading 标题

## 用途

页面与区块的三级衬线标题，可带副标题与旧金色装饰分隔线（对称细线加中央菱形）。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	display.Heading(c, "The Grand Hall", display.HeadingOptions{Subtitle: "Winter ceremonies", Divider: true})
	display.Heading(c, "Guest Ledger", display.HeadingOptions{Level: 2})
}
```

## 参数

`func Heading(c *ui.Context, title string, opts HeadingOptions) ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Level` | `int` | `0`（即 1） | 1、2、3 级，对应 28、22、18 DIP。其他值 panic。 |
| `Subtitle` | `string` | `""` | 标题下方的说明文字，使用 `TextMuted`。 |
| `Divider` | `bool` | `false` | 在下方绘制装饰分隔线，颜色为 `Ornament`。 |

## 状态

纯展示，仅默认状态。字号按系统文本缩放（`Preferences().TextScale`）放大。

## 事件

无。

## 键盘操作

不获取焦点。装饰线不进入无障碍树。

## 限制

- 标题字体为 `theme.SerifFont`（Georgia、宋体类回退）；系统缺字时由系统回退字体显示。
- MyGo v0.2.10 没有 Heading 角色（`ui/access.go` 的 Role 列表中没有），辅助技术把标题读作普通文本，无法按标题层级导航。
- 分隔线只在 Heading 内绘制，不依赖 Divider 组件。
