package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetNewFille1(oldPath string) (newPath string) {
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
	//pa := "/Users/gaotiansong/gohome/changedata/DAI1118.xlsx"
	//NewPath := GetNewFille1(pa)
	//fmt.Println("NewPath==", NewPath)
	fmt.Print("拉入要处理的文件:")
	scanner := bufio.NewScanner(os.Stdin)
	oldPath := ""
	for scanner.Scan() {
		oldPath = scanner.Text()
		break
	}
	fmt.Println("完成:", oldPath)

}
