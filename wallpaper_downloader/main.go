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

const (
	base_url            = "http://www.hdwallpapers.in"
	wallpaper_page_base = "http://www.hdwallpapers.in/nature__landscape-desktop-wallpapers/page/"
	download_resolution = "1920 x 1080"
)

var (
	wg sync.WaitGroup
)

//Download url: http://www.hdwallpapers.in/download/valley_reflections-1920x1080.jpg

func downloadImage(url string, name string) (err error) {
	log.Println("Downloading image url =", url)
	tmp := strings.LastIndex(url, ".")
	image_format := url[tmp:len(url)]
	image_name := name + string(image_format)
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error: downloadImage url =", url)
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error: reading image url =", url)
		return err
	}

	ioutil.WriteFile(image_name, data, 0666)
	log.Println("Sucess!. Downloaded image = ", image_name)
	return err
}

func fetchDownloadableUrl(url string, image_name string) (err error) {
	log.Println("\nGetting needed resolution image from url = ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Error while fetching downloadable url form url = ", url)
		return err
	}
	curr_selction := doc.Find(".wallpaper-resolutions a")
	curr_selction.Each(func(i int, s *goquery.Selection) {
		resolution := s.Text()
		if resolution == download_resolution {
			href, _ := s.Attr("href")
			title, _ := s.Attr("title")
			fmt.Println("Found resolution = ", resolution, " href = ", href, " title = ", title)
			go downloadImage(base_url+href, image_name)
			wg.Add(1)
		}
	})
	return err
}

func extractUrls(url string) (err error) {
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
	curr_sel := doc.Find(".wallpapers .wall .thumb a")
	curr_sel.Each(func(i int, s *goquery.Selection) {
		hrefs, is_present := s.Attr("href")
		if is_present == true {
			image_name := hrefs
			image_name = strings.TrimLeft(image_name, "/")
			image_name = strings.TrimRight(image_name, ".html")
			log.Println("Image name =", image_name, " url = ", base_url+hrefs)
			go fetchDownloadableUrl(base_url+hrefs, image_name)
			wg.Add(1)
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
		go extractUrls(wallpaper_page_base + string(i))
		wg.Add(1)
	}

	wg.Wait()

}
