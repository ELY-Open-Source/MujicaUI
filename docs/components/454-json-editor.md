# 454 JsonEditor · JSON 编辑器

## 用途

编辑 JSON、格式化并显示解析错误行列与上下文。

## 最小示例

```go
text := `{"port":8080,"ready":true}`
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := devtools.JsonEditor(c, &text, devtools.JsonEditorOptions{})
	if r.Changed {
		fmt.Println("valid", r.Valid)
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Label 默认 JSON；MinHeight 默认 160，MaxHeight 默认 400 DIP；Disabled 禁用；text 必填。

## 状态

有效 JSON、解析失败、错误行列定位、禁用。

## 事件

Changed 报告文本变化；Valid 为解析结果；Line、Column 为从 1 开始的错误位置。

## 键盘操作

Tab 到格式化和文本框；标准文本编辑与剪贴板快捷键。

## 限制

编辑单个 JSON 值；格式化使用两个空格缩进。
