package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("henet", "Scrapeing bgp.he.ne")

	search     = app.Command("search", "Search assets by word.")
	searchWord = search.Flag("word", "Word to search.").Required().String()

	asn       = app.Command("asn", "Search assets by ASN.")
	asnNumber = asn.Flag("number", "ASN to search.").String()
	asnList   = asn.Flag("list", "ASNs to search.").File()

	prefix       = app.Command("prefix", "Search assets by prefix.")
	prefixPrefix = prefix.Flag("prefix", "Prefix to search.").String()
	prefixList   = prefix.Flag("list", "Prefixes to search.").File()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case search.FullCommand():
		prov := SearchScraperProvider[SearchResult]{}
		s := SearchClient{searchWord, nil, SearchResult{}, &prov}
		s.Search()

		bytes, err := json.Marshal(s.Result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bytes))
		return

	case asn.FullCommand():
		sprov := ASNScraperProvider[ASNResult]{}
		fprov := ASNFileReaderProvider{}

		if asnNumber != nil {
			a := ASNClient{asnNumber, nil, ASNResult{}, &sprov, &fprov}
			a.Search()

			bytes, err := json.Marshal(a.Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
			return
		}

		if asnList != nil {
			a := ASNClient{nil, nil, ASNResult{}, &sprov, &fprov}
			asns := a.R.ReadFromFile(*asnList)
			a.ASNs = &asns
			a.SearchMulti()

			bytes, err := json.Marshal(a.Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
			return
		}

		fmt.Println(app.Help)

	case prefix.FullCommand():
		sprov := PrefixScraperProvider[PrefixResult]{}
		fprov := PrefixFileReaderProvider{}

		if prefixPrefix != nil {
			p := PrefixClient{prefixPrefix, nil, PrefixResult{}, &sprov, &fprov}
			p.Search()

			bytes, err := json.Marshal(p.Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
			return
		}

		if prefixList != nil {
			p := PrefixClient{nil, nil, PrefixResult{}, &sprov, &fprov}
			prefixes := p.R.ReadFromFile(*prefixList)
			p.Prefixes = &prefixes
			p.Result = p.SearchMulti()

			bytes, err := json.Marshal(p.Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
			return
		}

		fmt.Println(app.Help)

	default:
		fmt.Println(app.Help)
	}
}
