package structs

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LogEvent represents a line of the log file
type LogEvent struct {
	Host       string
	User       string
	Date       time.Time
	Verb       string
	Section    string
	Path       string
	StatusCode int
	ByteSize   int
	Error      bool
}

// findSectionDetail finds the index of the matching section (used by GroupBySection)
func findSectionDetail(details []SectionDetail, section string) int {
	for i, detail := range details {
		if detail.Section == section {
			return i
		}
	}
	return -1
}

/*
TrailingEvents iterates through all of the logEvents appending any that occurred less than
lastSeconds seconds ago to the filteredEvents and then returns filteredEvents
*/
func TrailingEvents(logEvents []LogEvent, lastSeconds int64) []LogEvent {
	now := time.Now()
	filteredEvents := make([]LogEvent, 0)

	for _, event := range logEvents {
		if int64(now.Sub(event.Date).Seconds()) <= lastSeconds {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents
}

/*
ParseLogEvent takes the log string and returns a LogEvent

A log line is of the format:
127.0.0.1 - frank [23/Mar/2019:18:44:53 +0000] "DELETE /config/update HTTP/1.0" 401 491
*/
func ParseLogEvent(line string) (LogEvent, error) {
	// if we get a blank line, we return an empty LogEvent and an error
	if line == "" {
		return LogEvent{}, errors.New("Empty String")
	}

	// double check that we don't have any newlines (tail *should* help us with this)
	line = strings.ReplaceAll(line, "\n", "")

	//	The heart of the program - loads up a big regex to match on the log line and capture necessary tokens
	re, _ := regexp.Compile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - (.*) \[(.*)\] \"((.*) (\/.*) .*)\" (\d{3}) (\d*)$`)
	result := re.FindStringSubmatch(line)

	// We have 9 capture places, so we have to get that many back
	if len(result) == 9 {
		host := result[1]
		user := result[2]

		// parse the date back from the log file format
		dateString := result[3]
		const longForm = "02/Jan/2006:15:04:05 -0700"

		/*
			We are swallowing this error.  if the log has a date that doesn't match, it *shouldn't* get through the regex,
			but if it does, we will blow up here
		*/
		date, _ := time.Parse(longForm, dateString)
		verb := result[5]

		// this comes in as something like /path or /section/path so we split and try to get the pieces separately
		path := result[6]
		section := path
		pieces := strings.Split(path, "/")
		if len(pieces) > 2 {
			section = "/" + pieces[1]
		}

		// we consider it an error if it is not informational or success https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
		status, _ := strconv.Atoi(result[7])

		// convert string to integer
		size, _ := strconv.Atoi(result[8])

		return LogEvent{
			Verb:       verb,
			Host:       host,
			User:       user,
			Date:       date,
			Section:    section,
			Path:       path,
			StatusCode: status,
			ByteSize:   size,
			Error:      status >= 400,
		}, nil
	}
	return LogEvent{}, errors.New("Bad regex")
}

// GroupBySection iterates through the logEvents generating a slice of SectionDetails grouped by section
func GroupBySection(logEvents []LogEvent) []SectionDetail {
	groupedDetails := make([]SectionDetail, 0)
	for _, v := range logEvents {
		index := findSectionDetail(groupedDetails, v.Section)

		if index >= 0 {
			errors := groupedDetails[index].Errors
			if v.Error {
				errors++
			}
			events := append(groupedDetails[index].Events, v)

			groupedDetails[index] = SectionDetail{
				Section: v.Section,
				Events:  events,
				Hits:    len(events),
				Errors:  errors,
			}
		} else {
			errors := 0
			if v.Error {
				errors++
			}
			events := append(make([]LogEvent, 0), v)
			groupedDetails = append(groupedDetails, SectionDetail{
				Section: v.Section,
				Events:  events,
				Hits:    1,
				Errors:  errors,
			})
		}
	}
	return groupedDetails
}
