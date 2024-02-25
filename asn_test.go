package main

import (
	"regexp"
	"testing"
)

type TestASNSingle struct {
	Name     string
	Args     TestASNSingleArgs
	Expected TestASNSingleExpected
}

type TestASNSingleArgs struct {
	ASN string
}

type TestASNSingleExpected struct {
	ExistsV4       bool
	ExistsV6       bool
	MatchV4Pattern bool
	MatchV6Pattern bool
}

var testsASNSingle = []TestASNSingle{
	TestASNSingle{"should get v4 and v6", TestASNSingleArgs{"AS11251"}, TestASNSingleExpected{true, true, true, true}},
	TestASNSingle{"should get v4", TestASNSingleArgs{"AS399490"}, TestASNSingleExpected{true, false, true, false}},
	TestASNSingle{"should get nothing", TestASNSingleArgs{"AS400265"}, TestASNSingleExpected{false, false, false, false}},
}

func TestASNClientSingle(t *testing.T) {
	sprov := ASNScraperProvider[ASNResult]{}
	fprov := ASNFileReaderProvider{}

	for _, tt := range testsASNSingle {
		t.Run(tt.Name, func(t *testing.T) {
			c := ASNClient{&tt.Args.ASN, nil, nil, &sprov, &fprov}
			got := c.Search()

			if (len(got.IPV4Prefix) > 0) != tt.Expected.ExistsV4 {
				t.Errorf("TestASNClientSingle(), %v, ExistsV4 error, %v", tt.Name, got.IPV4Prefix)
				return
			}

			if (len(got.IPV6Prefix) > 0) != tt.Expected.ExistsV6 {
				t.Errorf("TestASNClientSingle() %v ExistsV6 error %v", tt.Name, got.IPV6Prefix)
				return
			}

			for _, v4Prefix := range got.IPV4Prefix {
				pat := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
				if (pat.Match([]byte(v4Prefix.Prefix))) != tt.Expected.MatchV4Pattern {
					t.Errorf("TestASNClientSingle(), %v, MatchV4Pattern error, %v", tt.Name, got.IPV4Prefix)
					return
				}
			}

			for _, v6Prefix := range got.IPV6Prefix {
				pat := regexp.MustCompile(`([0-9,a-d]|::)`)
				if (pat.Match([]byte(v6Prefix.Prefix))) != tt.Expected.MatchV6Pattern {
					t.Errorf("TestASNClientSingle(), %v, MatchV6Pattern error, %v", tt.Name, got.IPV6Prefix)
					return
				}
			}
		})
	}
}
