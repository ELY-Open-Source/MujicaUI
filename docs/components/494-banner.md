# 494 Banner 页面横幅

## 用途

页面顶部的通栏提示：信息、成功、警告、错误；可带操作和关闭按钮。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.Banner(c, &open, "Maintenance tonight at 22:00.", mujica.BannerOptions{
		Severity: mujica.SeverityWarning,
		Action:   "Details",
	})
	if r.Action {
		// show the details
	}
}
```

签名：`func Banner(c *ui.Context, open *bool, message string, opts BannerOptions) AlertResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `open` | `*bool` | `nil` | 非 nil 时显示关闭按钮；为 false 时不显示。 |
| `message` | `string` | 必填 | 文字。 |
| `Severity` | `Severity` | `SeverityInfo` | 严重程度；未知值 panic。 |
| `Action` | `string` | `""` | 操作按钮文字。 |

## 状态

图标、底边色条与浅色底同时表达严重程度；警告与错误出现或变化时播报。

## 事件

`AlertResult.Action`、`AlertResult.Closed`。

## 键盘操作

Tab 到操作与关闭按钮；Enter/Space 触发。

## 限制

无。
