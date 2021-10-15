package main

import (
	"github.com/flopp/go-findfont"
	"log"
	"os"
	"strings"
)

var head = ``

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//宋体：Songti.ttc
		if strings.Contains(path, "Songti.ttc") {
			err := os.Setenv("FYNE_FONT", path)
			if err != nil {
				log.Println("设置环境变量出错", err)
			}
			return
		}
	}
	log.Fatal("未查找到中文字体")
}

func main() {
	CreatClient()
}
