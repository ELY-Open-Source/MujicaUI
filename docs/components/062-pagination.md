# 062 Pagination 分页

## 用途

在分页数据之间移动：页码、上一页、下一页、页容量选择与跳页。总数变化后，超出范围的页码自动修正到最后一页并报告变化。

## 最小示例

```go
if mujica.Pagination(c, &page, &size, mujica.PaginationOptions{Total: 420, PageSizes: []int{10, 20, 50}, Jump: true}).Changed() {
	load(page, size)
}
```

签名：`func Pagination(c *ui.Context, page, pageSize *int, opts PaginationOptions) PaginationResult`

## 参数

`page`：当前页，从 1 开始；每帧被限制在 `[1, 总页数]`。`pageSize`：每页条数，必须为正，否则 panic。

`PaginationOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Total` | `int` | 0 | 记录总数；为负 panic。为 0 时仍显示 1 页。 |
| `PageSizes` | `[]int` | nil | 可选的页容量；nil 隐藏选择器。非 nil 时须包含 `*pageSize`，否则 panic。 |
| `Jump` | `bool` | false | 显示跳页输入框。 |
| `Siblings` | `int` | 0（即 1） | 当前页两侧显示的页码数；为负 panic。 |
| `Disabled` | `bool` | false | 整体禁用。 |

返回 `PaginationResult`，内嵌 `*ui.Element`。

## 状态

- 当前页：Accent 填充、OnAccent 文字、600 字重，名称附“当前”。
- 页码按钮：Border 细边框，悬停 SurfaceHover 与 ControlBorder 边框；首页时“上一页”禁用，末页时“下一页”禁用。
- 页码多时首尾保留，中间以“…”省略，按钮数固定为 `2×Siblings+5`，位置稳定。
- 跳页错误：输入非数字或超出范围时边框变为 Danger，并在旁显示有效范围；输入保留，可直接修改。
- 修改页容量时，保持原页第一条记录可见。

## 事件

- `Changed()`：本帧页码或页容量改变（用户操作或总数变化导致的修正）。

## 键盘操作

- 页码按钮组是一个 Tab 停靠点：←/→ 移动，Home / End 到两端，Enter 或 Space 跳到该页。
- 页容量：↑/↓ 打开并移动，Enter 选择，Escape 关闭。
- 跳页框：输入页码后按 Enter。

## 限制

- 页容量选择基于 MyGo `SelectBase`，弹层由 MujicaUI 绘制。
- 弹层打开时在 160 ms 内淡入，减少动画时直接显示。关闭时立即消失：弹层关闭后 MyGo 不再构建其元素，公开接口无法为已移除的浮层做淡出。
