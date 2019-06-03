package main

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	mainPicDirTpl   string = "/主图/"
	detailPicDirTpl string = "/详情图/"
	skuPicDirTpl    string = "/款式图/"
)

type imageStruct struct {
	Preview  string `json:"preview"`
	Original string `json:"original"`
}

type skuStruct struct {
	Prop  string              `json:"prop"`
	Value []DownloadImgStruct `json:"value"`
}

type DownloadImgStruct struct {
	Url  string `json:"imageUrl"`
	Name string `json:"name"`
}

func DownloadImgs(url string) (string, error) {
	mainImgUrls, deatilUrlStr, skuImgs, err := spliderMainPic(url)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}

	t := time.Now()
	dir := fmt.Sprintf("./%s-%d", t.Format("20060102"), t.Unix())

	var files []string
	fmt.Println("采集主图....")
	mfiles := downloadMainImgs(mainImgUrls, dir)
	files = append(files, mfiles...)

	fmt.Println("采集详情图....")
	dfiles := downloadDetailImgs(deatilUrlStr, dir)
	files = append(files, dfiles...)

	fmt.Println("采集款式图....")
	sfiles := downloadSkuImgs(skuImgs, dir)
	files = append(files, sfiles...)

	zipName := fmt.Sprintf("%s.zip", dir)
	generateZip(zipName, files)
	fmt.Println("采集成功 ", dir)

	err = os.RemoveAll(dir)
	if err != nil {
		logrus.Error(err.Error())
	}

	return zipName, nil
}

func generateZip(zipName string, files []string) {

	// 创建一个缓冲区用来保存压缩文件内容
	buf := new(bytes.Buffer)

	// 创建一个压缩文档
	w := zip.NewWriter(buf)

	// 将文件加入压缩文档
	for _, file := range files {
		if file == "" {
			continue
		}

		f, err := w.Create(file)
		if err != nil {
			log.Printf("[generateZip]  %s %s", file, err.Error())
			continue
		}
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("[generateZip]  %s %s", file, err.Error())
			continue
		}

		_, err = f.Write(fileContent)
		if err != nil {
			log.Printf("[generateZip]  %s %s", file, err.Error())
			continue
		}
	}

	// 关闭压缩文档
	err := w.Close()
	if err != nil {
		log.Printf("[generateZip] close %s", err.Error())
		return
	}

	// 将压缩文档内容写入文件
	f, err := os.OpenFile(zipName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("[generateZip] close %s", err.Error())
		return
	}
	buf.WriteTo(f)

}

func downloadMainImgs(mainImgUrls []string, dir string) []string {
	mainPicDir := dir + mainPicDirTpl
	if err := os.MkdirAll(mainPicDir, os.ModeDir); err != nil {
		log.Println(err.Error())
		return nil
	}
	var mainImgs []DownloadImgStruct
	for j, mPic := range mainImgUrls {
		fn := fmt.Sprintf("main-%d", j)
		mainImgs = append(mainImgs, DownloadImgStruct{
			Url:  mPic,
			Name: fn,
		})

	}
	return saveImages(mainPicDir, mainImgs)
}

func downloadDetailImgs(deatilUrlStr string, dir string) []string {

	detailUrls, err := getDetailImgUrl(deatilUrlStr)
	if err != nil {
		return nil
	}

	detailPicDir := dir + detailPicDirTpl
	if err := os.MkdirAll(detailPicDir, os.ModeDir); err != nil {
		log.Println(err.Error())
		return nil
	}
	var detailImgs []DownloadImgStruct
	for j, dPic := range detailUrls {
		fn := fmt.Sprintf("detail-%d", j)
		detailImgs = append(detailImgs, DownloadImgStruct{
			Url:  dPic,
			Name: fn,
		})
	}
	return saveImages(detailPicDir, detailImgs)

}

func downloadSkuImgs(skuImgs []DownloadImgStruct, dir string) []string {

	skuPicDir := dir + skuPicDirTpl
	if err := os.MkdirAll(skuPicDir, os.ModeDir); err != nil {
		log.Println(err.Error())
		return nil
	}

	return saveImages(skuPicDir, skuImgs)
}

func spliderMainPic(url string) (mainImgUrl []string, detailContentUrl string, skuImgs []DownloadImgStruct, err error) {
	resp := GetHtml(url)
	defer resp.Body.Close()

	content, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.Error(err.Error())
		return mainImgUrl, detailContentUrl, skuImgs, err
	}

	skuImagesStr := content.Find("script").Text()
	skuImgs = getSkuImgUrl(skuImagesStr)

	content.Find(`.tab-trigger`).Each(func(i int, s *goquery.Selection) {
		var img imageStruct
		imgJsonStr, exist := s.Attr("data-imgs")
		if !exist {
			return
		}
		json.Unmarshal([]byte(imgJsonStr), &img)

		mainImgUrl = append(mainImgUrl, img.Original)
	})

	deatailUrl, _ := content.Find(".desc-lazyload-container").Attr("data-tfs-url")

	return mainImgUrl, deatailUrl, skuImgs, err
}

func getSkuImgUrl(skuImagesStr string) []DownloadImgStruct {
	exp := regexp.MustCompile(`"skuProps":[\s\S]*?]}],`)
	matchs := exp.FindStringSubmatch(skuImagesStr)
	if len(matchs) == 0 {
		return nil
	}

	var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(matchs[0]))
	str := string(decodeBytes)

	var skuDatas []skuStruct
	err := json.Unmarshal([]byte(str[11:len(str)-1]), &skuDatas)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var resp []DownloadImgStruct
	for _, sdata := range skuDatas {
		resp = append(resp, sdata.Value...)
	}

	return resp
}

