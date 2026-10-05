# 042 FormField

表单字段。

## 用途

为单个控件提供标签、必填标记、说明与错误文字。标签通过 MyGo `ui.Field` 关联控件：屏幕阅读器以标签命名控件，并读出说明与错误；点击标签聚焦控件。

## 最小示例

```go
email := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FormField(c, "Email", mujica.FormFieldOptions{Required: true, Description: "For receipts."}, func() {
		mujica.InputGroup(c, &email, mujica.InputGroupOptions{})
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `label` | `string` | 必填 | 标签文字；空值 panic。 |
| `FormFieldOptions.Required` | `bool` | false | 在标签后加“ *”。 |
| `FormFieldOptions.Description` | `string` | "" | 控件下方的说明。 |
| `FormFieldOptions.Error` | `string` | "" | 控件下方的错误，并把控件标记为无效。 |
| `FormFieldOptions.Horizontal` | `bool` | false | 标签在左、控件在右。 |
| `control` | `func()` | 必填 | 构建控件；其中第一个可获取焦点的元素即被关联的控件。 |

## 状态

默认、必填、带说明、错误；禁用由控件自身决定，标签随之变淡。

## 事件

无自身事件；事件来自内部控件。

## 键盘操作

Tab 直接到达控件；点击标签聚焦控件，或切换复选框、开关。

## 限制

- 标签、说明与错误文字由 MyGo `ui.Field` 绘制，使用主题颜色（文字、TextMuted、Danger）；字号与间距跟随 MyGo。
- 错误只把边框染红的效果对 MyGo 自带输入框生效；MujicaUI 的 InputGroup 等需同时传入 `Invalid: true`。
- 左右布局时多个字段需放在同一个 `ui.Form` 或 MujicaUI `Form(Horizontal)` 中才能对齐标签宽度；单独使用时只按自身标签宽度排列。
