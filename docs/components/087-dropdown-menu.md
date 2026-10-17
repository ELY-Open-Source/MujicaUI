# 087 DropdownMenu 下拉菜单

## 用途

按钮触发的应用内菜单，由 MujicaUI 自绘（不是系统菜单）。支持分组、子菜单、勾选项与禁用项；选择后按规则关闭。

## 最小示例

```go
pinned := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	overlay.DropdownMenu(c, "Actions", overlay.DropdownMenuOptions{}, func(m *overlay.PopupMenu) {
		m.Group("Letter", func() {
			if m.Item("Rename", overlay.PopupMenuItemOptions{}) {
				fmt.Println("rename")
			}
			m.Check("Pinned", &pinned, overlay.PopupMenuItemOptions{KeepOpen: true})
		})
		m.Submenu("Move to", overlay.PopupMenuItemOptions{}, func(m *overlay.PopupMenu) {
			m.Item("Archive", overlay.PopupMenuItemOptions{})
		})
		m.Separator()
		m.Item("Burn", overlay.PopupMenuItemOptions{Disabled: true})
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func DropdownMenu(c *ui.Context, label string, opts DropdownMenuOptions, build func(m *PopupMenu)) ui.Element`，返回触发按钮。`build` 为 nil、或 `label` 与 `opts.Label` 都为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Icon` | `*ui.SVG` | `nil` | 标签前的图标。 |
| `Disabled` | `bool` | `false` | 禁用触发按钮并关闭已打开的菜单。 |
| `Label` | `string` | `""` | 仅图标按钮的无障碍名称。 |

`PopupMenu` 方法：

| 方法 | 说明 |
| --- | --- |
| `Item(label, opts) bool` | 普通项；本帧被选择时返回 true。 |
| `Check(label, *bool, opts) bool` | 勾选项，切换 `*bool`；本帧变化时返回 true。 |
| `Submenu(label, opts, build)` | 子菜单项；`build` 为 nil 时 panic。 |
| `Group(label, fn)` | 带说明文字的分组，前面有内容时自动加分隔线。 |
| `Separator()` | 分隔线。 |

`PopupMenuItemOptions`：`Disabled bool`（不可选、不进入方向键序列），`KeepOpen bool`（选择后不关闭菜单）。

## 状态

- 触发按钮：次按钮外观加下拉箭头，打开时保持按下底色。
- 菜单项：悬停或键盘焦点时显示 Selection 底色和左侧 2 DIP 强调条；勾选项显示勾号；禁用项变暗。
- 菜单面板最高为可用空间，超出时滚动；子菜单在右侧打开，右侧放不下时翻到左侧，仍放不下时贴窗口边缘，宽度不超过窗口；下方放不下时上移。

## 事件

`Item` / `Check` 的返回值；每次选择只报告一次。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Enter / Space / ↓ | 在触发按钮上打开菜单，焦点在第一项。 |
| ↑ / ↓ / Home / End | 在当前层移动，跳过禁用项，首尾循环。 |
| → / Enter | 在子菜单项上打开子菜单并进入。 |
| ← / Escape | 在子菜单中关闭子菜单，回到父项。 |
| Escape / Tab | 关闭菜单，焦点回到触发按钮。 |

## 限制

- MyGo 没有 menu / menuitem 无障碍角色（`ui/access.go:22` Role 列表），菜单面板报告为 popup，普通项为 button，勾选项为 checkbox，触发按钮为 menu button。
- 打开期间外部点击只用于关闭，不传给下方控件：`PopoverBase` 用覆盖整窗的背板接收点击（`ui/base.go:373`）。
- Escape 只关闭最上层的 MujicaUI 浮层（窗口内按打开顺序记录）。MyGo 自带的 `ui.Popover`、`ui.Select` 弹层不在这份记录中；它们打开在非模态 MujicaUI 浮层之上时，Escape 可能先关闭下面的 MujicaUI 浮层。
