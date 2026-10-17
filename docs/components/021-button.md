# 021 Button 按钮

## 用途

触发一个操作。提供主按钮、次按钮、描边、无底色和危险五种变体，可带图标，可显示加载与禁用。基于 MyGo `ui.ButtonBase`。

## 最小示例

```go
if input.Button(c, "Save", input.ButtonOptions{Icon: icons.Must("check")}).Clicked() {
	fmt.Println("saved")
}
```

签名：`func Button(c *ui.Context, label string, opts ButtonOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `label` | `string` | 必填（仅图标时可为空） | 按钮文字；为空时必须提供 `Label`，否则 panic。 |
| `Variant` | `ButtonVariant` | `Primary` | `Primary`、`Secondary`、`Outline`、`Ghost`、`Danger`。 |
| `Icon` | `*ui.SVG` | `nil` | 文字前的 16 DIP 图标。 |
| `Loading` | `bool` | `false` | 加载标记出现在前置 16 DIP 图标位置，文字保持可见；有图标时替换图标，纯文字时从左右内边距中让出该位置，宽度不变；阻止点击。 |
| `Disabled` | `bool` | `false` | 不可点击、不进入 Tab 序列。 |
| `Label` | `string` | `""` | 无障碍名称；`label` 为空时必填。 |

## 状态

默认、悬停（底色加深或出现浅底）、按下、键盘聚焦（2 DIP 酒红焦点环，外留 2 DIP）、加载、禁用（MyGo 统一减淡）。高度不低于当前密度的控件高度（28/36/44 DIP）。

## 事件

`Clicked()`：鼠标在按钮内释放、Enter 或 Space 时为真，一次操作只报告一次。加载与禁用时始终为假。

## 键盘操作

Tab / Shift+Tab 聚焦；Enter 或 Space 触发。

## 限制

- 加载标记在减少动画模式下静止显示。
- 危险变体的悬停与按下色由 `DangerFill` 混合黑色得到，未单独定义变量。
