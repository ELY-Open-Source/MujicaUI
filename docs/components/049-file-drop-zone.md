# 049 FileDropZone

文件拖放区。

## 用途

接收从系统拖入或通过“浏览…”选择的文件，检查类型、数量和大小，显示文件列表与错误，可逐项移除。不上传，不读取内容。

## 最小示例

```go
var files []mujica.DroppedFile
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FileDropZone(c, &files, mujica.FileDropZoneOptions{
		Extensions: []string{"png", "jpg"}, MaxFiles: 3, MaxSize: 5 << 20,
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `files` | `*[]DroppedFile` | 必填 | 已接受的文件（Path、Size）。 |
| `Extensions` | `[]string` | nil（不限） | 允许的扩展名，不带点，不区分大小写。 |
| `MaxFiles` | `int` | 0（不限） | 列表最多文件数。 |
| `MaxSize` | `int64` | 0（不限） | 单个文件的最大字节数。 |
| `Hint` | `string` | 内置“Drop files here/将文件拖到此处” | 提示文字。 |
| `Disabled` | `bool` | false | 禁用拖放、浏览与移除。 |

## 状态

默认、文件拖到上方（选中底色与强调色虚线边框）、检查中（列表显示“Checking…”）、已接受、被拒绝（危险色图标与具体原因）、禁用。

## 事件

`FileDropZoneParts.Changed()` 添加或移除了文件。被拒绝的条目只显示在组件内，可单独关闭。

## 键盘操作

Tab 到达“浏览…”与每个移除按钮；拖放需要指针，键盘用户通过“浏览…”添加文件。

## 限制

- 文件大小通过后台 goroutine 中的 `os.Stat` 获取，不阻塞界面线程；文件夹、无法读取的路径会被拒绝。
- 重复路径静默忽略。
- MyGo 的无窗口测试器不能模拟系统拖放，自动测试通过“浏览…”路径覆盖同一检查流程。
