# 048 FilePicker

文件选择。

## 用途

调用 MyGo 的系统文件或文件夹选择对话框（`mygo.Dialog.Open`），支持单选、多选与扩展名筛选，只返回所选路径。

## 最小示例

```go
var paths []string
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.FilePicker(c, &paths, mujica.FilePickerOptions{
		Multiple: true,
		Filters:  []mygo.FileFilter{{Name: "Scores", Extensions: []string{"pdf", "mid"}}},
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `paths` | `*[]string` | 必填 | 所选路径；取消时不变。 |
| `Directory` | `bool` | false | 选择文件夹。 |
| `Multiple` | `bool` | false | 允许多选。 |
| `Filters` | `[]mygo.FileFilter` | nil | 扩展名筛选，扩展名不带点；带点时 panic。 |
| `Title` | `string` | "" | 对话框标题。 |
| `Label` | `string` | 内置“Choose file…/选择文件…” | 按钮文字。 |
| `Parent` | `*mygo.Window` | nil | 附着的窗口（macOS 显示为表单）。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、对话框打开中（按钮显示加载并阻止重复打开）、已选择（列出文件名与所在目录）、对话框失败（危险色说明）、禁用。

## 事件

`FilePickerParts.Changed()` 选中了新路径；`Canceled()` 用户取消，不视为错误。

## 键盘操作

Enter/Space 打开系统对话框；对话框内的键盘操作由系统提供。

## 限制

- 对话框为系统样式，不换肤。
- 对话框在后台 goroutine 中等待，结果经 `Context.Invalidate` 回到界面线程；需要应用已通过 `mygo.App.Run` 启动，无窗口测试中无法打开真实对话框。
- 不读取文件内容。
