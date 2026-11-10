package core

func init() {
	addMessages(map[string][2]string{
		"mermaid.label":  {"Diagram", "图表"},
		"mermaid.failed": {"Diagram failed", "图表出错"},
	})
}
