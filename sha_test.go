package main

import (
	"fmt"
	"strings"
	"testing"
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
