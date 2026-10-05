# 304 ConflictResolver 冲突解决

## 用途

读取带 git 冲突标记的文本，列出冲突块，逐块采纳当前、传入或两者，并显示已解决计数。

## 最小示例

```go
var state git.ConflictResolverState
text := "a\n<<<<<<< HEAD\nmine\n=======\ntheirs\n>>>>>>> feature\nb"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if git.ConflictResolver(c, &state, text, git.ConflictResolverOptions{}).Changed() {
		_ = state.Result() // save the file once every block is resolved
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ConflictResolver(c *ui.Context, s *ConflictResolverState, text string, opts ConflictResolverOptions) ChoiceResult`

## 参数

`ConflictResolverOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Disabled | bool | false | 禁用采纳。 |

`ConflictResolverState.Result()` 返回替换后的文本，未解决的块保留原标记；`Resolved()` 返回已解决与总数。只在文本改变时重新解析。

## 状态

- 顶部显示“已解决 n / m 个冲突”与 Success 色进度条。
- 未解决块标题为警示图标与“未解决”，已解决为对勾与“已解决”。
- 被采纳的一侧显示浅色底与对勾，按钮带对勾；已解决块提供“撤销”。
- 标记未闭合时显示 Danger 色提示与起始行号，不列出冲突。

## 事件

- `Changed()`：某块的选择改变。

## 键盘操作

- Tab 依次经过各块按钮，Enter 或 Space 采纳。

## 限制

- 支持 diff3 样式的 ||||||| 基准段，但不显示它。
- 冲突块不虚拟化；适合单个文件的冲突数量。
