# 070 List · 列表

## 用途

展示任意数量的条目，只构建可见行（加预取两行），10 万条记录也能流畅滚动、选择。支持列表与网格两种模式、单选与多选、稳定条目键、按标签输入定位，以及首次加载、追加加载、失败与空状态。

## 最小示例

```go
names := []string{"Matins", "Lauds", "Prime", "Terce"}
var state mujica.ListState[string]
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	v := mujica.List(c, &state, len(names), mujica.ListOptions[string]{
		Key:   func(i int) string { return names[i] },
		Label: func(i int) string { return names[i] },
	}, func(i int) { ui.Text(c, names[i]) })
	v.Size(240, 200)
	if v.Changed() {
		key, _ := state.Selected()
		fmt.Println("chose", key)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func List[K comparable](c *ui.Context, s *ListState[K], n int, opts ListOptions[K], item func(i int)) ListView`。`ListView` 嵌入 `*ui.Element`；`ListState` 由调用方持有，零值表示未选择。`state.Selected()` 返回所选键，多选时 `state.Selection` 保存全部键；`state.ScrollTo(i)` 滚动到第 i 项。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Key` | `func(i int) K` | 必填 | 条目的稳定键；选择与行状态跟随键，而不是下标。 |
| `Label` | `func(i int) string` | 必填 | 条目文字：无障碍名称与输入定位。 |
| `Multiple` | `bool` | false | 多选：Cmd/Ctrl 点击增减、Shift 连选、Cmd/Ctrl+A 全选。 |
| `Grid` | `bool` | false | 网格模式：按 `ItemWidth` 自动分列。 |
| `ItemWidth`, `ItemHeight` | `float32` | 160, 120 | 网格条目最小宽度与高度（DIP）。 |
| `Status` | `DataStatus` | `DataReady` | 加载状态：`DataLoading` 首次加载（无行时显示加载标记）、`DataLoadingMore` 追加加载（行下方显示“正在加载更多”）、`DataFailed` 失败（显示错误与“重试”）。 |
| `Error` | `string` | 内置“加载失败” | `DataFailed` 时显示的具体错误。 |
| `Empty` | `string` | 内置“暂无数据” | 无数据时的说明文字。 |

## 状态

默认、悬停（SurfaceHover）、选中（Selection 背景 + Text 文字 + AccentText 勾选标记）、聚焦（列表外 2 DIP 焦点环）、网格选中（Selection 背景 + AccentText 边框与勾选）、加载、追加加载、失败、空。

## 事件

`ListView.Changed()`：选择变化；`Submitted()`：Enter 或双击；`Retried()`：点击失败状态的“重试”。

## 键盘操作

- ↑/↓、Home/End，Linux 与 Windows 上 PageUp/PageDown：移动选择；网格模式下方向键二维移动。
- 输入字母：按 `Label` 定位。
- Enter：提交。
- 多选：Shift+方向键连选；Cmd/Ctrl+A 全选；Linux/Windows 上 Ctrl+方向键只移动、Ctrl+Space 切换。

## 限制

- 虚拟化由 MyGo `ui.List`/`ui.GridView` 提供：每帧构建的元素数由可见区域和预取范围决定；程序跳转（`ScrollTo`）时目标两侧各构建一屏。测试以 10 万条记录统计“最后一帧调用条目构建函数的次数”验证上限（list_test.go）。
- 单选按键保存：每帧先检查原位置的键，不符时扫描全部条目重新定位（MyGo 自身只在附近 1000 行内查找，ui/list.go:30），所以反转、排序或在网格前端插入后选择仍指向同一条目；条目被删除时选择清除。
- Shift 连选的锚点由 MyGo 决定（ui/list.go `extend`），与平台习惯一致。
- 网格模式没有 `ListState.Visible` 等滚动位置读取接口。
