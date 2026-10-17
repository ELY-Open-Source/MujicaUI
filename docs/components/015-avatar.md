# 015 Avatar 头像

## 用途

圆形人物头像。依次回退：图片 → 姓名首字母 → 默认人物图标。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Row(c).Gap(8).Children(func() {
		display.Avatar(c, "Sakiko Togawa", display.AvatarOptions{Size: 40})
		display.Avatar(c, "", display.AvatarOptions{})
	})
}
```

## 参数

`func Avatar(c *ui.Context, name string, opts AvatarOptions) ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Image` | `*ui.Bitmap` | `nil` | 头像图片，按圆形裁切填充。 |
| `Size` | `float32` | `0`（即 32） | 直径（DIP）。负数 panic。 |

`name` 用作无障碍名称和首字母来源；为空时无障碍名称是“未命名”。

## 状态

- 图片：圆形裁切。
- 首字母：取第一个和最后一个词的首字母（大写），衬线字体，`AccentText` 文字配 `Selection` 底色。
- 无可用首字母：`SurfaceHover` 底色加人物图标。

## 事件

无。

## 键盘操作

不获取焦点。辅助技术把头像读作以 `name` 命名的图片。

## 限制

- 首字母字号随头像直径而定，不随系统文本缩放放大，以免溢出圆形。

- 中文姓名没有空格时按一个词处理，只取第一个字，例如“张伟”显示为“张”。
- 头像不使用 MyGo `ui.Avatar`，因为它的尺寸固定为字号的 2.25 倍，无法配置（`ui/indicators.go:383`）。
