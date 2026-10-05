# 234 TokenCounter Token计数

## 用途

当前输入 token 数与上限；估算值由调用方提供。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TokenCounter(c, 1234, 8000)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`used` 为调用方估算的非负 token 数；`limit` 为正数；达到 90% 时警告，超过上限时显示超额。

## 状态

展示由输入数据驱动；边界值与警告状态带文字说明。

## 事件

展示组件，调用方更新估算值。

## 键盘操作

展示组件，沿用外层 Tab 顺序。

## 限制

token 分词与模型计算由调用方完成。
