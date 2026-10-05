# 315 PullRequestCard 拉取请求卡片

## 用途

概览一个拉取请求：标题、作者、状态、检查、评审人与差异统计。

## 最小示例

```go
pr := mujica.PullRequest{Number: 42, Title: "Add git components", Author: "Sakiko", Head: "feature", Base: "main",
	ChecksPassed: 3, Reviewers: []mujica.Reviewer{{Name: "Uika", State: mujica.ReviewApproved}}, Added: 120, Deleted: 8, Files: 6}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.PullRequestCard(c, pr, mujica.PullRequestCardOptions{}).Changed() {
		// open the pull request
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func PullRequestCard(c *ui.Context, pr PullRequest, opts PullRequestCardOptions) ChoiceResult`

## 参数

`PullRequest`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Number | int | 必填 | 编号，须为正。 |
| Title | string | 必填 | 标题。 |
| Author | string | "" | 作者。 |
| State | PullRequestState | PullRequestOpen | Open、Draft、Merged、Closed；未知值 panic。 |
| Head, Base | string | "" | 源分支与目标分支，都有值时显示。 |
| ChecksPassed, ChecksFailed, ChecksPending | int | 0 | 检查计数；为负 panic。 |
| Reviewers | []Reviewer | nil | 评审人与状态（Pending、Approved、ChangesRequested）。 |
| Added, Deleted, Files | int | 0 | 差异统计。 |

`PullRequestCardOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Disabled | bool | false | 禁用打开。 |

返回 `ChoiceResult`。

## 状态

- 状态徽章同时使用图标、文字与描边颜色。
- 检查以图标加数字显示，朗读完整说明；为 0 的类别不显示。
- 评审人以头像加角标图标表示状态，名称与状态写入提示和朗读文本。

## 事件

- `Changed()`：点击标题或在标题上按 Enter，表示打开。

## 键盘操作

- Tab 到标题链接，Enter 或 Space 打开。

## 限制

- 数据由调用方获取；组件不刷新检查状态。
