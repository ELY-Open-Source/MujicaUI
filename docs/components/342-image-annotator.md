# 342 ImageAnnotator 图片标注

## 用途

在图片上画矩形、箭头和文字标注，可选择、移动和删除，并导出标注数据。

## 最小示例

```go
photo := ui.NewBitmap(examplePhoto)
var notes []mujica.Annotation
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ImageAnnotator(c, photo, &notes, mujica.ImageAnnotatorOptions{}).Changed() {
		data, _ := mujica.ExportAnnotations(notes)
		_ = data // save the JSON
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ImageAnnotator(c *ui.Context, src ui.ImageSource, notes *[]Annotation, opts ImageAnnotatorOptions) ChoiceResult`

`func ExportAnnotations(notes []Annotation) ([]byte, error)` 把标注写成 JSON（kind、x0、y0、x1、y1、text）。

`Annotation` 的坐标是图片像素；文字标注以 `From` 为左上角，`Text` 不能为空。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Height` | `float32` | 0（即 320） | 舞台高度。 |
| `Disabled` | `bool` | false | 禁用绘制、选择和删除。 |
| `Label` | `string` | "Annotations" | 舞台的无障碍名称。 |

## 状态

- 工具：选择、矩形、箭头、文字，用分段控件切换。
- 选中：虚线框加菱形标记，不只靠颜色。
- 标注用暗红墨线加浅色描边，在任何图片上都看得清。

## 事件

`Changed()` 在新增、移动、改字或删除标注后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| [ / ] | 选中上一个 / 下一个标注 |
| 方向键 | 移动选中的标注，每次图片宽度的 1% |
| Delete / Backspace | 删除选中的标注 |
| Escape | 取消选择 |

## 限制

- 文字工具：先在工具栏输入文字，再点击放置；选中文字标注时同一输入框可改字。
- 标注不能缩放，只能移动；要改大小请删除后重画。
- 文字的命中区域按每字 8 像素估算。
