# 229 ContextChips 上下文标签

## 用途

已附加上下文列表；移除；溢出折叠。

## 最小示例

```go
items := []mujica.ContextItem{{ID: "f1", Label: "guests.csv", Kind: mujica.ContextFile}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if id, ok := mujica.ContextChips(c, items, mujica.ContextChipsOptions{}).Removed(); ok {
		fmt.Println("detach", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`items` 为已附加上下文，`ID`、`Label` 必填，`Kind` 指定图标；`Max=0` 默认显示 4 项，超出时可展开。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Removed() 返回待移除 ID；展开与折叠保留在组件状态。

## 键盘操作

Tab 到移除与展开；Enter、Space 激活。

## 限制

上下文文件内容由调用方管理。
