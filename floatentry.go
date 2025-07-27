package midget

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var _ fyne.Widget = (*FloatEntry)(nil)

type FloatEntry struct {
	ttw.ToolTipWidget
	value          binding.Float
	valueString    binding.String
	entry          *widget.Entry
	entryMaxWidth  float32
	FormatString   string
	miniButtonPair *MiniButtonPair
	Min, Max, Step float64
	validator      func(float64) bool
}

func NewFloatEntry() *FloatEntry {
	return NewFloatEntryWithSpecs(0, math.MaxFloat64, 0.1)
}
func NewFloatEntryWithData(i binding.Float) *FloatEntry {
	ie := NewFloatEntryWithDataAndSpecs(i, 0, math.MaxFloat64, 0.1)
	ie.Bind(i)
	return ie
}
func NewFloatEntryWithDataAndSpecs(i binding.Float, min, max, step float64) *FloatEntry {
	ie := NewFloatEntryWithSpecs(0, math.MaxFloat64, 1)
	ie.Bind(i)
	return ie
}
func NewFloatEntryWithSpecs(min, max, step float64) *FloatEntry {
	ie := &FloatEntry{
		value: binding.NewFloat(),
	}
	ie.entry = widget.NewEntry()
	ie.entryMaxWidth = 100
	ie.Min = min
	ie.Max = max
	ie.Step = step
	ie.FormatString = "%.2f"

	ie.validator = func(i float64) bool {
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
	ie.valueString = binding.FloatToStringWithFormat(ie.value, ie.FormatString)
	ie.entry.Bind(ie.valueString)
	return ie
}
func (ie *FloatEntry) Bind(bi binding.Float) {
	ie.value = bi
	intBinding := binding.FloatToStringWithFormat(ie.value, ie.FormatString)
	ie.entry.Bind(intBinding)
}
func (ie *FloatEntry) Unbind() {
	ie.entry.Unbind()
}
func (ie *FloatEntry) Validate(v float64) bool {
	return ie.validator(v)
}

func (ie *FloatEntry) CreateRenderer() fyne.WidgetRenderer {
	ie.ExtendBaseWidget(ie)
	return &floatEntryRenderer{ie}
}

var _ fyne.WidgetRenderer = (*floatEntryRenderer)(nil)

type floatEntryRenderer struct {
	e *FloatEntry
}

// Destroy implements fyne.WidgetRenderer.
func (r *floatEntryRenderer) Destroy() {
}

// Layout implements fyne.WidgetRenderer.
func (r *floatEntryRenderer) Layout(s fyne.Size) {
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
func (r *floatEntryRenderer) MinSize() fyne.Size {
	return fyne.NewSize(
		(r.e.entry.MinSize().Width + r.e.miniButtonPair.MinSize().Width),
		fyne.Max(r.e.entry.MinSize().Height, r.e.miniButtonPair.MinSize().Height),
	)
}

// Objects implements fyne.WidgetRenderer.
func (r *floatEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.e.entry, r.e.miniButtonPair}
}

// Refresh implements fyne.WidgetRenderer.
func (r *floatEntryRenderer) Refresh() {
	r.e.entry.Refresh()
	r.e.miniButtonPair.Refresh()
}
