# 286 ToolRegistryPanel 工具注册面板

## 用途

列出智能体可用的工具与启停开关，并预览所选工具的参数模式。

## 最小示例

```go
var state data.ListState[string]
tools := []feedback.AgentTool{{Name: "web_search", Description: "Search the web", Enabled: true,
	Params: []feedback.ToolParam{{Name: "query", Type: "string", Required: true}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.ToolRegistryPanel(c, &state, tools, feedback.ToolRegistryPanelOptions{})
	r.Element.Height(260)
	if name, on, ok := r.Toggled(); ok && name == tools[0].Name {
		tools[0].Enabled = on
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ToolRegistryPanel(c *ui.Context, s *ListState[string], tools []AgentTool, opts ToolRegistryPanelOptions) ToolRegistryResult`

## 参数

`AgentTool`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Name | string | 必填 | 唯一名称；重复 panic。 |
| Description | string | "" | 说明。 |
| Enabled | bool | false | 是否启用，由调用方持有。 |
| Params | []ToolParam | nil | 参数：名称、类型、必填、说明。 |

`ToolRegistryPanelOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Disabled | bool | false | 禁用开关。 |
| Version | int | 0 | 原地修改数据时递增；新切片或长度变化自动刷新。 |

给返回的 `Element` 一个高度；`List.Changed()` 报告工具选择。名称索引按数据版本缓存，参数预览独立滚动，按可见行虚拟化。

## 状态

- 开关使用 Switch，名称为“启用 工具名”。
- 参数以等宽名称、Info 色类型与“必填”描边标签显示。

## 事件

- `Toggled()`：被切换的工具与新状态。
- 列表选择变化时右侧预览随之更新。

## 键盘操作

- 列表 ↑/↓ 选择工具；开关用 Tab 到达，Space 切换。

## 限制

- 参数只做平铺预览，不展开嵌套对象。
