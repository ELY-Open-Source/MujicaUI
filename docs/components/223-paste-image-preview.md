# 223 PasteImagePreview 图片预览

## 用途

粘贴或拖入图片的预览缩略图；移除。

## 最小示例

```go
pic := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="3"/>`))
images := []mujica.PastedImage{{ID: "p1", Source: pic, Name: "screenshot"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if id, ok := mujica.PasteImagePreview(c, images).Removed(); ok {
		fmt.Println("remove", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`images` 为调用方解码的缩略图列表。每项 `ID` 为唯一稳定标识，`Source` 为必填图片，`Name` 为操作名称；缩略图边长为 72 DIP。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Removed() 返回待移除的稳定 ID 与事件标志。

## 键盘操作

Tab 到移除按钮；Enter、Space 移除。

## 限制

粘贴和拖入的图片解码、获取及去重由调用方完成。MyGo ui/headless.go:377–380 仅提供文本剪贴板；ui/drop.go:11 返回文件路径，原生图片剪贴板接口仍需上游扩展。
