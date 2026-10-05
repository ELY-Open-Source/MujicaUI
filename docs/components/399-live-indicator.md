# LiveIndicator

## 用途

显示实时同步状态、离线说明和重试。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.LiveIndicator(c, files.LiveOffline, files.LiveIndicatorOptions{})
}
ui.Render(view, 500, 100, 1)
```

## 参数

LastSynced 可选；Now 默认帧时间；Error 补充失败原因。

## 状态

已同步、同步中、重连、离线、失败。

## 事件

ChoiceResult.Changed() 表示用户请求重试。

## 键盘操作

Tab 到重试；Enter / Space 激活。

## 限制

同步由调用方执行；状态变化会朗读，减少动画使用静态说明。
