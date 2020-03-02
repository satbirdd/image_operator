package image_operator

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/fogleman/gg"
	// "github.com/golang/freetype"
	// "github.com/golang/freetype/truetype"
)

const (
	// 叠加在上部
	PositionOverlapTop = iota
	// 叠加在下部
	PositionOverlapBottom
	// 文字不叠加在图片上，位置在图片上面
	PositionConcateTop
	// 文字不叠加在图片上，位置在图片下面
	PositionConcateBottom

	defaultMarginTop    = 10.0
	defaultMarginBottom = 10.0
	defaultMarginLeft   = 10.0
	defaultMarginRight  = 10.0
)

var (
	ErrInvalidPosition = errors.New("无效的文字叠加方式")
	white              = color.RGBA{255, 255, 255, 255}
	black              = color.RGBA{0, 0, 0, 255}
)

type Bounds struct {
	X     int // 总宽坐标
	Y     int // 总高坐标
	ImgY  int // 图片开始坐标
	ImgH  int // 图片高度
	TextY int // 文字开始y坐标
	TextH int // 文字高度
}

func (b Bounds) ImgBounds() image.Rectangle {
	return image.Rect(0, b.ImgY, b.X, b.ImgY+b.ImgH)
}

type TextSetting struct {
	// Font     *truetype.Font
	FontPath     string
	FontSize     float64
	LineSpacing  float64
	MarginTop    *float64
	MarginBottom *float64
	MarginLeft   *float64
	MarginRight  *float64
	Dpi          float64
	Position     int
	Color        *color.RGBA
	BgColor      *color.RGBA
}

func (st *TextSetting) Valid() (bool, error) {
	if st.FontPath == "" {
		return false, fmt.Errorf("请指定字体路径")
	}

	if st.FontSize == 0 {
		return false, fmt.Errorf("请指定字体大小")
	}

	if st.LineSpacing <= 0 {
		return false, fmt.Errorf("请指定行间距，行间距不能小于零")
	}

	if st.Dpi == 0 {
		return false, fmt.Errorf("请指定dpi参数")
	}

	if st.Position != PositionOverlapTop &&
		st.Position != PositionOverlapBottom &&
		st.Position != PositionConcateTop &&
		st.Position != PositionConcateBottom {

		return false, ErrInvalidPosition
	}

	if st.Color == nil {
		return false, fmt.Errorf("请指定文字颜色")
	}

	if st.BgColor == nil {
		return false, fmt.Errorf("请指定背景颜色")
	}

	if st.MarginBottom == nil {
		marginB := defaultMarginBottom
		st.MarginBottom = &marginB
	}

	if st.MarginTop == nil {
		marginT := defaultMarginTop
		st.MarginTop = &marginT
	}

	if st.MarginLeft == nil {
		marginL := defaultMarginLeft
		st.MarginLeft = &marginL
	}

	return true, nil
}

func AddTextLinesToImage(path string, textLines []string, setting TextSetting) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return AddTextLines(img, textLines, setting)
}

func AddTextLines(img image.Image, textLines []string, setting TextSetting) (image.Image, error) {
	if _, err := setting.Valid(); err != nil {
		return nil, err
	}

	if len(textLines) == 0 {
		return nil, fmt.Errorf("叠加的文字不能为空")
	}

	dc := gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy())

	if err := dc.LoadFontFace(setting.FontPath, setting.FontSize); err != nil {
		return nil, err
	}

	bounds, textLines, err := calcualteBounds(dc, img, textLines, setting)
	if err != nil {
		return nil, err
	}

	dc = gg.NewContext(bounds.X, bounds.Y)
	if err := dc.LoadFontFace(setting.FontPath, setting.FontSize); err != nil {
		return nil, err
	}

	dc.DrawImage(img, 0, bounds.ImgY)

	err = drawTextLines(dc, textLines, bounds, setting)
	if err != nil {
		return nil, err
	}

	return dc.Image(), nil
}

