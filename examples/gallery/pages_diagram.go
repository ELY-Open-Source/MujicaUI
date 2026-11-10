package main

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/diagram"
	"github.com/ELY-Open-Source/MujicaUI/navigation"
)

var diagramPages = []page{
	{
		ID: "501", Name: "Mermaid", ZhName: "Mermaid 图 流程图 时序图", Category: "Diagrams",
		Demo: mermaidDemo,
		Code: `source := "flowchart LR\n  Score[Write score] --> Rehearse[Rehearse]\n  Rehearse --> Stage([Take the stage])"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	diagram.Mermaid(c, source, diagram.MermaidOptions{Label: "Concert plan"})
}`,
		Keys: "Static picture; assistive technology reads its label. A source Pele rejects shows an alert with its message.",
	},
}

// mermaidSamples are a band's week in Mermaid; the last is broken.
var mermaidSamples = []struct{ name, source string }{
	{"Flowchart", "flowchart LR\n  Score[Write score] --> Rehearse[Rehearse]\n  Rehearse --> Ready{Ready?}\n" +
		"  Ready -- no --> Rehearse\n  Ready -- yes --> Stage([Take the stage])\n  style Stage fill:#7D2034,stroke:#511020,color:#FFF8EE"},
	{"Sequence", "sequenceDiagram\n  participant Sakiko\n  participant Uika\n  participant Hall as Concert hall\n" +
		"  Sakiko->>Uika: Send the new score\n  Uika-->>Sakiko: Vocals recorded\n  Sakiko->>Hall: Book 12 December\n  Hall-->>Sakiko: Confirmed"},
	{"Pie", "pie title Set list by tempo\n  \"Slow\" : 4\n  \"Mid\" : 6\n  \"Fast\" : 3"},
	{"Broken", "flowchart LR\n  Score -->"},
}

func mermaidDemo(c *ui.Context, log func(string)) {
	at := ui.Local(c.Root(), "501", func() int { return 0 })
	ui.Column(c).Gap(16).Width(520).Children(func() {
		segments := make([]navigation.Segment, len(mermaidSamples))
		for i, s := range mermaidSamples {
			segments[i] = navigation.Segment{Label: s.name}
		}
		if navigation.SegmentedControl(c, at, segments, navigation.SegmentedControlOptions{Label: "Sample"}).Changed() {
			log("drew " + mermaidSamples[*at].name)
		}
		diagram.Mermaid(c, mermaidSamples[*at].source, diagram.MermaidOptions{Label: mermaidSamples[*at].name})
	})
}
