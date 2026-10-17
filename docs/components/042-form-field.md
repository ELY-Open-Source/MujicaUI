# 042 FormField

表单字段。

## 用途

为单个控件提供标签、必填标记、说明与错误文字。标签、说明与错误显式关联到 control 返回的元素：屏幕阅读器以标签命名该元素并读出说明、必填与错误；点击标签聚焦该元素。

## 最小示例

```go
email := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.FormField(c, "Email", input.FormFieldOptions{Required: true, Description: "For receipts."}, func() ui.Element {
		return input.InputGroup(c, &email, input.InputGroupOptions{}).Input
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `label` | `string` | 必填 | 标签文字；空值 panic。 |
| `FormFieldOptions.Required` | `bool` | false | 在标签后显示危险色“*”，并向辅助技术说明“必填”。 |
| `FormFieldOptions.Description` | `string` | "" | 控件下方的说明。 |
| `FormFieldOptions.Error` | `string` | "" | 控件下方的错误，并把控件标记为无效。 |
| `FormFieldOptions.Horizontal` | `bool` | false | 标签在左、控件在右。 |
| `control` | `func() ui.Element` | 必填 | 构建控件并返回要关联的可聚焦元素（如 `InputGroupParts.Input`）；返回零值 `ui.Element{}` 时 panic。 |

## 状态

默认、必填、带说明、错误；禁用由控件自身决定，标签随之变淡。

## 事件

无自身事件；事件来自内部控件。

## 键盘操作

Tab 直接到达控件；点击标签聚焦控件，或切换复选框、开关。

## 限制

- control 返回的元素的 Label 会被字段标签覆盖。
- Error 只把元素标记为无效；InputGroup 等的红色边框需同时传入 `Invalid: true`。
- 左右布局时多个字段需放在 MujicaUI `Form(Horizontal)` 中才能对齐标签宽度；单独使用时只按自身标签宽度排列。
