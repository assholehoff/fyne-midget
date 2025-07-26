package midget

import (
	// mytheme "UppSpar/internal/ui/theme"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var (
	_ fyne.Widget            = (*Label)(nil)
	_ fyne.Tappable          = (*Label)(nil)
	_ fyne.DoubleTappable    = (*Label)(nil)
	_ fyne.SecondaryTappable = (*Label)(nil)
)

/* A Label widget with a smaller descriptive text, click, right click and tooltip on hover */
type Label struct {
	ttw.ToolTipWidget
	Text       string
	Alignment  fyne.TextAlign
	Wrapping   fyne.TextWrap
	TextStyle  fyne.TextStyle
	Truncation fyne.TextTruncation
	Importance widget.Importance
	SizeName   fyne.ThemeSizeName
	Selectable bool

	VerticalAlignment LabelVerticalAlignment
	scale             MidgetScale
	color             fyne.ThemeColorName
	subcolor          fyne.ThemeColorName
	textBinder        basicBinder
	subBinder         basicBinder
	tipBinder         basicBinder

	Subtext      string
	SubStyle     fyne.TextStyle
	SubAlignment fyne.TextAlign
	SubSizeName  fyne.ThemeSizeName
	SubPosition  SubtextPosition
	SubTopic     TextTopic
	subInvisible bool

	OnTapped          func(*fyne.PointEvent)
	OnDoubleTapped    func(*fyne.PointEvent)
	OnTappedSecondary func(*fyne.PointEvent)
}

func NewLabel(text, subtext, tooltip string) *Label {
	return NewLabelWithStyle(text, subtext, tooltip,
		fyne.TextAlignLeading, fyne.TextAlignLeading,
		fyne.TextStyle{}, fyne.TextStyle{},
		SubtextBelowText, LabelAlignCenter,
	)
}
func NewLabelWithData(text, subtext, tooltip binding.String) *Label {
	label := NewLabel("", "", "")
	label.Bind(text, subtext, tooltip)
	return label
}
func NewLabelWithStyle(
	text, subtext, tooltip string,
	textAlignment, subAlignment fyne.TextAlign,
	textStyle, subStyle fyne.TextStyle,
	subPosition SubtextPosition,
	verticalAlignment LabelVerticalAlignment,
) *Label {
	label := &Label{
		Text:              text,
		Subtext:           subtext,
		Alignment:         textAlignment,
		SubAlignment:      subAlignment,
		TextStyle:         textStyle,
		SubStyle:          subStyle,
		SubPosition:       subPosition,
		VerticalAlignment: verticalAlignment,
	}
	label.ExtendBaseWidget(label)
	label.SetToolTip(tooltip)
	label.OnTapped = func(pe *fyne.PointEvent) {}
	label.OnDoubleTapped = func(pe *fyne.PointEvent) {}
	label.OnTappedSecondary = func(pe *fyne.PointEvent) {}
	if subtext == "" {
		label.SetSubInvisible()
	}
	return label
}
func (l *Label) Bind(text, subtext, tooltip binding.String) {
	l.BindText(text)
	l.BindSubtext(subtext)
	l.BindTooltip(tooltip)
}
func (l *Label) BindText(text binding.String) {
	l.textBinder.SetCallback(l.updateTextFromData)
	l.textBinder.Bind(text)
}
func (l *Label) BindSubtext(subtext binding.String) {
	l.subBinder.SetCallback(l.updateSubtextFromData)
	l.subBinder.Bind(subtext)
}
func (l *Label) BindTooltip(tooltip binding.String) {
	l.tipBinder.SetCallback(l.updateTooltipFromData)
	l.tipBinder.Bind(tooltip)
}
func (l *Label) Unbind() {
	l.UnbindText()
	l.UnbindSubtext()
	l.UnbindTooltip()
}
func (l *Label) UnbindText() {
	l.textBinder.Unbind()
}
func (l *Label) UnbindSubtext() {
	l.subBinder.Unbind()
}
func (l *Label) UnbindTooltip() {
	l.tipBinder.Unbind()
}
func (l *Label) SetText(s string)                  { l.Text = s }
func (l *Label) SetSubtext(s string)               { l.Subtext = s }
func (l *Label) SetColor(c fyne.ThemeColorName)    { l.color = c }
func (l *Label) SetSubColor(c fyne.ThemeColorName) { l.subcolor = c }
func (l *Label) SetSubVisible()                    { l.subInvisible = false }
func (l *Label) SetSubInvisible()                  { l.subInvisible = true }
func (l *Label) SetScale(n MidgetScale)            { l.scale = n }
func (l *Label) SetBottom()                        { l.SubPosition = SubtextBelowText }
func (l *Label) SetTop()                           { l.SubPosition = SubtextAboveText }
func (l *Label) ToggleSubtext() {
	if l.subInvisible {
		l.SetSubVisible()
		l.Refresh()
	} else {
		l.SetSubInvisible()
		l.Refresh()
	}
}

/* CreateRenderer implements fyne.Widget */
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	text := &canvas.Text{
		Alignment: l.Alignment,
		Color:     theme.Color(l.color),
		Text:      l.Text,
		TextSize:  theme.TextSize(),
	}
	subtext := &canvas.Text{
		Alignment: l.SubAlignment,
		Color:     theme.Color(l.subcolor),
		Text:      l.Subtext,
		TextSize:  theme.CaptionTextSize(),
	}

	if l.subInvisible {
		subtext.Hide()
	}

	return &labelRenderer{
		text:    text,
		subtext: subtext,
		label:   l,
	}
}

