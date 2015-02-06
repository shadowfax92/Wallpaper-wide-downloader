package main

/*
References:
- HTML parsing using goquery example - https://www.socketloop.com/tutorials/golang-how-to-extract-links-from-web-page
- goquery documentation - http://godoc.org/github.com/PuerkitoBio/goquery#Selection
*/

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	//"strings"
)

const (
	base_url            = "http://wallpaperswide.com/"
	wallpaper_page_base = "http://wallpaperswide.com/mac-desktop-wallpapers/page/"
	download_resolution = "1920x1080"
)

func fetchDownloadableUrl(url string) (err error) {
	log.Println("\nGetting needed resolution image from url = ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalln("Error while fetching downloadable url form url = ", url)
	}
	curr_selction := doc.Find(".wallpaper-resolutions a")
	curr_selction.Each(func(i int, s *goquery.Selection) {
		resolution := s.Text()
		if resolution == download_resolution {
			href, _ := s.Attr("href")
			fmt.Println("Found resolution = ", resolution, " href = ", href)
		}
	})
	return err
}

func extractUrls(url string) (err error) {
	log.Println("\nExtracting from url = ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalln("Error parsing url = ", url)
	}

	/*
			<ul class="wallpapers">

			<li class="wall">
			<div class="thumb">
		    	<div class="mini-hud" id="hudtitle" align="center">
		        <a href="/apple_mac_os_x_blue-wallpapers.html" title="Apple MAC OS X Blue HD Wide Wallpaper for Widescreen">
		        <h1>Apple MAC OS X Blue</h1>
		        </a>
	*/

	// Try to get extract the above sub-html
	curr_sel := doc.Find(".wallpapers .wall .thumb .mini-hud a")
	curr_sel.Each(func(i int, s *goquery.Selection) {
		hrefs, is_present := s.Attr("href")
		if is_present == true {
			log.Println("Image url is ", base_url+hrefs)
			fetchDownloadableUrl(base_url + hrefs)
		}
	})
	return err
}

func main() {
	fmt.Println("Wallpaper downlaoder by Nsonti")
	var page_start int
	var page_end int
	fmt.Println("Enter page_start and page_end")
	_, _ = fmt.Scanf("%d", &page_start)
	_, _ = fmt.Scanf("%d", &page_end)

	for i := page_start; i <= page_end; i++ {
		extractUrls(wallpaper_page_base + string(i))
	}

}
