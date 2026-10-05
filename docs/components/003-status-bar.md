# 003 StatusBar 状态栏

## 用途

窗口底部的状态条，分左右两组条目：状态文字、图标、进度与可点击操作。宽度不足时自动隐藏次要条目，并显示剩余数量。

## 最小示例

```go
bar := mujica.StatusBar(c, mujica.StatusBarOptions{
	Left:  []mujica.StatusItem{{ID: "sync", Text: "Syncing", ShowProgress: true, Progress: 0.4}},
	Right: []mujica.StatusItem{{ID: "bell", Label: "Notifications", Icon: icons.Must("bell"), Action: true}},
})
if bar.Clicked == "bell" {
	opened++
}
```

签名：`func StatusBar(c *ui.Context, opts StatusBarOptions) StatusBarResult`

## 参数

`StatusBarOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Left` | `[]StatusItem` | nil | 左侧条目，从左到右。 |
| `Right` | `[]StatusItem` | nil | 右侧条目，从左到右。 |

`StatusItem`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `string` | 必填 | 稳定键，也是 `Clicked` 的值；为空或重复 panic。 |
| `Text` | `string` | "" | 状态文字（12 DIP，随系统文本缩放）。 |
| `Icon` | `*ui.SVG` | nil | 前置图标（16 DIP）。 |
| `Tone` | `StatusItemTone` | `StatusItemNeutral` | 语义色：Neutral、Success、Warning、Danger、Info。 |
| `ShowProgress` | `bool` | false | 显示 64 DIP 进度条。 |
| `Progress` | `float32` | 0 | 进度 0–1，超出 panic。 |
| `Action` | `bool` | false | 条目成为按钮。 |
| `Label` | `string` | "" | 无文字条目的无障碍名称；`Text` 与 `Label` 都为空时 panic。 |
| `Secondary` | `bool` | false | 空间不足时优先隐藏。 |

## 状态

- 默认：Surface 底色，顶部 Border 线，组内条目以细线分隔。
- 操作条目：悬停 SurfaceHover，按下 Selection，键盘焦点环。
- 折叠：按“右组末项 → 左组首项”的顺序隐藏 `Secondary` 条目，直到放得下（按实际宽度计算，操作条目至少 28 DIP，含条目间距）；显示 `+N` 按钮，无障碍名称为“另有 N 项”，点击或 Enter/Space 打开弹层，列出被隐藏的条目，其中的操作照常可用，执行后弹层关闭。主要条目永不隐藏。
- 文本放大时状态栏随文字增高。
- 焦点所在条目被折叠时，焦点移到收纳它的 `+N` 按钮上，不会丢失。
- 进度：AccentText 填充，Border 轨道，描述为百分比。

## 事件

- `Clicked`：本帧被点击的操作条目 ID，没有则为空。一次点击只报告一次。

## 键盘操作

Tab 到达操作条目和 `+N` 按钮，Enter 或 Space 触发。`+N` 打开的弹层中 Tab 依次到达隐藏的操作，Escape 关闭弹层并把焦点还给 `+N`。纯文字条目不进入 Tab 序列。

## 限制

- 折叠依据上一帧的宽度和文字测量估算，窗口缩放后下一帧生效。
- 不提供不确定进度动画；需要时在 `Text` 中说明。
