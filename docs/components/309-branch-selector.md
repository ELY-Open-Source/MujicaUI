# 309 BranchSelector 分支选择器

## 用途

按钮显示当前分支，点击弹出可搜索的本地与远端分支列表；可选“新建分支”入口。

## 最小示例

```go
branch := "main"
branches := []mujica.Branch{{Name: "main", Current: true}, {Name: "origin/main", Remote: true}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.BranchSelector(c, &branch, branches, mujica.BranchSelectorOptions{AllowCreate: true})
	if name, ok := r.Created(); ok {
		_ = name // run git switch -c name
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func BranchSelector(c *ui.Context, selected *string, branches []Branch, opts BranchSelectorOptions) BranchSelectorResult`

## 参数

`Branch`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Name | string | 必填 | 唯一名称，远端分支带远端名，如 origin/main；为空时 panic。 |
| Remote | bool | false | 归入远端组。 |
| Current | bool | false | 当前检出的分支。 |
| Upstream | string | "" | 上游分支。 |
| Ahead, Behind | int | 0 | 领先、落后的提交数；为负 panic。 |

`BranchSelectorOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| AllowCreate | bool | false | 搜索词不是已有分支名时，列表末尾提供“新建分支”。 |
| Disabled | bool | false | 禁用；打开中的弹层会关闭。 |

返回 `BranchSelectorResult`，含 `Element`。过滤结果在搜索词或分支切片变化时才重新计算。

## 状态

- 弹层使用统一浮层外观，160 ms 淡入，尺寸限制在窗口内，内容滚动。
- 本地、远端分组标题为衬线小字；当前分支加粗并带对勾。
- 高亮项：Selection 背景，键盘移动时滚入视野。

## 事件

- `Changed()`：选择了另一个分支，`*selected` 已更新。
- `Created()`：用户要求新建的分支名（已去除首尾空白）。

## 键盘操作

- Enter 或 Space 打开；焦点进入搜索框。
- ↑/↓ 移动高亮并循环，Enter 选择，Escape 关闭并把焦点还给按钮。

## 限制

- 不执行 git：选择与新建都由调用方处理。
- 搜索为大小写不敏感的子串匹配。
