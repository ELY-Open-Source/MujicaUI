package diagram_test

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/diagram"
)

func ExampleMermaid() {
	source := "flowchart LR\n  Score[Write score] --> Rehearse[Rehearse]\n  Rehearse --> Stage([Take the stage])"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		diagram.Mermaid(c, source, diagram.MermaidOptions{Label: "Concert plan"})
	}
	ui.Render(view, 480, 200, 1)
}
