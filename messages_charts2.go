package mujica

func init() {
	addMessages(map[string][2]string{
		"chart_legend.more":  {"+%d more", "另有 %d 项"},
		"chart_legend.fewer": {"Show fewer", "收起"},
	})
}
