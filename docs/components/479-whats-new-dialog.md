# 479 WhatsNewDialog

更新内容。

## 用途

以带饰边的对话框列出各版本的更新（新增、改进、修复），提供“不再显示”。

## 最小示例

```go
open, hide := true, false
releases := []account.Release{{Version: "1.1", Changes: []account.ReleaseChange{{Kind: account.ChangeNew, Text: "500 components."}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	account.WhatsNewDialog(c, &open, &hide, releases, account.WhatsNewDialogOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `open` | `*bool` | 必填 | 是否显示。 |
| `hide` | `*bool` | 必填 | “不再显示”的值，由调用方保存。 |
| `releases` | `[]Release` | 必填 | Version 必填；Date、Changes。 |
| `Title` | `string` | 内置“What's new” | 标题。 |

## 状态

显示、勾选不再显示。

## 事件

`ChoiceResult.Changed()` 在“不再显示”变化时。

## 键盘操作

“知道了”自动聚焦；Escape 关闭。

## 限制

- 无。
