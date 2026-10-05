# 351 VolumeControl 音量控制

## 用途

音量滑块与静音按钮。

## 最小示例

```go
volume, muted := 0.8, false
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.VolumeControl(c, &volume, &muted, mujica.VolumeControlOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func VolumeControl(c *ui.Context, volume *float64, muted *bool, opts VolumeControlOptions) ChoiceResult`

`*volume` 取 0 到 1，超出 panic；步长 5%。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

静音或音量为 0 时显示静音图标；静音按钮开启时有选中底色和下划线。静音时调高音量会自动取消静音。

## 事件

`Changed()` 在音量或静音变化后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| ← / → | 音量 ±5% |
| Home / End | 0% / 100% |
| M | 组件内切换静音 |

## 限制

- 组件不调节系统音量，只改调用方的值。
