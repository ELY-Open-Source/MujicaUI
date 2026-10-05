# 061 Breadcrumb 面包屑

## 用途

显示当前位置的层级路径。最后一项为当前页（非链接），前面各项可点击返回；层级过多时中间项折叠到“…”菜单；过长名称以省略号截断，悬停显示完整名称。

## 最小示例

```go
if i, ok := mujica.Breadcrumb(c, path, mujica.BreadcrumbOptions{}).Chosen(); ok {
	path = path[:i+1]
}
```

签名：`func Breadcrumb(c *ui.Context, items []BreadcrumbItem, opts BreadcrumbOptions) BreadcrumbResult`

## 参数

`items` 为空时 panic。

`BreadcrumbItem`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 必填 | 名称；为空时 panic。 |
| `Icon` | `*ui.SVG` | nil | 名称前的 14 DIP 图标。 |

`BreadcrumbOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `MaxItems` | `int` | 0（即 4） | 最多显示的项数；超出时保留第一项和最后 `MaxItems-2` 项，其余折叠。小于 3 时 panic。 |
| `MaxItemWidth` | `float32` | 0（即 160 DIP） | 单项文字最大宽度，超出以省略号截断。 |
| `Label` | `string` | 内置“路径导航” | 整组的无障碍名称。 |

返回 `BreadcrumbResult`，内嵌 `*ui.Element`。

## 状态

- 链接：TextMuted 文字，悬停 SurfaceHover 背景与 Text 文字，焦点环 2 DIP。
- 当前项：Text 文字、600 字重，无障碍名称附“当前页”。
- 分隔：旧金（Ornament）小箭头，不进入无障碍树。
- 折叠按钮“…”：打开弹层列出被折叠的项，弹层打开时焦点在第一项。

## 事件

- `Chosen() (int, bool)`：本帧被激活的项下标（含折叠菜单中的项）。组件不修改路径，由调用方处理。

## 键盘操作

- Tab / Shift+Tab：在各链接与折叠按钮之间移动；Enter 或 Space 激活。
- 折叠菜单：Tab 在项之间移动，Enter 选择，Escape 关闭并把焦点还给“…”按钮。

## 限制

- 完整名称提示对每一项都提供，MyGo 不报告文字是否被截断。
- 未复用 MyGo `ui.Breadcrumbs`：它不支持折叠与提示，本组件用 `ButtonBase`、`PopoverBase` 组合实现。
