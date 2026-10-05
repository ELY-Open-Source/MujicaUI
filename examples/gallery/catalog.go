package main

import "github.com/egoist/mygo/ui"

// page is one component's gallery entry.
type page struct {
	ID       string
	Name     string
	ZhName   string
	Category string
	// Demo builds the interactive example; log records events.
	Demo func(c *ui.Context, log func(string))
	Code string
	Keys string
}

func catalog() []page {
	var all []page
	for _, g := range [][]page{layoutPages, displayPages, inputPages, choicePages, formPages, datetimePages, navigationPages, dataPages, overlayPages, feedbackPages,
		charts1Pages, charts2Pages, mktchartPages, tradingPages, marketsPages, chatmsgPages, chatinputPages, chatsessionPages, agentPages, motionPages, gitPages, codePages, mediaPages, projectsPages, filesPages, messagingPages, calendarPages, devtoolsPages, accountPages, essentialsPages} {
		all = append(all, g...)
	}
	return all
}
