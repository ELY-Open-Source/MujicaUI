package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// MessageAction is one action on a chat message.
type MessageAction int

const (
	ActionCopy MessageAction = iota + 1
	ActionRegenerate
	ActionEdit
	ActionLike
	ActionDislike
	ActionSpeak
	ActionShare
)

// MessageFeedback is the user's rating of an assistant message.
type MessageFeedback int

const (
	FeedbackNone MessageFeedback = iota
	FeedbackLiked
	FeedbackDisliked
)

// MessageActionsOptions configure MessageActions.
type MessageActionsOptions struct {
	// Feedback holds the rating; nil hides like and dislike.
	Feedback *MessageFeedback
}

// MessageActionsResult is the toolbar's element and the action chosen this frame.
type MessageActionsResult struct {
	Element *ui.Element
	action  MessageAction
}

// Action returns the action the user chose this frame.
func (r MessageActionsResult) Action() (MessageAction, bool) { return r.action, r.action != 0 }

// messageActionIcons names each action's icon and copy key.
var messageActionIcons = map[MessageAction][2]string{
	ActionCopy:       {"copy", "message_actions.copy"},
	ActionRegenerate: {"refresh-cw", "message_actions.regenerate"},
	ActionEdit:       {"pencil", "message_actions.edit"},
	ActionLike:       {"thumbs-up", "message_actions.like"},
	ActionDislike:    {"thumbs-down", "message_actions.dislike"},
	ActionSpeak:      {"volume-2", "message_actions.speak"},
	ActionShare:      {"share-2", "message_actions.share"},
}

// messageActionsFor lists the actions a role's messages offer.
func messageActionsFor(role MessageRole, rated bool) []MessageAction {
	switch role {
	case MessageUser:
		return []MessageAction{ActionCopy, ActionEdit}
	case MessageAssistant:
		if rated {
			return []MessageAction{ActionCopy, ActionRegenerate, ActionLike, ActionDislike, ActionSpeak, ActionShare}
		}
		return []MessageAction{ActionCopy, ActionRegenerate, ActionSpeak, ActionShare}
	case MessageSystem:
		return []MessageAction{ActionCopy}
	case MessageError:
		return []MessageAction{ActionCopy, ActionRegenerate}
	}
	panic("mujica: unknown MessageRole")
}

// MessageActions is a toolbar of the actions a role's message offers; arrows move between them.
func MessageActions(c *ui.Context, role MessageRole, opts MessageActionsOptions) MessageActionsResult {
	var r MessageActionsResult
	r.Element = ui.Row(c).Gap(Density(c).Unit()).AlignItems(ui.Center).FocusGroup(ui.Horizontal).
		Role(ui.RoleToolbar).Label(msg(c, "message_actions.label"))
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		for _, a := range messageActionsFor(role, opts.Feedback != nil) {
			ic := messageActionIcons[a]
			hit := false
			if a == ActionLike || a == ActionDislike {
				want := FeedbackLiked
				if a == ActionDislike {
					want = FeedbackDisliked
				}
				on := *opts.Feedback == want
				if codeBlockToggle(c, ic[0], msg(c, ic[1]), &on) && !disabled {
					hit = true
					if on {
						*opts.Feedback = want
					} else {
						*opts.Feedback = FeedbackNone
					}
				}
			} else {
				hit = iconAction(c, icons.Must(ic[0]), msg(c, ic[1])).Tooltip(msg(c, ic[1])).Clicked() && !disabled
			}
			if hit {
				r.action = a
				slog.Debug("mujica: message_actions action", "action", a)
			}
		}
	})
	return r
}
