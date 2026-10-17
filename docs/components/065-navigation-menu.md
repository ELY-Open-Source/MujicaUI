# 065 NavigationMenu 导航菜单

## 用途

站点式的顶层导航：直接链接与下拉分组（分组内可再分小节），用应用已有的 `*ui.Router` 标记当前路径并跳转；可选显示 MyGo 的后退、前进按钮。

## 最小示例

```go
navigation.NavigationMenu(c, router, []navigation.NavigationLink{
	{Label: "Home", Path: "/"},
	{Label: "Library", Items: []navigation.NavigationLink{
		{Label: "Poems", Path: "/library/poems", Description: "Verse of the court"},
	}},
}, navigation.NavigationMenuOptions{History: true})
router.View(c, func(r *ui.Route) { ui.Text(c, r.Path()) })
```

签名：`func NavigationMenu(c *ui.Context, router *ui.Router, links []NavigationLink, opts NavigationMenuOptions) ui.Element`

## 参数

`router` 为 nil 或 `links` 为空时 panic。路由状态由应用持有，组件只读 `Path()` 并调用 `Push`。

`NavigationLink`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 必填 | 名称；为空时 panic。 |
| `Path` | `string` | 叶子必填 | 点击后 `router.Push(Path)`；既无 `Path` 也无 `Items` 时 panic。 |
| `Icon` | `*ui.SVG` | nil | 名称前图标。 |
| `Description` | `string` | "" | 下拉面板中名称下方的说明。 |
| `Items` | `[]NavigationLink` | nil | 非 nil 时为分组：顶层显示为下拉按钮，面板内显示为小节标题。 |

`NavigationMenuOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `History` | `bool` | false | 在前面显示后退、前进按钮（MyGo `BackButton` / `ForwardButton`）。 |
| `Label` | `string` | 内置“导航菜单” | 链接组的无障碍名称。 |

## 状态

- 当前路径：路径等于链接，或位于其下（`/library` 覆盖 `/library/poems`；`/` 只匹配自身）。顶层项以 600 字重与底部 AccentText 线标记；面板中的当前项为 Selection 背景、左侧边线与 600 字重。
- 悬停或展开：SurfaceHover 背景。焦点环 2 DIP。
- 历史按钮在历史两端禁用；右键列出可前往的页面（MyGo 行为）。

## 事件

- 无返回事件：跳转直接作用于 `router`，调用方读取 `router.Path()`。

## 键盘操作

- Tab：进入顶层链接组；←/→ 在顶层项之间移动。
- 在分组上按 ↓、Enter 或 Space 打开面板，焦点在第一个链接；↑/↓ 在面板内移动，Enter 跳转并关闭，Escape 关闭并把焦点还给分组按钮。
- 历史：Cmd+[ / Cmd+]（macOS），Alt+←/→（其他平台），以及鼠标侧键，由 MyGo 路由处理。

## 限制

- MyGo 历史按钮的提示文字固定为英文 “Back” / “Forward”（`ui/router.go:713`）；本组件覆盖了无障碍名称，但提示仍为英文。
- 历史按钮的外观来自 MyGo 主题按钮样式（经 `core.Use` 映射为 MujicaUI 颜色）。
- 弹层打开时在 160 ms 内淡入，减少动画时直接显示。关闭时立即消失：弹层关闭后 MyGo 不再构建其元素，公开接口无法为已移除的浮层做淡出。
- 下拉面板最高为锚点上方或下方可用空间，超出时在面板内滚动；方向键移动焦点时自动滚入视野。
