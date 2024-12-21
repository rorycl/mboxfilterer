package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/rorycl/letters"
)

type filterFunc func(EmailWithSource) (string, bool)

// Filters contains a set of filters for excluding email headers from
// consideration together with stats on both ok emails and those that
// have been filtered out by any filter (identified by name).
type Filters struct {
	stats      map[string]int
	filters    []filterFunc
	filterChan chan string
	start, end time.Time // processing time
}

func NewFilters(funcs ...filterFunc) *Filters {
	f := Filters{
		filters:    funcs,
		stats:      map[string]int{},
		start:      time.Now(),
		filterChan: make(chan string),
	}
	f.stats["ok"] = 0

	// start running stats collector
	go func() {
		for name := range f.filterChan {
			f.stats[name]++
		}
	}()
	return &f
}

// Filter filters an EmailWithSource through each filter exiting on
// first false match or falling through to "ok". This function is is
// designed for concurrent access.
func (f *Filters) Filter(e EmailWithSource) bool {
	for _, fn := range f.filters {
		if name, ok := fn(e); !ok {
			f.filterChan <- name
			return false
		}
	}
	f.filterChan <- "ok"
	return true
}

// Stats shows how many times particular filters or the fallthrough "ok"
// condition have been called during processing of emails in mboxes.
func (f *Filters) Stats() string {
	f.end = time.Now()
	d := f.end.Sub(f.start)
	dRound := d.Round(time.Millisecond * 10)

	tpl := "%-25s: %4d\n"
	t := fmt.Sprintf("%-25s: %4s\n", "time processing", dRound)
	t += fmt.Sprintf(tpl, "OK", f.stats["ok"])
	t += "skipped\n"
	names := []string{}
	for k := range f.stats {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		if k == "ok" {
			continue
		}
		v := f.stats[k]
		t += fmt.Sprintf(tpl, k, v)
	}
	return t
}

// newFilterByID filters out emails not matching the provided IP fragment
func newFilterIP(name string, ipFragment string) filterFunc {
	return func(e EmailWithSource) (string, bool) {
		return name, strings.Contains(
			strings.Join(e.Headers.ExtraHeaders["Received"], " "),
			ipFragment,
		)
	}
}

// newFilterByID filters out emails outside of the report period
func newFilterByReportDate(name string, reportStart, reportEnd time.Time) filterFunc {
	return func(e EmailWithSource) (string, bool) {
		return name, e.Date.After(reportStart) && e.Date.Before(reportEnd)
	}
}

// Holiday is a holiday type for use in config and newFilterByHoliday
type Holiday struct {
	Start time.Time
	End   time.Time
}

func (h Holiday) String() string {
	return fmt.Sprintf("%s : %s", h.Start.Format("2006-01-02"), h.End.Format("2006-01-02"))
}

// newFilterByID filters out emails during holiday periods
func newFilterByHoliday(name string, holidays []Holiday) filterFunc {
	return func(e EmailWithSource) (string, bool) {
		for _, h := range holidays {
			if e.Date.After(h.Start) && e.Date.Before(h.End) {
				return name, false
			}
		}
		return name, true
	}
}

// newFilterBySender filters out emails from non matching senders
func newFilterBySender(name string, validSenderRegexp *regexp.Regexp) filterFunc {
	return func(e EmailWithSource) (string, bool) {
		return name, validSenderRegexp.MatchString(e.From[0].Address)
	}
}

// newFilterByID filters out emails with duplicate IDs
func newFilterByID(name string) filterFunc {
	idHash := map[letters.MessageId]struct{}{}
	return func(e EmailWithSource) (string, bool) {
		if _, ok := idHash[e.MessageID]; ok {
			return name, false
		}
		idHash[e.MessageID] = struct{}{}
		return name, true
	}
}
