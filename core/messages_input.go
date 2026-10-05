package core

func init() {
	addMessages(map[string][2]string{
		"text_input.clear":      {"Clear", "清空"},
		"number_input.increase": {"Increase", "增加"},
		"number_input.decrease": {"Decrease", "减少"},
		"number_input.invalid":  {"Enter a number", "请输入数字"},
		"number_input.range":    {"Enter a number from %s to %s", "请输入 %s 到 %s 之间的数字"},
		"password_input.show":   {"Show password", "显示密码"},
		"password_input.hide":   {"Hide password", "隐藏密码"},
		"search_input.search":   {"Search", "搜索"},
		"search_input.find":     {"Find", "查找"},
		"search_input.previous": {"Previous match", "上一项"},
		"search_input.next":     {"Next match", "下一项"},
		"search_input.count":    {"%d of %d", "第 %d 项，共 %d 项"},
		"search_input.none":     {"No matches", "无匹配"},
		"checkbox.mixed":        {"Partly selected", "部分选中"},
		"text_area.count":       {"%d of %d characters", "%d / %d 字"},
		"text_area.over":        {"%d characters over the limit", "超出 %d 字"},
	})
}
