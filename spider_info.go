package main

import (
	"github.com/anaskhan96/soup"
)

type WebBody struct {
	body string
}

type ShopValue struct {
	Name    string `名字`
	Address string `地址`
	Phone   string `电话`
}

func NewWebBody(content string) WebBody {
	return WebBody{body: content}
}

func (wb WebBody) SpiderAddress() string {
	htmlInfo := soup.HTMLParse(wb.body)
	address := htmlInfo.Find("body").Find("div", "class", "main")
	return address.FullText()
}

func (wb WebBody) SpiderInfo() ShopValue {
	htmlInfo := soup.HTMLParse(wb.body)
	shopSoup := htmlInfo.Find("body").Find("div", "class", "main")
	address := shopSoup.Find("div", "id", "J_map-show").FullText()
	phone := shopSoup.Find("p").FullText()
	return ShopValue{
		Address: address,
		Phone:   phone,
	}
}
