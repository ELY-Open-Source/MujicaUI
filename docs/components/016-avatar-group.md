# 016 AvatarGroup 头像组

## 用途

把多个头像重叠排列。超过最大数量时显示“+N”按钮，点开后在弹层中列出完整名单。

## 最小示例

```go
people := []mujica.AvatarItem{
	{Name: "Sakiko Togawa"}, {Name: "Uika Misumi"}, {Name: "Mutsumi Wakaba"},
	{Name: "Umiri Yahata"}, {Name: "Nyamu Yutenji"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AvatarGroup(c, people, mujica.AvatarGroupOptions{Max: 3})
}
```

## 参数

`func AvatarGroup(c *ui.Context, items []AvatarItem, opts AvatarGroupOptions) AvatarGroupView`

`AvatarItem`：`Name string`、`Image *ui.Bitmap`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Max` | `int` | `0`（即 4） | 最多直接显示的头像数。负数 panic。 |
| `Size` | `float32` | `0`（即 32） | 每个头像的直径。负数 panic。 |

返回值 `AvatarGroupView`：`Element`；`Open` 表示完整名单正在显示；`Toggled` 表示本帧打开或关闭了名单。

## 状态

- 默认：头像之间重叠四分之一直径，用 2 DIP 窗口背景色描边分隔。
- “+N”按钮：基于 `ui.CollapsibleBase` 的 Trigger，辅助技术读作展开/收起按钮，并报告是否已展开。命中区域至少 28×28 DIP（`Size` 很小时也如此）。悬停、按下改变背景；聚焦时显示 2 DIP 焦点环。
- 展开：按钮下方弹出 6 DIP 圆角面板，列出全部头像与姓名，超过 240 DIP 时滚动；遵循减少动画设置。

## 事件

- `AvatarGroupView.Toggled` / `Open`：名单展开或关闭（点击按钮、Enter、Space、Escape、点击外部）。同时记录 `slog.Debug("mujica: avatargroup toggle")`。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦“+N”按钮（读作“显示全部 N 人”） |
| Enter / Space | 展开或收起完整名单 |
| Escape | 关闭名单 |

## 限制

- 分隔描边使用窗口 `Background` 颜色，放在 `Surface` 卡片上时描边颜色会与卡片背景不同。
- 名单中的条目只用于展示，不可选择。
