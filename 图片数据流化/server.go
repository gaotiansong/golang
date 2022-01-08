package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ReadImg 实现File和FileInfo接口的类
type ReadImg struct {
	buf      *bytes.Reader
	fileUrl  string
	fileData []byte
}

// ReadImgData 获取C的图片数据
func ReadImgData(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return pix
}

// Close 实现File和FileInfo接口
func (r *ReadImg) Close() (err error) {
	return nil
}

func (r *ReadImg) Read(p []byte) (n int, err error) {
	return r.buf.Read(p)
}

func (r *ReadImg) Readdir(count int) ([]os.FileInfo, error) {
	var i os.FileInfo = &ReadImg{buf: bytes.NewReader(r.fileData), fileUrl: r.fileUrl, fileData: r.fileData}
	return []os.FileInfo{i}, nil
}

func (r *ReadImg) Seek(offset int64, whence int) (int64, error) {
	return r.buf.Seek(offset, whence)
}

func (r *ReadImg) Stat() (os.FileInfo, error) {
	var i os.FileInfo = &ReadImg{buf: bytes.NewReader(r.fileData), fileUrl: r.fileUrl, fileData: r.fileData}
	return i, nil
}

func (r *ReadImg) Name() string {
	return filepath.Base(r.fileUrl)[:len(filepath.Base(r.fileUrl))-4]
}

func (r *ReadImg) Size() int64 {
	return (int64)(len(r.fileData))
}

func (r *ReadImg) Mode() os.FileMode {
	return os.ModeSetuid
}

func (r *ReadImg) ModTime() time.Time {
	return time.Now()
}

func (r *ReadImg) IsDir() bool {
	return false
}

func (r *ReadImg) Sys() interface{} {
	return nil
}

// HttpDealImg 处理请求
type HttpDealImg struct{}

func (self HttpDealImg) Open(name string) (http.File, error) {
	imgName := name[1:]
	fmt.Println(imgName)
	imgUrl := "https://img.lianzhixiu.com/uploads/allimg/210111/" + name //C(文件服务器地址)
	imgData := ReadImgData(imgUrl)                        //向服务器气球图片数据
	if len(imgData) == 0 {
		fmt.Println("file access forbidden:", name)
		return nil, os.ErrNotExist
	}
	decode, err := imaging.Decode(bytes.NewReader(imgData))
	//图片解析失败
	if err != nil {
		log.Println("图片解析失败")
		return nil, os.ErrNotExist
	}
	log.Println(decode.Bounds(), err)
	bounds := decode.Bounds()
	cropBounds := decode.Bounds()
	cropBounds.Max.Y += rand.Intn(5) - 6
	cropBounds.Max.X += rand.Intn(5) - 6
	bounds.Max.Y += rand.Intn(10) - 8
	bounds.Max.X += rand.Intn(10) - 7
	fmt.Println(bounds.Max.X,bounds.Max.Y)
	fmt.Println(cropBounds.Max.X,cropBounds.Max.Y)
	//先随机缩放大小
	resize := imaging.Resize(decode, bounds.Max.X, bounds.Max.Y, imaging.Lanczos)
	//resize := imaging.Resize(decode, 100, 100, imaging.Lanczos)
	//再随机裁剪一下
	crop := imaging.Crop(resize, cropBounds)
	var ImgByt bytes.Buffer
	writer := bufio.NewWriter(&ImgByt)
	//编码为png图片
	_ = imaging.Encode(writer, crop, imaging.PNG)
	fmt.Println("get img file:", imgUrl)
	var f http.File = &ReadImg{
		buf:      bytes.NewReader(ImgByt.Bytes()),
		fileUrl:  imgName,
		fileData: imgData,
	}
	return f, nil
}

func InitHttpImgFileServ() {
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(HttpDealImg{})))
}

func main() {
	InitHttpImgFileServ()
	http.ListenAndServe(":8000", nil)
}
