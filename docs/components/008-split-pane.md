# 008 SplitPane 分栏

## 用途

将区域分为左右（或上下）两栏，中间分隔条可拖动、可用键盘调整；两栏有各自最小尺寸，可指定一栏可折叠。

## 最小示例

```go
split := mujica.SplitPaneState{Size: 160}
mujica.SplitPane(c, &split, mujica.SplitPaneOptions{MinFirst: 120, Collapsible: mujica.SplitPaneFirst},
	func() { ui.Text(c, "Files") },
	func() { ui.Text(c, "Editor") },
).Element.Fill()
```

签名：`func SplitPane(c *ui.Context, state *SplitPaneState, opts SplitPaneOptions, first, second func()) SplitPaneResult`

## 参数

`SplitPaneState`（应用持有，nil 时 panic）

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Size` | `float32` | 0（首次布局后取一半） | 第一栏宽度（纵向时为高度），DIP。每帧被限制在最小尺寸之间。 |
| `Collapsed` | `bool` | false | 可折叠栏是否折叠；折叠时保留 `Size` 以便恢复。未设置 `Collapsible` 却为 true 时 panic。 |

`SplitPaneOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Vertical` | `bool` | false | 第一栏在上。 |
| `MinFirst`, `MinSecond` | `float32` | 10 单位（默认 40 DIP） | 两栏最小尺寸；负数 panic。 |
| `Collapsible` | `SplitPaneSide` | `SplitPaneNeither` | 可折叠的栏：`SplitPaneFirst` 或 `SplitPaneSecond`。 |

返回 `SplitPaneResult`，`Element` 字段是组件元素。

## 状态

- 分隔条为 1 DIP ControlBorder 线，命中区域 1.5 单位宽；悬停或拖动时变为 AccentText；键盘焦点时在命中区域外画焦点环。
- 折叠时该栏不构建，分隔条贴边并显示一段短标记，提示可从此处拉开。
- 无障碍：角色为 splitter，名称“分隔条”，描述为第一栏所占百分比或“已折叠”。

## 事件

- `Resized`：本帧尺寸改变（拖动或键盘）。
- `Toggled`：本帧折叠或展开。

## 键盘操作

Tab 聚焦分隔条后：

- ←/→（纵向为 ↑/↓）：移动 4 单位，到达最小尺寸后停止。
- 可折叠栏已在最小尺寸时继续朝它按方向键：折叠；折叠后按反方向键：展开。
- Home / End：移到最小 / 最大位置。
- Enter：切换可折叠栏（未设置 `Collapsible` 时不占用 Enter）。

指针：拖动分隔条；拖到可折叠栏最小尺寸的一半以内即折叠，拖回超过一半即展开；双击切换折叠。

## 限制

- 未直接使用 MyGo `ui.Split`：它的最小栏宽固定为 10 个间距单位（`ui/split.go:20`），没有折叠能力，且分隔条元素不对外暴露，无法挂接 Enter、Home、End。本组件用 MyGo 公开的 `Focusable`、`Dragged`、`Shortcut`、`RoleSplitter` 重新组合，行为与 MyGo 分隔条一致并补齐上述功能。
- MyGo 没有公开的无障碍数值范围接口，当前位置通过描述文字提供。
- 键盘或拖动产生的尺寸在下一帧生效（组件会立即请求下一帧）。
