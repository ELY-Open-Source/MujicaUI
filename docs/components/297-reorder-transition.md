# 297 ReorderTransition 重排过渡

## 用途

按稳定键排列子项，顺序变化时子项平滑移动到新位置；新增淡入，移除淡出。

## 最小示例

```go
order := []string{"Sakiko", "Uika", "Mutsumi"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.ReorderTransition(c, order, feedback.ReorderTransitionOptions{Gap: 4}, func(name string) {
		ui.Text(c, name)
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func ReorderTransition[K comparable](c *ui.Context, keys []K, opts ReorderTransitionOptions, item func(key K)) ui.Element`

## 参数

`ReorderTransitionOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Duration | time.Duration | 0（即 120 ms） | 移动时长；为负 panic。 |
| Gap | float32 | 0 | 子项间距。 |

键重复时 panic。基于 MyGo `Transition` 的 Position、Enter 与 Exit。

## 状态

- 减少动画：直接到位。动画结束后停止重绘。

## 事件

- 无。

## 键盘操作

- 焦点随键跟随子项移动。

## 限制

- 只做纵向列；虚拟化列表（List）的行不支持此过渡（MyGo 文档 transitions.md）。
