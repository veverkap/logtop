package main

import (
	"log"
	"os"

	"github.com/hpcloud/tail"

	"github.com/veverkap/logtop/reader/helpers"
)

func main() {
	helpers.ParseFlags()
	tail := loadTail(helpers.LogFileLocation)
	helpers.LoopUI(tail)
}

// loadTail loads up a pointer to the tail object used to get updates from inotify
func loadTail(logFileLocation string) *tail.Tail {
	tail, err := tail.TailFile(
		logFileLocation,
		tail.Config{
			Follow:    true, // Continue looking for new lines (tail -f)
			MustExist: true, // Fail early if the file does not exist
			Location: &tail.SeekInfo{ // Seek to this location before tailing
				Offset: 0,
				Whence: os.SEEK_END,
			},
		},
	)
	if err != nil {
		log.Fatalf("Could not open log file at %s", logFileLocation)
	}
	return tail
}
