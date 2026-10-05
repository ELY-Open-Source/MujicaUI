# 311 TagList 标签列表

## 用途

标签列表：名称、注释与目标提交，每行可跳转到提交。

## 最小示例

```go
var state mujica.ListState[string]
tags := []mujica.GitTag{{Name: "v1.1.0", Target: "a1b2c3d4", Annotation: "Release"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.TagList(c, &state, tags, mujica.TagListOptions{})
	r.Element.Height(200)
	if _, hash, ok := r.Action(); ok {
		_ = hash // show the tagged commit
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func TagList(c *ui.Context, s *ListState[string], tags []GitTag, opts TagListOptions) GitListResult`

## 参数

`GitTag`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Name | string | 必填 | 标签名，唯一。 |
| Target | string | 必填 | 目标提交哈希。 |
| Annotation | string | "" | 注释；轻量标签留空。 |

`TagListOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Status | DataStatus | DataReady | 加载、追加、失败状态。 |
| Error | string | "" | 失败说明。 |

Name 或 Target 为空时 panic。给返回的 `Element` 一个高度。

## 状态

- 注释为 TextMuted 单行，超长省略。

## 事件

- `Action()` 返回 `"jump"` 与目标哈希。
- `Changed()`、`Submitted()` 同 CommitList。

## 键盘操作

- ↑/↓、Home/End 移动选择；“跳转到提交”按钮用 Tab 到达。

## 限制

- 不读取标签：数据由调用方提供。
