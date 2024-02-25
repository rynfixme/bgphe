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
		cli := SearchClient{searchWord, nil, nil, &prov}
		result := cli.Search()

		cli.Result = &result
		if cli.Result == nil {
			log.Fatalln("Search fetching has not completed.")
			return
		}

		bytes, err := json.Marshal(cli.Result)
		if err != nil {
			log.Fatalln(err)
			return
		}
		fmt.Println(string(bytes))
		return

	case asn.FullCommand():
		var cli ASNClient
		var result ASNResult
		sprov := ASNScraperProvider[ASNResult]{}
		fprov := ASNFileReaderProvider{}

		if asnNumber != nil && *asnNumber != "" {
			cli = ASNClient{asnNumber, nil, nil, &sprov, &fprov}
			result = cli.Search()
			cli.Result = &result
			if cli.Result == nil {
				log.Fatalln("ASN fetching has not completed.")
				return
			}

			bytes, err := json.Marshal(cli.Result)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(bytes))
			return
		}

		if asnList != nil {
			fmt.Println("--list", *asnList)
			cli = ASNClient{nil, nil, nil, &sprov, &fprov}
			asns := cli.R.ReadFromFile(*asnList)
			cli.ASNs = &asns

			result = cli.SearchMulti()
			cli.Result = &result
			if cli.Result == nil {
				log.Fatalln("ASN fetching has not completed.")
				return
			}

			bytes, err := json.Marshal(cli.Result)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(bytes))
			return
		}

	case prefix.FullCommand():
		var cli PrefixClient
		var result PrefixResult
		sprov := PrefixScraperProvider[PrefixResult]{}
		fprov := PrefixFileReaderProvider{}

		if prefixPrefix != nil && *prefixPrefix != "" {
			cli = PrefixClient{prefixPrefix, nil, nil, &sprov, &fprov}

			result = cli.Search()
			cli.Result = &result
			if cli.Result == nil {
				log.Fatalln("Prefix fetching has not completed.")
				return
			}

			bytes, err := json.Marshal(cli.Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
		}

		if prefixList != nil {
			cli = PrefixClient{nil, nil, nil, &sprov, &fprov}
			prefixes := cli.R.ReadFromFile(*prefixList)
			cli.Prefixes = &prefixes

			result = cli.SearchMulti()
			cli.Result = &result
			if cli.Result == nil {
				log.Fatalln("Prefix fetching has not completed.")
				return
			}

			bytes, err := json.Marshal(cli.Result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
		}

	default:
		fmt.Println(app.Help)
	}
}
