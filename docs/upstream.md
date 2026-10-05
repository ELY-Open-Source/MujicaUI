# 上游归并与上游缺口

## 归并（PRD 4.11）

MyGo 已有能力纳入下列组件，不增加计数。

| MyGo 能力 | MujicaUI 组件 | 文件 |
| --- | --- | --- |
| Primary button、Menu button | Button 主按钮变体；DropdownMenu | `button.go`、`dropdown_menu.go` |
| Stepper、Number input | NumberInput 仅步进模式与输入模式 | `number_input.go` |
| Search field、Find bar | SearchInput 搜索模式与查找栏模式 | `search_input.go` |
| Autocomplete、Combobox | Combobox 自由输入与选项选择 | `combobox.go` |
| Token field、Editable text | TagsInput；Text 的 `EditableText` | `tags_input.go`、`text.go` |
| List、Grid view | List 列表模式与网格模式 | `list.go` |
| Table、Outline、Tree | DataTable；TreeTable；Tree | `data_table.go`、`tree_table.go`、`tree.go` |
| Date input、Time input | DatePicker；TimePicker | `date_picker.go`、`time_picker.go` |
| Back and forward buttons | NavigationMenu 历史导航 | `navigation_menu.go` |
| Meter、Progress bar | Progress 计量模式与进度模式 | `progress.go` |

## 上游缺口（待提交 MyGo）

PRD 第 2 节要求公开接口不足时向上游提交扩展。以下缺口在 MyGo v0.2.10 中核实，组件文档的“限制”节写明影响。本仓库不访问 MyGo internal，也未修改依赖。

| 缺口 | 位置 | 受影响组件 |
| --- | --- | --- |
| 自定义元素无法设置 expanded、mixed、数值范围与值 | `ui/element.go:282-295`、`ui/base.go:128,297` | Checkbox 半选、MultiSelect、TreeSelect、RangeSlider、Rating、全部日期时间弹出按钮 |
| 无 menu、menuitem 角色 | `ui/access.go:22-79` | Menubar、DropdownMenu、ContextMenu |
| 无 heading 角色 | `ui/access.go` | Heading |
| 输入法组合状态不公开 | `ui/editor.go:73`、`ui/input.go:689` | TextInput、TextArea、SearchInput、Dialog、AlertDialog |
| 打开链接失败不回传 | `ui/context.go:270`、`ui/window.go:583` | Link |
| 浮层焦点恢复与浮层快捷键不公开 | `ui/scope.go:109,123,137` | Dialog、Drawer、Overlay |
| `PopoverBase` 背板吞掉外部点击；只有上下两个位置 | `ui/base.go:373`、`ui/widgets.go:689` | Popover、Popconfirm、DropdownMenu、HoverCard |
| `activeDescendant` 不公开 | `ui/element.go:331` | Menubar、CommandPalette |
| `SliderBase` 方向键步长固定 1%，构建时处理指针 | `ui/base.go:111-117` | Slider |
| `SelectBase` 高亮无设置接口 | `ui/base.go:260` | Select |
| `TextInputBase` 无只读开关与占位符 | `ui/editor.go:139,734,811` | TextInput 只读模式、CommandPalette |
| `ui.Stepper` 禁用时仍可步进，固定 9 位取整 | `ui/indicators.go:216`、`ui/widgets.go` `snap` | NumberInput（已改为自绘） |
| 路径遮罩缓存按四分之一像素取键 | `ui/path.go:189-210` | 全部自绘路径（以 `internal.Snap` 对齐网格） |
