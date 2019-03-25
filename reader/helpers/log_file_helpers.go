package helpers

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/veverkap/logtop/reader/structs"
)

// LogEvents represents a slice of events
var LogEvents = make([]structs.LogEvent, 0)

var previousOffset int64

func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// LoadExistingLogFile will:
// 1. Read existing log file
// 2. Parse each line into a LogEvent
// 3. Set the previousOffset to the current size of the file
func LoadExistingLogFile() {
	file, err := os.Open(LogFileLocation)
	if err != nil {
		log.Fatalf("Could not load file %s", LogFileLocation)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		event, err := structs.ParseLogEvent(string(line))

		if err == nil {
			LogEvents = append(LogEvents, event)
		}
	}

	fileInfo, err := os.Stat(LogFileLocation)
	previousOffset = fileInfo.Size()
}

// LoadLogFileUpdates loads the last line
func LoadLogFileUpdates() ([]string, error) {
	fileInfo, err := os.Stat(LogFileLocation)
	handleError(err)
	file, err := os.Open(LogFileLocation)
	handleError(err)

	defer file.Close()

	buffer := make([]byte, 1024)
	offset := fileInfo.Size()

	updates := make([]string, 0)

	if previousOffset != offset {
		numRead, _ := file.ReadAt(buffer, previousOffset-1)
		// print out last line content
		buffer = buffer[:numRead]
		bufferString := string(buffer)

		for _, line := range strings.Split(bufferString, "\n") {
			logEvent, error := structs.ParseLogEvent(line)
			if error == nil {
				LogEvents = append(LogEvents, logEvent)
				updates = append(updates, line)
			}
		}
		previousOffset = offset
		return updates, nil
	}
	return updates, errors.New("No new line")
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
