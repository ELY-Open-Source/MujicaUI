# 482 CopyButton 复制按钮

## 用途

把文字放到剪贴板，按钮随即显示“已复制”与对勾两秒，并向屏幕阅读器播报。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.CopyButton(c, "https://example.com/invite/42", mujica.CopyButtonOptions{Label: "Copy link"})
}
```

签名：`func CopyButton(c *ui.Context, text string, opts CopyButtonOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `text` | `string` | 必填 | 要复制的文字；不进日志，只记录长度。 |
| `Label` | `string` | `""` | 可见文字；为空时是方形图标按钮。 |
| `Disabled` | `bool` | `false` | 不复制。 |

## 状态

默认、悬停、按下、焦点；复制后两秒显示“已复制”，宽度按两段文字中较长者预留，不跳动。

## 事件

`Changed()`：本帧完成复制。

## 键盘操作

Tab 聚焦；Enter 或 Space 复制。

## 限制

剪贴板写入由 MyGo 完成，无法得知系统是否拒绝写入。
