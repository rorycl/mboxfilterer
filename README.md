# mboxfilterer

version 0.0.1 : 21 December 2024

Filter emails in a set of mailboxes in unix mbox format by various
criteria to summarise these emails in a csv report. 

## Overview

The programme reads mbox files concurrently, using a patched version of
`github.com/mnako/letters` to only read email headers to speed up
processing. Email bodies are ignored by this programme.

The filters are presently:

```
newFilterIP           : sent from a specified ip range fragment
newFilterByReportDate : within the report date
newFilterByHoliday    : while not on holiday
newFilterBySender     : from specified senders only
newFilterByID         : a unique message id
```

The filters are initialised using information in a yaml configuration
file, for example:

```yaml
# report start/end dates; emails with dates outside of this range will
# be excluded
reportStart: "2022-01-01"
reportEnd:   "2023-03-12"

# ip range from which senders will be included, based on the email
# "Received" header.
receivedIPFragment: "10.1.99."

# legitimate senders is a regular expression string
validSenderRegexpStr: "(?i)(this|that|another.com)"

# holidays during which emails are ignored
holidayStrings: 
  -
    start: "2022-02-20"
    end: "2022-02-22"
  -
    start: "2022-02-25"
    end: "2022-02-27"

```

## Usage

```
./mboxfilterer -h

Usage:
  mboxfilterer [OPTIONS] MboxFiles...

Application Options:
  -c, --config=    yaml configuration file (required)
  -o, --output=    optional output csv file

Help Options:
  -h, --help       Show this help message

Arguments:
  MboxFiles:       one or more mbox files to process

```

## License

This project is licensed under the [MIT Licence](LICENCE).
