/*
mboxfilterer

This programme outputs a concise csv summarising unique emails in one or more mbox files which pass filtering.

These filters currently are:
newFilterIP           : ensuring the sender is from a specified ip range
newFilterByReportDate : within the report date
newFilterByHoliday    : while not on holiday
newFilterBySender     : from specified senders only
newFilterByID         : a unique message id

RCL 20 December 2024
*/
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"

	flags "github.com/jessevdk/go-flags"
)

// Options are flags options
type Options struct {
	// Verbose  bool `short:"v" long:"verbose"  description:"show verbose output\nthis presently does not do much"`
	Config  string `short:"c" long:"config" description:"yaml configuration file (required)" required:"yes"`
	Output  string `short:"o" long:"output" description:"optional output csv file"`
	shasums string
	Args    struct {
		MboxFiles []string `description:"one or more mbox files to process"`
	} `positional-args:"yes" required:"yes"`
}

// checkFileExists checks the existance of a file
func checkFileExists(file string) error {
	if _, err := os.Stat(file); errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

// makeOutputFile makes an output file for writing
func makeOutputFile(s string) (*os.File, error) {
	var fileName string
	if s != "" {
		if err := checkFileExists(s); err == nil {
			return nil, fmt.Errorf("file %s already exists", s)
		}
		fileName = s
		return os.Create(fileName)
	}
	timestamp := time.Now()
	fileName = timestamp.Format("20060102-150405.csv")
	return os.Create(fileName)
}

func main() {

	// process arguments
	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	// check input mbox files exists; better to error early
	for _, f := range options.Args.MboxFiles {
		err := checkFileExists(f)
		if err != nil {
			os.Exit(1)
		}
	}

	// start processing shasums for input files
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		options.shasums, err = sha256Summarize("input files", options.Args.MboxFiles...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// load configuration
	filer, err := os.ReadFile(options.Config)
	if err != nil {
		fmt.Printf("could not load file: %s", err)
		os.Exit(1)
	}
	config, err := LoadYaml(filer)
	if err != nil {
		fmt.Println("yaml loading error", err)
		os.Exit(1)
	}

	// initialise report output files
	wfile, err := makeOutputFile(options.Output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	writer := csv.NewWriter(wfile)

	// init Emails container
	emails := NewEmails()

	// init Filters
	filters := NewFilters(
		newFilterIP("ip invalid", config.ReceivedIPFragment),
		newFilterByReportDate("outside daterange", config.ReportStart, config.ReportEnd),
		newFilterByHoliday("on holiday", config.Holidays),
		newFilterBySender("invalid sender", config.ValidSenderRegexp),
		newFilterByID("duplicate id"),
	)

	// process files
	emailChan, errorChan := process(options.Args.MboxFiles, filters)

	// drain the error chan, exiting on first error
	go func() {
		for err := range errorChan {
			if err != nil {
				fmt.Println(err)
				fmt.Println("exiting...")
				os.Exit(1)
			}
		}
	}()

	// add emails
	for e := range emailChan {
		emails.Add(e)
	}

	// write out emails with a subject max length of 10 chars
	err = emails.Write(writer, 10)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// show stats
	fmt.Print(filters.Stats())

	// show parameters
	fmt.Println(config)

	// show input files processed in goroutine above
	wg.Wait()
	fmt.Println(options.shasums)

	// show output file sha256
	fileName := wfile.Name()
	wfile.Close()
	outputFileSummary, err := sha256Summarize("output file", fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(outputFileSummary)

}
