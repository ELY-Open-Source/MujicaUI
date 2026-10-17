package datetime

import (
	"github.com/egoist/mygo/ui"
)

// DateTimeResult is a date or time component and what happened to it this frame.
type DateTimeResult struct {
	Element ui.Element
	// Input is the text field inside, for labels and focus; nil without one.
	Input    ui.Element
	changed  bool
	canceled bool
}

// Changed reports whether the component wrote a new value this frame.
func (r DateTimeResult) Changed() bool { return r.changed }

// Canceled reports whether the user discarded a pending edit this frame.
func (r DateTimeResult) Canceled() bool { return r.canceled }

// NewDateTimeResult reports a date or time control built from parts in another component.
func NewDateTimeResult(e, input ui.Element, changed, canceled bool) DateTimeResult {
	return DateTimeResult{e, input, changed, canceled}
}
