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
	fDebug    bool
	fDetailed bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}
}

func main() {
	var client *cryptcheck.Client

	site := flag.Arg(0)

	if fDebug {
		client = cryptcheck.NewClient(cryptcheck.Config{Log: 2})
	} else {
		client = cryptcheck.NewClient()
	}

	if fDetailed {
		report, err := client.GetDetailedReport(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s'\n", site)
		}

		// Just dump the json
		jr, err := json.Marshal(report)
		fmt.Printf("%s\n", jr)
	} else {
		fmt.Printf("%s Wrapper: %s API version %s\n\n",
			MyName, cryptcheck.MyVersion, cryptcheck.APIVersion)
		grade, err := client.GetScore(site)
		if err != nil {
			log.Fatalf("impossible to get grade for '%s': %v\n", site, err)
		}
		fmt.Printf("Grade for '%s' is %s\n", site, grade)
	}
}
