# 013 Icon 图标

## 用途

按名称显示 `icons` 包内置的 Lucide SVG 图标，可设尺寸与语义颜色。名称不存在时显示可见的回退标记并写警告日志。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Row(c).Gap(8).Children(func() {
		display.Icon(c, "star", display.IconOptions{})
		display.Icon(c, "circle-check", display.IconOptions{Size: 20, Tone: display.IconSuccess, Label: "Done"})
	})
}
```

## 参数

`func Icon(c *ui.Context, name string, opts IconOptions) ui.Element`

`name` 必填，空字符串 panic。可用名称见 `icons.Names()`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Size` | `float32` | `0`（即 16） | 宽高（DIP），建议 16、20、24。负数 panic。 |
| `Tone` | `IconTone` | `IconInherit` | `IconInherit` 沿用文字颜色；另有 `IconMuted`、`IconAccent`、`IconSuccess`、`IconWarning`、`IconDanger`、`IconInfo`。未知值 panic。 |
| `Label` | `string` | `""` | 有含义的图标的无障碍名称；为空时图标视为装饰，不进入无障碍树。 |

## 状态

- 默认：按语义颜色显示图标。
- 加载失败（名称不存在）：在同尺寸方框中用 `Danger` 色绘制虚线框与斜线，提示文字为“缺少图标：名称”，无障碍名称相同。每个元素对每个缺失名称记录一次 `slog.Warn("mujica: icon missing", "name", ...)`。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 只能使用 `icons` 包内置的图标。要显示任意 SVG，可直接用 `ui.Icon`。
- 回退标记带悬停提示，但图标本身不获取焦点，键盘用户看不到这个提示，只能依赖无障碍名称。
