package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

var cityList = []string{"上海", "北京", "重庆", "成都", "广州"}
var cityMap = map[string]int{"上海": 1, "北京": 2, "广州": 4, "重庆": 9, "成都": 8}
var cityNameMap = map[int]string{1: "上海", 2: "北京", 9: "重庆", 8: "成都", 4: "广州"}
var CrawlCity = 1

var pageList = []string{"第一页", "第二页", "第三页", "第四页", "第五页", "第六页", "第七页", "第八页", "第九页", "第十页"}
var pageMap = map[string]int{"第一页": 1, "第二页": 2, "第三页": 3, "第四页": 4, "第五页": 5, "第六页": 6, "第七页": 7, "第八页": 8, "第九页": 9, "第十页": 10}
var CrawlPage = 1

func CreatClient() {
	a := app.New()
	var work fyne.Window

	work = a.NewWindow("点评爬虫工具")

	keyWorkEntry := widget.NewEntry()
	keyWorkEntry.SetPlaceHolder("输入查询关键字")

	cookieEntry := widget.NewMultiLineEntry()
	cookieEntry.SetPlaceHolder("输入请求头")

	pageSelect := widget.NewSelect(pageList, func(s string) {
		CrawlPage = pageMap[s]
	})
	pageSelect.SetSelected("第一页")
	citySelect := widget.NewSelect(cityList, func(s string) {
		CrawlCity = cityMap[s]
	})
	citySelect.SetSelected("上海")

	spiderButton := widget.NewButton("开始爬取", func() {
		log.Println("爬取参数", keyWorkEntry.Text, CrawlPage, CrawlCity)
		fileData, err := SpiderData(CrawlCity, CrawlPage, keyWorkEntry.Text)
		if err != nil {
			log.Println("爬取出错", err)
		}
		path := fmt.Sprintf("Desktop/%s%s.xlsx", keyWorkEntry.Text, cityNameMap[CrawlCity])
		err = DownLoadFile(fileData, path, CrawlPage)
		if err != nil {
			log.Println("数据存储失败", err)
		}
		log.Println("数据爬取成功。")
	})

	pauseButton := widget.NewButton("暂停", func() {

	})
	continueButton := widget.NewButton("继续", func() {

	})

	dropButton := widget.NewButton("删除", func() {

	})

	cookieButton := widget.NewButton("更新请求头", func() {
		head = cookieEntry.Text
	})

	//输入输出窗口和按钮布局
	buttonContain := container.NewGridWithColumns(4, spiderButton, pauseButton, continueButton, dropButton)
	spiderChoice := container.NewGridWithColumns(3, keyWorkEntry, pageSelect, citySelect)
	content := container.NewVBox(spiderChoice, buttonContain, cookieEntry, cookieButton)

	work.Resize(fyne.NewSize(800, 800))
	work.SetContent(content)
	work.ShowAndRun()
}
