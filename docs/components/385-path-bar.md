# PathBar

## 用途

分段跳转，长路径折叠中间段。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PathBar(c, "/court/chronicles/seals", mujica.PathBarOptions{})
}
ui.Render(view, 400, 100, 1)
```

## 参数

Separator 默认 /；RootLabel 默认分隔符；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Navigated 返回目标祖先路径。

## 键盘操作

Tab 进入；左右移动；Enter 激活；Escape 关闭折叠菜单。

## 限制

路径必须至少包含一个分段；调用方处理导航。
