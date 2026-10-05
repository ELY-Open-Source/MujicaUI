# 059 Tabs 标签页

## 用途

在同一区域内切换若干面板。标签列表与当前标签的面板一起构建；支持横向与纵向布局、禁用项和溢出滚动。Tabs 是视觉回归基准组件，快照覆盖三种密度。

## 最小示例

```go
if mujica.Tabs(c, &tab, []mujica.Tab{
	{Label: "General", Panel: func() { ui.Text(c, "General settings") }},
	{Label: "Fonts", Icon: icons.Must("italic"), Panel: func() { ui.Text(c, "Fonts") }},
	{Label: "Archive", Disabled: true},
}, mujica.TabsOptions{Label: "Settings"}).Changed() {
	fmt.Println("tab", tab)
}
```

签名：`func Tabs(c *ui.Context, selected *int, tabs []Tab, opts TabsOptions) ChoiceResult`

## 参数

`selected`：当前标签下标（从 0 开始）。越界或指向禁用标签时 panic。`tabs` 为空时 panic。

`Tab`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 必填 | 标签文字，也是面板的无障碍名称；为空时 panic。 |
| `Icon` | `*ui.SVG` | nil | 标签文字前的 16 DIP 图标。 |
| `Disabled` | `bool` | false | 禁用：不可点击、不进入方向键序列。 |
| `Panel` | `func()` | nil | 选中时构建的面板内容；nil 显示空面板。 |

`TabsOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Vertical` | `bool` | false | 标签竖排在左侧，面板在右侧。 |
| `Label` | `string` | "" | 标签列表的无障碍名称。 |

返回 `ChoiceResult`，`Element` 字段是包含标签与面板的元素。

## 状态

- 选中：Selection 背景、2 DIP AccentText 边线（横向在底边，纵向在右边）、600 字重，三者同时标记。
- 悬停：SurfaceHover 背景；按下：Selection 背景。背景与文字颜色在 120 ms 内过渡，减少动画时直接切换。
- 焦点：2 DIP Focus 焦点环，外留 2 DIP。
- 禁用：MyGo 统一降低不透明度，不响应指针与键盘。
- 溢出：标签条横向（纵向时竖向）滚动；键盘或点击切换后，选中标签滚入视野。

## 事件

- `Changed()`：本帧用户选择了新标签（点击或方向键）。

## 键盘操作

- Tab：进入标签列表，落在当前标签；再按 Tab 离开。
- ←/→、↑/↓：选择上一个/下一个可用标签，跳过禁用项，首尾循环；焦点随选择移动。
- Home / End：第一个 / 最后一个可用标签。

## 限制

- MyGo `TabsBase` 的方向键不知道禁用项（`ui/base.go:192` 直接选择 i±1）。本组件只把可用标签交给 `TabsBase`，禁用标签以普通 `RoleTab` 行显示，因此禁用标签不在 MyGo 的标签计数内。
- MyGo 没有 tabpanel 角色与 aria-controls 式关联；面板以 `RoleGroup` 呈现，名称取自所选标签。
- 不提供标签关闭、拖动排序。
