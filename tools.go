package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 调用os.MkdirAll递归创建文件夹
func CreateDir(dir string) error {
	dir = filepath.Dir(dir)
	if !Exists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		return err
	}
	return nil
}

func HttpMethod(uri string, method string, header map[string]interface{}, data string) (http.Header, []byte) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 60 * 5, //超时时间
	}
	payload := strings.NewReader(data)
	req, _ := http.NewRequest(method, uri, payload)
	for k, v := range header {
		req.Header.Add(k, v.(string))
	}
	resp, err := httpClient.Do(req)
	if err == nil {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return resp.Header, body
	}
	//fmt.Println(err)
	return nil, nil
}

func HttpPut(uri string, header map[string]interface{}, data string) []byte {
	_, resp := HttpMethod(uri, "PUT", header, data)
	return resp
}

func HttpPost(uri string, header map[string]interface{}, data string) []byte {
	_, resp := HttpMethod(uri, "POST", header, data)
	return resp
}

func HttpGet(uri string) []byte {
	header := map[string]interface{}{}
	_, resp := HttpMethod(uri, "GET", header, "")
	return resp
}

func Download(url string, file string) bool {
	header := map[string]interface{}{}
	header[`UserAgent`] = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1"
	header[`AcceptLanguage`] = "zh-CN,zh;q=0.9"
	header[`Platform`] = "iPhone"

	_header, body := HttpMethod(url, "GET", header, "")
	if _header != nil && (strings.Contains(_header.Get(`Content-Type`), `image`) == false &&
		strings.Contains(_header.Get(`Content-Type`), `audio`) == false) {
		log.Printf("%s  %s", url, body)
		return false
	}
	//创建文件
	if err := CreateDir(file); err != nil {
		log.Println(err)
		return false
	}
	out, err := os.Create(file)
	if err != nil {
		log.Println(err)
		return false
	}
	// defer延迟调用 关闭文件，释放资源
	defer out.Close()
	//添加缓冲 bufio 是通过缓冲来提高效率。
	wt := bufio.NewWriter(out)
	_, _ = io.Copy(wt, bytes.NewReader(body))
	//将缓存的数据写入到文件中
	_ = wt.Flush()
	return true
}
