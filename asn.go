package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

type ASNScraperProvider[T ASNResult] struct{}

func (a *ASNScraperProvider[ASNResult]) Scrape(asn string) ASNResult {
	ipv4_prefix := []ASNIPV4Prefix{}
	ipv6_prefix := []ASNIPV6Prefix{}
	col := colly.NewCollector()
	col.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	col.OnHTML(`table#table_prefixes4 > tbody > tr`, func(ele *colly.HTMLElement) {
		ch := ele.DOM.Children()
		prefix := ch.Eq(0).Children().Eq(-1).Text()
		desc := ch.Eq(1).Text()
		re := regexp.MustCompile(`(\t|\n| {2,})`)
		desc = re.ReplaceAllString(desc, "")

		ipv4_prefix = append(ipv4_prefix, ASNIPV4Prefix{prefix, desc})
	})

	col.OnHTML(`table#table_prefixes6 > tbody > tr`, func(ele *colly.HTMLElement) {
		ch := ele.DOM.Children()
		prefix := ch.Eq(0).Children().Eq(-1).Text()
		desc := ch.Eq(1).Text()
		re := regexp.MustCompile(`(\t|\n| {2,})`)
		desc = re.ReplaceAllString(desc, "")

		ipv6_prefix = append(ipv6_prefix, ASNIPV6Prefix{prefix, desc})
	})

	url := "https://bgp.he.net/" + url.QueryEscape(asn)
	err := col.Visit(url)
	if err != nil {
		log.Fatalln(err, url)
	}

	return ASNResult{ipv4_prefix, ipv6_prefix}
}

type ASNFileReaderProvider struct{}

func (a *ASNFileReaderProvider) ReadFromFile(f *os.File) []string {
	result := []string{}
	asnPat := regexp.MustCompile(`AS\d{1,12}`)
	skipPat := regexp.MustCompile(`\s*#.*AS\d{1,12}`)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		t := scanner.Text()

		if skipPat.Match([]byte(t)) {
			continue
		}

		if !asnPat.Match([]byte(t)) {
			log.Fatalln("Unsupported pattern", scanner)
		}
		result = append(result, t)
	}

	return result
}

type ASNIPV4Prefix struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
}

type ASNIPV6Prefix struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
}

type ASNResult struct {
	IPV4Prefix []ASNIPV4Prefix `json:"ipv4_prefix"`
	IPV6Prefix []ASNIPV6Prefix `json:"ipv6_prefix"`
}

type ASNClient struct {
	ASN    *string
	ASNs   *[]string
	Result *ASNResult
	S      Scraper[ASNResult]
	R      Reader
}

func (a *ASNClient) Search() ASNResult {
	if a.ASN == nil {
		log.Fatalln("ASN not specified.")
	}

	return a.S.Scrape(*a.ASN)
}

func (a *ASNClient) SearchMulti() ASNResult {
	if a.ASNs == nil {
		log.Fatalln("ASN list is not specified.")
	}

	if len(*a.ASNs) == 0 {
		log.Fatalln("Empty ASN list is not allowed.")
	}

	result := ASNResult{[]ASNIPV4Prefix{}, []ASNIPV6Prefix{}}
	for _, asn := range *a.ASNs {
		r := a.S.Scrape(asn)
		result.IPV4Prefix = append(result.IPV4Prefix, r.IPV4Prefix...)
		result.IPV6Prefix = append(result.IPV6Prefix, r.IPV6Prefix...)
		time.Sleep(100 * time.Millisecond)
	}

	return result
}
