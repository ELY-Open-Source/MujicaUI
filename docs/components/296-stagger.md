# 296 Stagger 依次进场

## 用途

子项按间隔依次进场，只在首次出现时执行一次。

## 最小示例

```go
guests := []string{"Sakiko", "Uika", "Mutsumi"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.Stagger(c, len(guests), feedback.StaggerOptions{}, func(i int) {
		ui.Text(c, guests[i])
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Stagger(c *ui.Context, n int, opts StaggerOptions, item func(i int)) ui.Element`

## 参数

`StaggerOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Interval | time.Duration | 0（即 60 ms） | 相邻子项间隔。 |
| Duration | time.Duration | 0（即 PopupDuration） | 每项进场时长。 |
| Horizontal | bool | false | 横向排列。 |

n、Interval、Duration 为负时 panic。基于 MyGo `Transition` 的 Enter。

## 状态

- 每项从 8 DIP 外淡入。
- 序列结束后新增的项照常进场，不再等待。
- 减少动画：一次全部显示。按定时请求帧。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 尚未出现的项不占位，容器随之增长。
