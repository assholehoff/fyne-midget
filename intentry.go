package midget

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var _ fyne.Widget = (*IntEntry)(nil)

type IntEntry struct {
	ttw.ToolTipWidget
	value          binding.Int
	valueString    binding.String
	entry          *widget.Entry
	entryMaxWidth  float32
	FormatString   string
	miniButtonPair *MiniButtonPair
	Min, Max, Step int
	validator      func(int) bool
}

func NewIntEntry() *IntEntry {
	return NewIntEntryWithSpecs(math.MinInt, math.MaxInt, 1)
}
func NewIntEntryWithData(i binding.Int) *IntEntry {
	ie := NewIntEntryWithDataAndSpecs(i, math.MinInt, math.MaxInt, 1)
	ie.Bind(i)
	return ie
}
func NewIntEntryWithDataAndSpecs(i binding.Int, min, max, step int) *IntEntry {
	ie := NewIntEntryWithSpecs(0, math.MaxInt, 1)
	ie.Bind(i)
	return ie
}
func NewIntEntryWithSpecs(min, max, step int) *IntEntry {
	ie := &IntEntry{
		value: binding.NewInt(),
	}
	ie.entry = widget.NewEntry()
	ie.entryMaxWidth = 100
	ie.Min = min
	ie.Max = max
	ie.Step = step
	ie.FormatString = "%d"

	ie.validator = func(i int) bool {
		if i >= ie.Min && i <= ie.Max {
			return true
		}
		return false
	}
	onIncrement := func() {
		v, _ := ie.value.Get()
		if ie.Validate(v + ie.Step) {
			ie.value.Set(v + step)
		}
	}
	onDecrement := func() {
		v, _ := ie.value.Get()
		if ie.Validate(v - ie.Step) {
			ie.value.Set(v - step)
		}
	}

	ie.miniButtonPair = NewMiniButtonPair("▲", "▼", onIncrement, onDecrement)
	ie.valueString = binding.IntToStringWithFormat(ie.value, ie.FormatString)
	ie.entry.Bind(ie.valueString)
	return ie
}
func (ie *IntEntry) Bind(bi binding.Int) {
	ie.value = bi
	intBinding := binding.IntToStringWithFormat(ie.value, ie.FormatString)
	ie.entry.Bind(intBinding)
}
func (ie *IntEntry) Unbind() {
	ie.entry.Unbind()
}
func (ie *IntEntry) Validate(v int) bool {
	return ie.validator(v)
}

func (ie *IntEntry) CreateRenderer() fyne.WidgetRenderer {
	ie.ExtendBaseWidget(ie)
	return &intEntryRenderer{ie}
}

var _ fyne.WidgetRenderer = (*intEntryRenderer)(nil)

type intEntryRenderer struct {
	e *IntEntry
}

// Destroy implements fyne.WidgetRenderer.
func (r *intEntryRenderer) Destroy() {
}

// Layout implements fyne.WidgetRenderer.
func (r *intEntryRenderer) Layout(s fyne.Size) {
	entrySize := s
	entrySize.Width -= r.e.miniButtonPair.MinSize().Width
	entrySize.Width = fyne.Min(entrySize.Width, r.e.entryMaxWidth)
	entrySize.Height = r.e.entry.MinSize().Height
	r.e.entry.Resize(entrySize)

	pos := fyne.NewSquareOffsetPos(0)
	pos.Y = (s.Height - r.e.entry.Size().Height) / 2
	r.e.entry.Move(pos)

	pos.X += (r.e.entry.Size().Width + 1)
	pos.Y = (s.Height - r.e.miniButtonPair.MinSize().Height) / 2
	r.e.miniButtonPair.Move(pos)
	r.e.miniButtonPair.Resize(r.e.miniButtonPair.MinSize())
}

// MinSize implements fyne.WidgetRenderer.
func (r *intEntryRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		(r.e.entry.MinSize().Width + r.e.miniButtonPair.MinSize().Width),
		fyne.Max(r.e.entry.MinSize().Height, r.e.miniButtonPair.MinSize().Height),
	)
}

// Objects implements fyne.WidgetRenderer.
func (r *intEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.e.entry, r.e.miniButtonPair}
}

// Refresh implements fyne.WidgetRenderer.
func (r *intEntryRenderer) Refresh() {
	r.e.entry.Refresh()
	r.e.miniButtonPair.Refresh()
}
