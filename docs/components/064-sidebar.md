# 064 Sidebar 侧边栏

## 用途

应用的源列表导航：分组、嵌套、图标、计数和选中项，分组与嵌套组可展开折叠；紧凑模式只显示图标。基于 MyGo `ui.Sidebar`，保留其键盘、首字母跳转与树形无障碍结构。

## 最小示例

```go
mujica.Sidebar(c, &selected, []mujica.SidebarSection{
	{Title: "Mail", Items: []mujica.SidebarItem{
		{ID: "inbox", Label: "Inbox", Icon: icons.Must("inbox"), Count: 3},
		{Label: "Folders", Items: []mujica.SidebarItem{{ID: "work", Label: "Work", Icon: icons.Must("folder")}}},
	}},
}, mujica.SidebarOptions{})
```

签名：`func Sidebar(c *ui.Context, selected *string, sections []SidebarSection, opts SidebarOptions) SidebarResult`

## 参数

`selected`：选中项的 ID。

`SidebarSection`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 分组标题，点击展开折叠；为空时 panic（紧凑模式不显示）。 |
| `Items` | `[]SidebarItem` | nil | 分组内的项。 |

`SidebarItem`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `string` | 叶子项必填 | 选择时写入 `*selected`；叶子项为空时 panic。 |
| `Label` | `string` | 必填 | 名称；为空时 panic。 |
| `Icon` | `*ui.SVG` | nil | 名称前图标；紧凑模式下必填。 |
| `Count` | `int` | 0 | 大于 0 时在右侧显示计数，超过 999 显示 “999+”。 |
| `Items` | `[]SidebarItem` | nil | 非 nil 时该项成为可展开折叠的嵌套组，子项缩进。 |

`SidebarOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Compact` | `bool` | false | 图标栏：隐藏文字与组标题，组之间以旧金细线分隔，名称作为悬停提示。 |
| `Width` | `float32` | 0（240，紧凑时 56） | 宽度，DIP。高度填满父容器。 |

返回 `SidebarResult`，内嵌 `*ui.Element`。展开折叠状态由组件保存。

## 状态

- 选中：Selection 背景、左侧 3 DIP AccentText 边线、600 字重，三者同时标记。
- 悬停：SurfaceHover 背景。
- 焦点：侧边栏获得焦点时，在选中项上画焦点环。
- 计数：等宽数字，名称附计数（如 “Inbox (12)”）。

## 事件

- `Changed()`：本帧用户选择了新项。

## 键盘操作

- Tab：进入侧边栏（一个停靠点）。
- ↑/↓：上一项/下一项；Home / End：首项/末项；输入名称开头字母跳到匹配项。
- 分组与嵌套组标题用指针点击展开折叠。

## 限制

- 复用 MyGo `ui.SidebarItem`，它用主题 Accent 绘制图标（`ui/sidebar.go:182`）；深色主题 Accent 对比度不足，本组件在构建侧边栏时临时把 MyGo 主题的 Accent 换成 AccentText。
- MyGo `SidebarItem` 总是追加名称文字（`ui/sidebar.go:212`）；紧凑模式传入空名称，名称改由 `Label` 与提示提供，因此紧凑模式下首字母跳转不可用。
- 嵌套组用 MyGo `SidebarSection` 实现，组标题本身不能被选中，且只能用指针展开折叠（MyGo 不给标题注册按键）。
