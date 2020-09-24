package main

import (
	"log"
	"os"
)

const (
	version = "1.1.0"
)

// scan will try to fetch the hostname of the node running the agent.
func getHostname() string {
	h, err := os.Hostname()
	if err != nil {
		log.Panic("error acquiring hostname: ", err)
	}
	return h
}

// scan will start a scanning event, and go though all sources, and at
// completion it will send out the event to the configured destination.
func scan(f string) {
	var e event
	var c conf

	loadConfig(f, &c)
	scanFileSystemCertificates(&e, c)
	scanF5Certificates(&e, c)

	e.SenderInfo.Hostname = getHostname()

	sendEvent(e, c)

	return
}
