# 475 ApiKeyManager

API 密钥管理。

## 用途

密钥列表只显示前缀与掩码；创建后明文只显示一次并可复制，关闭即从内存清除；撤销需确认。

## 最小示例

```go
var state mujica.ApiKeyManagerState
keys := []mujica.APIKey{{ID: "k1", Name: "Deploy bot", Prefix: "mk_live_9f2c"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if name, ok := mujica.ApiKeyManager(c, &state, keys, mujica.ApiKeyManagerOptions{}).Create(); ok {
		fmt.Println("create", name)
		state.Reveal("secret from your key service")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*ApiKeyManagerState` | 必填 | 对话框与一次性明文。 |
| `keys` | `[]APIKey` | 必填 | ID、Name、Prefix 必填；不含明文。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

列表、空列表、输入名称、创建中、显示明文、失败、撤销确认、禁用。

## 事件

`Create()` 返回新密钥名称：调用方生成后调用 `s.Reveal(secret)` 或 `s.Fail(msg)`；`Revoked()` 返回确认撤销的 ID。

## 键盘操作

Tab 到创建与撤销按钮；对话框中 Enter 创建；复制按钮写入剪贴板并播报“已复制”。

## 限制

大集合使用虚拟列表，仅构建当前可见行。

- 验证码与密钥通过调用方状态接入；组件事件与日志仅包含完成或失败标记。
- 明文作为可选择文字显示一次，屏幕阅读器可读出文字本身。
- 创建进行中等待调用方结果；禁用关闭全部弹层并清除已显示明文。
- 待处理请求跨禁用保留；失败后恢复名称与错误，支持重试。
