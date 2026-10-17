# 085 Tooltip 文字提示

## 用途

指针停留或键盘聚焦一段时间后，在元素旁显示简短说明。提示不含交互内容，不拦截指针，并作为元素的无障碍描述。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	b := input.Button(c, "Seal", input.ButtonOptions{})
	overlay.Tooltip(c, b, "Seal the letter with wax", overlay.TooltipOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Tooltip(c *ui.Context, anchor ui.Element, text string, opts TooltipOptions) ui.Element`，返回 `anchor`。`anchor` 为 nil 时 panic；`text` 为空时什么也不做。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Delay` | `time.Duration` | 500 ms | 悬停或键盘聚焦持续多久后显示；负数 panic。 |

## 状态

- 炭黑底（深色主题为象牙底）、2 DIP 圆角、12+1 DIP 文字，最宽 240 DIP 后换行。
- 默认在元素下方居中，下方放不下时翻到上方，左右限制在窗口内。

## 事件

无业务事件。`text` 同时写入 `anchor` 的无障碍描述。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Tab 聚焦 | 键盘焦点停留 `Delay` 后显示。 |
| Escape | 隐藏，直到指针离开或焦点移走。 |

## 限制

- Escape 只关闭最上层的 MujicaUI 浮层（窗口内按打开顺序记录）。MyGo 自带的 `ui.Popover`、`ui.Select` 弹层不在这份记录中；它们打开在非模态 MujicaUI 浮层之上时，Escape 可能先关闭下面的 MujicaUI 浮层。
- 只在键盘焦点（可见焦点）时显示，鼠标点击获得的焦点不触发。
- 不使用 MyGo 的 `Element.Tooltip`：它只跟随指针、不响应焦点（`ui/widgets.go:549`）。
