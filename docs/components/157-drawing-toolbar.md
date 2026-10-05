# 157 DrawingToolbar

## 用途

画线工具栏；当前工具；锁定与清除；键盘快捷键提示。

## 最小示例

```go
	var drawings mujica.DrawingTools
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DrawingToolbar(c, &drawings, mujica.DrawingToolbarOptions{Horizontal: true})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

传入 DrawingTools；Horizontal=true 横向，默认纵向；工具名称附快捷键提示。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化。

## 键盘操作

Tab 进入工具栏；方向键移动焦点，Enter 激活；Alt+1 至 Alt+6 选择工具。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
