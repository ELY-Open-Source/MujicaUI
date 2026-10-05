package mujica

func init() {
	addMessages(map[string][2]string{
		"key_value_input.key":    {"Key", "键"},
		"key_value_input.value":  {"Value", "值"},
		"key_value_input.row":    {"Row %d", "第 %d 行"},
		"key_value_input.remove": {"Remove row %d", "删除第 %d 行"},
		"key_value_input.add":    {"Add row", "添加一行"},
	})
}
