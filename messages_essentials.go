package mujica

func init() {
	addMessages(map[string][2]string{
		"copy_button.copy":   {"Copy", "复制"},
		"copy_button.copied": {"Copied", "已复制"},

		"split_button.more": {"More options", "更多选项"},

		"hold_to_confirm.hint": {"Press and hold, or hold Space, to confirm", "按住鼠标或空格键以确认"},
		"hold_to_confirm.left": {"Keep holding… %.1f s", "继续按住… %.1f 秒"},
		"hold_to_confirm.done": {"Confirmed", "已确认"},

		"bulk_action_bar.selected": {"%d selected", "已选 %d 项"},
		"bulk_action_bar.clear":    {"Clear selection", "清除选择"},

		"color_palette.custom": {"Custom color", "自定义颜色"},

		"qr_code.label":    {"QR code", "二维码"},
		"qr_code.too_long": {"Text too long for a QR code", "文字过长，无法生成二维码"},
	})
}
