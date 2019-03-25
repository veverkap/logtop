package helpers

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	"github.com/veverkap/logtop/reader/structs"
)

// AccessLog represents the default AccessLog location
var AccessLog = "/tmp/access.log"

// LogEvents represents a slice of events
var LogEvents = make([]structs.LogEvent, 0)

var previousOffset int64

func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// LoadExistingLogFile will return the structs
func LoadExistingLogFile() {
	file, err := os.Open(AccessLog)
	if err != nil {
		log.Fatalf("Could not load file %s", AccessLog)
	}
	defer file.Close()

	LogEvents = parseStructs(file)

	fileInfo, err := os.Stat(AccessLog)
	previousOffset = fileInfo.Size()
}

func parseStructs(file io.Reader) []structs.LogEvent {
	reader := bufio.NewReader(file)
	events := make([]structs.LogEvent, 0)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		event, err := structs.ParseLogEvent(string(line))
		if err == nil {
			events = append(events, event)
		}
	}
	return events
}

// LogFileLastLine loads the last line
func LogFileLastLine() (string, error) {
	fileInfo, err := os.Stat(AccessLog)
	handleError(err)
	file, err := os.Open(AccessLog)
	handleError(err)

	defer file.Close()
	buffer := make([]byte, 1024)

	// +1 to compensate for the initial 0 byte of the line
	// otherwise, the initial character of the line will be missing

	// instead of reading the whole file into memory, we just read from certain offset
	offset := fileInfo.Size()
	numRead, err := file.ReadAt(buffer, previousOffset-1)

	if previousOffset != offset {
		// print out last line content
		buffer = buffer[:numRead]
		logEvent, error := structs.ParseLogEvent(string(buffer))
		if error == nil {
			LogEvents = append(LogEvents, logEvent)
		}
		previousOffset = offset
		return string(buffer), nil
	}
	return "", errors.New("No new line")
}
