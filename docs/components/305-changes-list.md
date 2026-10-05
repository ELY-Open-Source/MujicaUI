# 305 ChangesList 更改列表

## 用途

把工作区更改分为“已暂存”与“更改”两组，逐文件或整组暂存、取消暂存，并选择文件查看差异。只构建可见行。

## 最小示例

```go
var state git.ChangesListState
files := []git.ChangedFile{
	{Path: "main.go", Status: git.GitModified, Staged: true},
	{Path: "notes.txt", Status: git.GitUntracked},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.ChangesList(c, &state, files, git.ChangesListOptions{})
	r.Element.Height(240)
	if action, path, ok := r.Action(); ok {
		_, _ = action, path // run git add or git restore --staged
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ChangesList(c *ui.Context, s *ChangesListState, files []ChangedFile, opts ChangesListOptions) GitListResult`

## 参数

`ChangedFile`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Path | string | 必填 | 文件路径；为空时 panic。 |
| Status | GitStatus | GitModified | 状态，见 GitStatusBadge；未知值 panic。 |
| Staged | bool | false | 放入已暂存组。 |

`ChangesListOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Disabled | bool | false | 禁用：不暂存、不改变选择。 |

`ChangesListState` 由调用方持有，`Selected()` 返回所选文件在 `files` 中的下标，无则 -1。`files` 换成新切片或长度变化时重新分组，分组不在每帧计算。给返回的 `Element` 一个高度。

## 状态

- 组标题为衬线小字与计数，空组的整组按钮禁用。
- 每行显示 GitStatusBadge（字母、边框与颜色）、文件名与目录。
- 选中行：Selection 背景加左侧 AccentText 竖线。
- 悬停与按下：图标按钮 SurfaceHover / Selection 背景。

## 事件

- `Action()` 返回 `"stage"`、`"unstage"`（附路径）或 `"stage_all"`、`"unstage_all"`。
- `Changed()` 选择变化；`Submitted()` 双击或 Enter。

## 键盘操作

- Tab 进入列表；↑/↓、Home/End 移动选择，跳过组标题；输入文件名开头字母跳转；Enter 提交。
- 行内暂存按钮可用 Tab 到达，Enter 或 Space 执行。

## 限制

- 不读取仓库：文件列表与暂存操作都由调用方完成。
- 重命名的旧路径不单独显示。
