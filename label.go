package midget

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"

	mtheme "github.com/assholehoff/fyne-theme"
)

var (
	_ fyne.Widget            = (*Label)(nil)
	_ fyne.Tappable          = (*Label)(nil)
	_ fyne.DoubleTappable    = (*Label)(nil)
	_ fyne.SecondaryTappable = (*Label)(nil)
)

/* A ridiculously overspecced Label widget with a smaller descriptive text, click, right click and tooltip on hover */
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

	VerticalAlignment VerticalAlignment
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
		SubtextBelowText, AlignCenter,
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
	verticalAlignment VerticalAlignment,
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

		color:    theme.ColorNameForeground,
		subcolor: mtheme.ColorNameDiscrete,
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
func (l *Label) SetText(s string)    { l.Text = s }
func (l *Label) SetSubtext(s string) { l.Subtext = s }
func (l *Label) SetColor(c fyne.ThemeColorName) {
	l.color = c
	l.Refresh()
}
func (l *Label) SetSubColor(c fyne.ThemeColorName) {
	l.subcolor = c
	l.Refresh()
}
func (l *Label) SetSubVisible()         { l.subInvisible = false }
func (l *Label) SetSubInvisible()       { l.subInvisible = true }
func (l *Label) SetScale(n MidgetScale) { l.scale = n }
func (l *Label) SetBottom()             { l.SubPosition = SubtextBelowText }
func (l *Label) SetTop()                { l.SubPosition = SubtextAboveText }
func (l *Label) ToggleSubtext() {
	if l.subInvisible {
		l.SetSubVisible()
		l.Refresh()
	} else {
		l.SetSubInvisible()
		l.Refresh()
	}
}

func (l *Label) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)
	return l.BaseWidget.MinSize()
}

/* CreateRenderer implements fyne.Widget */
func (l *Label) CreateRenderer() fyne.WidgetRenderer {
	// l.textProvider = widget.NewRichTextFromMarkdown(l.Text)
	// l.subProvider = widget.NewRichTextFromMarkdown(l.Subtext)
	l.ExtendBaseWidget(l)

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
	l.OnTapped(event)
}

/* DoubleTapped implements fyne.DoubleTappable. */
func (l *Label) DoubleTapped(event *fyne.PointEvent) {
	l.OnDoubleTapped(event)
}

/* TappedSecondary implements fyne.SecondaryTappable. */
func (l *Label) TappedSecondary(event *fyne.PointEvent) {
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
	l.Refresh()
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
	l.Refresh()
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
	l.Refresh()
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
	// text := r.label.textProvider
	// subtext := r.label.subProvider
	text := r.text
	subtext := r.subtext

	size := s
	size.Width -= theme.InnerPadding()
	size.Height -= theme.InnerPadding()

	primSize := size
	if !r.label.subInvisible {
		primSize.Height -= subtext.MinSize().Height
	}

	pos := fyne.NewSquareOffsetPos(theme.InnerPadding())

	text.Resize(primSize)

	switch r.label.VerticalAlignment {
	case AlignTop:
		pos.Y = theme.InnerPadding()
		if !r.label.subInvisible {
			subSize := subtext.MinSize()
			subSize.Height -= text.Size().Height
			subtext.Resize(subSize)
			if r.label.SubPosition == SubtextAboveText {
				subtext.Move(pos)
				pos.Y += subSize.Height
			} else {
				pos.Y += text.Size().Height
				subtext.Move(pos)
				pos.Y -= text.Size().Height
			}
		}
	case AlignCenter:
		if !r.label.subInvisible {
			subSize := subtext.MinSize()
			subtext.Resize(subSize)
			if r.label.SubPosition == SubtextAboveText {
				pos.Y = (s.Height - text.Size().Height - subSize.Height) / 2
				subtext.Move(pos)
				pos.Y += subSize.Height
			} else {
				pos.Y = (s.Height - text.Size().Height - subSize.Height) / 2
				pos.Y += text.Size().Height
				subtext.Move(pos)
				pos.Y -= text.Size().Height
			}
		} else {
			pos.Y = (s.Height - text.Size().Height) / 2
		}
	case AlignBottom:
		pos.Y = size.Height
		if !r.label.subInvisible {
			subSize := subtext.MinSize()
			subSize.Height -= text.Size().Height
			subtext.Resize(subSize)
			if r.label.SubPosition == SubtextAboveText {
				pos.Y -= (subtext.Size().Height + text.Size().Height)
				subtext.Move(pos)
				pos.Y += subtext.Size().Height
			} else {
				pos.Y -= subtext.Size().Height
				subtext.Move(pos)
				pos.Y -= text.Size().Height
			}
		} else {
			pos.Y -= text.Size().Height
		}
	}

	text.Move(pos)
}

/* MinSize implements fyne.WidgetRenderer. */
func (r *labelRenderer) MinSize() fyne.Size {
	// text := r.label.textProvider
	// subtext := r.label.subProvider
	text := r.text
	subtext := r.subtext

	size := text.MinSize()
	size.Width = fyne.Max(size.Width, subtext.MinSize().Width)

	if !r.label.subInvisible {
		size.Height += subtext.MinSize().Height
	}

	size.Width += theme.InnerPadding()
	size.Height += theme.InnerPadding()

	return size
}

/* Objects implements fyne.WidgetRenderer. */
func (r *labelRenderer) Objects() []fyne.CanvasObject {
	// text := r.label.textProvider
	// subtext := r.label.subProvider
	text := r.text
	subtext := r.subtext
	return []fyne.CanvasObject{text, subtext}
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
	r.text.Color = theme.Color(r.label.color)
	r.subtext.Color = theme.Color(r.label.subcolor)
	r.text.Refresh()
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
func alignedLabelPosition(align VerticalAlignment, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case AlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case AlignTop:
		pos.X = padding.Width / 2
	case AlignBottom:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}
