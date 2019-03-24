package structs

import (
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
func TrailingEvents(logEvents []LogEvent, lastSeconds float64) []LogEvent {
	now := time.Now()

	return Filter(logEvents, func(v LogEvent) bool {
		diff := now.Sub(v.Date)
		seconds := diff.Seconds()

		return (seconds <= lastSeconds)
	})
}

// ParseLogEvent takes the log string and returns a LogEvent struct
//  A log line is of the format:
// 127.0.0.1 - frank [23/Mar/2019:18:44:53 +0000] "DELETE /config/update HTTP/1.0" 401 491
func ParseLogEvent(line string) LogEvent {
	re, _ := regexp.Compile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - (.*) \[(.*)\] \"((.*) (\/.*) .*)\" (\d{3}) (\d*)$`)
	result := re.FindStringSubmatch(line)
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
	}
}
