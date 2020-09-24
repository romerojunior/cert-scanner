package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// A f5CertResponse represents the unchanged REST API response (`JSON` formated)
// from the F5 endpoint, it contains all certificates under `Items` property.
type f5CertResponse struct {
	Kind     string          `json:"kind"`
	SelfLink string          `json:"selfLink"`
	Items    []f5Certificate `json:"items"`
}

// A f5Certificate represents an x509 certificate as reported by F5 REST API.
type f5Certificate struct {
	Kind                    string                    `json:"kind"`
	Name                    string                    `json:"name"`
	FullPath                string                    `json:"fullPath"`
	Generation              int                       `json:"generation"`
	SelfLink                string                    `json:"selfLink"`
	APIRawValues            f5APIRawValues            `json:"apiRawValues"`
	CommonName              string                    `json:"commonName"`
	Country                 string                    `json:"country"`
	Fingerprint             string                    `json:"fingerprint"`
	Organization            string                    `json:"organization"`
	OU                      string                    `json:"ou"`
	SubjectAlternativeName  string                    `json:"subjectAlternativeName"`
	CertValidatorsReference f5CertValidatorsReference `json:"certValidatorsReference"`
}

type f5APIRawValues struct {
	CertificateKeySize string `json:"certificateKeySize"`
	Expiration         string `json:"expiration"`
	Issuer             string `json:"issuer"`
	PublicKeyType      string `json:"publicKeyType"`
}

type f5CertValidatorsReference struct {
	Link            string `json:"link"`
	IsSubcollection bool   `json:"isSubcollection"`
}

// ParsedFingerprint returns a parsed fingerprint value of type `string`, it is
// consistent with the definition of fingerprint (sha256, lowercase, no special
// characteres) across the entire application.
func (e f5Certificate) ParsedFingerprint() string {
	var fp string
	fp = strings.ReplaceAll(e.Fingerprint, "SHA256/", "")
	fp = strings.ReplaceAll(fp, ":", "")
	fp = strings.ToLower(fp)
	return fp
}

// ParsedExpiration returns a `time.Time` type representing the expiration date
// of the certificate.
func (e f5APIRawValues) ParsedExpiration() time.Time {
	layout := "Jan _2 15:04:05 2006 MST"
	t, _ := time.Parse(layout, e.Expiration)
	return t
}

// scanF5Certificates will load the result of a scan on a F5 endpoint inside
// `event` type. Information about the F5 endpoint such as password and URL are
// passed as the second argument as a `conf` type.
//
// Returns `int` representing the amount of certificates scanned.
func scanF5Certificates(e *event, cfg conf) int {
	if cfg.Sources.F5.URL == "" {
		return 0
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	p := "/mgmt/tm/sys/crypto/cert"

	req, _ := http.NewRequest(http.MethodGet, cfg.Sources.F5.URL+p, nil)
	req.SetBasicAuth(cfg.Sources.F5.User, cfg.Sources.F5.Password)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("error requesting f5: ", err)
	}

	var certs f5CertResponse
	err = json.NewDecoder(res.Body).Decode(&certs)
	if err != nil {
		log.Fatal("error decoding f5 api response body: ", err)
	}

	var c []certificate

	for _, cert := range certs.Items {

		c = append(c,
			certificate{
				Path:   cert.FullPath,
				Source: "f5",
				Validity: validity{
					ValidUntil: cert.APIRawValues.ParsedExpiration(),
				},
				Issuer: issuer{
					CommonName: cert.APIRawValues.Issuer,
				},
				Fingerprint: cert.ParsedFingerprint(),
				Subject: subject{
					Country:            cert.Country,
					Organization:       cert.Organization,
					OrganizationalUnit: cert.OU,
					CommonName:         cert.CommonName,
				},
			})
	}

	(*e).Sources.Certificates = append((*e).Sources.Certificates, c...)

	return len(c)
}
