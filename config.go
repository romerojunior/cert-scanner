package main

import (
	"encoding/json"
	"log"
	"os"
)

// A `conf` represents the contents of the decoded `JSON` configuration file.
type conf struct {
	Sources     confSources
	Destination confDestination
}

type confSources struct {
	Filesystem confSourceFilesystem
	F5         confSourceF5
}

type confDestination struct {
	Aws confDestinationAws
}

type confSourceFilesystem struct {
	ScanPaths []string
}

type confSourceF5 struct {
	User     string
	Password string
	URL      string
	Port     int
}

type confDestinationAws struct {
	Profile string
	URL     string
}

// loadConfig reads a configuration file (declared in `JSON` format) and decode
// its parameters intro a `conf` data type using a `JSON` decoder, it exits in
// case of failure while reading or decoding the configuration file.
func loadConfig(filepath string, c *conf) {
	cf, err := os.Open(filepath)
	if err != nil {
		log.Fatal("error loading config: ", err)
	}
	defer cf.Close()

	d := json.NewDecoder(cf)

	err = d.Decode(c)
	if err != nil {
		log.Fatal("error decoding json config: ", err)
	}
}
