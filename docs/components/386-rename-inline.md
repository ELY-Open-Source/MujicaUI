# RenameInline

## 用途

编辑文件名，校验后确认或取消。

## 最小示例

```go
name, editing := "seal.txt", true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RenameInline(c, &name, &editing, mujica.RenameInlineOptions{})
}
ui.Render(view, 400, 100, 1)
```

## 参数

name、editing 必填；Folder 默认 false；Validate 返回错误文案；Disabled 默认 false。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Renamed 表示已写回名称；Canceled 表示取消。

## 键盘操作

Enter 确认；Escape 取消；双击开始编辑；失焦确认有效名称。

## 限制

默认文件名主体选择采用视觉标记和首次输入替换。MyGo ui/editor.go:133 的选择接口为私有；剪贴板和撤销沿用原生编辑器的实际选区。组合输入期间 Enter 由 MyGo ui/editor.go:485、688 处理，组件无法查询组合状态；使用输入法提交文本后再确认。
