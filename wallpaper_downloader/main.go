package main

/*
References:
- HTML parsing using goquery example - https://www.socketloop.com/tutorials/golang-how-to-extract-links-from-web-page
- goquery documentation - http://godoc.org/github.com/PuerkitoBio/goquery#Selection
- saving image to disk - http://stackoverflow.com/questions/8648682/reading-image-from-http-requests-body-in-go
*/

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

/*const (
	base_url            = "http://www.hdwallpapers.in"
	wallpaper_page_base = "http://www.hdwallpapers.in/nature__landscape-desktop-wallpapers/page/"
	download_resolution = "1920 x 1080"
)*/

const (
	base_url            = "http://wallpaperswide.com"
	wallpaper_page_base = "http://wallpaperswide.com/3840x2160-wallpapers-r/page/"
	//wallpaper_page_base = "http://wallpaperswide.com/nature-desktop-wallpapers/page/"
	download_resolution = "2560x1440"
)

var (
	wg                sync.WaitGroup
	check_resolutions = []string{"1920x1080", "2048x1152", "2400x1350", "2560x1440", "2880x1620", "3554x1999", "3840x2160"}
)

//Download url: http://www.hdwallpapers.in/download/valley_reflections-1920x1080.jpg

func downloadImage(url string, name string) (err error) {
	wg.Add(1)
	log.Println("Downloading image url =", url)
	tmp := strings.LastIndex(url, ".")
	image_format := url[tmp:len(url)]
	image_name := name + string(image_format)

	/*res, err := http.Get(url)
	if err != nil {
		log.Println("Error: downloadImage url =", url)
		return err
	}
	data, err := ioutil.ReadAll(res.Body)*/

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:35.0) Gecko/20100101 Firefox/35.0")
	req.Header.Add("Host", "wallpaperswide.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:35.0) Gecko/20100101 Firefox/35.0")
	req.Header.Add("Referer", "http://wallpaperswide.com/os_x_yosemite_2-wallpapers.html")
	req.Header.Add("Cookie", "__qca=P0-629664164-1423208545996; __utma=30129849.1884121915.1423208546.1423214454.1423922331.3; __utmz=30129849.1423208546.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); ae74935a9f5bd890e996f9ae0c7fe805=q5vS1ldKBFw%3DRXHN2qLD5gI%3DuvhKVtnz6aQ%3D6yAA0QoLSpo%3Daa0wj%2BrGoS4%3DlopdREWA8%2B4%3DvuEblRbQplU%3D8NgLP0uGZcM%3D; __utmb=30129849.1.10.1423922331; __utmc=30129849; __utmt=1")
	req.Header.Add("Connection", "keep-alive")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error: reading image url =", url)
		return err
	}

	ioutil.WriteFile("./wallpapers/"+image_name, data, 0666)
	log.Println("Sucess!. Downloaded image = ", image_name)
	return err
}

func fetchDownloadableUrl(url string, image_name string) (err error) {
	wg.Add(1)
	log.Println("\nGetting needed resolution image from url = ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Error while fetching downloadable url form url = ", url)
		return err
	}

	var found_resolution bool = false
	curr_selction := doc.Find(".wallpaper-resolutions a")
	curr_selction.Each(func(i int, s *goquery.Selection) {
		resolution := s.Text()
		for i := (len(check_resolutions) - 1); i >= 0 && found_resolution != true; i-- {
			//log.Println("test ", check_resolutions[i])
			if resolution == check_resolutions[i] {
				href, _ := s.Attr("href")
				title, _ := s.Attr("title")
				fmt.Println("Found resolution = ", resolution, " href = ", href, " title = ", title)
				//go downloadImage(base_url+href, image_name)
				found_resolution = true
			}
		}

	})
	return err
}

func extractUrls(url string) (err error) {
	wg.Add(1)
	log.Println("\nExtracting from url = ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Error parsing url = ", url)
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
			log.Println("Image name =", image_name, " url = ", base_url+hrefs)
			go fetchDownloadableUrl(base_url+hrefs, image_name)
		}
	})
	return err
}

func main() {
	fmt.Println("Wallpaper downlaoder by Nsonti")
	defer wg.Done()

	var page_start int
	var page_end int
	fmt.Println("Enter page_start and page_end")
	_, _ = fmt.Scanf("%d", &page_start)
	_, _ = fmt.Scanf("%d", &page_end)

	for i := page_start; i <= page_end; i++ {
		wg.Add(1)
		fmt.Println(i)
		go extractUrls(fmt.Sprint(wallpaper_page_base, i))
	}

	wg.Wait()

}
