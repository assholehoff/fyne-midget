package midget

import (
	"fyne.io/fyne/v2"
	xwidget "fyne.io/x/fyne/widget"
)

type NumericEntry struct {
	xwidget.CompletionEntry
	onIncrement func()
	onDecrement func()
}

func NewNumericEntry() *NumericEntry {
	return NewNumericEntryWithCompletions([]string{})
}
func NewNumericEntryWithCompletions(options []string) *NumericEntry {
	e := &NumericEntry{}
	e.Options = options
	e.ExtendBaseWidget(e)
	return e
}

func (e *NumericEntry) FocusGained() {
	e.CompletionEntry.FocusGained()
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
}
func (e *NumericEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyDown:
		e.onDecrement()
	case fyne.KeyUp:
		e.onIncrement()
	default:
		e.CompletionEntry.TypedKey(key)
	}
}
func (e *NumericEntry) SetOnIncrement(f func()) {
	e.onIncrement = f
}
func (e *NumericEntry) SetOnDecrement(f func()) {
	e.onDecrement = f
}
