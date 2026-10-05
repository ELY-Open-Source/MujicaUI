# 179 RiskMeter

风险等级。

## 用途

分段刻度显示风险等级，当前段上方有菱形标记，并显示等级名称与说明。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.RiskMeter(c, 4, finance.RiskMeterOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `level` | `int` | 必填 | 从 1 开始，超出范围 panic。 |
| `Levels` | `[]RiskLevel` | 内置五级 | 各级名称与说明，至少两级。 |
| `Label` | `string` | 内置“Risk level” | 名称。 |

## 状态

各等级。

## 事件

无。

## 键盘操作

不获取焦点；以计量角色读出“第 n 级，共 N 级”和名称。

## 限制

- 颜色只辅助：位置标记与等级文字同时表达等级。
