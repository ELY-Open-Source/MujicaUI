# 031 Switch · 开关

## 用途

表示立即生效的开/关设置，如“自动备份”。组件只翻转值并回传结果，不调用任何业务接口；调用方在 `Changed()` 时执行操作，可用 `Loading` 表示“正在应用”。

## 最小示例

```go
wifi := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Switch(c, &wifi, "Wi-Fi", input.SwitchOptions{}).Changed() {
		log.Println("Wi-Fi:", wifi)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Switch(c *ui.Context, on *bool, label string, opts SwitchOptions) ui.Element`。`label` 为空时请在返回值上调用 `Label` 命名。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Disabled` | `bool` | false | 禁用：不响应点击与键盘，退出 Tab 序列，由 MyGo 置灰。 |
| `Loading` | `bool` | false | 调用方正在应用结果：滑块内显示旋转标记，保留尺寸与标签，暂不接受输入。 |
| `Description` | `string` | "" | 辅助说明，屏幕阅读器在标签后读出。 |

## 状态

关（象牙底、褐色方块滑块）、开（酒红轨道、象牙滑块，并在左侧显示竖线标记，不只依赖颜色）、悬停（边框与底色加深）、聚焦（轨道外 2 DIP 焦点环）、禁用、加载。切换动画 120 ms，减少动画模式下直接切换。

## 事件

`Changed()` 在用户切换的那一帧返回 true，`*on` 已是新值。若业务失败，调用方自行把值改回并提示错误。

## 键盘操作

- Tab / Shift+Tab：聚焦或离开。
- Space：切换。

## 限制

- 基于 `ui.SwitchBase`，无障碍角色为 Switch，开关状态由 MyGo 报告。
- 组件不保存“请求中”状态；`Loading` 由调用方根据自己的异步任务设置。
