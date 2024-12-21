package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func sha256Summarize(label string, files ...string) (string, error) {
	output := ""
	if label != "" {
		output += fmt.Sprintf("%s\n", label)
	}
	for _, f := range files {
		s, err := fileCalcSHA(f)
		if err != nil {
			return "", err
		}
		output += fmt.Sprintf("file                     : %s\n", f)
		output += fmt.Sprintf("sha256sum                : %s\n", s)
	}
	return output, nil
}

func fileCalcSHA(f string) (string, error) {
	o, err := os.Open(f)
	if err != nil {
		return "", err
	}
	defer o.Close()
	s, err := SHA256(o)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", s), nil
}

func SHA256(r io.Reader) ([]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return []byte{}, err
	}
	return hash.Sum(nil), nil
}
