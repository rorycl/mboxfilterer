package main

import (
	"fmt"
	"net/mail"
	"regexp"
	"testing"
	"time"

	"github.com/rorycl/letters"
)

func discardName(name string, result bool) bool {
	return result
}

func TestReceivedFilter(t *testing.T) {

	nf := newFilterIP("ip filter", "10.99.1.")

	tests := []struct {
		received string
		ok       bool
	}{
		{
			received: "from 23.188.159.143.dyn.plus.net ([10.99.1.23] helo=robertosmith-t470s) by smythersbrown.net with esmtpsa (TLS1.3:ECDHE_RSA_AES_256_GCM_SHA384:256) (Exim 97.65) (envelope-from <robertosmith@smythersbrown.net>) id 1lpydR-0001ll-Ru; Sun, 06 Jun 2021 19:40:37 +0000",
			ok:       true,
		},
		{
			received: "by smythersbrown.net with esmtpsa (TLS1.3:ECDHE_RSA_AES_256_GCM_SHA384:256) (Exim 97.65) (envelope-from <robertosmith@smythersbrown.net>) id 1lpydR-0001ll-Ru; Sun, 06 Jun 2021 19:40:37 +0000 from 23.188.159.143.dyn.plus.net ([10.99.1.23] helo=robertosmith-t470s)",
			ok:       true,
		},
		{
			received: "from [217.138.52.158] (helo=robertosmith-t470s) by smythersbrown.net with esmtpsa  (TLS1.3) tls TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 (Exim 4.94.2) (envelope-from <robertosmith@smythersbrown.net>) id 1o15hj-003Xom-Ud; Tue, 14 Jun 2022 12:31:31 +0000 from robertosmith by robertosmith-t470s with local (Exim 4.94.2) (envelope-from <robertosmith@robertosmith-t470s>) id 1o15hj-003dRx-5f; Tue, 14 Jun 2022 13:31:31 +0100",
			ok:       false,
		},
		{
			received: "from x.y.z.dyn.plus.net ([142.159.188.23] helo=fail) by smythersbrown.net with esmtpsa (TLS1.3:ECDHE_RSA_AES_256_GCM_SHA384:256) (Exim 97.65) (envelope-from <robertosmith@smythersbrown.net>) id 1lpydR-0001ll-Ru; Sun, 06 Jun 2021 19:40:37 +0000",
			ok:       false,
		},
		{
			received: "from 155.178.125.91.dyn.plus.net ([91.125.178.155] helo=smtpclient.apple) by smythersbrown.net with esmtpsa  (TLS1.3) tls TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 (Exim 4.96) (envelope-from <robertosmith@smythersbrown.net>) id 1snuhs-004Id5-2D; Tue, 10 Sep 2024 06:50:32 +0000",
			ok:       false,
		},
		{
			received: "from 92.40.177.66.threembb.co.uk ([92.40.177.66] helo=smtpclient.apple) by smythersbrown.net with esmtpsa  (TLS1.3) tls TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 (Exim 4.94.2) (envelope-from <robertosmith@smythersbrown.net>) id 1oY83v-00Fl82-HI for savimbi@smythersbrown.net; Tue, 13 Sep 2022 15:42:59 +0000",
			ok:       false,
		},
		{
			// studios ip range
			received: "Received: from 155.178.125.91.dyn.plus.net ([91.125.178.155] helo=smtpclient.apple) by smythersbrown.net with esmtpsa  (TLS1.2) tls TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 (Exim 4.96) (envelope-from <savimbi@smythersbrown.net>) id 1sp75q-004VlS-2e for robertosmith@smythersbrown.net; Fri, 13 Sep 2024 14:16:15 +0000",
			ok:       false,
		},
	}
	for i, tt := range tests {
		e := EmailWithSource{letters.Headers{}, "test"}
		e.ExtraHeaders = map[string][]string{}
		e.ExtraHeaders["Received"] = []string{tt.received}
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got, want := discardName(nf(e)), tt.ok; got != want {
				t.Errorf("received %s\ngot %t want %t", tt.received, got, want)
			}
		})
	}
}

