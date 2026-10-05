# 470 ProfileEditor

个人资料编辑。

## 用途

头像、名称与简介；与已保存版本不同时显示未保存标记；保存、放弃与更换头像由调用方处理。

## 最小示例

```go
saved := account.Profile{Name: "Sakiko", Bio: "Keyboard."}
profile := saved
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if account.ProfileEditor(c, &profile, account.ProfileEditorOptions{Saved: saved}).Save() {
		saved = profile
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `p` | `*Profile` | 必填 | Name、Bio。 |
| `Saved` | `Profile` | 零值 | 已保存的资料。 |
| `Avatar` | `*ui.Bitmap` | nil | 头像；nil 显示首字母。 |
| `BioLimit` | `int` | 160 | 可保存的简介字数上限；超限保留输入并阻止保存。 |
| `Saving` | `bool` | false | 保存中，阻止编辑。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

已保存、未保存（菱形与文字）、名称为空、简介超限、保存中、禁用。

## 事件

`Save()`、`ChangeAvatar()`、`Discarded()`。

## 键盘操作

Tab 依次经过更换头像、名称、简介、放弃、保存。

## 限制

- 头像选择与上传由调用方完成，可配合 FilePicker。