func getSkuProps(url string) (skuProps []skuStruct, err error) {
	resp := GetHtml(url)
	defer resp.Body.Close()

	content, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.Errorf("url: %s, error: %s", url, err.Error())
		return skuProps, err
	}

	skuImagesStr := content.Find("script").Text()
	exp := regexp.MustCompile(`"skuProps":[\s\S]*?]}],`)
	matchs := exp.FindStringSubmatch(skuImagesStr)
	if len(matchs) == 0 {
		return skuProps, nil
	}

	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(matchs[0]))
	str := string(decodeBytes)

	err = json.Unmarshal([]byte(str[11:len(str)-1]), &skuProps)
	if err != nil {
		logrus.Error(err.Error())
		return skuProps, err
	}

	return skuProps, nil
}

func getDetailImgUrl(url string) (detailImgUrl []string, err error) {
	//	fmt.Println(url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("NewRequest", err.Error())
		return
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	reg := regexp.MustCompile(`src=\\"(.*?)\\"`)

	urlRaws := reg.FindAllString(string(content), -1)

	for _, urlRaw := range urlRaws {
		tmp := strings.Split(urlRaw, "\\")

		detailImgUrl = append(detailImgUrl, tmp[1][1:])
	}
	return detailImgUrl, nil
}

func saveImages(path string, imgUrls []DownloadImgStruct) []string {
	var files []string
	for i, img := range imgUrls {
		if img.Url == "" {
			continue
		}

		response, err := http.Get(img.Url)
		if err != nil {
			log.Println("get img_url failed:", err)
			return nil
		}

		defer response.Body.Close()

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("read data failed:", img.Url, err)
			return nil
		}

		urlSlice := strings.Split(img.Url, ".")
		var filename string
		if img.Name == "" {
			filename = fmt.Sprintf("%d.%s", i, urlSlice[len(urlSlice)-1])
		} else {
			filename = fmt.Sprintf("%s.%s", img.Name, urlSlice[len(urlSlice)-1])
		}

		file := path + filename

		image, err := os.Create(file)
		if err != nil {
			log.Println("create file failed:", filename, err)
			continue
		}
		defer image.Close()
		image.Write(data)

		files = append(files, file)
	}

	return files
}

func GetHtml(bossurl string) *http.Response {
	netTransport := &http.Transport{ //要管理代理、TLS配置、keep-alive、压缩和其他设置，可以创建一个Transport
		//Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 2, //超时设置
	}

	client := &http.Client{ //要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client
		Timeout:   time.Second * 2,
		Transport: netTransport,
	}
	req, err := http.NewRequest("GET", bossurl, nil) //NewRequest使用指定的方法、网址和可选的主题创建并返回一个新的*Request。

	if err != nil {
		logrus.Errorf(err.Error())
	}
	req = addHeader(req)
	resp, err := client.Do(req) //Do方法发送请求，返回HTTP回复
	if err != nil {
		logrus.Errorf(err.Error())
	}

	return resp //返回网页响应
}

func addHeader(req *http.Request) *http.Request {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36") //模拟浏览器User-Agent
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("cookie", `UM_distinctid=169b9663a1024c-0610d13d3f621d-5f1d3a17-100200-169b9663a1181; cna=zPIcFRPem3wCAQ6bnNGsvpWL; ali_ab=218.18.229.179.1553593679505.0; hng=CN%7Czh-CN%7CCNY%7C156; h_keys="%u91d1%u521a%u53d8%u5f62%u98de%u673a#%u6298%u53e0%u7bb1"; ad_prefer="2019/05/24 10:16:03"; cookie2=1b27ead2f1df197cd602757c68f5289b; t=cdcef7f2c7e9abfb623db7fc4d2ed773; _tb_token_=ed5365e111e53; __wapcsf__=1; alicnweb=homeIdttS%3D72684858764667525296251236198019265397%7ChomeIdttSAction%3Dtrue%7Ctouch_tb_at%3D1559530186745%7Clastlogonid%3Dtb21014790; cookie1=BvXn%2FQU7MfUXazIubvQL%2B%2BCkvgZm9HAdPqkuP4MZy4k%3D; cookie17=UNDVc8fWLWFBIA%3D%3D; sg=151; csg=3e0f88c9; lid=counting111; unb=3011023735; __cn_logon__=true; __cn_logon_id__=counting111; ali_apache_track=c_mid=b2b-3011023735808f2|c_lid=counting111|c_ms=1|c_mt=3; ali_apache_tracktmp=c_w_signed=Y; _nk_=counting111; last_mid=b2b-3011023735808f2; _csrf_token=1559530207141; _is_show_loginId_change_block_=b2b-3011023735808f2_false; _show_force_unbind_div_=b2b-3011023735808f2_false; _show_sys_unbind_div_=b2b-3011023735808f2_false; _show_user_unbind_div_=b2b-3011023735808f2_false; CNZZDATA1253659577=277548610-1553592267-https%253A%252F%252Fs.1688.com%252F%7C1559527168; __rn_alert__=false; l=bBxXyQq7vAuqf3MSBOCZCQhfhk79jIRxjuSJcRxMi_5Iq9L6fabOlUUtShp6Vj5R_qTH4keiqTy9-etkx; isg=BMnJNn-F8YwbDI2VLr6ZX08s2PWPArdLNj31WWs-RrDvsunEsWWPGJTo9Fah2VWA`)
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en;q=0.7")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")

	return req
}
