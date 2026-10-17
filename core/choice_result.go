package core

import (
	"github.com/egoist/mygo/ui"
)

// ChoiceResult is a choice control's element and whether its value changed this frame.
type ChoiceResult struct {
	Element ui.Element
	changed bool
}

// Changed reports whether the user changed the value in this frame.
func (r ChoiceResult) Changed() bool { return r.changed }

// NewChoiceResult reports e and whether its value changed this frame.
func NewChoiceResult(e ui.Element, changed bool) ChoiceResult { return ChoiceResult{e, changed} }
