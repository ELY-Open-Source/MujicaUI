# 308 CommitList 提交列表

## 用途

显示提交列表：标题、引用、作者、短哈希与时间，选择后查看。只构建可见行。

## 最小示例

```go
var state data.ListState[string]
commits := []git.Commit{{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Subject: "Add the commit list"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := git.CommitList(c, &state, commits, git.CommitListOptions{})
	r.Element.Height(240)
	if hash, ok := state.Selected(); ok && r.Submitted() {
		_ = hash // show the commit
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func CommitList(c *ui.Context, s *ListState[string], commits []Commit, opts CommitListOptions) GitListResult`

## 参数

`CommitListOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| TimeFormat | string | "2006-01-02 15:04" | 时间格式。 |
| Status | DataStatus | DataReady | 加载、追加、失败状态。 |
| Error | string | "" | 失败说明。 |

`Commit` 见 CommitGraph。选择保存在 `ListState`，以哈希为键。给返回的 `Element` 一个高度。

## 状态

- 引用以 AccentText 描边小标签显示；短哈希为等宽 AccentText。
- 空列表显示“没有提交”；加载与失败使用列表的统一占位。
- 选中行：Selection 背景与对勾。

## 事件

- `Changed()` 选择变化；`Submitted()` 双击或 Enter。

## 键盘操作

- Tab 进入列表；↑/↓、Home/End 移动选择；输入标题开头字母跳转；Enter 打开。

## 限制

- 时间按调用方给定格式显示，不做相对时间。
