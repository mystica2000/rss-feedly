package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)


func setUpCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Methods","POST")
	(*w).Header().Set("Access-Control-Allow-Origin","http://127.0.0.1:5173")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func matchForRSS(link string) bool {
	if(len(link)>0) {
	matched, _ := regexp.Match(`(feed|rss|feed.xml|Feed|RSS)`, []byte(link))
	comments, _ := regexp.Match(`(comments/feed)`, []byte(link))
	if(comments) {
		return false;
	}
	return matched;
	}
	return false
}

func FindRSSFeed(webURL string) string {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/6.0"),
		colly.IgnoreRobotsTxt(),
	)

	// if p, err := proxy.RoundRobinProxySwitcher(
	// ); err == nil {
	// 	c.SetProxyFunc(p)
	// }

	str := "Cannot be retreived"; // Some sites don't allow to scrape content

	c.OnRequest(func(r *colly.Request){{
		fmt.Println("visiting", r.URL.String())
	}});

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "\nError:", err)
	})

	c.OnHTML("link[href]",func(h *colly.HTMLElement) {
		link := h.Attr("href")
		res := matchForRSS(link)
		if(res) {
		   fmt.Println("RSS FEED LINK = ",h.Request.AbsoluteURL(link))
		   str = h.Request.AbsoluteURL(link)
		}
	});

	c.Visit(webURL)
	return str;
}


func RSSFinderHandler(w http.ResponseWriter, r *http.Request) {

	setUpCORS(&w,r)
	bytes, _ := ioutil.ReadAll(r.Body)
	webURL := string(bytes)

	str := FindRSSFeed(webURL)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	res,_ := json.Marshal(str)

	w.Write(res)

}



func main() {
	r := mux.NewRouter();

	r.HandleFunc("/rss-finder",RSSFinderHandler).Methods("POST")



	log.Fatal(http.ListenAndServe(":8080",r))

}