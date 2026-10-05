# 432 TimezoneSelect · 时区选择

## 用途

按城市名或偏移搜索 IANA 时区，常用时区置顶，显示当前偏移。

## 最小示例

```go
zone := "Asia/Shanghai"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if datetime.TimezoneSelect(c, &zone, datetime.TimezoneSelectOptions{}).Changed() {
		fmt.Println("zone", zone)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TimezoneSelect(c *ui.Context, zone *string, opts TimezoneSelectOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Zones` | `[]string` | 内置世界列表 | 可选时区。 |
| `Common` | `[]string` | 内置常用列表 | 置顶的常用时区。 |
| `Now` | `time.Time` | 零值（系统时钟） | 计算偏移所用的时刻（夏令时会改变偏移）。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

关闭、打开（搜索框 + 分组列表）、高亮、已选（勾）、无结果、禁用。

## 事件

`Changed()` 报告选中的时区变化。

## 键盘操作

- Enter、Space 或 ↓ 打开；输入搜索；↑/↓ 移动高亮；Enter 选择；Escape 关闭。

## 限制

- 内嵌 `time/tzdata`，离线可用；未知时区名 panic。
