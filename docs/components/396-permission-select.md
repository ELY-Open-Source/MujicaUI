# PermissionSelect

## 用途

选择可查看、可评论、可编辑。

## 最小示例

```go
value := files.PermissionView
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.PermissionSelect(c, &value, files.PermissionSelectOptions{})
}
ui.Render(view, 300, 100, 1)
```

## 参数

PermissionView 为零值；Label 默认权限；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

ChoiceResult.Changed() 表示值已变化。

## 键盘操作

Enter 打开；上下移动；Enter 确认；Escape 关闭。

## 限制

权限枚举越界会 panic；授权操作由调用方处理。
