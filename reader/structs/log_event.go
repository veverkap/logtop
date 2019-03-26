package structs

import (
	"errors"
	"regexp"
	"sort"
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

// LogEventBy is the type of a "less" function that defines the ordering of its arguments.
type LogEventBy func(p1, p2 *LogEvent) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by LogEventBy) Sort(events []LogEvent) {
	es := &eventSorter{
		events: events,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(es)
}

type eventSorter struct {
	events []LogEvent
	by     func(p1, p2 *LogEvent) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *eventSorter) Len() int {
	return len(s.events)
}

// Swap is part of sort.Interface.
func (s *eventSorter) Swap(i, j int) {
	s.events[i], s.events[j] = s.events[j], s.events[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *eventSorter) Less(i, j int) bool {
	return s.by(&s.events[i], &s.events[j])
}

// SortBySection returns the log events sorted by section
func SortBySection(events []LogEvent) []LogEvent {
	section := func(p1, p2 *LogEvent) bool {
		return p1.Section < p2.Section
	}

	LogEventBy(section).Sort(events)
	return events
}

// SortByDateAsc returns the log events sorted by date
func SortByDateAsc(events []LogEvent) []LogEvent {
	date := func(p1, p2 *LogEvent) bool {
		return p1.Date.Before(p2.Date)
	}

	LogEventBy(date).Sort(events)
	return events
}

// SortByDateDesc returns the log events sorted by date
func SortByDateDesc(events []LogEvent) []LogEvent {
	date := func(p1, p2 *LogEvent) bool {
		return p1.Date.After(p2.Date)
	}

	LogEventBy(date).Sort(events)
	return events
}

// Filter returns the slice of LogEvents matching filter
func Filter(vs []LogEvent, f func(LogEvent) bool) []LogEvent {
	vsf := make([]LogEvent, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// TrailingEvents returns the events the occurred in the lastSeconds
func TrailingEvents(logEvents []LogEvent, lastSeconds int64) []LogEvent {
	now := time.Now()

	events := Filter(logEvents, func(v LogEvent) bool {
		return (int64(now.Sub(v.Date).Seconds()) <= lastSeconds)
	})

	return events
}

// ParseLogEvent takes the log string and returns a LogEvent struct
//  A log line is of the format:
// 127.0.0.1 - frank [23/Mar/2019:18:44:53 +0000] "DELETE /config/update HTTP/1.0" 401 491
func ParseLogEvent(line string) (LogEvent, error) {
	if line == "" {
		return LogEvent{}, errors.New("Empty String")
	}
	line = strings.ReplaceAll(line, "\n", "")
	re, _ := regexp.Compile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - (.*) \[(.*)\] \"((.*) (\/.*) .*)\" (\d{3}) (\d*)$`)
	result := re.FindStringSubmatch(line)
	if len(result) == 9 {
		host := result[1]
		user := result[2]
		dateString := result[3]
		const longForm = "02/Jan/2006:15:04:05 -0700"
		date, _ := time.Parse(longForm, dateString)
		verb := result[5]
		path := result[6]
		section := path

		pieces := strings.Split(path, "/")
		if len(pieces) > 2 {
			section = "/" + pieces[1]
		}

		status, _ := strconv.Atoi(result[7])
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
			Error:      status != 200,
		}, nil
	}
	return LogEvent{}, errors.New("Bad regex")
}

// GroupBySection returns the slice of LogEvents matching filter
func GroupBySection(vs []LogEvent) []SectionDetail {
	vsf := make([]SectionDetail, 0)
	for _, v := range vs {
		index := findSectionDetail(vsf, v.Section)

		if index >= 0 {
			errors := vsf[index].Errors
			if v.Error {
				errors++
			}
			events := append(vsf[index].Events, v)

			vsf[index] = SectionDetail{
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
			vsf = append(vsf, SectionDetail{
				Section: v.Section,
				Events:  events,
				Hits:    1,
				Errors:  errors,
			})
		}
	}
	return vsf
}
