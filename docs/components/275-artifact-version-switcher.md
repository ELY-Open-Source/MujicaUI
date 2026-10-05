# 275 ArtifactVersionSwitcher 产物版本切换

## 用途

在产物的多个版本之间切换，并提供与上一版本比较差异的入口。

## 最小示例

```go
current := 1
versions := []agent.ArtifactVersion{{Label: "v1", Note: "First draft"}, {Label: "v2", Note: "Added music"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if agent.ArtifactVersionSwitcher(c, &current, versions, agent.ArtifactVersionSwitcherOptions{}).Compare {
		fmt.Println("diff v1 and v2")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ArtifactVersionSwitcher(c *ui.Context, current *int, versions []ArtifactVersion, opts ArtifactVersionSwitcherOptions) ArtifactVersionResult`。没有版本、`*current` 越界或版本 `Label` 为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Disabled` | `bool` | `false` | 禁止切换与比较。 |

## 状态

- 显示版本名、“N / 总数”和说明；第一个版本时“较早版本”禁用，最后一个时“较新版本”禁用。
- 非第一个版本显示“与上一版本比较”。

## 事件

- `Changed`：`*current` 改变。
- `Compare`：请求比较当前版本与上一版本。

## 键盘操作

- 焦点在组件内时 ← → 切换版本；Tab 到各按钮，Enter / Space 触发。

## 限制

- 差异视图由调用方打开，可用 FileChangeCard 或 MultiFileDiffReview 呈现。
