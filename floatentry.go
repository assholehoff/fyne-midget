package midget

import (
	"errors"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var _ fyne.Widget = (*FloatEntry)(nil)

type FloatEntry struct {
	ttw.ToolTipWidget
	value          binding.Float
	valueString    binding.String
	entry          *NumericEntry
	entryMaxWidth  float32
	FormatString   string
	miniButtonPair *MiniButtonPair
	Min, Max, Step float64
	validator      func(float64) bool
}

func NewFloatEntry() *FloatEntry {
	min := math.MaxFloat64 * -1
	return NewFloatEntryWithSpecs(min, math.MaxFloat64, 0.1)
}
func NewFloatEntryWithData(i binding.Float) *FloatEntry {
	min := math.MaxFloat64 * -1
	fe := NewFloatEntryWithDataAndSpecs(i, min, math.MaxFloat64, 0.1)
	fe.Bind(i)
	return fe
}
func NewFloatEntryWithDataAndSpecs(i binding.Float, min, max, step float64) *FloatEntry {
	fe := NewFloatEntryWithSpecs(min, max, step)
	fe.Bind(i)
	return fe
}
func NewFloatEntryWithSpecs(min, max, step float64) *FloatEntry {
	fe := &FloatEntry{
		value: binding.NewFloat(),
	}
	fe.ExtendBaseWidget(fe)
	fe.entry = NewNumericEntry()
	fe.entryMaxWidth = 100
	fe.Min = min
	fe.Max = max
	fe.Step = step
	fe.FormatString = "%.2f"

	fe.validator = func(i float64) bool {
		if i >= fe.Min && i <= fe.Max {
			return true
		}
		return false
	}
	onIncrement := func() {
		v, _ := fe.value.Get()
		if fe.Validate(v + fe.Step) {
			fe.value.Set(v + step)
		}
	}
	onDecrement := func() {
		v, _ := fe.value.Get()
		if fe.Validate(v - fe.Step) {
			fe.value.Set(v - step)
		}
	}

	fe.miniButtonPair = NewMiniButtonPair("▲", "▼", onIncrement, onDecrement)
	fe.valueString = binding.FloatToStringWithFormat(fe.value, fe.FormatString)

	fe.entry.onIncrement = onIncrement
	fe.entry.onDecrement = onDecrement

	fe.entry.Bind(fe.valueString)
	fe.entry.Validator = func(s string) error {
		f, _ := strconv.ParseFloat(s, 64)
		if !fe.Validate(f) {
			return errors.New("invalid")
		}
		return nil
	}
	fe.entry.OnChanged = func(s string) {
		f, _ := strconv.ParseFloat(s, 64)
		if f > fe.Max {
			fe.value.Set(fe.Max)
		}
		if f < fe.Min {
			fe.value.Set(fe.Min)
		}
	}
	return fe
}
func (fe *FloatEntry) Bind(bi binding.Float) {
	fe.value = bi
	intBinding := binding.FloatToStringWithFormat(fe.value, fe.FormatString)
	fe.entry.Bind(intBinding)
}
func (fe *FloatEntry) Unbind() {
	fe.entry.Unbind()
}
func (fe *FloatEntry) Validate(v float64) bool {
	return fe.validator(v)
}
func (fe *FloatEntry) Disable() {
	fe.entry.Disable()
	fe.miniButtonPair.Disable()
}
func (fe *FloatEntry) Enable() {
	fe.entry.Enable()
	fe.miniButtonPair.Enable()
}
func (fe *FloatEntry) Increment() {
	fe.miniButtonPair.OnTappedTop()
}
func (fe *FloatEntry) Decrement() {
	fe.miniButtonPair.OnTappedBottom()
}

func (fe *FloatEntry) CreateRenderer() fyne.WidgetRenderer {
	fe.ExtendBaseWidget(fe)
	return &floatEntryRenderer{fe}
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
