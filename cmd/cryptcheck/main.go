// main.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

/*
This is just a very short example.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/keltia/cryptcheck"
	"log"
	"os"
	"path/filepath"
)

var (
	fDetailed bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}
}

func main() {
	site := flag.Arg(0)

	if fDetailed {
		report, err := cryptcheck.NewClient().GetDetailedReport(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s'\n", site)
		}

		// Just dump the json
		jr, err := json.Marshal(report)
		fmt.Printf("%s\n", jr)
	} else {
		fmt.Printf("%s Wrapper: %s API version %s\n\n",
			MyName, cryptcheck.MyVersion, cryptcheck.APIVersion)
		grade, err := cryptcheck.NewClient().GetScore(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s'\n", site)
		}
		fmt.Printf("Grade for '%s' is %s\n", site, grade)
	}
}
