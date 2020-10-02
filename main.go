package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Image ...
type Image struct {
	URL      string
	FileName string
}

// Download Image
func (image *Image) Download(foldPath string) {
	bytes := GetSource(image.URL, map[string]string{
		"content-type": "image/gif",
	})
	err := ioutil.WriteFile(path.Join(foldPath, image.FileName), bytes, 0666)

	HandlerErr(err)
}

// GetSource from url with spec headers
func GetSource(url string, headers map[string]string) []byte {
	tr := &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	HandlerErr(err)
	for k := range headers {
		req.Header.Set(k, headers[k])
	}

	resp, err := client.Do(req)
	HandlerErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body

}

// HandlerErr handler the err
func HandlerErr(err error) {
	if err != nil {
		panic(err)
	}

}

//GithubUserSpider
func GithubUserSpider(githubUrl string) []Image {

	var imgs []Image
	html := string(GetSource(githubUrl, map[string]string{
		"User-Agent": "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
	}))
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	HandlerErr(err)
	dom.Find("div.user-list-item").Each(
		func(i int, selection *goquery.Selection) {
			nickName := selection.Find("a.text-gray").First().Text()
			fmt.Println(nickName + "sss")
			imgURL, IsExists := selection.Find("img").First().Attr("src")
			if IsExists && strings.HasPrefix(imgURL, "http") {
				fileName := fmt.Sprintf("%s.png", nickName)
				imgs = append(imgs, Image{URL: strings.Split(imgURL, "?")[0], FileName: fileName})

			}
		})
	fmt.Println(imgs)
	return imgs
}

func main() {

	var imageDownloadFold = "imgs"
	var wg sync.WaitGroup
	city := "melbourne"

	for i := 1; i < 10; i++ {
		githubURL := fmt.Sprintf("https://github.com/search?p=%d&q=location%%3A%s&type=Users", i, city)
		fmt.Println(githubURL)
		imgs := GithubUserSpider(githubURL)
		fmt.Println(imgs)
		wg.Add(len(imgs))
		for _, im := range imgs {

			go func(img Image) {
				img.Download(imageDownloadFold)
				wg.Done()
			}(im)

		}
		wg.Wait()
		time.Sleep(5 * time.Second)
	}

}
