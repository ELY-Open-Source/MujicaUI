package mujica

func init() {
	addMessages(map[string][2]string{
		"range_slider.min":        {"minimum", "最小值"},
		"range_slider.max":        {"maximum", "最大值"},
		"select.placeholder":      {"Select…", "请选择…"},
		"select.empty":            {"No options", "暂无选项"},
		"select.clear":            {"Clear", "清空"},
		"multi_select.search":     {"Search", "搜索"},
		"multi_select.all":        {"Select all", "全选"},
		"multi_select.no_results": {"No results", "无匹配结果"},
		"multi_select.limit":      {"%d of %d chosen", "已选 %d / %d"},
		"multi_select.more":       {"+%d", "+%d"},
		"multi_select.remove":     {"Remove %s", "移除 %s"},
		"combobox.loading":        {"Loading…", "加载中…"},
		"combobox.no_results":     {"No matches", "无匹配项"},
		"combobox.use":            {"Use “%s”", "使用“%s”"},
		"cascader.loading":        {"Loading…", "加载中…"},
		"cascader.failed":         {"Couldn't load", "加载失败"},
		"cascader.retry":          {"Retry", "重试"},
		"cascader.empty":          {"Nothing here", "暂无数据"},
		"tree_select.search":      {"Search", "搜索"},
		"tree_select.no_results":  {"No results", "无匹配结果"},
		"tree_select.expand":      {"Expand %s", "展开 %s"},
		"tree_select.collapse":    {"Collapse %s", "收起 %s"},
	})
}
