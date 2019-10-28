package image_operator

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

// 排列方式为二维整形数组，数组中的元素为图片编号（下标+1）.图片编号不能重复
// 例如：
// [[1, 2],[3, 5, 4]]
// 表示图片按如下顺序排列
// [[1, 2]
//  [3, 5, 4]]
type Permutation [][]int

func (perm Permutation) String() string {
	var strs []string

	for _, outs := range perm {
		var innerStr []string
		for _, inner := range outs {
			innerStr = append(innerStr, fmt.Sprint(inner))
		}

		strs = append(strs, strings.Join(innerStr, ","))
	}

	return "[" + strings.Join(strs, "],[") + "]"
}

func (perm Permutation) IsBegin(index int) bool {
	for _, outs := range perm {
		if len(outs) > 0 && outs[0] == index {
			return true
		}
	}

	return false
}

func (perm Permutation) IsEnd(index int) bool {
	for _, outs := range perm {
		if len(outs) > 0 && outs[len(outs)-1] == index {
			return true
		}
	}

	return false
}

func (perm Permutation) Include(index int) bool {
	for _, outs := range perm {
		for _, inner := range outs {
			if inner == index {
				return true
			}
		}
	}

	return false
}

// 合并几张图片为一张
func CombineImages(paths []string, perm Permutation) (image.Image, error) {
	images, err := getImagesByPath(paths)
	if err != nil {
		return nil, err
	}

	return Combine(images, perm)
}

// 合并几张图片为一张
func Combine(images []image.Image, perm Permutation) (image.Image, error) {
	switch len(images) {
	case 0:
		return nil, fmt.Errorf("用于合并的图片集合是空的")
	case 1:
		return images[0], nil
	default:
		return combine(images, perm)
	}
}

type CombineStatus struct {
	X          int // 下一张图片绘制起点的x值
	UpY        int // 同排下一张图片绘制起点的y值
	DownY      int // 下一排图片绘制起点的y值
	OutIndex   int // 已经绘到第几排
	InnerIndex int // 已经绘到这一排的第几张
}

func (status *CombineStatus) Change(rowBegin, rowEnd bool, x, y int) {
	// if rowBegin {
	status.NextCol(x, y)
	// }

	if rowEnd {
		status.NextRow(x, y)
	}
}

func (status *CombineStatus) NextCol(x, y int) {
	status.InnerIndex += 1
	status.X += x

	if status.DownY == 0 {
		status.DownY = y
	}

	if status.UpY+y < status.DownY {
		status.DownY = status.UpY + y
	}
}

func (status *CombineStatus) NextRow(x, y int) {
	status.OutIndex += 1
	status.X = 0

	status.UpY = status.DownY
	status.DownY += y
}

func (status CombineStatus) Current() (int, int) {
	return status.X, status.UpY
}

func combine(images []image.Image, perm Permutation) (image.Image, error) {
	var (
		status CombineStatus
	)

	w, h := calculateWH(images, perm)
	dc := gg.NewContext(w, h)

	for _, outs := range perm {
		for _, inner := range outs {
			// 如果有超出范围的index,忽略它
			if inner <= len(images) {
				img := images[inner-1]

				err := addImage(dc, inner, img, perm, &status)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return dc.Image(), nil
}

func calculateWH(images []image.Image, perm Permutation) (int, int) {
	var maxWidth, height int

	for _, outs := range perm {
		var (
			maxHeight int
			width     int
		)

		for _, inner := range outs {
			// 如果有超出范围的index,忽略它
			if len(images) >= inner {
				img := images[inner-1]
				width += img.Bounds().Dx()

				h := img.Bounds().Dy()
				if maxHeight < h {
					maxHeight = h
				}
			}
		}

		height += maxHeight
		if maxWidth < width {
			maxWidth = width
		}
	}

	return maxWidth, height
}

func addImage(dc *gg.Context, index int, img image.Image, perm Permutation, status *CombineStatus) error {
	// 未设置的图片直接不绘制
	included := perm.Include(index)
	if !included {
		return nil
	}

	x, y := status.Current()
	dc.DrawImage(img, x, y)

	status.Change(perm.IsBegin(index), perm.IsEnd(index), img.Bounds().Dx(), img.Bounds().Dy())

	return nil
}

func getImagesByPath(paths []string) ([]image.Image, error) {
	var (
		images []image.Image
	)

	for _, path := range paths {
		img, err := readImage(path)
		if err != nil {
			return images, fmt.Errorf("地址 %v 对应的文件不存在或者不是图片", path)
		}

		images = append(images, img)
	}

	return images, nil
}

func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
