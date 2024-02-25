package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/gocolly/colly"
)

type SearchScraperProvider[T SearchResult] struct{}

func (s *SearchScraperProvider[SearchResult]) Scrape(word string) SearchResult {
	v4Prefix := []SearchIPV4Prefix{}
	v6Prefix := []SearchIPV6Prefix{}
	asn := []SearchASN{}
	asnPat := regexp.MustCompile(`AS\d{1,12}`)
	ipv4Pat := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	ipv6Pat := regexp.MustCompile(`([0-9,a-d]|::)`)

	col := colly.NewCollector()

	col.OnHTML("div#search > table > tbody > tr", func(ele *colly.HTMLElement) {
		ch := ele.DOM.Children()
		result := ch.Eq(0).Children().Eq(0).Text()
		rtype := ch.Eq(1).Text()
		desc := ch.Eq(2).Text()

		if asnPat.Match([]byte(result)) {
			asn = append(asn, SearchASN{result, rtype, desc})
		} else if ipv4Pat.Match([]byte(result)) {
			v4Prefix = append(v4Prefix, SearchIPV4Prefix{result, rtype, desc})
		} else if ipv6Pat.Match([]byte(result)) {
			v6Prefix = append(v6Prefix, SearchIPV6Prefix{result, rtype, desc})
		} else {
			// do nothing
			fmt.Println("ignored", result, rtype, desc)
		}
	})

	url := "https://bgp.he.net/search?search%5Bsearch%5D=" + url.QueryEscape(word) + "&commit=Search"
	err := col.Visit(url)
	if err != nil {
		log.Fatalln(err, url)
	}

	return SearchResult{v4Prefix, v6Prefix, asn}
}

func (s *SearchScraperProvider[SearchResult]) ScrapeMulti(words []string) SearchResult {
	return SearchResult{
		[]SearchIPV4Prefix{},
		[]SearchIPV6Prefix{},
		[]SearchASN{},
	}
}

type SearchIPV4Prefix struct {
	Result      string `json:"result"`
	SearchType  string `json:"search_type"`
	Description string `json:"description"`
}

type SearchIPV6Prefix struct {
	Result      string `json:"result"`
	SearchType  string `json:"search_type"`
	Description string `json:"description"`
}

type SearchASN struct {
	Result      string `json:"result"`
	SearchType  string `json:"search_type"`
	Description string `json:"description"`
}

type SearchResult struct {
	IPV4Prefix []SearchIPV4Prefix `json:"ipv4_prefix"`
	IPV6Prefix []SearchIPV6Prefix `json:"ipv6_prefix"`
	ASN        []SearchASN        `json:"asn"`
}

type SearchClient struct {
	Word   *string
	Words  *[]string
	Result SearchResult
	S      Scraper[SearchResult]
}

func (s *SearchClient) Search() SearchResult {
	if s.Word == nil {
		log.Fatalln("Search word is not specified.")
	}
	return s.S.Scrape(*s.Word)
}
