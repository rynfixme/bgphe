package main

import "testing"

type TestPrefixSingle struct {
	Name     string
	Args     TestPrefixSingleArgs
	Expected TestPrefixSingleExptected
}

type TestPrefixSingleArgs struct {
	Prefix string
}

type TestPrefixSingleExptected struct {
	ExistsPTR     bool
	ExistsNullPTR bool
	ExistsA       bool
	ExistsNullA   bool
}

var testsPrefixSingle = []TestPrefixSingle{
	TestPrefixSingle{"should get PTR and A", TestPrefixSingleArgs{"91.204.192.0/22"}, TestPrefixSingleExptected{true, true, true, true}},
	TestPrefixSingle{"should get no PTR and A", TestPrefixSingleArgs{"91.204.195.0/24"}, TestPrefixSingleExptected{false, false, false, false}},
}

func TestPrefixClientSingle(t *testing.T) {
	sprov := PrefixScraperProvider[PrefixResult]{}
	fprov := PrefixFileReaderProvider{}

	for _, tt := range testsPrefixSingle {
		t.Run(tt.Name, func(t *testing.T) {
			c := PrefixClient{&tt.Args.Prefix, nil, nil, &sprov, &fprov}
			got := c.Search()

			var prt = false
			var prtnl = false
			var a = false
			var anl = false

			for _, prefix := range got.DNS {
				prt = prt || (prefix.PTR != nil)
				prtnl = prtnl || (prefix.PTR == nil)
				a = a || (prefix.A != nil)
				anl = anl || (prefix.A == nil)
			}

			if prt != tt.Expected.ExistsPTR {
				t.Errorf("TestPrefixClientSingle(), %v, ExistsPTR error, %v", tt.Name, got.DNS)
				return
			}

			if prtnl != tt.Expected.ExistsNullPTR {
				t.Errorf("TestPrefixClientSingle(), %v, ExistsNullPTR error, %v", tt.Name, got.DNS)
				return
			}

			if a != tt.Expected.ExistsA {
				t.Errorf("TestPrefixClientSingle(), %v, ExistsA error, %v", tt.Name, got.DNS)
				return
			}

			if anl != tt.Expected.ExistsNullA {
				t.Errorf("TestPrefixClientSingle(), %v, ExistsNullPTR error, %v", tt.Name, got.DNS)
				return
			}
		})
	}
}
