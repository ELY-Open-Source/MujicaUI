# 066 Menubar 菜单栏

## 用途

应用窗口内由 MujicaUI 自绘的菜单栏（不是系统菜单）：多级子菜单、勾选项、单选项、分隔线、快捷键提示和完整键盘导航。

## 最小示例

```go
res := mujica.Menubar(c, []mujica.MenubarMenu{
	{Label: "File", Items: []mujica.MenubarItem{
		{ID: "new", Label: "New", Shortcut: mujica.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyN}},
		{Label: "Export", Items: []mujica.MenubarItem{{ID: "pdf", Label: "PDF"}}},
	}},
	{Label: "View", Items: []mujica.MenubarItem{
		{ID: "wrap", Label: "Word wrap", Kind: mujica.MenubarCheck, Checked: wrap},
	}},
}, mujica.MenubarOptions{})
if id, ok := res.Chosen(); ok && id == "wrap" {
	wrap = !wrap
}
```

签名：`func Menubar(c *ui.Context, menus []MenubarMenu, opts MenubarOptions) MenubarResult`

## 参数

`menus` 为空或标题为空时 panic。

`MenubarItem`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `string` | 必填（子菜单与分隔线除外） | `Chosen` 报告的标识。 |
| `Label` | `string` | 必填（分隔线除外） | 名称。 |
| `Kind` | `MenubarItemKind` | `MenubarAction` | `MenubarAction`、`MenubarCheck`、`MenubarRadio`、`MenubarSeparator`。子菜单必须为 `MenubarAction`。 |
| `Checked` | `bool` | false | 勾选项与单选项的当前值，由调用方持有。 |
| `Disabled` | `bool` | false | 禁用：不可选择，方向键跳过。 |
| `Shortcut` | `MenubarShortcut` | 无 | 仅显示的快捷键提示：macOS 显示 ⌘⇧ 等符号，其他平台显示 Ctrl+Shift+。键名无法表示时 panic。 |
| `Items` | `[]MenubarItem` | nil | 非 nil 时为子菜单。 |

`MenubarOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 内置“菜单栏” | 菜单栏的无障碍名称。 |

返回 `MenubarResult`，`Element` 字段是菜单栏元素。

## 状态

- 打开的菜单标题：Selection 背景；高亮项：Selection 背景。
- 勾选项：AccentText 勾号；单选项：AccentText 菱形标记；名称附“已勾选”。
- 禁用项：降低不透明度，指针与键盘均跳过。
- 弹层：6 DIP 圆角、Border 边框；靠近窗口边缘时左移，子菜单放不下时翻到父菜单左侧，过高时限高并在菜单内滚动，键盘高亮项自动滚入视野。

## 事件

- `Chosen() (string, bool)`：本帧被选择的项 ID。选择后菜单关闭；勾选与单选由调用方更新 `Checked`。

## 键盘操作

- Tab：进入菜单栏（一个停靠点）；←/→ 在标题间移动；↓、↑、Enter 或 Space 打开菜单。
- 菜单中：↑/↓ 移动（跳过分隔线与禁用项并循环），Home / End 到首末项。
- →：打开子菜单；在普通项上则打开右侧下一个菜单。←：关闭子菜单；在顶层则打开左侧菜单。
- Enter / Space：选择或打开子菜单。Escape：关闭一级。
- 指针：点击标题打开；菜单打开时移到其他标题即切换；悬停子菜单项展开；点击外部关闭。

## 限制

- 快捷键只做提示，不注册；用 `c.Shortcut` 处理。
- MyGo 没有 menu / menuitem 角色（`ui/access.go:25-79`），项以按钮呈现；勾选状态无法通过公开接口写入（`ui/base.go:56` 的 checked 为内部字段），因此写入名称。
- 指针点击非可聚焦元素会清除焦点（`ui/input.go:238`）；选择后组件把焦点还给对应标题。
- 菜单打开期间，焦点保持在标题上，高亮项由组件管理。驱动菜单的按键注册在标题上，而 MyGo 的 active descendant 不公开，因此高亮变化时用 `c.Announce` 朗读当前项，并把它写入标题的 `Description`。
- 弹层打开时在 160 ms 内淡入，减少动画时直接显示。关闭时立即消失：弹层关闭后 MyGo 不再构建其元素，公开接口无法为已移除的浮层做淡出。
- 点击菜单内任意位置（包括禁用项）或用指针打开子菜单后，焦点回到所属标题，方向键、Enter、Escape 继续有效。
- 菜单打开时用 Tab 或 Shift+Tab 离开菜单栏，菜单关闭，焦点留在 Tab 目标上；焦点在菜单栏内任意位置时 Escape 关闭菜单。
