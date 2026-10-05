# 240 SearchProgress 搜索进度

## 用途

搜索中关键词与已读来源数；阶段列表；完成后折叠。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	stages := []chat.SearchStage{{Label: "Searching the web", Done: true}, {Label: "Reading sources"}}
	chat.SearchProgress(c, chat.SearchProgressOptions{Query: "Versailles", Read: 3, Total: 8, Stages: stages})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`Query` 必填；`Read`、`Total` 满足 0 ≤ Read ≤ Total；`Stages` 每项提供 `Label` 与 `Done`；`Done=false`，完成时自动折叠一次。

## 状态

搜索中展开，完成后自动折叠；Done 从 true 切换为 false 时展开下一次搜索；支持键盘焦点与祖先禁用。

## 事件

ChoiceResult.Changed() 表示用户展开或折叠。

## 键盘操作

Tab 到折叠标题；Enter、Space 展开或收起。

## 限制

网络搜索由调用方执行，阶段与来源数量由调用方更新。
