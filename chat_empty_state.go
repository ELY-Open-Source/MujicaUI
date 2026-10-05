package mujica

import (
	"cmp"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// ChatEmptyKind is why a chat is empty; zero is ChatEmptyNew.
type ChatEmptyKind int

const (
	// ChatEmptyNew: a new conversation with no messages.
	ChatEmptyNew ChatEmptyKind = iota
	// ChatEmptyNoResults: a conversation search matched nothing.
	ChatEmptyNoResults
	// ChatEmptyOffline: the connection is down.
	ChatEmptyOffline
)

// ChatEmptyStateOptions configure a ChatEmptyState; empty texts use built-in copy.
type ChatEmptyStateOptions struct {
	Kind ChatEmptyKind
	// Query, for ChatEmptyNoResults, is named in the description.
	Query       string
	Title       string
	Description string
	// Action replaces the kind's built-in button text.
	Action string
}

// ChatEmptyState explains an empty chat and offers the next step.
func ChatEmptyState(c *ui.Context, opts ChatEmptyStateOptions) EmptyStateResult {
	switch opts.Kind {
	case ChatEmptyNew:
		return EmptyState(c, EmptyStateOptions{
			Title:       cmp.Or(opts.Title, msg(c, "chat_empty_state.new.title")),
			Description: cmp.Or(opts.Description, msg(c, "chat_empty_state.new.description")),
			Icon:        icons.Must("star"),
			Action:      cmp.Or(opts.Action, msg(c, "chat_empty_state.new.action")),
		})
	case ChatEmptyNoResults:
		return EmptyState(c, EmptyStateOptions{
			Kind: EmptyNoResults, Query: opts.Query, Title: opts.Title, Description: opts.Description,
			Action: cmp.Or(opts.Action, msg(c, "chat_empty_state.no_results.action")),
		})
	case ChatEmptyOffline:
		return EmptyState(c, EmptyStateOptions{
			Title:       cmp.Or(opts.Title, msg(c, "chat_empty_state.offline.title")),
			Description: cmp.Or(opts.Description, msg(c, "chat_empty_state.offline.description")),
			Icon:        icons.Must("wifi-off"),
			Action:      cmp.Or(opts.Action, msg(c, "chat_empty_state.offline.action")),
		})
	}
	panic("mujica: ChatEmptyState Kind out of range")
}
