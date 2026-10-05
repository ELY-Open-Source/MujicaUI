package mujica

func init() {
	addMessages(map[string][2]string{
		"qr_code.label":    {"QR code", "二维码"},
		"qr_code.too_long": {"Text too long for a QR code", "文字过长，无法生成二维码"},
	})
}
