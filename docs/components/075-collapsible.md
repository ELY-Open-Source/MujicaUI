# 075 Collapsible · 折叠面板

## 用途

一个触发区控制下方内容的显示。展开状态由调用方的 `*bool` 控制，可在外部打开或关闭；关闭后的内容不构建，因此退出焦点序列。

## 最小示例

```go
open := false
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Collapsible(c, &open, mujica.CollapsibleOptions{Title: "Details"}, func() {
		ui.Text(c, "Candles are lit at the seventh bell.")
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Collapsible(c *ui.Context, open *bool, opts CollapsibleOptions, content func()) CollapsibleView`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 空 | 触发区文字；与 `Trigger` 至少设置一个。 |
| `Trigger` | `func()` | nil | 自定义触发区内容（箭头之后）。 |
| `Disabled` | `bool` | false | 禁止展开或关闭，退出 Tab 序列。 |

## 状态

默认、悬停、按下、聚焦、展开（箭头旋转 90°，120 ms）、禁用。

## 事件

`CollapsibleView.Changed()`：用户点击或按键切换了展开状态。

## 键盘操作

- Enter/Space：展开或关闭。
- 关闭（包括关闭动画期间）的内容不在 Tab 序列中。

## 限制

- MyGo `CollapsibleBase` 自带 200 ms 动画计时（ui/collapsible.go），组件使用自己的 120 ms 进度绘制，但 MyGo 的计时仍会在 200 ms 内请求帧。
