# 100 LoadingOverlay 加载遮罩

## 用途

在一个容器范围内显示加载遮罩：被覆盖的内容不接收指针、按键和焦点，防止重复操作；窗口其余部分照常可用。显示说明文字与可选取消操作。

## 最小示例

```go
r := mujica.LoadingOverlay(c, loading, mujica.LoadingOverlayOptions{Cancel: "Cancel"}, func() {
	ui.Text(c, "Quarterly report")
	mujica.Button(c, "Export", mujica.ButtonOptions{})
})
if r.Canceled {
	loading = false
}
```

签名：`func LoadingOverlay(c *ui.Context, loading bool, opts LoadingOverlayOptions, content func()) LoadingOverlayResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `loading` | `bool` | 必填 | 是否覆盖。 |
| `content` | `func()` | 必填 | 构建被覆盖的内容；nil 时 panic。 |
| `Label` | `string` | 内置“Loading… / 正在加载…” | 说明文字。 |
| `Cancel` | `string` | `""` | 取消按钮文字；为空时不显示。 |

返回 `LoadingOverlayResult{Element, Canceled}`。`Element` 即遮罩所覆盖的范围，可在其上设置宽高。

## 状态

- 加载时内容经 MyGo `Disabled` 阻止输入（MyGo 统一减淡），上方为半透明背景与居中面板。
- 开始加载时播报说明；焦点原在内容内时移到取消按钮（无取消时移到可聚焦的说明面板），焦点不会离开该区域。
- 遮罩淡入淡出，减少动画时直接切换。

## 事件

`Canceled`：点击取消，或焦点在面板内时按 Escape。

## 键盘操作

Tab 到达取消按钮；Enter、Space 或 Escape 取消。被覆盖的控件不进入 Tab 序列。

## 限制

- 只遮盖容器自身范围；整窗阻塞请使用模态对话框。
