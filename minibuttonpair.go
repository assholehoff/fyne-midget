package midget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var (
	_ fyne.Widget = (*MiniButtonPair)(nil)
)

type MiniButtonPair struct {
	ttw.ToolTipWidget
	topButton    *miniButton
	bottomButton *miniButton

	LabelTop    string
	LabelBottom string

	OnTappedTop    func()
	OnTappedBottom func()

	disabled bool
	unsquare bool
}

func NewMiniButtonPair(top, bottom string, ftop, fbot func()) *MiniButtonPair {
	// "▲"
	// "▼"
	i := &MiniButtonPair{
		LabelTop:    top,
		LabelBottom: bottom,

		OnTappedTop:    ftop,
		OnTappedBottom: fbot,

		topButton:    newMiniButton(top, ftop),
		bottomButton: newMiniButton(bottom, fbot),
	}

	i.topButton.OnTapped = func() { i.OnTappedTop() }
	i.bottomButton.OnTapped = func() { i.OnTappedBottom() }

	return i
}
func (i *MiniButtonPair) Square()   { i.unsquare = false }
func (i *MiniButtonPair) Unsquare() { i.unsquare = true }
func (i *MiniButtonPair) Disable() {
	i.topButton.Disable()
	i.bottomButton.Disable()
}
func (i *MiniButtonPair) Enable() {
	i.topButton.Enable()
	i.bottomButton.Enable()
}
func (i *MiniButtonPair) Disabled() bool {
	return i.disabled
}

/* CreateRenderer implements fyne.Widget. */
func (i *MiniButtonPair) CreateRenderer() fyne.WidgetRenderer {
	i.ExtendBaseWidget(i)
	r := &miniButtonPairRenderer{i}
	return r
}

/* MinSize implements fyne.Widget. */
func (i *MiniButtonPair) MinSize() fyne.Size {
	size := i.topButton.MinSize()
	size.Height += 1.0 // space between buttons
	size.Height += i.bottomButton.MinSize().Height
	return size
}

/* Refresh implements fyne.Widget. */
func (i *MiniButtonPair) Refresh() {
	i.topButton.Label = i.LabelTop
	i.topButton.OnTapped = i.OnTappedTop
	i.topButton.Refresh()

	i.bottomButton.Label = i.LabelBottom
	i.bottomButton.OnTapped = i.OnTappedBottom
	i.bottomButton.Refresh()
}

var _ fyne.WidgetRenderer = (*miniButtonPairRenderer)(nil)

type miniButtonPairRenderer struct {
	i *MiniButtonPair
}

/* Destroy implements fyne.WidgetRenderer. */
func (r *miniButtonPairRenderer) Destroy() {
}

/* Layout implements fyne.WidgetRenderer. */
func (r *miniButtonPairRenderer) Layout(s fyne.Size) {
	minSize := r.i.topButton.MinSize()
	minSize.Height += (r.i.bottomButton.MinSize().Height + 1.0)
	minSize.Width += 2.0
	pos := fyne.NewPos(
		(s.Width-minSize.Width)/2,
		(s.Height-minSize.Height)/2,
	)
	size := minSize
	size.Height -= (r.i.bottomButton.MinSize().Height + 1.0)
	if !r.i.unsquare {
		size = fyne.NewSquareSize(fyne.Min(size.Width, size.Height))
	}
	r.i.topButton.Move(pos)
	r.i.topButton.Resize(size)
	pos.Y += r.i.topButton.Size().Height
	pos.Y += 1.0 // space between buttons
	r.i.bottomButton.Move(pos)
	r.i.bottomButton.Resize(size)
}

/* MinSize implements fyne.WidgetRenderer. */
func (r *miniButtonPairRenderer) MinSize() fyne.Size {
	size := fyne.NewSquareSize(
		fyne.Max(
			r.i.topButton.MinSize().Width,
			r.i.topButton.MinSize().Height,
		),
	)
	size.Height += 1.0
	size.Height += fyne.Max(
		r.i.bottomButton.MinSize().Width,
		r.i.bottomButton.MinSize().Height,
	)
	return size
}

/* Objects implements fyne.WidgetRenderer. */
func (r *miniButtonPairRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.i.topButton, r.i.bottomButton}
}

/* Refresh implements fyne.WidgetRenderer. */
func (r *miniButtonPairRenderer) Refresh() {
	r.i.topButton.Refresh()
	r.i.bottomButton.Refresh()
}

var (
	_ desktop.Hoverable = (*miniButton)(nil)
	_ fyne.Disableable  = (*miniButton)(nil)
	_ fyne.Focusable    = (*miniButton)(nil)
	_ fyne.Tappable     = (*miniButton)(nil)
)

type miniButton struct {
	ttw.ToolTipWidget
	// ▲▼
	Icon             fyne.Resource
	Label            string
	Alignment        widget.ButtonAlign
	Importance       widget.Importance
	OnTapped         func()
	disabled         bool
	hovered, focused bool
	unsquare         bool
	tapAnim          *fyne.Animation
}

