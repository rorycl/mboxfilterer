package main

import (
	"fmt"
	"testing"
)

func TestConfigFail(t *testing.T) {
	yaml := []byte(`
reportStart: "2022-01-01"
reportEnd:   "2023-03-12"
# receivedIPFragment: "10.1.99."
validSenderRegexpStr: "(?i)(this|that|another.com)"
holidayStrings: 
  -
    start: "2022-02-20"
    end: "2022-02-22"
  -
    start: "2022-02-25"
    end: "2022-02-27"
`)

	_, err := LoadYaml(yaml)
	if err == nil {
		t.Fatalf("expected receivedIPFragment error")
	}
	fmt.Println(err)
}

func TestConfigFail2(t *testing.T) {
	yaml := []byte(`
reportStart: "2022-01-01"
reportEnd:   "2023-03-12"
receivedIPFragment: "10.1.99."
validSenderRegexpStr: "(?i)(this|that|another.com)"
holidayStrings: 
  -
    start: "2022-02-23"
    end: "2022-02-22"
  -
    start: "2022-02-25"
    end: "2022-02-27"
`)

	_, err := LoadYaml(yaml)
	if err == nil {
		t.Fatalf("expected holiday start > end error")
	}
	fmt.Println(err)
}

func TestConfigOK(t *testing.T) {
	yaml := []byte(`
reportStart: "2022-01-01"
reportEnd:   "2023-03-12"
receivedIPFragment: "10.1.99."
validSenderRegexpStr: "(?i)(this|that|another.com)"
holidayStrings: 
  -
    start: "2022-02-21"
    end: "2022-02-22"
  -
    start: "2022-02-25"
    end: "2022-02-27"
`)

	config, err := LoadYaml(yaml)
	if err != nil {
		t.Fatalf("got unexpected error %s", err)
	}
	if got, want := config.ReceivedIPFragment, "10.1.99."; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
