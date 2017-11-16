package main

import (
	"flag"

	"github.com/reinventer/fb2epub/converter"
	"log"
)

func main() {
	var (
		inFile          string
		outFile         string
		translit        bool
		sectionsPerPage int
	)
	flag.StringVar(&inFile, "f", "-", `fb2 file, use "-" for STDIN`)
	flag.StringVar(&outFile, "t", "-", `epub file, use "-" for STDOUT`)
	flag.BoolVar(&translit, "translit", false, "transliterate header information")
	flag.IntVar(&sectionsPerPage, "sections", 10, "sections per page")
	flag.Parse()

	c, err := converter.New(inFile, sectionsPerPage)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Convert(outFile, translit); err != nil {
		log.Fatal(err)
	}
}
