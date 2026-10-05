# 303 ThreeWayMerge 三方合并

## 用途

以本地、基准、远端三栏和结果栏显示三方合并，逐块采纳本地、远端、两者或基准。

## 最小示例

```go
var state git.ThreeWayMergeState
base, local, remote := "a\nb\nc", "a\nB1\nc", "a\nB2\nc"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.ThreeWayMerge(c, &state, base, local, remote, git.ThreeWayMergeOptions{})
	r.Element.Height(280)
	if done, total := state.Resolved(); r.Changed() && done == total {
		_ = state.Result() // write the merged file
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ThreeWayMerge(c *ui.Context, s *ThreeWayMergeState, base, ours, theirs string, opts ThreeWayMergeOptions) ChoiceResult`

## 参数

`ThreeWayMergeOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Disabled | bool | false | 禁用采纳。 |

`ThreeWayMergeState.Result()` 返回合并文本，未解决的冲突保留冲突标记；`Resolved()` 返回已解决与总冲突数。分块只在三段文本改变时计算。给返回的 `Element` 一个高度。

## 状态

- 只有一侧修改或两侧相同修改的块自动合并，来源一侧以浅色底标出。
- 冲突块在结果栏显示 Danger 色“未解决”与警示图标。
- 已采纳的按钮为 Selection 底、AccentText 边框与对勾。
- 较长的未变更块只显示首尾各两行。

## 事件

- `Changed()`：某个冲突的选择改变。

## 键盘操作

- 列表可获得焦点并滚动；采纳按钮用 Tab 到达，Enter 或 Space 选择。

## 限制

- 按行合并，不做词级合并。
- 相邻（中间没有未变更行）的两侧修改视为冲突。
