# 092 Toast 轻提示

## 用途

在窗口右下角短暂显示操作结果，可带一个操作（如撤销）与关闭按钮。每个窗口一条队列：按 ID 去重，同时最多显示 3 条，其余排队。

## 最小示例

```go
if mujica.Button(c, "Delete note", mujica.ButtonOptions{Variant: mujica.Danger}).Clicked() {
	mujica.ShowToast(c, mujica.ToastMessage{ID: "note-1", Message: "Note deleted", Action: "Undo"})
}
// Last in the view:
if r := mujica.Toaster(c); r.Action == "note-1" {
	// restore the note
}
```

接口：

- `func ShowToast(c *ui.Context, m ToastMessage)`：加入队列，须在界面线程调用（后台任务经 `Window.Update` 回到界面线程）。
- `func DismissToast(c *ui.Context, id string)`：关闭显示中或排队中的消息。
- `func Toaster(c *ui.Context) ToasterResult`：每帧在视图末尾调用一次，绘制消息并报告事件。

## 参数

`ToastMessage`：

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `string` | `Message` | 去重键。相同 ID 再次显示时替换内容并重新计时，不新增一条。 |
| `Severity` | `Severity` | `SeverityInfo` | 图标、色条与颜色。 |
| `Message` | `string` | 必填 | 文字；为空时 panic。 |
| `Action` | `string` | `""` | 操作按钮文字。 |
| `Duration` | `time.Duration` | `ToastDuration`（5 秒） | 显示时长；负数 panic。错误消息在为 0 时一直保留到用户关闭。 |

常量：`ToastDuration = 5 * time.Second`，`ToastVisible = 3`。

## 状态

- 排队消息在显示时才开始计时。
- 指针悬停或焦点位于消息内时暂停计时，离开后从剩余时间继续。
- 出现时向读屏播报文字；进出场 160 ms 位移与淡入，减少动画时直接切换。

## 事件

`ToasterResult{Action, Closed}`：本帧被点击操作或被关闭的消息 ID；两者都会移除该消息。到时自动消失不报告。

## 键盘操作

消息从不抢夺焦点，正在输入的文本框保持焦点。Tab 在页面之后到达消息中的按钮；焦点在消息内时 Escape 关闭该消息。

## 复用与限制

MyGo 自带 `c.Toast` / `c.ToastAction`（docs/ui/toast.md），但它按文字去重、没有 ID、不限制条数、无关闭按钮、错误不能常驻、外观不可替换，计时只在悬停时暂停而不在聚焦时暂停。这些都是 PRD §5.5 的要求，因此 MujicaUI 自建队列，只复用 MyGo 的 `Overlay`、`Transition`、`Announce` 与 `After` 定时唤醒。

- 所有消息（不只是含操作的消息）都在悬停或聚焦时暂停，这是 PRD 要求的超集，便于阅读。
- 位置固定在右下角；切换到未调用 `Toaster` 的视图时消息仍在队列中，回来后继续显示。
