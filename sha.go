package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
)

func sha256Summarize(label string, files ...string) (string, error) {
	output := ""
	if label != "" {
		output += fmt.Sprintf("%s\n", label)
	}

	type sumErr struct {
		sum string
		err error
	}
	results := make(chan sumErr)
	var err error

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, f := range files {
		go func(ff string) {
			defer wg.Done()
			sum, err := fileCalcSHA(ff)
			if err != nil {
				results <- sumErr{err: err}
				return
			}
			s := ""
			s += fmt.Sprintf("file                     : %s\n", f)
			s += fmt.Sprintf("sha256sum                : %s\n", sum)
			results <- sumErr{s, nil}
		}(f)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	for r := range results {
		if r.err != nil {
			err = r.err
			break
		}
		output += r.sum
	}
	return output, err
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
