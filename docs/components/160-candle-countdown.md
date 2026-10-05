# 160 CandleCountdown

## 用途

当前 K 线收盘倒计时；随周期变化；隐藏时停止计时重绘。

## 最小示例

```go
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandleCountdown(c, "1h", mujica.CandleCountdownOptions{Location: time.UTC})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

interval 为 Intervals 中的周期；Location 默认 UTC；Now 零值读取帧时钟，指定值冻结时间。

## 状态

随周期和时间更新；指定 Now 显示冻结读数；移出可见区域停止安排计时帧。

## 事件

返回 Element；时间由帧时钟或调用方指定。

## 键盘操作

展示型组件；周期由调用方选择控件调整。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
