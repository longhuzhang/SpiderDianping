package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var woffType = []string{"address", "shopNum", "tagName", "reviewTag", "num", "dishname", "shopdesc",
	"review", "hours"}

//校验所需要的破解密钥是否都存在
func EntryFront(resource string) (map[string]string, error) {
	fontBaseUrl := regexp.MustCompile(`href="(//s3plus.meituan.net/v1/.*?)">`)
	urlList := fontBaseUrl.FindAllStringSubmatch(resource, -1)
	if len(urlList) < 1 {
		return nil, errors.New(fmt.Sprintf("Did not match font url. resource is:%s", resource))
	}
	aimUrl := "https:" + urlList[0][1]
	response, err := http.Get(aimUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	fontData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	woffUrls := regexp.MustCompile(`,url\("(.*?\.woff"\).*?\{)`)
	woffList := woffUrls.FindAllStringSubmatch(string(fontData), -1)

	var entryMap = make(map[string]string)
	for _, v := range woffList {
		catch := catchType(v[1])
		if catch == "" {
			continue
		}
		addressCompile := regexp.MustCompile(`(//.*?woff)`)
		addressFont := addressCompile.FindAllStringSubmatch(v[1], -1)
		addressUrl := "https:" + addressFont[0][1]
		length := len(addressUrl)
		name := addressUrl[length-13 : length-5]
		addressJsonName := "tmp/" + name + ".json"
		entryMap[catch] = addressJsonName
		if FileStore[addressJsonName] {
			continue
		}
		return nil, errors.New(fmt.Sprintf("the woff json name(%s) is not exit", addressJsonName))
	}

	return entryMap, nil
}

//解密加密的点评网页
func DecryptFront(text string, decryptKey map[string]string) (string, error) {
	for k, v := range decryptKey {
		decrypt, err := jsonFile.ReadFile(v)
		if err != nil {
			return "", err
		}
		var content map[string]string
		err = json.Unmarshal(decrypt, &content)
		if err != nil {
			return "", err
		}
		for jsonKey, jsonValue := range content {
			key := strings.Replace(jsonKey, "uni", "&#x", -1)
			key = `"` + k + `">` + key + `;`
			value := `"` + k + `">` + jsonValue
			text = strings.Replace(text, key, value, -1)
		}
	}
	return text, nil
}

func catchType(woff string) string {
	for _, v := range woffType {
		if strings.Contains(woff, v) {
			return v
		}
	}
	return ""
}

func ParesHead(str string) http.Header {
	headSlice := strings.Split(str, "\n")
	head := http.Header{}
	var sliceStr []string
	for _, v := range headSlice {
		sliceStr = strings.Split(v, ": ")
		if len(sliceStr) < 2 {
			continue
		}
		sliceStr[0] = strings.Replace(sliceStr[0], " ", "", -1)
		sliceStr[1] = strings.Replace(sliceStr[1], " ", "", -1)
		head.Set(sliceStr[0], sliceStr[1])
	}
	return head
}

func DownLoadFile(shopList []ShopValue, path string, page int) error {
	var f *excelize.File
	var err error
	if FileExit(path) {
		f, err = excelize.OpenFile(path)
		if err != nil {
			return err
		}
	} else {
		f = excelize.NewFile()
	}

	var i int
	if err = f.SetCellValue("Sheet1", "A1", "序号"); err != nil {
		return err
	}
	if err = f.SetCellValue("Sheet1", "B1", "店名"); err != nil {
		return err
	}
	if err = f.SetCellValue("Sheet1", "C1", "地址"); err != nil {
		return err
	}
	if err = f.SetCellValue("Sheet1", "D1", "电话"); err != nil {
		return err
	}
	for k, v := range shopList {
		i = k + 2 + (page-1)*15
		if err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), i-1); err != nil {
			return err
		}
		if err = f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), v.Name); err != nil {
			return err
		}
		if err = f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), v.Address); err != nil {
			return err
		}
		if err = f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), v.Phone); err != nil {
			return err
		}
	}
	// Set active sheet of the workbook.
	// Save spreadsheet by the given path.
	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}

func ExportDataToExcel(data interface{}, path string, page int) error {
	var f *excelize.File
	var err error
	if FileExit(path) {
		f, err = excelize.OpenFile(path)
		if err != nil {
			return err
		}
	} else {
		f = excelize.NewFile()
	}

	//获取字段名和备注名字
	point := reflect.TypeOf(data).Elem()
	var fields = make([]string, 0)
	if page == 1 {
		if err = f.SetCellValue("Sheet1", "A1", "序号"); err != nil {
			return err
		}
	}

	for i := 0; i < point.NumField(); i++ {
		elem := point.FieldByIndex([]int{i})
		fields = append(fields, elem.Name)
		if page > 1 {
			continue
		}
		if err = f.SetCellValue("Sheet1", string(byte('A'+i+1))+strconv.Itoa(1), elem.Tag); err != nil {
			return err
		}
	}

	//fmt.Println(fields)
	//数据内容写入到excel表格中
	value := reflect.Indirect(reflect.ValueOf(data))
	var index int
	for i := 0; i < value.Len(); i++ {
		index = i + 2 + (page-1)*25
		if err = f.SetCellValue("Sheet1", "A"+strconv.Itoa(index), strconv.Itoa(index-1)); err != nil {
			return err
		}
		for k, _ := range fields {
			if err = f.SetCellValue("Sheet1", string(byte(('A'+k+1)))+strconv.Itoa(index), value.Index(i).FieldByName(fields[k])); err != nil {
				return err
			}
			//fmt.Println("index",string(byte(('A'+k+1)))+strconv.Itoa(i+2),"value:",value.Index(i).FieldByName(fields[k]))
		}
	}
	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}

func FileExit(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
