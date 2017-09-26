package main

import (
	"log"
	"net/http"
	"util"

	"os"
	"process"
	"strings"
)

func main() {
	//接受URL
	//校验URL
	//缓存代理
	//输出结果
	//ScanThread("HEAD")
	//CacheUrl()
	//println(CheckUrl("www/s/3?a=1#b=2"))
	//println(CheckUrl("www.me/s/3?a=1#b=2"))
	//println(CheckUrl("changwei.me/s/3?a=1#b=2"))
	//println(CheckUrl("..me/s/3?a=1#b=2"))
	//println(CheckUrl("http://c.me/s/3?a=1#b=2"))
	//println(CheckUrl("https://c.me/s/3?a=1#b=2"))
	//println(CheckUrl("https//c.me/s/3?a=1#b=2"))
	//println(CheckUrl("https:/c.me/s/3?a=1#b=2"))
	//println(len(process.CacheUrl("directories", "txt")))
	//println()
	//util.Dump(util.ReadFile("directories\\配置文件\\DIR.txt"))
	//util.Dump(process.CacheHeader())
	//var headers map[string]string
	//headers := make(map[string]string)
	//headers := process.CacheHeader()
	//println(headers)
	scan("127.0.0.1", "head", nil, 100)
}

func scan(urlStr string, methodStr string, headers map[string]string, threadNum int) {
	log.Println("threadNum:", threadNum)
	urlStr = process.FormatUrl(urlStr)
	if !process.CheckUrl(urlStr) {
		log.Fatalln("check url fail")
	}

	//urlList := process.CacheUrl("directories", "txt")
	urlList := process.CacheUrl("directories/配置文件", "txt")
	log.Println("urlList", urlList)
	urlList = util.RemoveDuplicateMultiThread(urlList)
	log.Println("urlList", urlList)
	os.Exit(0)
	headerMap := process.CacheHeader()
	log.Println("headerList", len(headerMap))

	methodStr = strings.ToUpper(methodStr)

	urlQueue := make(chan string)

	for i := 0; i < threadNum; i++ {
		go scanThread(urlQueue, methodStr, headerMap, urlList)
	}

	process.AppendUrl(urlQueue, urlStr, urlList)
}

//扫描thread
func scanThread(urlQueue chan string, methodStr string, headers map[string]string, urlList []string) {
	// get url from urlQueue channel
	for urlStr := range urlQueue {
		//log.Println(urlStr)

		client := &http.Client{}

		req, err := http.NewRequest(methodStr, urlStr, nil)
		if err != nil {
			log.Println(" 扫描URL： " + urlStr + " 出错，可能为Bad Request。\n" + err.Error())
			continue
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		/*body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(body))*/
		//将结果定位到标准输出 也可以直接打印出来 或者定位到其他地方进行相应的处理
		/*stdout := os.Stdout
		_, err = io.Copy(stdout, response.Body)*/

		//返回的状态码
		StatusCode := resp.StatusCode
		status := resp.Status
		if StatusCode != 404 {
			log.Println(urlStr, status)
		}

		if !strings.Contains(urlStr, ".") && (StatusCode == 200 || StatusCode == 403 || strings.HasPrefix(status, "30")) {
			process.AppendUrl(urlQueue, urlStr, urlList)
		}
	}
}
