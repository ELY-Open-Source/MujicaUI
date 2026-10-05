# 472 WorkspaceSwitcher

工作区切换。

## 用途

显示当前工作区，弹层中搜索并切换。

## 最小示例

```go
current := "ave"
spaces := []mujica.Workspace{{ID: "ave", Name: "Ave Mujica"}, {ID: "mygo", Name: "MyGO!!!!!"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WorkspaceSwitcher(c, &current, spaces, mujica.WorkspaceSwitcherOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `current` | `*string` | 必填 | 当前工作区 ID。 |
| `spaces` | `[]Workspace` | 必填 | ID、Name 必填，Detail 可选。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

关闭、展开、搜索、无结果、当前项勾选、禁用。

## 事件

`ChoiceResult.Changed()`。

## 键盘操作

打开后焦点在搜索框；↑/↓ 移动，Enter 切换，Escape 关闭。

## 限制

大集合使用虚拟列表，构建当前可见行；搜索与分组在数据或查询变化时更新。

- 无。
