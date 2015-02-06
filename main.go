package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	//"strings"
)

const (
	base_url = "http://wallpaperswide.com/"
)

func main() {
	doc, err := goquery.NewDocument("http://wallpaperswide.com/mac-desktop-wallpapers.html")
	if err != nil {
		log.Fatalf("Error querying url", err)
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
		hrefs, _ := s.Attr("href")
		fmt.Println(hrefs)
	})
}
