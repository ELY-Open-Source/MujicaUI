# 328 Minimap · 代码缩略图

## 用途

以细条绘制代码缩略图，框出编辑器可见区域；点击或拖动跳转，键盘可移动。

## 最小示例

```go
lines := []string{"package court", "", "func Seal() {}", "func Herald() {}"}
view := mujica.MinimapView{Top: 0, Visible: 2}
draw := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Minimap(c, lines, &view, mujica.MinimapOptions{Language: "go"}).Changed() {
		fmt.Println("scroll editor to", view.Top)
	}
}
ui.Render(draw, 400, 300, 1)
```

## 参数

`func Minimap(c *ui.Context, lines []string, view *MinimapView, opts MinimapOptions) ChoiceResult`。`MinimapView{Top, Visible}` 是编辑器当前首行与可见行数，由调用方与编辑器同步（可用 `CodeViewerState.Visible` 与 `ScrollTo`）。每行 2 DIP、每字符 1 DIP；文件比缩略图高时缩略图随视图滚动，只绘制可见的行。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Language` | `string` | "" | 着色语言，同 CodeViewer。 |
| `Width` | `float32` | 96 | 宽度（DIP）；负数 panic。 |
| `Disabled` | `bool` | false | 不能跳转。 |

## 状态

可见区为 Selection 半透明填充加 AccentText 边框；聚焦时外侧焦点环；禁用。

## 事件

`Changed()`：用户移动了视图，读取 `view.Top`。

## 键盘操作

- Tab 聚焦；↑/↓ 移动一行，PageUp/PageDown 移动一屏，Home/End 到开头与结尾。

## 限制

- 缩略图按行独立着色，不跟踪跨行的块注释。
