# 279 CheckpointList 检查点列表

## 用途

列出 Agent 运行的检查点，标出当前点，回滚到其它检查点前要求确认。

## 最小示例

```go
current := "c2"
points := []mujica.Checkpoint{{ID: "c2", Label: "Menu drafted"}, {ID: "c1", Label: "Start"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.CheckpointList(c, &current, points, mujica.CheckpointListOptions{}).Changed() {
		fmt.Println("restore", current)
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func CheckpointList(c *ui.Context, current *string, points []Checkpoint, opts CheckpointListOptions) ChoiceResult`。`current` 为 nil、`ID` 为空或重复、`Label` 为空、`MaxHeight` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `MaxHeight` | `float32` | 360 | 列表最大高度。 |
| `Disabled` | `bool` | `false` | 禁止回滚。 |

## 状态

- 当前点用实心菱形、选中底色和“当前”文字标出；其它点为空心菱形。
- 回滚用 AlertDialog 确认：默认聚焦取消，Enter 不会直接回滚。
- 超过 200 个检查点时虚拟化。

## 事件

- `Changed()`：确认回滚后 `*current` 已改为目标 ID。

## 键盘操作

- Tab 到各“回滚”，Enter 打开确认；确认框内 Tab 到“回滚”再按 Enter。

## 限制

- 恢复文件与状态由调用方完成。
