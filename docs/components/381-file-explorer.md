# FileExplorer

## 用途

组合文件夹树、文件列表、导航历史与多选。

## 最小示例

```go
var state mujica.FileExplorerState
opts := mujica.FileExplorerOptions{
	Roots:   []string{"/court"},
	Folders: func(string) []string { return nil },
	Entries: func(string) []mujica.FileEntry {
		return []mujica.FileEntry{{Path: "/court/seal.go", Name: "seal.go", Size: 1024}}
	},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileExplorer(c, &state, opts).Element.Height(240)
}
ui.Render(view, 600, 300, 1)
```

## 参数

Roots、Folders、Entries 必填；Disabled 默认 false；Revision 用于原地更新条目。

## 状态

默认、焦点、禁用；空数据按组件内容展示。

## 事件

Navigated、Changed、Opened 返回本帧导航、选择和打开。

## 键盘操作

Tab 切换；方向键选择；Enter 打开；macOS Command+[ / ] / ↑，其他平台 Alt+← / → / ↑。

## 限制

文件及目录由调用方提供，读取与打开由调用方执行。条目 Path 必须唯一且非空；切片地址、长度或 Revision 变化时重建索引并清理已删除的选择。
