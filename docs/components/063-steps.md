# 063 Steps 步骤条

## 用途

显示多步流程的进度：已完成、当前、等待和出错四种状态；横向或纵向排列；每一步是否可点击由调用方决定。

## 最小示例

```go
mujica.Steps(c, &current, []mujica.Step{
	{Title: "Account", Clickable: true},
	{Title: "Profile", Description: "Name and seal"},
	{Title: "Confirm"},
}, mujica.StepsOptions{Label: "Sign-up"})
```

签名：`func Steps(c *ui.Context, current *int, steps []Step, opts StepsOptions) StepsResult`

## 参数

`current`：进行中的步骤下标；之前的步骤视为已完成。等于 `len(steps)` 表示全部完成；超出范围 panic。`steps` 为空时 panic。

`Step`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 步骤标题；为空时 panic。 |
| `Description` | `string` | "" | 标题下的说明文字。 |
| `Error` | `bool` | false | 标为出错，优先于位置推导的状态。 |
| `Clickable` | `bool` | false | 允许用户点击跳到该步。 |

`StepsOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Vertical` | `bool` | false | 纵向排列，连接线在标记下方。 |
| `Label` | `string` | "" | 整组的无障碍名称。 |

返回 `StepsResult`，内嵌 `*ui.Element`。

## 状态

- 已完成：Accent 圆形标记内 OnAccent 勾号；之后的连接线为 AccentText。
- 当前：AccentText 圆环与序号，标题 600 字重。
- 等待：ControlBorder 圆环与 TextMuted 序号、标题。
- 出错：Danger 圆环与叉号，标题 Danger 色。
- 每步的无障碍名称为“第 n 步，共 m 步：标题，状态”，状态不只靠颜色表达。

## 事件

- `Changed()`：本帧用户点击了可点击的步骤，`*current` 已更新。

## 键盘操作

- 可点击的步骤组成一个 Tab 停靠点，方向键在其间移动，Enter 或 Space 跳到该步。不可点击的步骤不获取焦点。

## 限制

- MyGo 没有 step 或 aria-current 角色；步骤以列表项（可点击时为按钮）呈现，状态写入名称。
