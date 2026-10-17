# 496 Masonry 瀑布流

## 用途

按可用宽度决定列数，把每一项放进当前最短的列；只构建可见范围内的项，以稳定键记住测得的高度。

## 最小示例

```go
titles := []string{"Sea", "Forest", "Desert", "Glacier"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.Masonry(c, len(titles), input.MasonryOptions{
		Key: func(i int) string { return titles[i] },
	}, func(i int) {
		ui.Text(c, titles[i])
	}).Size(380, 280)
}
```

签名：`func Masonry(c *ui.Context, n int, opts MasonryOptions, item func(i int)) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `n` | `int` | 必填 | 项数；负数 panic。 |
| `Key` | `func(int) string` | 必填 | 稳定且唯一的键；重复 panic。 |
| `ColumnWidth` | `float32` | 220 | 最窄列宽。 |
| `Gap` | `float32` | 3 个间距单位 | 间隔。 |
| `Estimate` | `float32` | 160 | 测量前的估计高度。 |

## 状态

自身是滚动容器，需由调用方给定尺寸。

## 事件

无。

## 键盘操作

滚动容器按 MyGo 规则滚动；项内控件按各自规则聚焦。

## 限制

排布只在项数、列数、宽度、测得高度或可见项的键变化时重算。
