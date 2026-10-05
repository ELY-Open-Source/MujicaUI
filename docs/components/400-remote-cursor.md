# RemoteCursor

## 用途

绘制远程光标、选区与姓名标签。

## 最小示例

```go
carets := []mujica.RemoteCaret{{ID: "ada", Name: "Ada", Offset: 4, Anchor: 4}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	ui.Box(c).Size(300, 100).Children(func() {
		mujica.RemoteCursor(c, carets, mujica.RemoteCursorOptions{
			Caret: func(offset int) ui.Rect { return ui.Rect{X: float32(offset) * 8, Y: 20, H: 20} },
		})
	})
}
ui.Render(view, 320, 120, 1)
```

## 参数

ID 唯一且必填；Offset、Anchor 非负；Caret 必填；Selection 返回选区矩形。

## 状态

光标、选区、姓名；覆盖层输入穿透。

## 事件

展示组件。

## 键盘操作

输入穿透至下方编辑器。

## 限制

坐标映射、远程状态由调用方提供；覆盖层放在编辑器后方。
