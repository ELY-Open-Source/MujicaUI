# 354 Playlist 播放列表

## 用途

播放列表：标出当前项，可拖动排序、移除，只构建可见行。

## 最小示例

```go
items := []mujica.PlaylistItem{{ID: "a", Title: "Overture"}, {ID: "b", Title: "Masquerade"}}
current := "a"
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.Playlist(c, &items, mujica.PlaylistOptions{Current: current}); r.Played != "" {
		current = r.Played
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Playlist(c *ui.Context, items *[]PlaylistItem, opts PlaylistOptions) PlaylistResult`

`PlaylistItem`：`ID`（唯一且非空，否则 panic）、`Title`、`Subtitle`、`Duration`。组件直接修改 `*items` 完成排序和移除。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Current` | `string` | "" | 正在播放项的 ID，显示播放标记和强调文字。 |
| `Height` | `float32` | 0（即 320） | 列表高度。 |
| `Disabled` | `bool` | false | 禁用选择、排序和移除。 |

## 状态

基于 List：选中行使用 Selection 底色加勾选标记；空列表显示提示。

## 事件

`PlaylistResult`：`Played` 为 Enter 或双击选择的 ID；`Removed` 为被移除的 ID；`Moved` 表示顺序变化。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| ↑ / ↓ | 移动选择 |
| Enter | 播放选中项 |
| Alt+↑ / Alt+↓ | 上移 / 下移选中项 |
| Delete / Backspace | 移除选中项 |

## 限制

- 拖动排序由 MyGo List 的 Reorder 提供，拖动时显示插入线。
