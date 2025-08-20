# m(y w)idget(s)

Various widgets for the Fyne GUI toolkit. Made or augmented by me.
```Go
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
```

## `NumericEntry[T Number]`

An entry with two small buttons for adjusting a bound number.

### `NumericEntry[float64]`

Bindable to a `binding.Float`

### `NumericEntry[int]`

Bindable to a `binding.Int`

### `HexColorEntry`

For setting and adjusting a hex color string (not done)