func TestSenderFilter(t *testing.T) {

	nf := newFilterBySender("sender filter", regexp.MustCompile("(?i)(isadorax|robertosmith|smythersbrown)"))

	for i, tt := range []struct {
		address string
		ok      bool
	}{
		{"J.Barrage@camden.co.uk", false},
		{"a.nother.notaire@xxx-notaires.fr", false},
		{"savimbi@smythersbrown.net", true},
		{"isadoraxcampbelllange@gmail.com", true},
		{"IZASMYTHERSBROWN@GMAIL.COM", true},
		{"isadoraxclhome@gmail.com", true},
		{"isadoraxclwork@gmail.com", true},
		{"bob@smythersbrown.net", true},
		{"noxious@smythersbrown.net", true},
		{"JoannaDaly@tfl.gov.uk", false},
		{"maja@allbuildingcontrol.com", false},
		{"robertosmith@smythersbrown.net", true},
		{"boB@smythersbrown.net", true},
		{"robertosmith@codata.ltd", true},
		{"robertosmith@CODATA.ltd", true},
		{"sid@bobthebuilder.com", false},
		{"sto@smythersbrown.net", true},
	} {
		e := EmailWithSource{letters.Headers{}, "test"}
		e.From = []*mail.Address{&mail.Address{Address: tt.address}}
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got, want := discardName(nf(e)), tt.ok; got != want {
				t.Errorf("for %s got %t want %t", tt.address, got, want)
			}
		})
	}
}

func TestHolidaysFilter(t *testing.T) {
	holidayStrings := []Holiday{
		Holiday{
			time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 8, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	nf := newFilterByHoliday("holiday filter", holidayStrings)

	// don't use time.Local otherwise there might be a UTC offset
	// problem eg with British Summer Time
	beforeDate := time.Date(2021, 7, 30, 0, 0, 0, 0, time.UTC)
	onDate1 := time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC)
	okDate := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	onDate2 := time.Date(2021, 8, 3, 0, 0, 0, 0, time.UTC)
	afterDate := time.Date(2021, 8, 4, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		date time.Time
		ok   bool
	}{
		{beforeDate, true},
		{onDate1, true},
		{okDate, false},
		{onDate2, true},
		{afterDate, true},
	}

	for i, tt := range tests {
		e := EmailWithSource{letters.Headers{}, "test"}
		e.Date = tt.date
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if got, want := discardName(nf(e)), tt.ok; got != want {
				t.Errorf("got %t != want %t for %s", got, want, tt.date)
			}
		})
	}
}

func TestReportDateFilter(t *testing.T) {

	reportStart := time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC)
	reportEnd := time.Date(2021, 8, 3, 0, 0, 0, 0, time.UTC)

	nf := newFilterByReportDate("report date filter", reportStart, reportEnd)

	// don't use time.Local otherwise there might be a UTC offset
	// problem eg with British Summer Time
	beforeDate := time.Date(2021, 7, 30, 0, 0, 0, 0, time.UTC)
	okDate := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	afterDate := time.Date(2021, 8, 4, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		date time.Time
		ok   bool
	}{
		{beforeDate, false},
		{okDate, true},
		{afterDate, false},
	}

	for i, tt := range tests {
		e := EmailWithSource{letters.Headers{}, "test"}
		e.Date = tt.date
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if got, want := discardName(nf(e)), tt.ok; got != want {
				t.Errorf("got %t != want %t for %s", got, want, tt.date)
			}
		})
	}
}

func TestIDFilter(t *testing.T) {

	nf := newFilterByID("id filter")
	tests := []struct {
		id letters.MessageId
		ok bool
	}{
		{letters.MessageId("a"), true},
		{letters.MessageId("b"), true},
		{letters.MessageId("a"), false},
		{letters.MessageId("c"), true},
		{letters.MessageId("b"), false},
	}
	for i, tt := range tests {
		e := EmailWithSource{letters.Headers{}, "test"}
		e.MessageID = tt.id
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if got, want := discardName(nf(e)), tt.ok; got != want {
				t.Errorf("got %t != want %t for %s", got, want, string(tt.id))
			}
		})
	}

}
