package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

//csv文件写入
func WriterCSV(path string, data []string) {

	if _, err := os.Stat(path); err != nil {
		log.Println("文件不存在", err)
		//OpenFile读取文件，不存在时则创建，使用追加模式
		F, e := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if e != nil {
			log.Println("创建文件失败", e)
		}
		defer F.Close()

		WC := csv.NewWriter(F)
		headdata := []string{
			"ID", "Type", "SKU", "Name", "Published",
			"Is featured?", "Visibility in catalog", "Short description", "Description", "Date sale price starts",
			"Date sale price ends", "Tax status", "Tax class", "In stock?", "Stock",
			"Backorders allowed?", "Sold individually?", "Weight (lbs)", "Length (in)", "Width (in)",
			"Height (in)", "Allow customer reviews?", "Purchase note", "Sale price", "Regular price",
			"Categories", "Tags", "Shipping class", "Images", "Download limit",
			"Download expiry days", "Parent", "Grouped products", "Upsells", "Cross-sells",
			"External URL", "Button text", "Position", "Attribute 1 name", "Attribute 1 value(s)",
			"Attribute 1 visible", "Attribute 1 global", "Attribute 2 name", "Attribute 2 value(s)", "Attribute 2 visible",
			"Attribute 2 global", "Meta: _wpcom_is_markdown", "Download 1 name", "Download 1 URL", "Download 2 name",
			"Download 2 URL",
		}
		e1 := WC.Write(headdata)
		if e1 != nil {
			log.Println("写入文件失败", e1)
		}
		WC.Flush()
	}
	log.Println("文件存在")
	//使用追加模式打开文件
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！", err)
	}
	defer File.Close()

	//创建写入接口
	WriterCsv := csv.NewWriter(File)
	//data := []string{"chen1", "hai1", "wei1"} //需要写入csv的数据，切片类型

	//写入一条数据，传入数据为切片(追加模式)

	newdata := []string{
		"", "simple", data[0], data[1], "1", //5
		"0", "visible", "", "", "", //10
		"", "taxable", "", "1", "", //15
		"0", "0", "", "", "", //20
		"", "1", "", "", data[2], //25
		"分类", "Tag", "", data[3], "", //30
		"", "", "", "", "", //35
		"", "", "0", "", "", //40
		"1", "1", "", "", "", //45
		"", "1", "", "", "", "", //51
	}

	err1 := WriterCsv.Write(newdata)
	if err1 != nil {
		log.Println("WriterCsv写入文件失败", err1)
	}
	WriterCsv.Flush() //刷新，不刷新是无法写入的
	log.Println("数据写入成功...")

}

func GetNewFille(oldPath string) (newPath string) {
	paths, _ := filepath.Split(oldPath)
	//获取文件名带后缀
	filenameWithSuffix := path.Base(oldPath)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	newPath = paths + "New_" + filenameOnly + ".csv"
	return
}

func main() {

	fmt.Print("拉入要处理的文件:")
	scanner := bufio.NewScanner(os.Stdin)
	oldPath := ""
	for scanner.Scan() {
		oldPath = strings.TrimSpace(scanner.Text()) //去掉字符串前后空格
		break
	}
	newPath := GetNewFille(oldPath)
	f, err := excelize.OpenFile(oldPath)
	if err != nil {
		println(err.Error())
		return
	}
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	for i, row := range rows {
		var data []string
		if i == 0 {
			continue
		}
		sku := row[2]
		name := row[4]
		price := row[5]
		imgs := row[7][1 : len(row[7])-1]
		nimgs := strings.Split(imgs, ",")
		var newimgs []string
		for i2 := 0; i2 < len(nimgs); i2++ {
			s := nimgs[i2]
			s2 := strings.Trim(s, " ")
			s22 := s2[1 : len(s2)-1]
			//println("s22==", s22)
			newimgs = append(newimgs, s22)
		}
		ims := strings.Join(newimgs, ",")
		println("sku==", sku)
		println("name==", name)
		println("价格==", price)
		println("图片==", ims)
		println()
		data = append(data, sku, name, price, ims)
		WriterCSV(newPath, data)
	}
	fmt.Println("处理完毕")
}
