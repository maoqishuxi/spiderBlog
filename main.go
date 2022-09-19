package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type article struct {
	fullText []byte
	url      []byte
	title    []byte
}

type page struct {
	title   []byte
	time    []byte
	context []byte
}

func sendHttpPage(url string) ([]article, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("http request page error: ", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("page parse data error: ", err)
		return nil, err
	}

	express := `<a title="" target="_blank" href="(.+?)">(.+?)</a>`
	data := regexp.MustCompile(express).FindAllSubmatch(body, -1)

	var essay []article

	for _, item := range data {
		essay = append(essay, article{
			fullText: item[0],
			url:      item[1],
			title:    item[2],
		})
	}

	return essay, nil
}

func getPage() []map[int][]article {
	//url := "https://blog.sina.com.cn/s/articlelist_1281503010_0_%d.html"
	url := "https://blog.sina.com.cn/s/articlelist_1222362564_0_%d.html"
	var essayPage []map[int][]article
	for i := 1; i < 30; i++ {
		for j := 0; j < 5; j++ {
			res, err := sendHttpPage(fmt.Sprintf(url, i))
			if err == nil {
				essayPage = append(essayPage, map[int][]article{i: res})
				break
			}
			fmt.Printf("请求次数：%d\n", j)
		}
	}

	return essayPage
}

func sendHttpContent(url string) ([]page, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("content http error: ", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("content parse error: ", err)
	}

	title := `<h2 id=".+?" class="titName SG_txta">(.+?)</h2>`
	time := `<span class="time SG_txtc">\((.+?)\)</span>`
	content := `正文开始(\W+?.+?)+?正文结束`
	titleData := regexp.MustCompile(title).FindAllSubmatch(body, -1)
	timeData := regexp.MustCompile(time).FindAllSubmatch(body, -1)
	contentData := regexp.MustCompile(content).FindAllSubmatch(body, -1)

	//fmt.Printf("length :%d ;%q\n\n", len(titleData[0]), titleData)
	//fmt.Printf("length :%d ;%q\n\n", len(timeData[0]), timeData)
	//fmt.Printf("length :%d ;%q\n\n", len(contentData[0]), contentData)

	var pageContent []page

	contentDataStr := contentData[0][0]
	pageContent = append(pageContent, page{
		title:   titleData[0][1],
		time:    timeData[0][1],
		context: contentDataStr[16 : len(contentDataStr)-17],
	})

	return pageContent, nil
}

func getPageContent(url string, code int) []map[int][]page {
	var pageContent []map[int][]page
	for i := 0; i < 5; i++ {
		res, err := sendHttpContent(url)
		if err == nil {
			pageContent = append(pageContent, map[int][]page{code: res})
			break
		}
		fmt.Printf("get page count %d\n", i)
	}
	return pageContent
}

func work() {
	result := getPage()

	for _, contents := range result {
		for page := 1; page <= len(result); page++ {
			fmt.Printf("当前爬取页面%v\n", page)
			for _, art := range contents[page] {
				url := "https:" + string(art.url)
				fmt.Printf("当前爬取链接%v\n", url)
				content := getPageContent(url, page)
				fmt.Printf("爬取page %v 链接 %v\n", page, url)
				for _, c := range content {
					for _, p := range c[page] {
						outputMarkdown(p)
						fmt.Printf("输出Markdown文档完成")
					}
				}
			}
		}
	}
}

func outputMarkdown(pageContent page) {
	title := strings.ReplaceAll(string(pageContent.title), "/", "")
	file, err := os.OpenFile("./pages/"+title+".md",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write([]byte("# "))
	file.Write(pageContent.title)
	file.Write([]byte("\n"))
	file.Write(pageContent.time)
	file.Write([]byte("\n"))
	file.Write(pageContent.context)
}

func main() {
	//fmt.Println("welcome to you")
	//url := "https://blog.sina.com.cn/s/blog_4c622f220100fi23.html"
	//url := "https://blog.sina.com.cn/s/blog_4c622f2201019gcr.html"
	//
	//result := getPageContent(url, 1)
	//fmt.Printf("%q", result)
	//page := result[0][1][0]
	//outputMarkdown(page)

	work()
}
