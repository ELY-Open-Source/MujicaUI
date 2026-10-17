# 318 GitStatusBadge Git 状态标记

## 用途

以字母 M、A、D、U、R、C 加描边与颜色标记文件状态。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.GitStatusBadge(c, git.GitModified, git.GitStatusBadgeOptions{ShowLabel: true})
}
ui.Render(view, 400, 300, 1)
```

签名：`func GitStatusBadge(c *ui.Context, status GitStatus, opts GitStatusBadgeOptions) ui.Element`

## 参数

`GitStatusBadgeOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| ShowLabel | bool | false | 在字母旁写出状态名称；否则作为提示。 |

`GitStatus`：GitModified、GitAdded、GitDeleted、GitUntracked、GitRenamed、GitCopied；`Letter()` 返回字母。未知值 panic。

## 状态

- 字母为等宽粗体，颜色依状态：修改 Warning、添加 Success、删除 Danger、未跟踪与重命名 Info、复制 AccentText。
- 名称写入朗读文本与提示，不只靠颜色。

## 事件

- 无。

## 键盘操作

- 无：纯展示。

## 限制

- 无。
