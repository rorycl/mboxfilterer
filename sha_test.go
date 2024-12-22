package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSHA256(t *testing.T) {
	j := strings.NewReader("Foo")
	s, err := SHA256(j)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := fmt.Sprintf("%x", s), "1cbec737f863e4922cee63cc2ebbfaafcd1cff8b790d8cfd2e6a5d550b648afa"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestFileCalcSHA(t *testing.T) {
	f := "testdata/golang.mbox"
	s, err := fileCalcSHA(f)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := s, "8716f9d26984405c068dc64427c43c021bfd7b3c7ac4129b338bbd2a1e33a703"; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestSHAFiles(t *testing.T) {
	want := `test
file                     : testdata/gonuts.mbox
sha256sum                : 683606049b9c7a4154f2a57a61ad1934f37a82df461c4d0e008f864cd15b25e4
file                     : testdata/golang.mbox
sha256sum                : 8716f9d26984405c068dc64427c43c021bfd7b3c7ac4129b338bbd2a1e33a703
`

	got, err := sha256Summarize("test", "testdata/golang.mbox", "testdata/gonuts.mbox")
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("values are not the same %s", cmp.Diff(got, want))
	}
	// fmt.Println(summary)
}
