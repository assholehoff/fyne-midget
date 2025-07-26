# m(y w)idget(s)

Various widgets for the Fyne GUI toolkit. Made or augmented by me.
```Go
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
```

## `NumericEntry[T Number]`

An entry with two small buttons for adjusting a bound number.

### `NumericEntry[float64]`

Bindable to a `binding.Float`

### `NumericEntry[int]`

Bindable to a `binding.Int`

### `HexColorEntry`

For setting and adjusting a hex color string (not done)
