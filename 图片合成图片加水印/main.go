package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path"
)

var (
	bg      = flag.String("bg", "/Users/gaotiansong/Desktop/bj.jpg", "背景图片")
	pt      = flag.String("pt", "/Users/gaotiansong/Desktop/005.png", "前景图片")
	offsetX = flag.Int("offsetX", 0, "x轴偏移值")
	offsetY = flag.Int("offsetY", 0, "y轴偏移值")
	prefix  = flag.String("prefix", "test_", "文件名前缀")
)

func main() {
	flag.Parse()
	mergeImage(*pt)
}

func mergeImage(file string) {
	// 开始处理背景
	ImgB, _ := os.Open(*bg)
	img, _ :=jpeg.Decode(ImgB)
	b1 := img.Bounds()
	//获得背景图像的范围

	//开始处理元素
	wmb, _ := os.Open(file)
	im2, _ := png.Decode(wmb)
	// 设置图片大小
	img2:=resize.Resize(uint(b1.Max.X*5/10), uint(b1.Max.Y*5/10),im2,resize.Lanczos3)
	b2 := img2.Bounds()
	//获得图像的范围
	fmt.Println(b1,b2)
	offset := image.Pt((b1.Max.X-b2.Max.X)/2+*offsetX, (b1.Max.Y-b2.Max.Y)/2+*offsetY)
	//返回Point{X , Y} 也就是一个点的坐标
	b := img.Bounds() //图像范围
	m := image.NewRGBA(b) //新建一个具有b范围的图像RGBA 即 m 和img一样大 背景更好装下背景图
	// m理解为画布比较好 把图片 img 放入 m中 ，他们两个是一样大的
	draw.Draw(m, b, img, image.Point{X: 0, Y: 0}, draw.Src)
	// m 背景  b 是背景图的绘图区域  img 要绘制的图
	//把图片 watermark 放入m这个画布中去
	draw.Draw(m, img2.Bounds().Add(offset), img2, image.Point{X: 0,Y: 0}, draw.Over)
	// b2.Add(offset) b2平移offset

	//保存图片
	ImgW, _ := os.Create(*prefix + path.Base(file))
	err := png.Encode(ImgW, m)
	if err != nil {
		fmt.Println("保存图片失败")
	}
}
