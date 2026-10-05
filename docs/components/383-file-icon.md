# FileIcon

## 用途

按扩展名和文件夹开合显示类型图标。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.FileIcon(c, "seal.go", files.FileIconOptions{Size: 32})
}
ui.Render(view, 100, 80, 1)
```

## 参数

Size 默认 20 DIP；Folder 表示目录；Open 需要 Folder。

## 状态

文件、目录关闭、目录打开。

## 事件

展示组件。

## 键盘操作

展示组件，沿用父级焦点。

## 限制

扩展名大小写无关；未知扩展名显示通用文件图标。
