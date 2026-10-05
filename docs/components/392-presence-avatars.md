# PresenceAvatars

## 用途

显示协作者、溢出人数及跟随入口。

## 最小示例

```go
people := []mujica.Collaborator{{ID: "ada", Name: "Ada Lovelace"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PresenceAvatars(c, people, mujica.PresenceAvatarsOptions{})
}
ui.Render(view, 300, 100, 1)
```

## 参数

Max 默认 4；Size 默认 28；Following 为跟随 ID；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Follow 返回 ID；Unfollow 请求停止跟随。

## 键盘操作

Tab 进入；左右移动；Enter 激活；Escape 关闭溢出名单。

## 限制

在线数据与跟随行为由调用方维护。
