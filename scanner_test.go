package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFileSystemCertificates(t *testing.T) {
	var e event
	c := conf{Sources: confSources{Filesystem: confSourceFilesystem{ScanPaths: []string{"."}}}}
	i := scanFileSystemCertificates(&e, c)
	if i < 1 {
		t.Errorf("expected at least 1 certificate but got %v", i)
	}
}

func TestAppendX509EncodedFile(t *testing.T) {
	var certs []certificate
	filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			switch ext {
			case ".pem", ".crt":
				err = appendX509EncodedFile(path, &certs)
			}
			if err != nil {
				t.Errorf("error appending certificate file: %v", err)
			}
			return nil
		},
	)
	if len(certs) < 1 {
		t.Errorf("expected at least 1 certificate but got %v", len(certs))
	}
}

func TestStringfiedJSON(t *testing.T) {
	var e event
	c := conf{Sources: confSources{Filesystem: confSourceFilesystem{ScanPaths: []string{"."}}}}
	i := scanFileSystemCertificates(&e, c)
	if i < 1 {
		t.Errorf("expected at least 1 certificate but got %v", i)
	}
	_, err := e.stringfiedJSON()
	if err != nil {
		t.Errorf("could not stringify json")
	}
}
