package main

import (
	"compress/gzip"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

func SpiderData(city, page int, keyWord string) ([]ShopValue, error) {
	var url string
	if page == 1 {
		url = fmt.Sprintf(`http://www.dianping.com/search/keyword/%d/%d_%s`, city, 0, keyWord)
	} else {
		url = fmt.Sprintf(`http://www.dianping.com/search/keyword/%d/%d_%s/p%d`, city, 0, keyWord, page)
	}
	data, err := GetData(url, head)
	if err != nil {
		log.Println(err)
	}
	titleAndAddress, titles, err := ParseTitleAndAddress(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("length:", len(titleAndAddress), "\ntitle address map:", titleAndAddress)
	time.Sleep(time.Second)
	rand.Seed(time.Now().Unix())
	//获取店家的网页源头数据
	var shopData = make(map[string]string)
	var response string
	var n time.Duration
	for k, v := range titleAndAddress {
		response, err = GetData(v, head)
		shopData[k] = response
		n = time.Duration(rand.Intn(5))
		fmt.Println("点击进入商家获取数据成功，睡眠", n, "秒")
		time.Sleep(time.Second * n)
	}
	//对获取到的网页源头数据进行解密处理
	var aimData = make(map[string]ShopValue)
	for title, resource := range shopData {
		entrpyt, err := EntryFront(resource)
		if err != nil {
			log.Println(err)
			continue
		}
		data, err := DecryptFront(resource, entrpyt)
		if err != nil {
			log.Println(err)
			continue
		}
		pares := NewWebBody(data)
		aimData[title] = pares.SpiderInfo()
		n = time.Duration(rand.Intn(4))
		fmt.Println("解析数据成功睡眠", n, "秒")
		time.Sleep(time.Second * n)
	}
	fileData := make([]ShopValue, len(titles))
	for k, v := range titles {
		fileData[k] = aimData[v]
		fileData[k].Name = v
	}
	return fileData, nil
}

func ParseTitleAndAddress(data string) (map[string]string, []string, error) {
	url := regexp.MustCompile(`<a onclick="LXAnalytics\('moduleClick', 'shoppic'\)" target="_blank" href="(.*?)"`)
	urlList := url.FindAllStringSubmatch(data, -1)
	title := regexp.MustCompile(`<img title="(.*?)"`)
	titleList := title.FindAllStringSubmatch(data, -1)
	if len(urlList) != len(titleList) {
		return nil, nil, errors.New("标题和路径对应不上")
	}
	var shopUrl = make(map[string]string)
	var titles = make([]string, len(urlList))
	for i := 0; i < len(urlList); i++ {
		titles[i] = titleList[i][1]
		shopUrl[titleList[i][1]] = urlList[i][1]
	}
	return shopUrl, titles, nil
}

func GetData(url, head string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "构造请求出错")
	}
	if len(head) == 0 {
		return "", errors.New("传入的请求头未空！！！去设置请求头")
	}
	request.Header = ParesHead(head)
	client := &http.Client{
		Timeout: time.Millisecond * 500,
	}
	response, err := client.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "请求出错")
	}
	var reader io.ReadCloser
	if response.StatusCode != http.StatusOK {
		return "", errors.New("Spider not ok")
	}
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return "", err
		}
	} else {
		reader = response.Body
	}
	data, err := ioutil.ReadAll(reader)
	return string(data), err
}
