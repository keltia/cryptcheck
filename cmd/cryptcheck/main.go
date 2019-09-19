// main.go
//
// Copyright 2018-2019 © by Ollivier Robert <roberto@keltia.net>

/*
This is just a very short example.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/keltia/cryptcheck"
)

var (
	fDebug    bool
	fDetailed bool
	fRaw      bool
	fRefresh  bool

	// MyName is the application name
	MyName = filepath.Base(os.Args[0])
)

func init() {
	// Flags
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fRefresh, "R", false, "Force a refresh")
	flag.BoolVar(&fDetailed, "d", false, "Get a detailed report")

	// Commands
	flag.BoolVar(&fRaw, "raw", false, "RAW JSON mode.")

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("You must give at least one site name!")
	}
}

func main() {
	var client *cryptcheck.Client

	site := flag.Arg(0)

	// -raw implies -d
	if fRaw {
		fDetailed = true
	}

	if fDebug {
		client = cryptcheck.NewClient(cryptcheck.Config{Log: 2, Refresh: fRefresh})
	} else {
		client = cryptcheck.NewClient(cryptcheck.Config{Refresh: fRefresh})
	}

	report, err := client.GetDetailedReport(site)
	if err != nil {
		log.Fatalf("impossible to get grade for '%s': %v\n", site, err)
	}

	if !fRaw {
		fmt.Printf("%s Wrapper: %s API version %s\n\n",
			MyName, cryptcheck.MyVersion, cryptcheck.APIVersion)
	}

	if fDetailed {
		// Just dump the json
		jr, _ := json.Marshal(report)
		fmt.Printf("%s\n", jr)
	} else {

		if len(report.Result.Hosts) == 0 {
			log.Fatalf("No endpoint for %s.", site)
		}

		grade := report.Result.Hosts[0].Grade.Rank
		fmt.Printf("Grade for '%s' is %s (Date: %s)\n", site, grade, report.Result.Date.Local())
	}
}
