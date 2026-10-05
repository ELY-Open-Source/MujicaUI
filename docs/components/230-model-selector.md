# 230 ModelSelector 模型选择器

## 用途

模型列表；能力标签与上下文长度；不可用项说明。

## 最小示例

```go
model := "herald"
models := []chat.ModelOption{{ID: "herald", Name: "Herald 2", Capabilities: []string{"Vision"}, ContextTokens: 200000}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.ModelSelector(c, &model, models, chat.ModelSelectorOptions{}).Changed() {
		fmt.Println("model", model)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`value` 为模型 ID 指针；`models` 每项 `ID`、`Name` 必填，`Capabilities` 为能力标签，`ContextTokens=0` 隐藏容量；`Unavailable` 提供不可用原因。`Label` 默认“模型”，`Disabled=false`。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

ChoiceResult.Changed() 表示模型 ID 已改变。

## 键盘操作

Tab 到选择器；Up/Down、Home/End 选择，Enter/Space 确认，Escape 关闭。

## 限制

模型服务与可用性由调用方提供。
