# 477 OnboardingWizard

引导向导。

## 用途

多步引导：进度条、步骤内容、上一步、下一步、跳过与完成。

## 最小示例

```go
step := 0
steps := []mujica.WizardStep{{Title: "Welcome"}, {Title: "Pick a theme"}, {Title: "Done"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.OnboardingWizard(c, &step, steps, mujica.OnboardingWizardOptions{}).Finished() {
		fmt.Println("onboarded")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `step` | `*int` | 必填 | 当前步骤，越界 panic。 |
| `steps` | `[]WizardStep` | 必填 | Title 必填；Description、Content、Ready。 |
| `NoSkip` | `bool` | false | 不显示跳过。 |

## 状态

各步骤、未就绪（下一步禁用）、最后一步（完成）。

## 事件

`Finished()`、`Skipped()`。

## 键盘操作

下一步按钮自动聚焦；Tab 在按钮与内容间移动。

## 限制

- 无。
