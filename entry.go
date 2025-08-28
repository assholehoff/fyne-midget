package midget

import (
	"fyne.io/fyne/v2"
	xwidget "fyne.io/x/fyne/widget"
)

type Entry struct {
	xwidget.CompletionEntry
}

func (e *Entry) FocusGained() {
	e.CompletionEntry.FocusGained()
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
}
