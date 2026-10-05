# 418 RecipientInput 收件人输入

## 用途

通过标签输入收件人，校验邮箱格式并提供联系人建议。

## 最小示例

```go
var to []string
contacts := []mujica.Contact{{Name: "Uika Misumi", Email: "uika@example.com"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RecipientInput(c, &to, mujica.RecipientInputOptions{Contacts: contacts, Label: "To"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`recipients *[]string` 保存邮箱标签。`Contacts []Contact` 提供 Name、Email；联系人数据变化时替换切片。`Label` 默认“收件人”；`Placeholder` 默认空；`Disabled` 默认 false。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

返回 `TagsInputParts`；`Changed()` 表示标签变化，`Input` 可用于焦点控制。

## 键盘操作

Enter 或分隔符提交标签；空输入 Backspace 删除末项；上下方向键移动联系人建议，Enter 选择，Escape 收起建议；Tab 可聚焦建议并显示焦点环。

## 限制

联系人姓名中的逗号和引号按邮件地址语法转义。建议最多显示 6 项；联系人匹配在查询或联系人、收件人数据变化时计算；业务邮箱可达性由调用方验证。
