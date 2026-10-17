# 017 Badge 徽标

## 用途

显示短文字或计数。计数超过上限时显示为“上限+”。徽标可以内联排列，也可以钉在内容的某个角上。辅助技术朗读完整数值。

## 最小示例

```go
unread := 120
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	display.Badge(c, display.BadgeOptions{Text: "New", Tone: display.BadgeNeutral}, nil)
	display.Badge(c, display.BadgeOptions{Count: unread, Label: "Unread", Position: display.BadgeTopRight}, func() {
		display.Avatar(c, "Sakiko Togawa", display.AvatarOptions{})
	})
}
```

## 参数

`func Badge(c *ui.Context, opts BadgeOptions, content func()) ui.Element`

`Position` 为 `BadgeInline` 时 `content` 必须为 `nil`；为角落位置时 `content` 必须不为 `nil`。否则 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Text` | `string` | `""` | 显示的文字；非空时忽略 `Count`。 |
| `Count` | `int` | `0` | 计数。`Text` 为空且计数为 0 时不显示徽标。负数 panic。 |
| `Max` | `int` | `0`（即 99） | 显示上限，超过时显示 `Max+`。负数 panic。 |
| `Tone` | `BadgeTone` | `BadgeAccent` | `BadgeAccent`、`BadgeNeutral`、`BadgeSuccess`、`BadgeWarning`、`BadgeDanger`、`BadgeInfo`。 |
| `Position` | `BadgePosition` | `BadgeInline` | `BadgeTopRight`、`BadgeTopLeft`、`BadgeBottomRight`、`BadgeBottomLeft`：徽标中心对准内容的对应角。 |
| `Label` | `string` | `""` | 计数的含义，朗读时放在完整数值前，如“Unread 1204”。 |

## 状态

纯展示。字号与最小尺寸随系统文本缩放放大。强调和危险色调使用实色填充（`OnAccent` 配 `Accent`，`OnDanger` 配 `DangerFill`）；成功、警告、信息色调只把语义色用在文字和边框上。

## 事件

无。

## 键盘操作

不获取焦点。徽标的无障碍名称为 `Label` 加完整数值。

## 限制

- 角落徽标的位置相对于 `content` 所在的盒子。内容本身有外边距时，徽标不会跟着偏移。
- 徽标使用 `RoleStatus`，部分屏幕阅读器可能在数值变化时主动播报。
