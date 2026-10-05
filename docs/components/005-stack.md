# 005 Stack 堆叠

## 用途

按横向或纵向排列子元素，统一间距、对齐与换行。用 `StackItem` 包裹可隐藏的子项，隐藏时不占位置也不留间距。

## 最小示例

```go
mujica.Stack(c, mujica.StackOptions{Horizontal: true, Gap: mujica.SpaceM}, func() {
	mujica.Button(c, "Save", mujica.ButtonOptions{})
	mujica.StackItem(!showDraft, func() { mujica.Button(c, "Draft", mujica.ButtonOptions{Variant: mujica.Secondary}) })
	mujica.Button(c, "Cancel", mujica.ButtonOptions{Variant: mujica.Ghost})
})
```

签名：

- `func Stack(c *ui.Context, opts StackOptions, children func()) *ui.Element`
- `func StackItem(hidden bool, build func())`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Horizontal` | `bool` | false | 横向排列；否则纵向。 |
| `Gap` | `Space` | `SpaceS`（2 单位） | 子项间距，随密度缩放。 |
| `Align` | `StackAlign` | `StackAlignAuto` | 交叉轴对齐：Auto（纵向拉伸、横向居中）、Start、Center、End、Stretch；非法值 panic。 |
| `Justify` | `ui.Align` | `ui.Start` | 主轴分布：Start、Center、End、SpaceBetween、SpaceAround、SpaceEvenly。 |
| `Wrap` | `bool` | false | 横向换行；纵向时设置会 panic。 |

`StackItem`：`hidden` 为 true 时不构建 `build`，因此既不占空间也不产生间距。`build` 为 nil 时 panic。

## 状态

布局组件，无交互状态。

## 事件

无。

## 键盘操作

不获取焦点；Tab 顺序即子项构建顺序。

## 限制

- 隐藏项不会被构建，其内部界面状态（如未提交的输入）随之丢弃；需要保留时把状态放在应用中。
- 只有横向支持换行。
