package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"strings"
	"text/scanner"

	log "github.com/sirupsen/logrus"
)

type workSpec struct {
	StoreID      string
	Transactions int
}

var (
	numberOfTransactions int
	controllerEndpoint   string
	storeNumbers         string
	dryRun               bool
)

func main() {
	flag.IntVar(&numberOfTransactions,
		"number-of-transactions",
		10,
		"The number of transactions to send per-store")
	flag.StringVar(&controllerEndpoint,
		"controller-endpoint",
		"https://n2mmzmrjb3.execute-api.ap-southeast-2.amazonaws.com/AddWorkload",
		"HTTP URI of the controller")
	flag.StringVar(&storeNumbers,
		"store-numbers",
		"399, 303",
		"A comma separated list of store numbers")
	flag.BoolVar(&dryRun,
		"dry-run",
		true,
		"If true no data is sent. Set with '--dryrun=true|false'")

	flag.Parse()
	log.Infof("Controller: dry-run = %v", dryRun)

	// Create client
	client := &http.Client{}

	var s scanner.Scanner
	s.Init(strings.NewReader(storeNumbers))
	s.Whitespace = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ' | 1<<','

	storeIDs := []string{}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		storeIDs = append(storeIDs, s.TokenText())
	}

	workSpecs := []*workSpec{}

	for _, storeID := range storeIDs {
		workSpecs = append(workSpecs, &workSpec{
			StoreID:      storeID,
			Transactions: numberOfTransactions,
		})
	}

	for _, ws := range workSpecs {
		jsonStr, jsonErr := json.Marshal(ws)
		if jsonErr != nil {
			log.WithError(jsonErr).Errorf("failed to marshal JSON: %v", ws)
			continue
		}
		req, err := http.NewRequest(http.MethodPost, controllerEndpoint, bytes.NewBuffer(jsonStr))
		if err != nil {
			log.WithError(err).Errorf("request failed: '%+v'", ws)
			continue
		}
		// Headers
		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		log.Infof("send request: '%v'", string(jsonStr))

		if dryRun {
			continue
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			log.WithError(err).Errorf("Failure : ")
			continue
		}

		// Display Result(s)
		log.Infof("response Status : %s", resp.Status)
	}
}
