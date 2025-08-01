package types

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

type Key interface{} // FIXME: find out how to get `gocui.Key | rune`

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Handler     func() error
	Key         Key
	Modifier    gocui.Modifier
	Description string
	// DescriptionFunc is used instead of Description if non-nil, and is useful for dynamic
	// descriptions that change depending on context. Important: this must not be an expensive call.
	// Note that you should still provide a generic, non-dynamic description in the Description field,
	// as this is used in the cheatsheet.
	DescriptionFunc func() string
	// If defined, this is used in place of Description when showing the keybinding
	// in the options view at the bottom left of the screen.
	ShortDescription string
	// ShortDescriptionFunc is used instead of ShortDescription if non-nil, and is useful for dynamic
	// descriptions that change depending on context. Important: this must not be an expensive call.
	ShortDescriptionFunc func() string
	Alternative          string
	Tag                  string // e.g. 'navigation'. Used for grouping things in the cheatsheet
	OpensMenu            bool

	// If true, the keybinding will appear at the bottom of the screen.
	// Even if set to true, the keybinding will not be displayed if it is currently
	// disabled. We could instead display it with a strikethrough, but there's
	// limited realestate to show all the keybindings we want, so we're hiding it instead.
	DisplayOnScreen bool
	// if unset, the binding will be displayed in the default color. Only applies to the keybinding
	// on-screen, not in the keybindings menu.
	DisplayStyle *style.TextStyle

	// to be displayed if the keybinding is highlighted from within a menu
	Tooltip string

	// Function to decide whether the command is enabled, and why. If this
	// returns an empty string, it is; if it returns a non-empty string, it is
	// disabled and we show the given text in an error message when trying to
	// invoke it. When left nil, the command is always enabled. Note that this
	// function must not do expensive calls.
	GetDisabledReason func() *DisabledReason
}

func (b *Binding) IsDisabled() bool {
	return b.GetDisabledReason != nil && b.GetDisabledReason() != nil
}

func (b *Binding) GetDescription() string {
	if b.DescriptionFunc != nil {
		return b.DescriptionFunc()
	}
	return b.Description
}

func (b *Binding) GetShortDescription() string {
	if b.ShortDescriptionFunc != nil {
		return b.ShortDescriptionFunc()
	}
	if b.ShortDescription != "" {
		return b.ShortDescription
	}
	return b.GetDescription()
}

// A guard is a decorator which checks something before executing a handler
// and potentially early-exits if some precondition hasn't been met.
type Guard func(func() error) func() error

type KeybindingGuards struct {
	OutsideFilterMode Guard
	NoPopupPanel      Guard
}

type ErrKeybindingNotHandled struct {
	DisabledReason *DisabledReason
}

func (e ErrKeybindingNotHandled) Error() string {
	return e.DisabledReason.Text
}

func (e ErrKeybindingNotHandled) Unwrap() error {
	return gocui.ErrKeybindingNotHandled
}
