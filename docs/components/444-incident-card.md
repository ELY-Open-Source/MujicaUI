# 444 IncidentCard · 事件卡片

## 用途

展示事件标题、状态、影响与按时间排列的更新。

## 最小示例

```go
incident := mujica.Incident{Title: "API latency", Status: mujica.IncidentMonitoring, Impact: mujica.ImpactMinor, Updates: []mujica.IncidentUpdate{{Status: mujica.IncidentIdentified, Text: "Slow query located"}, {Status: mujica.IncidentMonitoring, Text: "Recovery observed"}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.IncidentCard(c, incident)
}
ui.Render(view, 720, 480, 1)
```

## 参数

Incident.Title 必填；Status、Impact 使用枚举；Started 与 Updates 的时间由调用方提供。

## 状态

调查中、已定位、观察中、已解决；四级影响。

## 事件

展示组件，调用方更新 Incident。

## 键盘操作

展示内容随调用方更新。

## 限制

时间线按输入顺序展示；调用方排序。
