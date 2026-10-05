# 216 BranchNavigator 分支切换

## 用途

在一条消息的多个分支间切换，显示“< 2 / 3 >”。

## 最小示例

```go
if chat.BranchNavigator(c, &branch, 3).Changed() {
	fmt.Println("show branch", branch)
}
```

签名：`func BranchNavigator(c *ui.Context, index *int, count int) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `index` | `*int` | 必填 | 当前分支，0 起；不在 0 到 count−1 之间时 panic。 |
| `count` | `int` | 必填 | 分支数，至少 1。 |

## 状态

首末分支时对应箭头禁用；切换时向辅助技术播报“第 n 个分支，共 m 个”。

## 事件

`Changed()`：分支改变。

## 键盘操作

整个控件是一个 Tab 停靠点：Left/Right 切换，Home、End 到第一个与最后一个；箭头按钮只响应指针。

## 限制

- 箭头按钮不进入 Tab 序列，避免到达首末分支时焦点停在禁用按钮上。
