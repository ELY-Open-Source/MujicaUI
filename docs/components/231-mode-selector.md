# 231 ModeSelector 模式选择器

## 用途

对话、搜索、深度研究、Agent 模式切换。

## 最小示例

```go
mode := mujica.ModeChat
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ModeSelector(c, &mode, mujica.ModeSelectorOptions{}).Changed() {
		fmt.Println("mode", mode)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`mode` 指向 ModeChat、ModeSearch、ModeResearch 或 ModeAgent；`IconsOnly=false` 保留文字；`Disabled=false`。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

ChoiceResult.Changed() 表示模式已改变。

## 键盘操作

Tab 进入模式组；Left/Right 切换模式。

## 限制

模式执行逻辑由调用方实现。
