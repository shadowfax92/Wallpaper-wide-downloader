package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	//"strings"
)

func main() {
	doc, err := goquery.NewDocument("http://wallpaperswide.com/mac-desktop-wallpapers.html")
	if err != nil {
		log.Fatalf("Error querying url", err)
	}
	curr_sel := doc.Find(".wallpapers .wall .thumb .mini-hud a")
	//curr_sel := doc.Find(".wallpapers .wall .thumb")
	fmt.Println(curr_sel)
	curr_sel.Each(func(i int, s *goquery.Selection) {

		//Title := strings.TrimSpace(s.Find("onclick").Text())
		// convert string to array
		//fields := strings.Fields(Title)
		//fmt.Println(fields)
		//		fmt.Println(s)
		//		Link, _ := s.Attr("href")
		//		fmt.Println("Link = ", Link)

		// working-1
		//on_clicks, _ := s.Attr("onclick")
		//fmt.Println(on_clicks)

		hrefs, _ := s.Attr("href")
		fmt.Println(hrefs)
	})
}
