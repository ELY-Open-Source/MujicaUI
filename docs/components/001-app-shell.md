# 001 AppShell 应用外壳

## 用途

搭建桌面窗口的整体骨架：顶部标题栏、左侧侧栏、主内容区和底部状态栏四个插槽。侧栏可折叠为只含切换按钮的窄条；侧栏与主内容各自独立滚动。

## 最小示例

```go
var collapsed bool
mujica.AppShell(c, &collapsed, mujica.AppShellOptions{
	TitleBar: func() { mujica.TitleBar(c, "Court Records", mujica.TitleBarOptions{}) },
	Sidebar:  func() { ui.Text(c, "Chronicles") },
	Content:  func() { ui.Text(c, "The Grand Hall") },
})
```

签名：`func AppShell(c *ui.Context, collapsed *bool, opts AppShellOptions) AppShellResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `collapsed` | `*bool` | 必填 | 侧栏是否折叠，由应用持有；为 nil 时 panic。 |
| `TitleBar` | `func()` | nil | 顶部插槽，通常放 `TitleBar`。 |
| `Sidebar` | `func()` | nil | 侧栏内容，自带独立滚动；为 nil 时没有侧栏与折叠快捷键。 |
| `Content` | `func()` | 必填 | 主内容，自带独立滚动；为 nil 时 panic。 |
| `StatusBar` | `func()` | nil | 底部插槽，通常放 `StatusBar`。 |
| `SidebarWidth` | `float32` | 240 | 展开时侧栏宽度（DIP）；负数 panic。 |
| `ContentScroll` | `*ui.ScrollState` | nil | 保存主内容滚动位置；nil 时按位置保留。 |

返回 `AppShellResult`，内嵌 `*ui.Element`。

## 状态

- 展开：侧栏宽 `SidebarWidth`，顶部右侧为收起按钮。
- 折叠：侧栏收为 `28 + 4 单位` 的窄条，只保留展开按钮。宽度变化 120 ms，减少动画时直接切换。
- 切换按钮有悬停、按下和键盘焦点环；折叠不会重置主内容滚动位置。

## 事件

- `Toggled`：本帧侧栏被展开或折叠（按钮或快捷键），只触发一次。

## 键盘操作

- `Cmd+B`（macOS）/ `Ctrl+B`（其他平台）：切换侧栏。焦点控件先处理该键，富文本等若占用此键则不会触发。
- Tab 依次到达切换按钮、侧栏与内容中的控件；滚动区域内容溢出时可获得焦点，方向键、Page Up/Down、Home、End 滚动。

## 限制

- 窄窗口不会自动折叠，由应用根据窗口宽度设置 `*collapsed`。
- 折叠动画期间侧栏内容立即隐藏，只有宽度过渡。
- 快捷键固定为 Cmd/Ctrl+B，不可配置。
