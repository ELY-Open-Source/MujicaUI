# 005 Stack 堆叠

## 用途

按横向或纵向排列子元素，统一间距、对齐与换行。要隐藏子项就不构建它（`if !hidden { build() }`），未构建的子项不占位置也不留间距。

## 最小示例

```go
layout.Stack(c, layout.StackOptions{Horizontal: true, Gap: layout.SpaceM}, func() {
	input.Button(c, "Save", input.ButtonOptions{})
	if showDraft {
		input.Button(c, "Draft", input.ButtonOptions{Variant: input.Secondary})
	}
	input.Button(c, "Cancel", input.ButtonOptions{Variant: input.Ghost})
})
```

签名：

- `func Stack(c *ui.Context, opts StackOptions, children func()) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Horizontal` | `bool` | false | 横向排列；否则纵向。 |
| `Gap` | `Space` | `SpaceS`（2 单位） | 子项间距，随密度缩放。 |
| `Align` | `StackAlign` | `StackAlignAuto` | 交叉轴对齐：Auto（纵向拉伸、横向居中）、Start、Center、End、Stretch；非法值 panic。 |
| `Justify` | `ui.Align` | `ui.Start` | 主轴分布：Start、Center、End、SpaceBetween、SpaceAround、SpaceEvenly。 |
| `Wrap` | `bool` | false | 横向换行；纵向时设置会 panic。 |

## 状态

布局组件，无交互状态。

## 事件

无。

## 键盘操作

不获取焦点；Tab 顺序即子项构建顺序。

## 限制

- 隐藏项不会被构建，其内部界面状态（如未提交的输入）随之丢弃；需要保留时把状态放在应用中。
- 只有横向支持换行。
