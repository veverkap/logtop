package structs

import (
	"regexp"
	"strconv"
	"strings"
)

// LogEvent represents a line of the log file
type LogEvent struct {
	Value      string
	Host       string
	User       string
	Date       string
	Verb       string
	Section    string
	Path       string
	StatusCode int
	ByteSize   int
}

// A log line is of the format:
// 127.0.0.1 - frank [23/Mar/2019:18:44:53 +0000] "DELETE /config/update HTTP/1.0" 401 491
func parseLogEvent(line string) LogEvent {
	re, _ := regexp.Compile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - (.*) \[(.*)\] \"((.*) (\/.*) .*)\" (\d{3}) (\d*)$`)
	result := re.FindStringSubmatch(line)
	host := result[1]
	user := result[2]
	date := result[3]
	verb := result[5]
	path := result[6]
	section := path
	pieces := strings.Split(path, "/")
	if len(pieces) > 2 {
		section = "/" + pieces[1]
	}
	print(section)

	status, _ := strconv.Atoi(result[7])
	size, _ := strconv.Atoi(result[8])

	return LogEvent{
		Value:      line,
		Host:       host,
		User:       user,
		Date:       date,
		Verb:       verb,
		Section:    section,
		Path:       path,
		StatusCode: status,
		ByteSize:   size,
	}
}
