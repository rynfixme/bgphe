package main

import (
	"regexp"
	"testing"
)

type TestSearchSingle struct {
	Name     string
	Args     TestSearchSingleArgs
	Expected TestSearchSingleExpected
}

type TestSearchSingleArgs struct {
	Search string
}

type TestSearchSingleExpected struct {
	ExistsASN       bool
	ExistsV4        bool
	ExistsV6        bool
	MatchASNPattern bool
	MatchV4Pattern  bool
	MatchV6Pattern  bool
}

var testsSearchSingle = []TestSearchSingle{
	TestSearchSingle{"should get ASN, V4, V6", TestSearchSingleArgs{"red bull"}, TestSearchSingleExpected{true, true, true, true, true, true}},
	TestSearchSingle{"should get no ASN, V4, V6", TestSearchSingleArgs{"abcdefg"}, TestSearchSingleExpected{false, false, false, false, false, false}},
}

func TestSearchClientSingle(t *testing.T) {
	sprov := SearchScraperProvider[SearchResult]{}

	for _, tt := range testsSearchSingle {
		t.Run(tt.Name, func(t *testing.T) {
			c := SearchClient{&tt.Args.Search, nil, SearchResult{}, &sprov}
			got := c.Search()

			if (len(got.ASN) > 0) != tt.Expected.ExistsASN {
				t.Errorf("TestSearchClientSingle(), %v, ExistsASN error, %v", tt.Name, got.ASN)
				return
			}

			if (len(got.IPV4Prefix) > 0) != tt.Expected.ExistsV4 {
				t.Errorf("TestSearchClientSingle(), %v, ExistsV4 error, %v", tt.Name, got.IPV4Prefix)
				return
			}

			if (len(got.IPV6Prefix) > 0) != tt.Expected.ExistsV6 {
				t.Errorf("TestSearchClientSingle() %v ExistsV6 error %v", tt.Name, got.IPV6Prefix)
				return
			}

			for _, asn := range got.ASN {
				pat := regexp.MustCompile(`AS\d{1,12}`)
				if (pat.Match([]byte(asn.Result))) != tt.Expected.MatchASNPattern {
					t.Errorf("TestSearchClientSingle(), %v, MatchASNPattern error, %v", tt.Name, got.ASN)
					return
				}
			}

			for _, v4Prefix := range got.IPV4Prefix {
				pat := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
				if (pat.Match([]byte(v4Prefix.Result))) != tt.Expected.MatchV4Pattern {
					t.Errorf("TestSearchClientSingle(), %v, MatchV4Pattern error, %v", tt.Name, got.IPV4Prefix)
					return
				}
			}

			for _, v6Prefix := range got.IPV6Prefix {
				pat := regexp.MustCompile(`([0-9,a-d]|::)`)
				if (pat.Match([]byte(v6Prefix.Result))) != tt.Expected.MatchV6Pattern {
					t.Errorf("TestSearchClientSingle(), %v, MatchV6Pattern error, %v", tt.Name, got.IPV6Prefix)
					return
				}
			}
		})
	}
}
