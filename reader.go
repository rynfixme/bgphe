package main

import "os"

type Reader interface {
	ReadFromFile(f *os.File) []string
}
