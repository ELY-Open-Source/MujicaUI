# 446 ServiceStatus · 服务状态

## 用途

展示服务正常、降级、中断和总体统计。

## 最小示例

```go
services := []devtools.ServiceItem{{Name: "API", Health: devtools.ServiceOperational}, {Name: "Search", Health: devtools.ServiceDegraded, Note: "Slow queries"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	devtools.ServiceStatus(c, services, devtools.ServiceStatusOptions{})
}
ui.Render(view, 720, 480, 1)
```

## 参数

MaxHeight 默认 360 DIP；ServiceItem.Name 必填；Health 为有效枚举；Note 展示原因。

## 状态

正常、降级、中断、空列表；总体统计随数据更新。

## 事件

展示组件，调用方更新服务切片。

## 键盘操作

滚动列表浏览服务。

## 限制

状态来自调用方；每行使用图标和文字区分状态。更新服务数据时提交新切片。
