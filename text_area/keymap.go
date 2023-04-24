package textarea

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the key bindings for different actions within the textarea.
type KeyMap struct {
	// Insert mode bindings
	CharacterBackward       key.Binding
	CharacterForward        key.Binding
	DeleteAfterCursor       key.Binding
	DeleteBeforeCursor      key.Binding
	DeleteCharacterBackward key.Binding
	DeleteCharacterForward  key.Binding
	DeleteWordBackward      key.Binding
	DeleteWordForward       key.Binding
	InsertNewline           key.Binding
	LineEnd                 key.Binding
	LineNext                key.Binding
	LinePrevious            key.Binding
	LineStart               key.Binding
	Paste                   key.Binding
	WordBackward            key.Binding
	WordForward             key.Binding
	InputBegin              key.Binding
	InputEnd                key.Binding
}

// DefaultKeyMap is the default set of key bindings for navigating and acting
// upon the textarea.
var DefaultKeyMap = KeyMap{
	CharacterForward:        key.NewBinding(key.WithKeys("right", "ctrl+f")),
	CharacterBackward:       key.NewBinding(key.WithKeys("left", "ctrl+b")),
	WordForward:             key.NewBinding(key.WithKeys("alt+right", "alt+f")),
	WordBackward:            key.NewBinding(key.WithKeys("alt+left", "alt+b")),
	LineNext:                key.NewBinding(key.WithKeys("down", "ctrl+n")),
	LinePrevious:            key.NewBinding(key.WithKeys("up", "ctrl+p")),
	DeleteWordBackward:      key.NewBinding(key.WithKeys("alt+backspace", "ctrl+w")),
	DeleteWordForward:       key.NewBinding(key.WithKeys("alt+delete", "alt+d")),
	DeleteAfterCursor:       key.NewBinding(key.WithKeys("ctrl+k")),
	DeleteBeforeCursor:      key.NewBinding(key.WithKeys("ctrl+u")),
	InsertNewline:           key.NewBinding(key.WithKeys("enter", "ctrl+m")),
	DeleteCharacterBackward: key.NewBinding(key.WithKeys("backspace", "ctrl+h")),
	DeleteCharacterForward:  key.NewBinding(key.WithKeys("delete", "ctrl+d")),
	LineStart:               key.NewBinding(key.WithKeys("home", "ctrl+a")),
	LineEnd:                 key.NewBinding(key.WithKeys("end", "ctrl+e")),
	Paste:                   key.NewBinding(key.WithKeys("ctrl+v")),
	InputBegin:              key.NewBinding(key.WithKeys("alt+<", "ctrl+home")),
	InputEnd:                key.NewBinding(key.WithKeys("alt+>", "ctrl+end")),
}
