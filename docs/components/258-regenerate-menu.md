# 258 RegenerateMenu · 重新生成菜单

## 用途

重新生成回复：同模型重试、换模型重试，以及更长、更短等改写变体。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.RegenerateMenu(c, chat.RegenerateMenuOptions{Models: []string{"Opus", "Sonnet"}, Current: "Opus"}); v.Chosen {
		log.Println("retry", v.Choice.Kind, v.Choice.Model, v.Choice.Variant)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func RegenerateMenu(c *ui.Context, opts RegenerateMenuOptions) RegenerateMenuView`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Models` | `[]string` | nil | 可换用的模型；为空或只有当前模型时子菜单禁用。 |
| `Current` | `string` | "" | 当前模型，不出现在子菜单中。 |
| `Variants` | `[]string` | nil（内置） | 替换内置的更长、更短、更简单、更详细。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

触发按钮（刷新图标）、菜单打开、子菜单、分组标题、禁用项。

## 事件

`Chosen` 与 `Choice`：`Kind` 为 `RegenerateAgain`、`RegenerateModel`（`Model`）或 `RegenerateVariant`（`Variant`）。

## 键盘操作

- Enter、Space 或 ↓ 打开菜单；↑/↓ 移动；→ 打开子菜单，← 返回；Enter 选择；Escape 关闭。

## 限制

- 菜单在窗口内绘制（DropdownMenu），不使用系统菜单。
