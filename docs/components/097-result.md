# 097 Result 结果页

## 用途

任务完成后的整页结果：成功、失败或无权限，附说明与下一步操作。组件只报告选择了哪个操作，不执行跳转。标题与操作之间用旧金色细线与菱形分隔。

## 最小示例

```go
r := mujica.Result(c, mujica.ResultOptions{
	Kind:        mujica.ResultForbidden,
	Description: "You need the Editor role to change billing.",
	Actions:     []string{"Request access", "Go back"},
})
switch r.Action {
case 0: // send the request
case 1: // navigate back
}
```

签名：`func Result(c *ui.Context, opts ResultOptions) ResultResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Kind` | `ResultKind` | `ResultSuccess` | `ResultSuccess`、`ResultFailure`、`ResultForbidden`；未知值 panic。 |
| `Title` | `string` | 按类型的内置文案 | 衬线标题。 |
| `Description` | `string` | `""` | 说明。 |
| `Actions` | `[]string` | 无 | 下一步操作；第一个为主按钮，其余为次按钮。 |

## 状态

图标（对勾、叉、盾牌）与标题文字共同表达结果，不只靠颜色。按钮有悬停、按下与焦点状态。

## 事件

`ResultResult.Action`：被触发的操作序号，未触发为 -1。

## 键盘操作

Tab 按顺序聚焦操作按钮；Enter 或 Space 触发。

## 限制

- 不含路由；导航由调用方完成。
