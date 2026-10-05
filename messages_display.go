package mujica

func init() {
	addMessages(map[string][2]string{
		"editable_text.saving": {"Saving", "正在保存"},
		"icon.missing":         {"Missing icon: %s", "缺少图标：%s"},
		"image.placeholder":    {"No image", "暂无图片"},
		"image.loading":        {"Loading image", "正在加载图片"},
		"image.error":          {"Image failed to load", "图片加载失败"},
		"image.retry":          {"Retry", "重试"},
		"avatar.unnamed":       {"Unnamed", "未命名"},
		"avatar_group.more":    {"Show all %d people", "显示全部 %d 人"},
		"tag.remove":           {"Remove %s", "移除%s"},
		"link.invalid":         {"This link cannot be opened: %s", "无法打开此链接：%s"},
	})
}