func newMiniButton(s string, f func()) *miniButton {
	b := &miniButton{
		Label:    s,
		OnTapped: f,
	}
	b.ExtendBaseWidget(b)
	return b
}
func (b *miniButton) Square() {
	b.unsquare = false
}
func (b *miniButton) Unsquare() {
	b.unsquare = true
}
func (b *miniButton) Disabled() bool {
	return b.disabled
}
func (b *miniButton) Disable() {
	b.disabled = true
	b.Refresh()
}
func (b *miniButton) Enable() {
	b.disabled = false
	b.Refresh()
}
func (b *miniButton) FocusGained() {
	b.focused = true
	b.Refresh()
}
func (b *miniButton) FocusLost() {
	b.focused = false
	b.Refresh()
}
func (b *miniButton) TypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeySpace {
		b.Tapped(nil)
	}
}
func (b *miniButton) TypedRune(rune) {
}
func (b *miniButton) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}
func (b *miniButton) MouseMoved(*desktop.MouseEvent) {
}
func (b *miniButton) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *miniButton) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(b.Label, theme.Color(theme.ColorNameForeground))
	text.TextSize = theme.CaptionTextSize()
	bgrect := canvas.NewRectangle(theme.Color(theme.ColorNameButton))
	tprect := canvas.NewRectangle(color.Transparent)
	return &miniButtonRenderer{
		b:          b,
		text:       text,
		background: bgrect,
		tapbg:      tprect,
	}
}

func (b *miniButton) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	size := b.BaseWidget.MinSize()
	size.Width += 2.0
	size.Height += 2.0
	if b.unsquare {
		return size
	}
	square := fyne.NewSquareSize(
		fyne.Max(
			size.Width,
			size.Height,
		),
	)
	return square
}

/* Tapped implements fyne.Tappable. */
func (b *miniButton) Tapped(*fyne.PointEvent) {
	b.OnTapped()
}

var _ fyne.WidgetRenderer = (*miniButtonRenderer)(nil)

type miniButtonRenderer struct {
	b          *miniButton
	text       *canvas.Text
	background *canvas.Rectangle
	tapbg      *canvas.Rectangle
}

func (r *miniButtonRenderer) updateText() {
	r.text.Text = r.b.Label
	r.text.TextSize = theme.Size(theme.SizeNameCaptionText)
	r.text.Refresh()

}
func (r *miniButtonRenderer) applyTheme() {
	th := r.b.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	fgColorName, bgColorName, bgBlendName := r.buttonColorNames()
	if bg := r.background; bg != nil {
		bgColor := color.Color(color.Transparent)
		if bgColorName != "" {
			bgColor = th.Color(bgColorName, v)
		}
		if bgBlendName != "" {
			bgColor = blendColor(bgColor, th.Color(bgBlendName, v))
		}
		bg.FillColor = bgColor
		bg.Refresh()
	}
	r.text.Color = theme.Color(fgColorName)
	r.text.Refresh()
}
func (r *miniButtonRenderer) buttonColorNames() (foreground, background, backgroundBlend fyne.ThemeColorName) {
	foreground = theme.ColorNameForeground
	b := r.b
	if b.Disabled() {
		foreground = theme.ColorNameDisabled
		if b.Importance != widget.LowImportance {
			background = theme.ColorNameDisabledButton
		}
	} else if b.focused {
		backgroundBlend = theme.ColorNameFocus
	} else if b.hovered {
		backgroundBlend = theme.ColorNameHover
	}
	if background == "" {
		switch b.Importance {
		case widget.DangerImportance:
			foreground = theme.ColorNameForegroundOnError
			background = theme.ColorNameError
		case widget.HighImportance:
			foreground = theme.ColorNameForegroundOnPrimary
			background = theme.ColorNamePrimary
		case widget.LowImportance:
			if backgroundBlend != "" {
				background = theme.ColorNameButton
			}
		case widget.SuccessImportance:
			foreground = theme.ColorNameForegroundOnSuccess
			background = theme.ColorNameSuccess
		case widget.WarningImportance:
			foreground = theme.ColorNameForegroundOnWarning
			background = theme.ColorNameWarning
		default:
			background = theme.ColorNameButton
		}
	}
	return
}
func (r *miniButtonRenderer) padding() fyne.Size {
	return fyne.NewSquareSize(2)
}

// Destroy implements fyne.WidgetRenderer.
func (r *miniButtonRenderer) Destroy() {
}

// Layout implements fyne.WidgetRenderer.
func (r *miniButtonRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.tapbg.Resize(size)

	if r.text.Text == "" {
		return
	}

	r.text.Move(alignedButtonPosition(r.b.Alignment, r.padding(), r.text.MinSize(), size))
	r.text.Resize(r.text.MinSize())
}

// MinSize implements fyne.WidgetRenderer.
func (r *miniButtonRenderer) MinSize() fyne.Size {
	size := r.text.MinSize()
	size.Add(r.padding())
	return size
}

// Objects implements fyne.WidgetRenderer.
func (r *miniButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.tapbg, r.text}
}

// Refresh implements fyne.WidgetRenderer.
func (r *miniButtonRenderer) Refresh() {
	r.updateText()
	r.applyTheme()
	r.Layout(r.b.Size())
}

/* from Fyne/v2/widget/button.go */
func alignedButtonPosition(align widget.ButtonAlign, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case widget.ButtonAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case widget.ButtonAlignLeading:
		pos.X = padding.Width / 2
	case widget.ButtonAlignTrailing:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}

// func newButtonTapAnimation(bg *canvas.Rectangle, w fyne.Widget, th fyne.Theme) *fyne.Animation {
// 	v := fyne.CurrentApp().Settings().ThemeVariant()
// 	return fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
// 		mid := w.Size().Width / 2
// 		size := mid * done
// 		bg.Resize(fyne.NewSize(size*2, w.Size().Height))
// 		bg.Move(fyne.NewPos(mid-size, 0))

// 		r, g, bb, a := col.ToNRGBA(th.Color(theme.ColorNamePressed, v))
// 		aa := uint8(a)
// 		fade := aa - uint8(float32(aa)*done)
// 		if fade > 0 {
// 			bg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
// 		} else {
// 			bg.FillColor = color.Transparent
// 		}
// 		canvas.Refresh(bg)
// 	})
// }
