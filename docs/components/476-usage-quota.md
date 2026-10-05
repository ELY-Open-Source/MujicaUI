# 476 UsageQuota

用量配额。

## 用途

每项资源的用量与上限进度条；接近与超出上限以图标与文字提示；显示重置时间。

## 最小示例

```go
items := []mujica.QuotaItem{{Label: "API calls", Used: 4200, Limit: 10000}, {Label: "Storage", Used: 9, Limit: 10, Unit: "GB"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.UsageQuota(c, items, mujica.UsageQuotaOptions{ResetAt: time.Now().Add(72 * time.Hour)})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `items` | `[]QuotaItem` | 必填 | Label、Used、Limit（>0）、Unit。 |
| `ResetAt` | `time.Time` | 零值 | 重置时间；一天内显示小时，否则显示日期。 |
| `Warn` | `float64` | 0.8 | 达到该比例即提示接近上限。 |

## 状态

正常、接近上限、已达上限。

## 事件

无。

## 键盘操作

不获取焦点；每项以计量角色读出。

## 限制

- 无。
