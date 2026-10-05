# 359 DeviceSelector 设备选择

## 用途

选择麦克风、摄像头、扬声器，设备列表由调用方提供。

## 最小示例

```go
devices := mujica.MediaDevices{
	Microphones: []mujica.MediaDevice{{ID: "mic-1", Name: "Built-in Microphone"}},
	Speakers:    []mujica.MediaDevice{{ID: "out-1", Name: "Headphones"}},
}
sel := mujica.DeviceSelection{Microphone: "mic-1"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.DeviceSelector(c, &sel, devices, mujica.DeviceSelectorOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DeviceSelector(c *ui.Context, sel *DeviceSelection, devices MediaDevices, opts DeviceSelectorOptions) ChoiceResult`

`MediaDevices` 的某类列表为 nil 时不显示该类；为空切片时显示禁用的“未找到设备”。设备 ID 重复，或已选 ID 不在列表中时 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Disabled` | `bool` | false | 禁用全部选择。 |

## 状态

每类一行：图标、名称和 Select。

## 事件

`Changed()` 在任一类选择变化后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 在各选择框间移动 |
| Enter / Space | 打开列表 |
| ↑ / ↓ | 移动 |
| 输入字母 | 跳到匹配的设备 |

## 限制

- 组件不枚举设备，也不监听设备插拔；调用方更新列表即可。
