package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

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
