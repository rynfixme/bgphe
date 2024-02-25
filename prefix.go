package main

import (
	"bufio"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

type PrefixScraperProvider[T PrefixResult] struct{}

func (p *PrefixScraperProvider[PrefixResult]) Scrape(prefix string) PrefixResult {
	col := colly.NewCollector()
	dns := []PrefixDNSRecord{}

	col.OnHTML("div#dnsrecords > table > tbody > tr", func(ele *colly.HTMLElement) {
		var ptr *string
		var a *string
		ch := ele.DOM.Children()
		ip := ch.Eq(0).Children().Eq(0).Text()

		ptrtxt := ch.Eq(1).Text()
		ptr = &ptrtxt
		if *ptr == "" {
			ptr = nil
		}

		atxt := ch.Eq(2).Children().Eq(0).Text()
		a = &atxt
		if *a == "" {
			a = nil
		}

		dns = append(dns, PrefixDNSRecord{ip, ptr, a})
	})

	url := "https://bgp.he.net/net/" + prefix
	err := col.Visit(url)
	if err != nil {
		log.Fatalln(err, url)
	}

	return PrefixResult{dns}
}

func (p *PrefixScraperProvider[PrefixResult]) ScrapeMulti(prefixes []string) PrefixResult {
	col := colly.NewCollector()
	dns := []PrefixDNSRecord{}

	col.OnHTML("div#dnsrecords > table > tbody > tr", func(ele *colly.HTMLElement) {
		var ptr *string
		var a *string
		ch := ele.DOM.Children()
		ip := ch.Eq(0).Children().Eq(0).Text()

		ptrtxt := ch.Eq(1).Text()
		ptr = &ptrtxt
		if *ptr == "" {
			ptr = nil
		}

		atxt := ch.Eq(2).Children().Eq(0).Text()
		a = &atxt
		if *a == "" {
			a = nil
		}

		dns = append(dns, PrefixDNSRecord{ip, ptr, a})
	})

	for _, prefix := range prefixes {
		url := "https://bgp.he.net/net/" + prefix
		err := col.Visit(url)
		if err != nil {
			log.Fatalln(err, url)
		}
	}

	return PrefixResult{dns}
}

type PrefixFileReaderProvider struct{}

func (p *PrefixFileReaderProvider) ReadFromFile(f *os.File) []string {
	result := []string{}
	ipv4Pat := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	ipv6Pat := regexp.MustCompile(`([0-9,a-d]|::)`)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if !ipv4Pat.Match([]byte(t)) && !ipv6Pat.Match([]byte(t)) {
			log.Fatalln("Unsupported pattern", scanner)
		}
		result = append(result, t)
	}

	return result
}

type PrefixDNSRecord struct {
	IP  string  `json:"ip"`
	PTR *string `json:"ptr"`
	A   *string `json:"a"`
}

type PrefixResult struct {
	DNS []PrefixDNSRecord `json:"dns"`
}

type PrefixClient struct {
	Prefix   *string
	Prefixes *[]string
	Result   PrefixResult
	S        Scraper[PrefixResult]
	R        Reader
}

func (p *PrefixClient) Search() PrefixResult {
	if p.Prefix == nil {
		log.Fatalln("Prefix not specified.")
	}

	return p.S.Scrape(*p.Prefix)
}

func (p *PrefixClient) SearchMulti() PrefixResult {
	if p.Prefixes == nil {
		log.Fatalln("Prefix not specified.")
	}

	return p.S.ScrapeMulti(*p.Prefixes)
}
