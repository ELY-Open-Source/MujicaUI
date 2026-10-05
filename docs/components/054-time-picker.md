# 054 TimePicker · 时间选择器

## 用途

分段编辑时、分、秒的时间输入，可在弹层列表中选择；支持 12/24 小时制、步长和范围。

## 最小示例

```go
alarm := mujica.TimeOfDay{Hour: 7, Minute: 30}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.TimePickerOptions{Label: "Alarm", Hour12: true, MinuteStep: 5}
	if mujica.TimePicker(c, &alarm, opts).Changed() {
		fmt.Println("alarm", alarm)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TimePicker(c *ui.Context, value *TimeOfDay, opts TimePickerOptions) DateTimeResult`，返回值的 `Element` 字段是组件元素。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Seconds` | `bool` | false | 显示并编辑秒。 |
| `Hour12` | `bool` | false | 12 小时制，附加上午/下午段。 |
| `MinuteStep` | `int` | 0（即 1） | 分钟步长，必须整除 60。 |
| `SecondStep` | `int` | 0（即 1） | 秒步长，必须整除 60。 |
| `Min` | `TimeOfDay` | 00:00:00 | 最早时间。 |
| `Max` | `TimeOfDay` | 零值（即 23:59:59） | 最晚时间。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 空 | 组名；各段名称为“Label: 时/分/秒”。 |

## 状态

默认、悬停、聚焦段（酒红填充 + 键盘焦点环）、只读、错误（危险色边框 + 说明，草稿保留）、禁用。弹层列表中超出范围的项禁用。

## 事件

`Changed()`：编辑后的时间合法（在范围内且符合步长）并与原值不同时为 true。不合法时不写入，显示错误。

## 键盘操作

- Tab：进入时间段（整组一个 Tab 停靠点）。
- ←/→：切换时、分、秒、上午/下午。
- ↑/↓：按步长增减，整圈循环。
- 数字键：直接键入，时完成后自动跳到分。
- A/P：设置上午/下午。
- Alt+↓ 或时钟按钮：打开列表弹层；Escape 关闭。

## 限制

- 超出范围或不符合步长的时间会提示，不会被自动夹取或取整；键入超出段范围的数字（如小时 29）会保留显示并报错。
- 段顺序固定为“时:分:秒 上午/下午”，zh-CN 下上午/下午也位于末尾。
- MyGo 的 `ui.TimeInput` 只有时分且无 12 小时制、步长和范围（ui/timeinput.go:18），因此本组件自绘分段。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
