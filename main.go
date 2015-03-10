package main

/*
References:
- HTML parsing using goquery example - https://www.socketloop.com/tutorials/golang-how-to-extract-links-from-web-page
- goquery documentation - http://godoc.org/github.com/PuerkitoBio/goquery#Selection
- saving image to disk - http://stackoverflow.com/questions/8648682/reading-image-from-http-requests-body-in-go
*/

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	base_url            = "http://wallpaperswide.com"
	wallpaper_page_base = "http://wallpaperswide.com/nature-desktop-wallpapers/page/"
	download_resolution = "3840x2160"
	wg                  sync.WaitGroup
	check_resolutions   = []string{"3840x2160", "3554x1999", "2880x1620", "2560x1440", "2400x1350", "2048x1152", "1920x1080"}
	ticker              *time.Ticker

	errorNon200HTTPStatus = errors.New("Non-200 HTTP status code")
)

func resolutionProduct(res string) (proc int64, err error) {
	proc = -1
	strings.Replace(res, " ", "", -1) //remove white spaces
	log.Println(res)
	var res1, res2 int
	_, err = fmt.Scanf(res, "%dx%d", &res1, &res2)
	if err != nil {
		log.Println("Error: resolutionProduct() = ", err)
		return proc, err
	}
	proc = int64(res1 * res2)
	return proc, err
}

//Download url: http://www.hdwallpapers.in/download/valley_reflections-1920x1080.jpg

func getCookie(url string) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:35.0) Gecko/20100101 Firefox/35.0")
	req.Header.Add("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		Error.Println("getCookie client.Do failed!. Err = ", err)
		return err
	}
	Info.Println("response header = ", resp.Header)
	Info.Println("response length = ", resp.ContentLength)
	Info.Println("response Body = ", resp.Body)
	Info.Println("Cookies array length = ", len(resp.Cookies()))
	for _, val := range resp.Cookies() {
		Info.Println(val)
	}
	return err
}

func downloadImage(url string, name string) (err error) {
	/*select {
	case tm := <-ticker.C:
		{
			defer wg.Done()
			Debug.Println("ticker received = ", tm)
			Debug.Println("Downloading image url = ", url)
			tmp := strings.LastIndex(url, ".")
			image_format := url[tmp:len(url)]
			image_name := name + string(image_format)

			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:35.0) Gecko/20100101 Firefox/35.0")
			req.Header.Add("Host", "wallpaperswide.com")
			req.Header.Add("Cookie", "__gads=ID=7d4ffd543145afad:T=1425997114:S=ALNI_MafR0WD0K9fFNiTNH7aPBIIBaMT0w; __utma=30129849.257949811.1425997114.1425997114.1425997114.1; __utmb=30129849.3.10.1425997114; __utmc=30129849; __utmz=30129849.1425997114.1.1.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmt=1; __qca=P0-928709853-1425997115885")

			req.Header.Add("Connection", "keep-alive")
			resp, err := client.Do(req)
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Println("Error: reading image url =", url)
				return err
			}

			ioutil.WriteFile("./downloads/"+image_name, data, 0666)
			Info.Println("Sucess!. Downloaded image = ", image_name)
		}
	default:
		return err
	}*/

	defer wg.Done()
	//Debug.Println("ticker received = ", tm)
	Debug.Println("Downloading image url = ", url)
	tmp := strings.LastIndex(url, ".")
	image_format := url[tmp:len(url)]
	image_name := name + string(image_format)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:35.0) Gecko/20100101 Firefox/35.0")
	req.Header.Add("Host", "wallpaperswide.com")
	req.Header.Add("Cookie", "__gads=ID=7d4ffd543145afad:T=1425997114:S=ALNI_MafR0WD0K9fFNiTNH7aPBIIBaMT0w; __utma=30129849.257949811.1425997114.1425997114.1425997114.1; __utmc=30129849; __utmz=30129849.1425997114.1.1.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __qca=P0-928709853-1425997115885")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Referer", url)

	req.Header.Add("Connection", "keep-alive")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	http_status := resp.StatusCode
	Info.Println("HTTP status code = ", http_status, " Url = ", url)
	if http_status != 200 {
		return errorNon200HTTPStatus
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error: reading image url =", url)
		return err
	}

	ioutil.WriteFile("./downloads/"+image_name, data, 0666)
	Info.Println("Sucess!. Downloaded image = ", image_name)

	return err
}

func fetchDownloadableUrl(url string, image_name string) (err error) {
	defer wg.Done()
	Debug.Println("\nGetting needed resolution image from url = ", url)
	doc, err := goquery.NewDocument(url)

	if err != nil {
		Error.Println("Error while fetching downloadable url form url = ", url)
		return err
	}

	var max_resolution_found string = "0x0"
	var download_url string

	curr_selection := doc.Find(".wallpaper-resolutions a")
	curr_selection.Each(func(i int, s *goquery.Selection) {
		//		resolution, err := resolutionProduct(s.Text())
		resolution := s.Text()
		//if resolution > max_resolution_found {
		if resolution == download_resolution {
			href, _ := s.Attr("href")
			//title, _ := s.Attr("title")
			max_resolution_found = resolution
			download_url = base_url + href

			Info.Println("new_max_resolution = ", max_resolution_found, " image_name = ", image_name, " download_url = ", download_url)
			wg.Add(1)
			go downloadImage(download_url, image_name)
		}

	})
	return err
}

func extractUrls(url string) (err error) {
	defer wg.Done()
	Debug.Println("\nExtracting from url = ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		Error.Println("Error parsing url = ", url)
		return err
	}

	/*
			<ul class="wallpapers">

			<li class="wall">
			<div class="thumb">
		        <a href="/apple_mac_os_x_blue-wallpapers.html" title="Apple MAC OS X Blue HD Wide Wallpaper for Widescreen">
		        <h1>Apple MAC OS X Blue</h1>
		        </a>
	*/

	// Try to get extract the above sub-html
	curr_sel := doc.Find(".wallpapers .wall .thumb .mini-hud a")
	curr_sel.Each(func(i int, s *goquery.Selection) {
		hrefs, is_present := s.Attr("href")
		if is_present == true {
			image_name := hrefs
			image_name = strings.TrimLeft(image_name, "/")
			image_name = strings.TrimRight(image_name, ".html")
			Info.Println("Image name = ", image_name, " url = ", base_url+hrefs)
			wg.Add(1)
			go fetchDownloadableUrl(base_url+hrefs, image_name)
		}
	})
	return err
}

func main() {
	LoggingInit(os.Stdout, os.Stdout, os.Stderr, os.Stdout)
	//getCookie(base_url)
	//os.Exit(1)

	var page_start int
	var page_end int
	fmt.Println("Enter page_start and page_end")
	_, _ = fmt.Scanf("%d", &page_start)
	_, _ = fmt.Scanf("%d", &page_end)

	ticker = time.NewTicker(time.Duration(1))
	for i := page_start; i <= page_end; i++ {
		wg.Add(1)
		go extractUrls(fmt.Sprint(wallpaper_page_base, i))
	}

	wg.Wait()

}
