# 039 Cascader · 级联选择

## 用途

逐级选择一条路径，如“欧洲 / 法国 / 巴黎”。字段显示完整路径；子级可由调用方按需加载，并有加载中、失败与重试状态。

## 最小示例

```go
var path []string
level := func(p []string) input.CascaderLevel {
	if len(p) == 0 {
		return input.CascaderLevel{Nodes: []input.CascaderNode{{Value: "fr", Label: "France"}}}
	}
	return input.CascaderLevel{Nodes: []input.CascaderNode{{Value: "paris", Label: "Paris", Leaf: true}}}
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.Cascader(c, &path, input.CascaderOptions{Level: level, Label: "City"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Cascader(c *ui.Context, path *[]string, opts CascaderOptions) ChoiceResult`。`*path` 是各级 `Value` 组成的路径。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Level` | `func(path []string) CascaderLevel` | 必填 | 返回某路径下的子级（空路径为根）。为 nil 时 panic。 |
| `Load` | `func(path []string)` | nil | 子级状态为 `CascaderIdle` 时调用一次，点击重试时再次调用。应把该级标为 `CascaderLoading`，结果经 `Window.Update` 写回。有 Idle 级别但 `Load` 为 nil 时 panic。 |
| `AnyLevel` | `bool` | false | 允许选择非叶节点；默认只有叶节点（`Leaf: true`）可成为值。 |
| `Placeholder` | `string` | 内置“请选择…” | 无值时提示。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "" | 无障碍名称。 |

`CascaderLevel`：`Nodes []CascaderNode`、`Status CascaderStatus`（零值 `CascaderLoaded`；另有 `CascaderIdle`、`CascaderLoading`、`CascaderFailed`）、`Err error`（失败原因，显示在重试按钮旁）。

`CascaderNode`：`Value`、`Label`、`Leaf`（无子级）、`Disabled`。

## 状态

字段：占位、完整路径、展开、聚焦、禁用。弹层每列：加载中（旋转标记）、失败（错误文字 + 重试按钮）、空（“暂无数据”）、节点列表；展开路径上的节点用玫瑰底，已选路径用酒红文字，分支节点带右箭头。

## 事件

`Changed()` 在路径变化的那一帧返回 true。叶节点选择后弹层关闭。

## 键盘操作

- Enter / Space / ↓：打开弹层，焦点在当前路径或第一项。
- 列内 ↑/↓、Home、End：移动。
- →：展开分支并进入下一列；←：回到上一列的父节点。
- Enter / Space：叶节点选择并关闭；分支展开（`AnyLevel` 时同时选择）。
- Escape：关闭，焦点回到字段。

## 限制

- `Load` 在视图构建中同步调用，不得阻塞；异步任务完成后请用 `Window.Update` 回到界面线程更新数据。
- 列表未虚拟化，每列适合数百项以内。
- 已选路径上的节点通过 `Selected` 报告选中；当前展开的分支无法报告 expanded 状态（ui/base.go:297 等处仅内部设置），仅在画面上用底色表示。
