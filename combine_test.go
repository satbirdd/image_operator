package image_operator_test

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satbirdd/image_operator"
)

var imgPaths = []string{
	"./src/images/1.jpg",
	"./src/images/2.jpg",
	"./src/images/3.jpg",
	"./src/images/4.jpg",
	"./src/images/5.jpg",
	"./src/images/6.jpg",
	"./src/images/7.jpg",
	"./src/images/8.jpg",
	"./src/images/9.jpg",
}

var inited bool

func combineSetup() {
	if !inited {
		inited = true

		files, _ := filepath.Glob("./src/outs/cl_*")
		for _, f := range files {
			os.Remove(f)
		}

		files, _ = filepath.Glob("./src/outs/cf_*")
		for _, f := range files {
			os.Remove(f)
		}
	}
}

func doCombine(paths []string, perm image_operator.Permutation, indices ...int) error {
	combineSetup()

	img, err := image_operator.CombineImages(paths, perm)
	if err != nil {
		return err
	}

	fileName := perm.String()

	if len(indices) > 0 {
		var strs []string

		for _, index := range indices {
			strs = append(strs, fmt.Sprint(index))
		}

		fileName = "cl_" + fileName + "(" + strings.Join(strs, ",") + ")"
	} else {
		fileName = "cf_" + fileName
	}

	file := fmt.Sprintf("src/outs/%v.jpg", fileName)
	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		return err
	}

	return nil
}

func TestTwoImgsVertical(t *testing.T) {
	paths := imgPaths[:2]
	perm := image_operator.Permutation{
		{1},
		{2},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTwoImgsHorizon(t *testing.T) {
	paths := imgPaths[:2]
	perm := image_operator.Permutation{
		{1, 2},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestThreeImgs01(t *testing.T) {
	paths := imgPaths[:3]
	perm := image_operator.Permutation{
		{1, 2, 3},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestThreeImgs02(t *testing.T) {
	paths := imgPaths[:3]
	perm := image_operator.Permutation{
		{1, 2},
		{3},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestThreeImgs03(t *testing.T) {
	paths := imgPaths[:3]
	perm := image_operator.Permutation{
		{1, 3},
		{2},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestThreeImgs04(t *testing.T) {
	paths := imgPaths[:3]
	perm := image_operator.Permutation{
		{3, 2},
		{1},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFourImgs01(t *testing.T) {
	paths := imgPaths[:4]
	perm := image_operator.Permutation{
		{1, 2},
		{3, 4},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFourImgs02(t *testing.T) {
	paths := imgPaths[:4]
	perm := image_operator.Permutation{
		{1, 3},
		{2, 4},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFiveImgs01(t *testing.T) {
	paths := imgPaths[:5]
	perm := image_operator.Permutation{
		{1, 2, 3},
		{4, 5},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFiveImgs02(t *testing.T) {
	paths := imgPaths[:5]
	perm := image_operator.Permutation{
		{1, 4},
		{2, 5},
		{3},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSixImgs01(t *testing.T) {
	paths := imgPaths[:6]
	perm := image_operator.Permutation{
		{1, 2},
		{3, 4},
		{5, 6},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSixImgs02(t *testing.T) {
	paths := imgPaths[:6]
	perm := image_operator.Permutation{
		{1, 4},
		{2, 5},
		{3, 6},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSixImgs03(t *testing.T) {
	paths := imgPaths[:6]
	perm := image_operator.Permutation{
		{1, 2, 3},
		{4, 5, 6},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSixImgs04(t *testing.T) {
	paths := imgPaths[:6]
	perm := image_operator.Permutation{
		{4, 5, 6},
		{3, 2, 1},
	}

	err := doCombine(paths, perm)
	if err != nil {
		t.Fatal(err)
	}
}

// ---------------- 传入图片数量少于的配置 ------------
func TestSixImgs05(t *testing.T) {
	paths := imgPaths[:3]
	perm := image_operator.Permutation{
		{4, 5, 6},
		{3, 2, 1},
	}

	err := doCombine(paths, perm, 1, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSixImgs06(t *testing.T) {
	paths := imgPaths[:3]
	perm := image_operator.Permutation{
		{1, 4},
		{2, 5},
		{3, 6},
	}

	err := doCombine(paths, perm, 1, 2, 3)
	if err != nil {
		t.Fatal(err)
	}
}
