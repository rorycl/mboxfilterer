package main

import (
	"strings"

	"github.com/rorycl/letters"
)

// EmailWithSource represents the headers of an email with their mbox
// source
type EmailWithSource struct {
	letters.Headers
	source string // source mbox
}

var csvHeader = []string{"date", "from", "subj", "source", "id", "received"}

func (e EmailWithSource) subj(n int) string {
	if n == 0 || len(e.Subject) < n {
		return e.Subject
	}
	return e.Subject[:n]
}

func (e EmailWithSource) forCSV(subjLen int) []string {
	dater := e.Date.Format("2006-01-02")
	return []string{
		dater,
		e.From[0].Address,
		e.subj(subjLen),
		e.source,
		string(e.MessageID),
		strings.Join(e.ExtraHeaders["Received"], " "),
	}
}
