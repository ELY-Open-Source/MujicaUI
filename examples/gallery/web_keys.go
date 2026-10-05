package main

import (
	"github.com/egoist/mygo/ui"
)

// keyCodes maps KeyboardEvent.code to MyGo keys.
var keyCodes = map[string]ui.Key{
	"Numpad0": ui.Key0, "Numpad1": ui.Key1, "Numpad2": ui.Key2, "Numpad3": ui.Key3, "Numpad4": ui.Key4,
	"Numpad5": ui.Key5, "Numpad6": ui.Key6, "Numpad7": ui.Key7, "Numpad8": ui.Key8, "Numpad9": ui.Key9,
	"NumpadDecimal": ui.KeyPeriod, "NumpadSubtract": ui.KeyMinus, "NumpadAdd": ui.KeyEqual,
	"NumpadMultiply": ui.Key8, "NumpadDivide": ui.KeySlash, "NumpadEqual": ui.KeyEqual, "NumpadComma": ui.KeyComma,
	"Enter":          ui.KeyEnter,
	"NumpadEnter":    ui.KeyEnter,
	"Escape":         ui.KeyEscape,
	"Backspace":      ui.KeyBackspace,
	"Tab":            ui.KeyTab,
	"Space":          ui.KeySpace,
	"Delete":         ui.KeyDelete,
	"Insert":         ui.KeyInsert,
	"Home":           ui.KeyHome,
	"End":            ui.KeyEnd,
	"PageUp":         ui.KeyPageUp,
	"PageDown":       ui.KeyPageDown,
	"ArrowLeft":      ui.KeyLeft,
	"ArrowRight":     ui.KeyRight,
	"ArrowUp":        ui.KeyUp,
	"ArrowDown":      ui.KeyDown,
	"Minus":          ui.KeyMinus,
	"Equal":          ui.KeyEqual,
	"Comma":          ui.KeyComma,
	"Period":         ui.KeyPeriod,
	"Slash":          ui.KeySlash,
	"Semicolon":      ui.KeySemicolon,
	"Quote":          ui.KeyQuote,
	"BracketLeft":    ui.KeyBracketLeft,
	"BracketRight":   ui.KeyBracketRight,
	"Backslash":      ui.KeyBackslash,
	"Backquote":      ui.KeyBackquote,
	"ContextMenu":    ui.KeyContextMenu,
	"BrowserBack":    ui.KeyBack,
	"BrowserForward": ui.KeyForward,
	"F1":             ui.KeyF1,
	"F2":             ui.KeyF2,
	"F3":             ui.KeyF3,
	"F4":             ui.KeyF4,
	"F5":             ui.KeyF5,
	"F6":             ui.KeyF6,
	"F7":             ui.KeyF7,
	"F8":             ui.KeyF8,
	"F9":             ui.KeyF9,
	"F10":            ui.KeyF10,
	"F11":            ui.KeyF11,
	"F12":            ui.KeyF12,
	"Digit0":         ui.Key0,
	"Digit1":         ui.Key1,
	"Digit2":         ui.Key2,
	"Digit3":         ui.Key3,
	"Digit4":         ui.Key4,
	"Digit5":         ui.Key5,
	"Digit6":         ui.Key6,
	"Digit7":         ui.Key7,
	"Digit8":         ui.Key8,
	"Digit9":         ui.Key9,
	"KeyA":           ui.KeyA,
	"KeyB":           ui.KeyB,
	"KeyC":           ui.KeyC,
	"KeyD":           ui.KeyD,
	"KeyE":           ui.KeyE,
	"KeyF":           ui.KeyF,
	"KeyG":           ui.KeyG,
	"KeyH":           ui.KeyH,
	"KeyI":           ui.KeyI,
	"KeyJ":           ui.KeyJ,
	"KeyK":           ui.KeyK,
	"KeyL":           ui.KeyL,
	"KeyM":           ui.KeyM,
	"KeyN":           ui.KeyN,
	"KeyO":           ui.KeyO,
	"KeyP":           ui.KeyP,
	"KeyQ":           ui.KeyQ,
	"KeyR":           ui.KeyR,
	"KeyS":           ui.KeyS,
	"KeyT":           ui.KeyT,
	"KeyU":           ui.KeyU,
	"KeyV":           ui.KeyV,
	"KeyW":           ui.KeyW,
	"KeyX":           ui.KeyX,
	"KeyY":           ui.KeyY,
	"KeyZ":           ui.KeyZ,
}

// cursorCSS names a MyGo cursor in CSS.
func cursorCSS(c ui.Cursor) string {
	switch c {
	case ui.CursorDefault:
		return "default"
	case ui.CursorPointer:
		return "pointer"
	case ui.CursorText:
		return "text"
	case ui.CursorMove:
		return "move"
	case ui.CursorResizeEW:
		return "ew-resize"
	case ui.CursorResizeNS:
		return "ns-resize"
	case ui.CursorResizeNWSE:
		return "nwse-resize"
	case ui.CursorResizeNESW:
		return "nesw-resize"
	case ui.CursorNotAllowed:
		return "not-allowed"
	case ui.CursorCrosshair:
		return "crosshair"
	case ui.CursorGrab:
		return "grab"
	case ui.CursorGrabbing:
		return "grabbing"
	case ui.CursorResizeN:
		return "n-resize"
	case ui.CursorResizeE:
		return "e-resize"
	case ui.CursorResizeS:
		return "s-resize"
	case ui.CursorResizeW:
		return "w-resize"
	case ui.CursorResizeColumn:
		return "col-resize"
	case ui.CursorResizeRow:
		return "row-resize"
	case ui.CursorVerticalText:
		return "vertical-text"
	case ui.CursorCopy:
		return "copy"
	case ui.CursorAlias:
		return "alias"
	case ui.CursorContextMenu:
		return "context-menu"
	case ui.CursorNone:
		return "none"
	}
	panic("gallery: unmapped cursor")
}