/* Tapped implements fyne.Tappable. */
func (l *Label) Tapped(event *fyne.PointEvent) {
	log.Printf("Tapped()")
	l.OnTapped(event)
}

/* DoubleTapped implements fyne.DoubleTappable. */
func (l *Label) DoubleTapped(event *fyne.PointEvent) {
	log.Printf("DoubleTapped()")
	l.OnDoubleTapped(event)
}

/* TappedSecondary implements fyne.SecondaryTappable. */
func (l *Label) TappedSecondary(event *fyne.PointEvent) {
	log.Printf("TappedSecondary()")
	l.OnTappedSecondary(event)
}

func (l *Label) updateTextFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.String)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		log.Printf("error getting current data value")
		fyne.LogError("Error getting current data value", err)
		return
	}
	l.SetText(val)
}
func (l *Label) updateSubtextFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.String)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		log.Printf("error getting current data value")
		fyne.LogError("Error getting current data value", err)
		return
	}
	l.SetSubtext(val)
}
func (l *Label) updateTooltipFromData(data binding.DataItem) {
	if data == nil {
		return
	}
	textSource, ok := data.(binding.String)
	if !ok {
		return
	}
	val, err := textSource.Get()
	if err != nil {
		log.Printf("error getting current data value")
		fyne.LogError("Error getting current data value", err)
		return
	}
	l.SetToolTip(val)
}

var _ fyne.WidgetRenderer = (*labelRenderer)(nil)

type labelRenderer struct {
	text, subtext *canvas.Text
	label         *Label
}

/* Destroy implements fyne.WidgetRenderer. */
func (r *labelRenderer) Destroy() {
}

/* Layout implements fyne.WidgetRenderer. */
func (r *labelRenderer) Layout(s fyne.Size) {
	size := s
	size.Width -= theme.InnerPadding()
	size.Height -= theme.InnerPadding()

	psize := size
	if !r.label.subInvisible {
		psize.Height -= r.subtext.MinSize().Height
	}

	pos := fyne.NewSquareOffsetPos(theme.InnerPadding())

	r.text.Resize(psize)

	ssize := fyne.NewSize(0, 0)
	if !r.label.subInvisible {
		ssize = size
		ssize.Height -= r.text.Size().Height
		r.subtext.Resize(ssize)
		if r.label.SubPosition == SubtextAboveText {
			r.subtext.Move(pos)
			pos.Y += r.subtext.Size().Height
		} else {
			spos := pos
			spos.Y += r.text.Size().Height
			r.subtext.Move(spos)
		}
	}

	switch r.label.VerticalAlignment {
	case LabelAlignCenter:
		// pos.Y = ((s.Height - ssize.Height) - r.text.Size().Height) / 2
		if r.label.subInvisible {
			pos.Y = (s.Height - r.text.Size().Height) / 2
		} else {
			if r.label.SubPosition == SubtextAboveText {
				pos.Y = ((s.Height - ssize.Height) - r.text.Size().Height) / 2
				pos.Y += ssize.Height
			} else {
				pos.Y = ((s.Height - ssize.Height) - r.text.Size().Height) / 2
			}
		}
	case LabelAlignTop:
		// fix later
	case LabelAlignBottom:
		// fix later
	}

	r.text.Move(pos)
}

/* MinSize implements fyne.WidgetRenderer. */
func (r *labelRenderer) MinSize() fyne.Size {
	size := r.text.MinSize()
	size.Width = fyne.Max(size.Width, r.subtext.MinSize().Width)
	size.Height += r.subtext.MinSize().Height
	if !r.label.subInvisible {
		size.Height += r.subtext.MinSize().Height
	}
	size.Width += theme.InnerPadding()
	return size
}

/* Objects implements fyne.WidgetRenderer. */
func (r *labelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text, r.subtext}
}

/* Refresh implements fyne.WidgetRenderer. */
func (r *labelRenderer) Refresh() {
	if r.label.subInvisible {
		r.subtext.Hide()
	} else {
		r.subtext.Show()
	}

	r.updateText()
	r.applyTheme()
	r.Layout(r.label.Size())
}
func (r *labelRenderer) applyTheme() {
	fgColorName := r.labelColorName()
	r.text.Color = theme.Color(fgColorName)
	r.text.Refresh()
}
func (r *labelRenderer) labelColorName() fyne.ThemeColorName {
	foreground := theme.ColorNameForeground
	// background = theme.ColorNameBackground
	return foreground
}
func (r *labelRenderer) updateText() {
	r.text.Text = r.label.Text
	r.text.Alignment = r.label.Alignment
	// TODO: text size
	r.text.TextStyle = r.label.TextStyle
	r.text.Refresh()

	r.subtext.Text = r.label.Subtext
	r.subtext.Alignment = r.label.SubAlignment
	// r.subtext.TextSize = r.labelSubTextSize
	r.subtext.TextStyle = r.label.SubStyle
	r.subtext.Refresh()
}
func alignedLabelPosition(align LabelVerticalAlignment, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case LabelAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case LabelAlignTop:
		pos.X = padding.Width / 2
	case LabelAlignBottom:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}
