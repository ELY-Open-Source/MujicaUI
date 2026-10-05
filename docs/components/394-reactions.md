# Reactions

## 用途

显示回应计数、当前用户标记与添加选择器。

## 最小示例

```go
items := []mujica.Reaction{{Emoji: "\U0001f44d", Count: 3, Mine: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Reactions(c, items, mujica.ReactionsOptions{Choices: []string{"\U0001f44d", "\U0001f389"}})
}
ui.Render(view, 400, 100, 1)
```

## 参数

Emoji 为显示字符；Count 必须为正；Choices 必填；People 用于提示；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Toggled 为切换意图；Added 为选择器选项。

## 键盘操作

Tab 到回应与添加；Enter / Space 激活；选择器左右导航；Escape 关闭。

## 限制

计数更新与持久化由调用方执行；传入系统字体支持的表情。
