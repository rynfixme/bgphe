package main

type Scraper[T any] interface {
	Scrape(word string) T
	ScrapeMulti(words []string) T
}
