# 277 MCPServerList MCP 服务器列表

## 用途

列出 MCP 服务器：连接状态、传输方式、工具数，提供启动、停止与重新连接回调。

## 最小示例

```go
servers := []agent.MCPServer{{ID: "fs", Name: "filesystem", Transport: "stdio", State: agent.MCPConnected, Tools: 12}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := agent.MCPServerList(c, servers, agent.MCPServerListOptions{}); r.Action == agent.MCPStop {
		fmt.Println("stop", r.Server)
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func MCPServerList(c *ui.Context, servers []MCPServer, opts MCPServerListOptions) MCPServerListResult`。`ID` 为空或重复、`Name` 为空、状态未知、`Tools` 为负、`MaxHeight` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `MaxHeight` | `float32` | 360 | 列表最大高度，超出后滚动。 |
| `Disabled` | `bool` | `false` | 禁止所有操作。 |
| `Still` | `bool` | `false` | 停止执行中动画；系统减少动画时也自动静止。 |

## 状态

- `MCPServer` 字段：`ID`、`Name`、`Transport`、`State`（`MCPStopped`、`MCPConnecting`、`MCPConnected`、`MCPError`）、`Tools`、`Error`。
- 已连接时显示“N 个工具”；连接失败时显示错误文字与“重新连接”。
- 超过 200 台服务器时虚拟化。

## 事件

- `Action`：本帧请求 `MCPStart`、`MCPStop` 或 `MCPReconnect`，`Server` 为对应 ID。

## 键盘操作

- Tab 依次到各行按钮，Enter / Space 触发。

## 限制

- 组件只发出请求，连接管理由调用方完成。
