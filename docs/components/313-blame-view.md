# 313 BlameView 逐行追溯

## 用途

代码逐行显示最后修改的作者与提交；同一提交的连续行合并显示，悬停显示提交详情。只构建可见行。

## 最小示例

```go
var state git.BlameViewState
lines := []git.BlameLine{
	{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "package main"},
	{Hash: "a1b2c3d4", Author: "Sakiko", When: time.Now(), Summary: "Add main", Text: "func main() {}"},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.BlameView(c, &state, lines, git.BlameViewOptions{Language: "go"}).Height(240)
}
ui.Render(view, 400, 300, 1)
```

签名：`func BlameView(c *ui.Context, s *BlameViewState, lines []BlameLine, opts BlameViewOptions) ui.Element`

## 参数

`BlameLine`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Hash | string | 必填 | 提交哈希；为空时 panic。 |
| Author | string | "" | 作者。 |
| When | time.Time | 零值 | 提交时间。 |
| Summary | string | "" | 提交标题，悬停时显示。 |
| Text | string | "" | 该行代码。 |

`BlameViewOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Language | string | "" | 语法高亮语言，如 go；空为纯文本。 |
| TimeFormat | string | "2006-01-02" | 日期格式。 |

`BlameViewState` 由调用方持有；块注释状态在数据或语言变化时计算一次。给返回元素一个高度。

## 状态

- 每段同提交行只在首行显示短哈希、作者与日期，段之间有细分隔线。
- 左栏悬停为 SurfaceHover，600 ms 后显示提示：完整哈希、作者、日期与标题。
- 代码使用统一的语法颜色，注释为斜体。

## 事件

- 无。

## 键盘操作

- 列表可获得焦点并滚动；输入行首文字跳转到该行。

## 限制

- 长行不换行，超出部分裁剪。
- 悬停详情使用 MyGo 提示，只能用指针触发；每行左栏的名称同样朗读作者与提交。
