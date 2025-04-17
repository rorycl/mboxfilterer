package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/attilabuti/mbox"
	"github.com/rorycl/letters"
	"github.com/rorycl/letters/parser"
)

// process processes email mbox files, processing each file
// concurrently, reading each email by email, putting emails on an email
// chan and errors on an error chan. Processing should stop on first
// error.
func process(filers []string, filters *Filters) (<-chan EmailWithSource, <-chan error) {

	done := make(chan struct{})
	emailChan := make(chan EmailWithSource)
	errorChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(len(filers))

	for _, filer := range filers {
		go func() {
			defer wg.Done()
			f, err := os.Open(filer)
			if err != nil {
				errorChan <- fmt.Errorf("file opening error, %w", err)
				done <- struct{}{} // stop further processing
				return
			}
			defer f.Close()

			mboxReader := mbox.NewReader(f)

			for {
				// stop processing early on done signal, to stop
				// processing of concurrent mboxes after first error
				select {
				case <-done:
					return
				default:
				}
				msg, err := mboxReader.NextMessage()
				if err == io.EOF {
					return
				}
				if err != nil {
					errorChan <- fmt.Errorf("mboxReader NextMessage error for %s, %w", filer, err)
					done <- struct{}{} // stop further processing
					return
				}

				p := letters.NewParser(parser.WithHeadersOnly())
				message, err := p.Parse(msg)
				if err != nil {
					errorChan <- fmt.Errorf("letters parsing error for %s, %w", filer, err)
					done <- struct{}{} // stop further processing
					return
				}

				es := EmailWithSource{message.Headers, filer}

				// continue if any filters return false
				if !filters.Filter(es) {
					continue
				}

				// put email on email channel
				emailChan <- es
			}
		}()
	}
	go func() {
		wg.Wait()
		close(emailChan)
		close(errorChan)
	}()
	return emailChan, errorChan
}
