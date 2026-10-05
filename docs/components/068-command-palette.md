# 068 CommandPalette 命令面板

## 用途

用快捷键（默认 Cmd+K，非 macOS 为 Ctrl+K）打开的模态命令搜索框：输入过滤、分组显示、无结果提示、方向键选择，执行后关闭。

## 最小示例

```go
res := mujica.CommandPalette(c, &open, []mujica.Command{
	{ID: "new", Label: "New document", Group: "File"},
	{ID: "theme", Label: "Toggle dark mode", Group: "View", Keywords: []string{"night"}},
}, mujica.CommandPaletteOptions{})
if id, ok := res.Chosen(); ok {
	run(id)
}
```

签名：`func CommandPalette(c *ui.Context, open *bool, commands []Command, opts CommandPaletteOptions) CommandPaletteResult`

## 参数

`open`：是否显示；快捷键把它设为 true，执行、Escape 或点击遮罩设为 false。调用方也可直接设置。

`Command`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `string` | 必填 | `Chosen` 报告的标识；为空或重复时 panic。 |
| `Label` | `string` | 必填 | 显示名称；为空时 panic。 |
| `Group` | `string` | "" | 分组标题；分组按首次出现的顺序排列，组内保持原顺序。 |
| `Keywords` | `[]string` | nil | 额外的搜索词。 |
| `Icon` | `*ui.SVG` | nil | 名称前图标。 |

`CommandPaletteOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OpenMods` | `ui.Modifiers` | 0 | 打开快捷键的修饰键。 |
| `OpenKey` | `ui.Key` | 0（即 `ui.Cmd`+K） | 打开快捷键；为 0 时忽略 `OpenMods`。 |

返回 `CommandPaletteResult`；关闭时其 `Element` 为 nil。

## 状态

- 打开：遮罩变暗，面板在窗口上部，焦点在搜索框；每次打开都清空搜索。
- 搜索：不区分大小写，按空格分词，每个词须出现在名称、分组或关键词中。
- 高亮项：Selection 背景、左侧 AccentText 边线、600 字重；键盘移动时滚入视野。
- 无结果：显示“没有匹配的命令”，Enter 不执行任何命令。
- 底部提示栏说明 ↑↓、Enter、Esc。

## 事件

- `Chosen() (string, bool)`：本帧执行的命令 ID（Enter 或点击），同时关闭面板。

## 键盘操作

- Cmd+K / Ctrl+K：打开（窗口级快捷键）。
- 输入：过滤。↑/↓：移动高亮并循环。Enter：执行并关闭。Escape：关闭不执行，焦点回到打开前的元素。
- Tab 限制在面板内（MyGo 模态对话框行为）。

## 限制

- MyGo `TextInputBase` 没有占位文字（`ui/editor.go:734`）；为避免与输入法组合文字重叠，说明放在底部提示栏，搜索框以无障碍名称“搜索命令”标识。
- 搜索是子串匹配，不做模糊匹配或排序打分。
- 焦点留在搜索框以便继续输入；方向键移动高亮时用 `c.Announce` 朗读当前命令，并写入搜索框的 `Description`（MyGo 的 active descendant 不公开）。
- 弹层打开时在 160 ms 内淡入，减少动画时直接显示。关闭时立即消失：弹层关闭后 MyGo 不再构建其元素，公开接口无法为已移除的浮层做淡出。
- 搜索内容变化后高亮回到第一项，并滚入视野。
