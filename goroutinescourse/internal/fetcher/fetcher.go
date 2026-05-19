package fetcher

import (
	"encoding/json"
	"log"
	"net/http"
)

func FetchHackerNews(keyword string) any {
	r, err := http.Get("https://hn.algolia.com/api/v1/search?query=" + keyword + "&tags=story")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var response any = nil
	json.NewDecoder(r.Body).Decode(&response)
	return response

}

func FetchReddit(keyword string) any {
	r, err := http.Get("https://www.reddit.com/search.json?q=" + keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var response any = nil
	json.NewDecoder(r.Body).Decode(&response)
	return response

}

func FetchDevTo(keyword string) any {
	r, err := http.Get("https://dev.to/api/articles?tag=" + keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var response any = nil
	json.NewDecoder(r.Body).Decode(&response)
	return response

}

func FetchGitHub(keyword string) any {
	r, err := http.Get("https://api.github.com/search/repositories?q=" + keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	var response any = nil
	json.NewDecoder(r.Body).Decode(&response)
	return response

}
