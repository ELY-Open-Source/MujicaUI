# 437 CronEditor · Cron 编辑器

## 用途

编辑标准 5 字段 cron 表达式，逐字段解释、就地报告错误，并预览接下来的运行时间。

## 最小示例

```go
expr := "0 9 * * MON-FRI"
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.CronEditor(c, &expr, mujica.CronEditorOptions{Location: time.Local}).Changed() {
		if _, err := mujica.ParseCron(expr); err == nil {
			fmt.Println("schedule", expr)
		}
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CronEditor(c *ui.Context, expr *string, opts CronEditorOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Now` | `time.Time` | 零值（系统时钟） | 预览的起点。 |
| `Runs` | `int` | 0（即 5） | 预览的次数。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

有效（字段卡片 + 运行时间）、错误（出错字段红框 + 说明，不显示预览）、禁用。

## 事件

`Changed()` 报告表达式文本变化；`ParseCron` 可单独校验并计算 `Next`。

## 键盘操作

- 在表达式框中输入；解释与预览随输入更新。

## 限制

- 支持列表、范围、步长、`*` 与月份/星期英文缩写；日与星期同时受限时任一匹配即运行（cron 惯例）。
- 夏令时跳过的时间不运行；重复的时间取第一次并标注。
- 只搜索未来五年。
