package image_operator_test

import (
	"image/color"
	"image/jpeg"
	"os"
	"testing"

	image_op "github.com/satbirdd/image_operator"
)

var (
	utf8FontFile string  = "src/fonts/simsun.ttf"
	imagepath    string  = "src/images/blank.jpg"
	fontsize     float64 = 15
)

func TestOsd(t *testing.T) {
	textLines := []string{
		`设备编号:R1E6015B 图片防伪码:adf2a5e5dge4d4f 抓拍时间:2019-09-11 09:12:13.365`,
		`违法地点:长沙人民东路与万家丽路交汇口 限速信息 小车120km/h 大车100km/h`,
		`实速:137km/h 超速百分比:14% 违法代码:60500 违法行为:违法超速 采集方向:车尾`,
	}

	dstImg, _ := image_op.AddTextLinesToImage(imagepath, textLines, image_op.TextSetting{
		FontPath:    utf8FontFile,
		LineSpacing: 1.4,
		FontSize:    fontsize,
		Dpi:         float64(72),
		Position:    image_op.TextPositionOverlapTop,
		Color:       &color.RGBA{255, 255, 255, 255},
		BgColor:     &color.RGBA{0, 0, 0, 120},
	})

	file := "./src/outs/osd.jpg"
	outFile, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, dstImg, nil)
	if err != nil {
		t.Fatal(err)
	}
}
