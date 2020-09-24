package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// scanFileSystemCertificates will load the result of a filesystem scan on the
// `event` type. Information about the filesystem paths to scan are passed as
//  the second argument as a `conf` type.
//
// Returns `int` representing the amount of certificates scanned.
func scanFileSystemCertificates(e *event, cfg conf) int {
	var c []certificate
	for _, root := range cfg.Sources.Filesystem.ScanPaths {
		filepath.Walk(root, visit(&c))
	}
	(*e).Sources.Certificates = append((*e).Sources.Certificates, c...)
	return len(c)
}

// visit is executed for each file scanned, it filters out everything except for
// .pem, .crt and .cer files.
func visit(certs *[]certificate) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)

		switch ext {
		case ".pem", ".crt", ".cer":
			appendX509EncodedFile(path, certs)
		}
		return nil
	}
}

// getFirstItem receives an string array and return the first item, if the array
// is defined to empty, it returns an empty string.
func getFirstItem(item []string) string {
	if len(item) > 0 {
		return item[0]
	}
	return ""
}

// appendX509EncodedFile receives a path of type `string`, opens, reads,
// decodes (pem), and parses it into a `crypto/x509.Certificate` data structure.
// It then calculates the sha265 hash that structure, and finally it loads the
// certificate into a `certs` (second argument).
func appendX509EncodedFile(path string, certs *[]certificate) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return errors.New("failed to open file")
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return errors.New("failed to get file stats")
	}

	data := make([]byte, fInfo.Size())
	f.Read(data)

	p, _ := pem.Decode(data)
	if p == nil {
		return errors.New("no PEM data found")
	}
	c, err := x509.ParseCertificate(p.Bytes)

	if err != nil {
		return errors.New("error parsing certificate")
	}

	h := sha256.New()
	h.Write(c.Raw)

	if !c.IsCA {
		*certs = append(*certs,
			certificate{
				Path:   path,
				Source: "filesystem",
				Validity: validity{
					ValidUntil: c.NotAfter,
				},
				Issuer: issuer{
					Country:      getFirstItem(c.Issuer.Country),
					Organization: getFirstItem(c.Issuer.Organization),
					CommonName:   c.Issuer.CommonName,
				},
				Fingerprint: fmt.Sprintf("%x", h.Sum(nil)),
				Subject: subject{
					Country:            getFirstItem(c.Subject.Country),
					Organization:       getFirstItem(c.Subject.Organization),
					OrganizationalUnit: getFirstItem(c.Subject.OrganizationalUnit),
					Locality:           getFirstItem(c.Subject.Locality),
					Province:           getFirstItem(c.Subject.Province),
					StreetAddress:      getFirstItem(c.Subject.StreetAddress),
					PostalCode:         getFirstItem(c.Subject.PostalCode),
					SerialNumber:       c.Subject.SerialNumber,
					CommonName:         c.Subject.CommonName,
				},
			})
	}

	return nil
}
