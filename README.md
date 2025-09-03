# m(y w)idget(s)

Various widgets for the Fyne GUI toolkit. Made or augmented by me.

## Entry

On focus it selects all its text

```Go
type Entry struct {
    xwidget.CompletionEntry
}
```

## MiniButtonPair

A pair of vertically aligned small buttons

```Go
type MiniButtonPair struct {
    ttw.ToolTipWidget
    topButton    *miniButton
    bottomButton *miniButton

    TopText    string
    BottomText string

    OnTappedTop    func()
    OnTappedBottom func()

    disabled bool
    unsquare bool
}
```


## Numeric entries

Entry fields with two buttons (a `MiniButtonPair`) for increment and decrement, also bound to `fyne.KeyUp` and `fyne.KeyDown`

```Go
type NumericEntry struct {
    xwidget.CompletionEntry
    onIncrement func()
    onDecrement func()
}

type IntEntry struct {
    ttw.ToolTipWidget
    value          binding.Int
    valueString    binding.String
    entry          *NumericEntry
    entryMaxWidth  float32
    FormatString   string
    miniButtonPair *MiniButtonPair
    Min, Max, Step int
    validator      func(int) bool
}

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
```