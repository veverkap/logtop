package helpers

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/veverkap/logtop/reader/structs"
)

// AccessLog represents the default AccessLog location
var AccessLog = "/tmp/access.log"

// LogEvents represents a slice of events
var LogEvents = make([]structs.LogEvent, 0)

var previousOffset int64

// LoadExistingLogFile will return the structs
func LoadExistingLogFile() {
	file, err := os.Open(AccessLog)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	LogEvents = parseStructs(file)
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
	file, err := os.Open(AccessLog)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	// we need to calculate the size of the last line for file.ReadAt(offset) to work

	// NOTE : not a very effective solution as we need to read
	// the entire file at least for 1 pass :(

	lastLineSize := 0

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		lastLineSize = len(line)
	}

	fileInfo, err := os.Stat(AccessLog)

	// make a buffer size according to the lastLineSize
	buffer := make([]byte, lastLineSize)

	// +1 to compensate for the initial 0 byte of the line
	// otherwise, the initial character of the line will be missing

	// instead of reading the whole file into memory, we just read from certain offset

	offset := fileInfo.Size() - int64(lastLineSize+1)
	numRead, err := file.ReadAt(buffer, offset)

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
