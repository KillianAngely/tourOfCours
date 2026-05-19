package main

import (
	"fmt"
	"goroutinescourse/internal/fetcher"

	"time"
)

func main() {
	println("Goroutines courses")
	startR := time.Now()
	search("React")
	duration := time.Since(startR)
	fmt.Print(duration)

}

func search(keyword string) []any {
	var articles []any
	articles = append(articles, fetcher.FetchDevTo(keyword))
	articles = append(articles, fetcher.FetchGitHub(keyword))
	articles = append(articles, fetcher.FetchHackerNews(keyword))
	articles = append(articles, fetcher.FetchReddit(keyword))
	return articles
}
