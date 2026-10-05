# 412 EmojiPicker 表情选择器

## 用途

按分类、搜索与最近使用选择表情，并给支持肤色的表情添加肤色修饰符。

## 最小示例

```go
var recent []string
tone := 0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if r := mujica.EmojiPicker(c, mujica.EmojiPickerOptions{Recent: &recent, Tone: &tone}); r.Picked != "" {
		_ = r.Picked // insert it
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`EmojiPickerOptions.Recent *[]string` 保存最近选择，选择时去重并保留最新 24 个；`Tone *int` 可选，范围 0–5；`Disabled` 默认 false。返回 `EmojiResult` 的命名字段 `Element` 与 `Picked`。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Picked` 是本帧选择的表情；选择时同时更新 Recent，肤色按钮更新 Tone。

## 键盘操作

Tab 访问搜索、肤色、分类、表情；方向键在组内移动；Enter / Space 选择。

## 限制

内置常用表情子集；对支持肤色的条目应用修饰符。
