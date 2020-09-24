package main

import "time"

type certificate struct {
	Path        string   `json:"path"`
	Source      string   `json:"source"`
	Validity    validity `json:"validity"`
	Issuer      issuer   `json:"issuer"`
	Subject     subject  `json:"subject"`
	Fingerprint string   `json:"fingerprint"`
}

type issuer struct {
	Country      string `json:"country,omitempty"`
	Organization string `json:"organization,omitempty"`
	CommonName   string `json:"commonName,omitempty"`
}

type subject struct {
	Country            string `json:"country,omitempty"`
	Organization       string `json:"organization,omitempty"`
	OrganizationalUnit string `json:"organizationalUnit,omitempty"`
	Locality           string `json:"locality,omitempty"`
	Province           string `json:"province,omitempty"`
	StreetAddress      string `json:"streetAddress,omitempty"`
	PostalCode         string `json:"postalCode,omitempty"`
	SerialNumber       string `json:"serialNumber,omitempty"`
	CommonName         string `json:"commonName,omitempty"`
}

type validity struct {
	ValidUntil time.Time `json:"validUntil"`
}
