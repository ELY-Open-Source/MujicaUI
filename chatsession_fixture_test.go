package mujica

import (
	"time"

	"mujicaui/icons"
)

// chatNow is the fixed clock of the chat session tests.
var chatNow = time.Date(2026, time.October, 5, 15, 30, 0, 0, time.UTC)

func chatConversations() []ChatConversation {
	h := time.Hour
	return []ChatConversation{
		{ID: "a", Title: "Court etiquette notes", Updated: chatNow.Add(-2 * h)},
		{ID: "b", Title: "Ballroom seating plan", Updated: chatNow.Add(-26 * h)},
		{ID: "c", Title: "Harpsichord repair", Updated: chatNow.Add(-4 * 24 * h)},
		{ID: "d", Title: "Coronation menu", Updated: chatNow.Add(-40 * 24 * h)},
		{ID: "e", Title: "Royal letters", Updated: chatNow.Add(-90 * 24 * h), Pinned: true},
		{ID: "f", Title: "Garden maze", Updated: chatNow.Add(-1 * h)},
	}
}

func chatMessages() []SessionMessage {
	return []SessionMessage{
		{Role: SessionUser, Text: "How should guests be seated at the banquet?", Time: chatNow.Add(-10 * time.Minute)},
		{Role: SessionAssistant, Text: "Seat by rank: the host at the head, honored guests to the right.", Time: chatNow.Add(-9 * time.Minute)},
		{Role: SessionUser, Author: "Lady Ann", Text: "And the musicians?", Time: chatNow.Add(-5 * time.Minute)},
	}
}

func chatCards() []CapabilityCard {
	return []CapabilityCard{
		{Title: "Write", Description: "Letters, speeches and verse.", Icon: icons.Must("file-text"), Examples: []string{"Draft a royal decree"}},
		{Title: "Plan", Description: "Banquets, seating and schedules.", Icon: icons.Must("calendar"), Examples: []string{"Plan a ball for 200 guests"}},
		{Title: "Research", Description: "History, heraldry and music.", Icon: icons.Must("search"), Examples: []string{"Explain the minuet"}},
	}
}
