package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

// Config describes the result of the desired yaml load.
type Config struct {
	ReportStart        time.Time
	ReportEnd          time.Time
	ReceivedIPFragment string
	ValidSenderRegexp  *regexp.Regexp
	Holidays           []Holiday
}

// String describes a Config for printing.
func (c Config) String() string {
	t := `
ReportStart        %s
ReportEnd          %s
ReceivedIPFragment %s
ValidSenderRegexp  %s
`
	s := fmt.Sprintf(t,
		time.Time(c.ReportStart).Format("2006-01-02"),
		time.Time(c.ReportEnd).Format("2006-01-02"),
		c.ReceivedIPFragment,
		c.ValidSenderRegexp,
	)
	for _, h := range c.Holidays {
		s += fmt.Sprintf("   %s\n", h)
	}
	return s
}

// UnmarshalYAML is a custom unmarshaller, which uses an auxillary
// struct (auxConfig) to deal with time.Time and regexp items.
func (c *Config) UnmarshalYAML(value *yaml.Node) error {

	tp := func(s string) (time.Time, error) { return time.Parse("2006-01-02", s) }

	type auxConfig struct {
		ReportStart          string `yaml:"reportStart"`
		reportStart          time.Time
		ReportEnd            string `yaml:"reportEnd"`
		reportEnd            time.Time
		ReceivedIPFragment   string `yaml:"receivedIPFragment"`
		ValidSenderRegexpStr string `yaml:"validSenderRegexpStr"`
		validSenderRegexp    *regexp.Regexp
		HolidayStrings       []map[string]string `yaml:"holidayStrings"`
		holidayStrings       []Holiday
	}

	var ac auxConfig
	var err error
	err = value.Decode(&ac)
	if err != nil {
		return err
	}
	ac.reportStart, err = tp(ac.ReportStart)
	if err != nil {
		return err
	}
	ac.reportEnd, err = tp(ac.ReportEnd)
	if err != nil {
		return err
	}
	if ac.ReceivedIPFragment == "" {
		return errors.New("no ip fragment found in config")
	}
	if ac.ValidSenderRegexpStr == "" {
		return errors.New("no regex string found in config")
	}
	ac.validSenderRegexp, err = regexp.Compile(ac.ValidSenderRegexpStr)
	if err != nil {
		return err
	}
	for _, h := range ac.HolidayStrings {
		if len(h) != 2 {
			return fmt.Errorf("Expected 2 date arguments, got %s", h)
		}
		start, ok := h["start"]
		if !ok {
			return fmt.Errorf("no start key for holiday %s", h)
		}
		end, ok := h["end"]
		if !ok {
			return fmt.Errorf("no start key for holiday %s", h)
		}
		startDate, err := tp(start)
		if err != nil {
			return err
		}
		endDate, err := tp(end)
		if err != nil {
			return err
		}
		if endDate.Before(startDate) {
			return fmt.Errorf("holiday %s after %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		}
		ac.holidayStrings = append(ac.holidayStrings, Holiday{startDate, endDate})
	}
	*c = Config{
		ReportStart:        ac.reportStart,
		ReportEnd:          ac.reportEnd,
		ReceivedIPFragment: ac.ReceivedIPFragment,
		ValidSenderRegexp:  ac.validSenderRegexp,
		Holidays:           ac.holidayStrings,
	}
	return nil
}

// LoadYaml loads a Config from a slice of bytes
func LoadYaml(yamlByte []byte) (Config, error) {
	var config Config
	err := yaml.Unmarshal(yamlByte, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
