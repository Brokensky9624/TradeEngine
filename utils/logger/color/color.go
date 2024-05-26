package color

import (
	"fmt"
	"strings"
)

type FgColor uint8
type BgColor uint8

// fg colors.
const FG_None FgColor = 0
const (
	FG_Black FgColor = iota + 30
	FG_Red
	FG_Green
	FG_Yellow
	FG_Blue
	FG_Magenta
	FG_Cyan
	FG_White
)
const (
	FG_HiBlack FgColor = iota + 90
	FG_HiRed
	FG_HiGreen
	FG_HiYellow
	FG_HiBlue
	FG_HiMagenta
	FG_HiCyan
	FG_HiWhite
)

// Bg colors.
const BG_None BgColor = 0
const (
	BG_Black BgColor = iota + 40
	BG_Red
	BG_Green
	BG_Yellow
	BG_Blue
	BG_Magenta
	BG_Cyan
	BG_White
)
const (
	BG_HiBlack BgColor = iota + 100
	BG_HiRed
	BG_HiGreen
	BG_HiYellow
	BG_HiBlue
	BG_HiMagenta
	BG_HiCyan
	BG_HiWhite
)

type rgbColor struct {
	r uint8
	g uint8
	b uint8
}
type rgbFgColor struct {
	fg *rgbColor
}

type rgbBgColor struct {
	bg *rgbColor
}

type ColorText struct {
	text string
	rgbFgColor
	rgbBgColor
	isSfg bool
	isSbg bool
	sfg   FgColor
	sbg   BgColor
}

func New(s string) *ColorText {
	return NewColorText(s)
}

func NewColorText(s string) *ColorText {
	return &ColorText{text: s, sfg: 255, sbg: 255}
}

func (ct *ColorText) SetText(s string) *ColorText {
	ct.text = s
	return ct
}

func (ct *ColorText) SetFg(r, g, b uint8) *ColorText {
	ct.fg = &rgbColor{r, g, b}
	return ct
}

func (ct *ColorText) SetFgObj(fg *rgbColor) *ColorText {
	ct.fg = fg
	return ct
}

func (ct *ColorText) SetSfg(sfg FgColor) *ColorText {
	ct.isSfg = true
	ct.sfg = sfg
	return ct
}

func (ct *ColorText) SetBg(r, g, b uint8) *ColorText {
	ct.bg = &rgbColor{r, g, b}
	return ct
}

func (ct *ColorText) SetBgObj(bg *rgbColor) *ColorText {
	ct.bg = bg
	return ct
}

func (ct *ColorText) SetSbg(sbg BgColor) *ColorText {
	ct.isSbg = true
	ct.sbg = sbg
	return ct
}

func (ct *ColorText) Build() string {
	if ct.bg == nil && ct.fg == nil && ct.sfg == 255 && ct.sbg == 255 {
		return ct.text
	}
	final := []string{}
	if ct.isSfg {
		final = append(final, fmt.Sprintf("\x1b[%dm", ct.sfg))
	} else if ct.fg != nil {
		final = append(final, fmt.Sprintf("\x1b[38;2;%d;%d;%dm", ct.fg.r, ct.fg.g, ct.fg.b))
	}
	if ct.isSfg {
		final = append(final, fmt.Sprintf("\x1b[%dm", ct.sbg))
	} else if ct.bg != nil {
		final = append(final, fmt.Sprintf("\x1b[48;2;%d;%d;%dm", ct.bg.r, ct.bg.g, ct.bg.b))
	}
	final = append(final, ct.text)
	final = append(final, "\x1b[0m")
	return strings.Join(final, "")
}

func GetDefaultBgDebug() *rgbColor {
	return &rgbColor{255, 255, 255}
}

func GetDefaultFgDebug() *rgbColor {
	return &rgbColor{0, 0, 0}
}

func GetDefaultBgInfo() *rgbColor {
	return &rgbColor{71, 110, 245}
}

func GetDefaultFgInfo() *rgbColor {
	return &rgbColor{245, 233, 201}
}

func GetDefaultBgWarning() *rgbColor {
	return &rgbColor{255, 215, 0}
}

func GetDefaultFgWarning() *rgbColor {
	return &rgbColor{84, 78, 168}
}

func GetDefaultBgError() *rgbColor {
	return &rgbColor{173, 55, 55}
}

func GetDefaultFgError() *rgbColor {
	return &rgbColor{212, 205, 205}
}
