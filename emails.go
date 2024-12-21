package main

import (
	"encoding/csv"
	"fmt"
	"sort"
)

// Emails are a collection of email headers with their mbox sources
type Emails []EmailWithSource

// NewEmails is a constructor for Emails
func NewEmails() Emails {
	return []EmailWithSource{}
}

// Add adds an email to the emails slice
func (emails *Emails) Add(e EmailWithSource) {
	*emails = append(*emails, e)
}

// Write writes out the emails to a csv.Writer, using a maximum subject
// length of subjectLen (use 0 to write the whole subject)
func (e Emails) Write(writer *csv.Writer, subjectLen int) error {
	sort.Slice(e,
		func(i, j int) bool {
			return e[i].Date.Before(e[j].Date)
		})

	// write out
	if err := writer.Write(csvHeader); err != nil {
		return fmt.Errorf("csv header writing error, %w", err)
	}
	for _, em := range e {
		if err := writer.Write(em.forCSV(subjectLen)); err != nil {
			return fmt.Errorf("csv writing error, %w", err)
		}
	}
	writer.Flush()
	return nil
}
