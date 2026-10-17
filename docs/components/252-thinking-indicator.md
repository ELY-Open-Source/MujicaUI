# 252 ThinkingIndicator · 思考指示

## 用途

显示助手正在思考：三个菱形依次起伏，旁边是阶段文字；减少动画时静止并加省略号。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.ThinkingIndicator(c, chat.ThinkingIndicatorOptions{Phase: "Searching the archives"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ThinkingIndicator(c *ui.Context, opts ThinkingIndicatorOptions) ui.Element`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Phase` | `string` | 内置“正在思考” | 阶段文字，如“正在检索资料”。 |

## 状态

动画（1.2 秒一周）、减少动画（静止菱形 + “…”）。

## 事件

无。

## 键盘操作

不获取焦点；角色为 Status，辅助技术读出阶段文字。

## 限制

- 动画只在绘制时请求下一帧；组件隐藏后不再重绘。
