# 072 Tree · 树

## 用途

展示层级数据：节点展开、选择、键盘导航、外部加载子项与空状态。只构建可见行，10 万节点同样流畅；折叠的子树不创建任何界面元素。

## 最小示例

```go
kids := map[string][]string{"Nave": {"Pulpit", "Pews"}}
var state mujica.TreeState[string]
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Tree(c, &state, mujica.TreeOptions[string]{
		Roots:    []string{"Nave", "Sacristy"},
		Children: func(n string) []string { return kids[n] },
		Label:    func(n string) string { return n },
	}).Size(280, 240)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Tree[K comparable](c *ui.Context, s *TreeState[K], opts TreeOptions[K]) TreeView`。节点值本身即稳定键。`TreeState` 方法：`Selected()`、`SetOpen(item, open)`、`IsOpen(item)`；字段 `Selection ui.Selection[K]` 保存多选。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Roots` | `[]K` | 空 | 顶层节点；为空时显示空状态。 |
| `Children` | `func(K) []K` | 必填 | 子节点：nil 表示叶子；未加载的节点返回空的非 nil 切片。 |
| `Label` | `func(K) string` | 必填 | 节点文字。 |
| `ItemStatus` | `func(K) DataStatus` | nil | 子项加载状态：`DataUnloaded`、`DataLoading`（行尾加载标记）、`DataFailed`（行尾错误与重试）；nil 表示全部已加载。 |
| `Load` | `func(K)` | nil | 节点以 `DataUnloaded` 状态展开时调用一次，点击重试时再次调用。 |
| `Multiple` | `bool` | false | 多选。 |
| `Status` | `DataStatus` | `DataReady` | 加载状态：`DataLoading` 首次加载（无行时显示加载标记）、`DataLoadingMore` 追加加载（行下方显示“正在加载更多”）、`DataFailed` 失败（显示错误与“重试”）。 |
| `Error` | `string` | 内置“加载失败” | `DataFailed` 时显示的具体错误。 |
| `Empty` | `string` | 内置“暂无数据” | 无数据时的说明文字。 |

## 状态

默认、悬停、选中（Accent 填充 + 勾选标记）、聚焦、展开/折叠（箭头旋转）、子项加载中、子项加载失败、空、首次加载与失败。

## 事件

`Changed()`：选择变化；`Submitted()`：Enter 或双击；`Retried()`：整树失败时点击重试。子项重试直接调用 `Load`。

## 键盘操作

- ↑/↓：移动；→：展开或进入第一个子节点；←：折叠或回到父节点；macOS 上 Option、其他平台 Shift 配合 ←/→ 展开或折叠全部子孙。
- 输入字母按 `Label` 定位；Enter 提交。

## 限制

- 基于 MyGo `ui.Outline` 而不是 `ui.Tree`：`ui.Tree` 会构建全部已展开节点（docs/ui/tree.md “A tree builds every item it shows”），不满足大量节点虚拟化要求。
- `Outline` 每帧遍历已展开的节点以计算行（ui/outline.go `flatten`），元素数受可见区域限制，但遍历成本随展开节点数增长。
- 单选跟随节点的搜索范围为 1000 行（ui/list.go:30）。
- 异步加载由调用方在后台完成并通过 `Window.Update` 回到界面线程。
