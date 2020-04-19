package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type  film struct {
	Name string
	Pic  string
	Director string
	Point string
	Review string
}

//爬取网页
func spiderDB(Url string)(bodystr string,err error){
	//设置client
	//useragent
	useragent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"
	//ssl证书
	tr := &http.Transport{
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},

	}
	client := http.Client{Transport:tr}
		req, err := http.NewRequest("GET", Url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("User-Agent", useragent)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		bodystr = string(body)
		return

}
//正则匹配
func reDB(body string )[]film{
	//电影名称和图片
	re1 :=regexp.MustCompile(`<img width="100" alt="(?s:(.*?))" src="(?s:(.*?))"`)
	res :=re1.FindAllStringSubmatch(body,-1)
	names := make([]string,25)
	for i,name := range res{
		names[i] = name[1]
	}
	pics := make([]string,25)
	for i,pic := range res{
		pics[i] = pic[2]
	}
	//电影导演
	re2 :=regexp.MustCompile(`<p class="">(?s:(.*?))&nbsp`)
	res2 :=re2.FindAllStringSubmatch(body,-1)
	directors := make([]string,25)
	for i,director := range res2{
		director[1] = strings.Replace(director[1]," ","",-1)
		director[1] = strings.Replace(director[1],"导演:","",-1)
		directors[i] = director[1]
	}
	//电影评分
	re3 :=regexp.MustCompile(`<span class="rating_num" property="v:average">(?s:(.*?))</span>`)
	res3 :=re3.FindAllStringSubmatch(body,-1)
	pionts := make([]string,25)
	for i,piont := range res3{
		pionts[i] = piont[1]
	}
	//电影评价
	re4 :=regexp.MustCompile(`<span class="inq">(?s:(.*?))</span>`)
	res4 :=re4.FindAllStringSubmatch(body,-1)
	reviews := make([]string,25)
	for i,review := range res4{
		reviews[i] = review[1]
	}
	films :=make([]film,25)
	for i:=0 ;i<=24;i++{
		films[i].Name = names[i]
		films[i].Pic = pics[i]
		films[i].Director = directors[i]
		films[i].Point = pionts[i]
		films[i].Review = reviews[i]
	}
	return films
}

//序列化
func JsonDB(films []film){
	results,_ := json.Marshal(films)
	fmt.Printf("%s\n",results)
}

func main() {
	for i:=0;i<10 ;i++  {
		Url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter="
		body,err := spiderDB(Url)
		if err != nil {
			fmt.Println(err)
		}
		films := reDB(body)
		JsonDB(films)
	}
	//Url := "https://movie.douban.com/top250?start=0&filter="
	//body,_ := spiderDB(Url)
	//films := reDB(body)
	//fmt.Println(films)
}
