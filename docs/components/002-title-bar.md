# 002 TitleBar 标题栏

## 用途

窗口顶部的标题行：衬线体应用标题、前后工具插槽。空白处可拖动窗口，并为系统窗口按钮留出空间。配合 `mygo.TitleBarHidden` 使用。

## 最小示例

```go
layout.TitleBar(c, "Court Records", layout.TitleBarOptions{
	Trailing: func() { input.Button(c, "Share", input.ButtonOptions{Variant: input.Ghost}) },
})
```

签名：`func TitleBar(c *ui.Context, title string, opts TitleBarOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `title` | `string` | 必填 | 标题文字，过长时单行省略。 |
| `Leading` | `func()` | nil | 标题前的工具，如侧栏开关。 |
| `Trailing` | `func()` | nil | 标题后的工具。 |

## 状态

- 最小高度取 `c.TitleBar().Height` 与 `控件高度 + 2 单位` 中较大者（默认密度 44 DIP）；标题字号随系统文本缩放，放大后标题栏随之增高。
- 左右内边距在 `c.TitleBar().Left/Right` 之外再加 3 单位，不覆盖系统按钮；全屏时预留为 0。
- Surface 底色，底部 1 DIP Border 分隔线。

## 事件

TitleBar 本身无事件；插槽中控件照常报告 `Clicked` 等事件。

## 键盘操作

插槽控件照常接受 Tab、Enter、Space。标题栏本身不获取焦点。

## 限制

- 拖动、双击缩放由 MyGo `DragWindow` 实现，只作用于非交互区域；需真机验证，无窗口测试无法覆盖。
- 系统窗口按钮外观由系统决定，不换肤。
- 普通有标题栏的窗口中 `c.TitleBar()` 为零，组件只是一行普通标题。
