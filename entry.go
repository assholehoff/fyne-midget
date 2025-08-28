package midget

import (
	"fyne.io/fyne/v2"
	xwidget "fyne.io/x/fyne/widget"
)

type Entry struct {
	xwidget.CompletionEntry
}

func NewEntry() *Entry {
	return NewEntryWithCompletions([]string{})
}
func NewEntryWithCompletions(options []string) *Entry {
	e := &Entry{}
	e.Options = options
	e.ExtendBaseWidget(e)
	return e
}

func (e *Entry) FocusGained() {
	e.CompletionEntry.FocusGained()
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
}
