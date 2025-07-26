package midget

import "image/color"

/* package midget contains m(y w)idget(s) */

type ButtonIconPlacement int

const (
	ButtonIconAboveText ButtonIconPlacement = iota
	ButtonIconBelowText
	ButtonIconLeadingText
	ButtonIconTrailingText
)

type LabelVerticalAlignment int

const (
	LabelAlignTop LabelVerticalAlignment = iota
	LabelAlignCenter
	LabelAlignBottom
)

type SubtextPosition int

const (
	SubtextAboveText SubtextPosition = iota
	SubtextBelowText
)

type TextTopic int

const (
	NonspecificTopic TextTopic = iota
	DetailsTopic
	HelpTopic
)

type MidgetScale int

const (
	Large MidgetScale = iota
	Regular
	Small
)

func blendColor(under, over color.Color) color.Color {
	// This alpha blends with the over operator, and accounts for RGBA() returning alpha-premultiplied values
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()

	srcAlpha := float32(srcA) / 0xFFFF
	dstAlpha := float32(dstA) / 0xFFFF

	outAlpha := srcAlpha + dstAlpha*(1-srcAlpha)
	outR := srcR + uint32(float32(dstR)*(1-srcAlpha))
	outG := srcG + uint32(float32(dstG)*(1-srcAlpha))
	outB := srcB + uint32(float32(dstB)*(1-srcAlpha))
	// We create an RGBA64 here because the color components are already alpha-premultiplied 16-bit values (they're just stored in uint32s).
	return color.RGBA64{R: uint16(outR), G: uint16(outG), B: uint16(outB), A: uint16(outAlpha * 0xFFFF)}
}