func calcualteBounds(dc *gg.Context, img image.Image, textLines []string, setting TextSetting) (Bounds, []string, error) {
	var (
		textHeight float64
		newLines   []string
	)

	for _, line := range textLines {
		lines := dc.WordWrap(line, float64(dc.Width()))
		newLines = append(newLines, lines...)
	}

	textHeight = float64(len(newLines)) * dc.FontHeight() * setting.LineSpacing
	textHeight -= (setting.LineSpacing - 1) * dc.FontHeight()
	textHeight = math.Ceil(textHeight)

	if setting.MarginTop != nil {
		textHeight += *setting.MarginTop
	}

	if setting.MarginBottom != nil {
		textHeight += *setting.MarginBottom
	}

	switch setting.Position {
	case PositionOverlapTop:
		return Bounds{
			X:     img.Bounds().Dx(),
			Y:     img.Bounds().Dy(),
			ImgY:  0,
			TextY: 0,
			ImgH:  img.Bounds().Dy(),
			TextH: int(textHeight),
		}, newLines, nil
	case PositionOverlapBottom:
		return Bounds{
			X:     img.Bounds().Dx(),
			Y:     img.Bounds().Dy(),
			ImgY:  0,
			TextY: img.Bounds().Dy() - int(textHeight),
			ImgH:  img.Bounds().Dy(),
			TextH: int(textHeight),
		}, newLines, nil
	case PositionConcateTop:
		return Bounds{
			X:     img.Bounds().Dx(),
			Y:     img.Bounds().Dy() + int(textHeight),
			ImgY:  int(textHeight),
			TextY: 0,
			ImgH:  img.Bounds().Dy(),
			TextH: int(textHeight),
		}, newLines, nil
	case PositionConcateBottom:
		return Bounds{
			X:     img.Bounds().Dx(),
			Y:     img.Bounds().Dy() + int(textHeight),
			ImgY:  0,
			TextY: img.Bounds().Dy(),
			ImgH:  img.Bounds().Dy(),
			TextH: int(textHeight),
		}, newLines, nil
	default:
		return Bounds{}, newLines, ErrInvalidPosition
	}
}

func drawTextLines(dc *gg.Context, textLines []string, bounds Bounds, setting TextSetting) error {
	var (
		fontColor color.RGBA
		bgColor   color.RGBA
	)

	if setting.Color == nil {
		fontColor = white
	} else {
		fontColor = *setting.Color
	}

	if setting.BgColor == nil {
		bgColor = black
	} else {
		bgColor = *setting.BgColor
	}

	rect := image.Rect(0, 0, bounds.X, bounds.TextH)
	bg := image.NewRGBA(rect)
	mask := image.NewAlpha(rect)

	// 读取alpha通道值
	_, _, _, a := bgColor.RGBA()

	for y := 0; y <= bg.Bounds().Dy(); y++ {
		for x := 0; x <= bg.Bounds().Dx(); x++ {
			// background.SetRGBA(x, bounds.TextY+y, bgColor)
			bg.SetRGBA(x, bounds.TextY+y, bgColor)
			mask.SetAlpha(x, y, color.Alpha{uint8(a)})
		}
	}

	dc.SetMask(mask)
	dc.DrawImage(bg, 0, 0)

	r, g, b, a := fontColor.RGBA()
	dc.SetRGBA255(int(r), int(g), int(b), int(a))

	textY := float64(bounds.TextY)
	if setting.MarginTop != nil {
		textY += *setting.MarginTop
	}

	textX := 0.0
	if setting.MarginLeft != nil {
		textX += *setting.MarginLeft
	}

	for _, str := range textLines {
		dc.DrawStringWrapped(str, textX, textY, 0, 0, float64(bounds.X), setting.LineSpacing, gg.AlignLeft)
		textY += dc.FontHeight() * setting.LineSpacing
	}

	return nil
}
