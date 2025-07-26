package midget

import (
	"log"
	"math"
	"reflect"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type Number interface {
	~float64 | ~int | ~int64
}

type NumericEnterface[T Number] interface {
	SetMax(v T)
	SetMin(v T)
	SetStep(v T)
	Bind(di binding.DataItem)
	Unbind()
}

type NumericEntry[T Number] struct {
	ttw.ToolTipWidget
	value           binding.DataItem
	valString       binding.String
	valEntry        *widget.Entry
	valFormatString string
	miniButtonPair  *MiniButtonPair
	min, max, step  T
	validator       func(T) bool
}

func (ne *NumericEntry[T]) SetMin(v T) {
	ne.min = v
}
func (ne *NumericEntry[T]) SetMax(v T) {
	ne.max = v
}
func (ne *NumericEntry[T]) SetStep(v T) {
	ne.step = v
}
func (ne *NumericEntry[T]) Bind(di binding.DataItem) {
	ne.value = di
	switch reflect.TypeOf(di).String() {
	case "*binding.Item[float64]":
		floatBinding := ne.value.(binding.Float)
		ne.valString = binding.FloatToStringWithFormat(floatBinding, ne.valFormatString)
	case "*binding.Item[int]":
		intBinding := ne.value.(binding.Int)
		ne.valString = binding.IntToStringWithFormat(intBinding, ne.valFormatString)
	default:
		log.Printf("cannot bind type %s", reflect.TypeOf(di).String())
	}
	ne.valEntry.Bind(ne.valString)
	ne.validator = func(v T) bool {
		if v >= ne.min && v <= ne.max {
			return true
		}
		return false
	}
}
func (ne *NumericEntry[T]) Unbind() {
	ne.valEntry.Unbind()
}
func (ne *NumericEntry[T]) Validate(v T) bool {
	return ne.validator(v)
}

func NewFloatEntryWithDataX(f binding.Float) *NumericEntry[float64] {
	return NewFloatEntryWithDataAndSpecsX(f, 0, math.MaxFloat64, 1.0)
}
func NewFloatEntryWithDataAndSpecsX(f binding.Float, min, max, step float64) *NumericEntry[float64] {
	ne := &NumericEntry[float64]{
		valFormatString: "%.2f",
		min:             min,
		max:             max,
		step:            step,
	}
	ne.Bind(f)
	ne.validator = func(f float64) bool {
		if f >= ne.min && f <= ne.max {
			return true
		}
		return false
	}
	onIncrement := func() {
		v, _ := ne.value.(binding.Float).Get()
		if ne.Validate(v + ne.step) {
			ne.value.(binding.Float).Set(v + ne.step)
		}
	}
	onDecrement := func() {
		v, _ := ne.value.(binding.Float).Get()
		if ne.Validate(v - ne.step) {
			ne.value.(binding.Float).Set(v - ne.step)
		}
	}
	ne.miniButtonPair = NewMiniButtonPair("▲", "▼", onIncrement, onDecrement)
	return ne
}
func NewIntEntryWithDataX(i binding.Int) *NumericEntry[int] {
	return NewIntEntryWithDataAndSpecsX(i, 0, math.MaxInt, 1)
}
func NewIntEntryWithDataAndSpecsX(i binding.Int, min, max, step int) *NumericEntry[int] {
	ne := &NumericEntry[int]{
		valFormatString: "%d",
		min:             min,
		max:             max,
		step:            step,
	}
	ne.Bind(i)
	ne.validator = func(i int) bool {
		if i >= ne.min && i <= ne.max {
			return true
		}
		return false
	}
	onIncrement := func() {
		v, _ := ne.value.(binding.Int).Get()
		if ne.Validate(v + ne.step) {
			ne.value.(binding.Int).Set(v + ne.step)
		}
	}
	onDecrement := func() {
		v, _ := ne.value.(binding.Int).Get()
		if ne.Validate(v - ne.step) {
			ne.value.(binding.Int).Set(v - ne.step)
		}
	}
	ne.miniButtonPair = NewMiniButtonPair("▲", "▼", onIncrement, onDecrement)
	return ne
}
