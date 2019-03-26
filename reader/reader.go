package main

import (
	"log"
	"os"

	"github.com/hpcloud/tail"
	"github.com/veverkap/logtop/reader/helpers"
)

func main() {
	helpers.ParseFlags()
	tail := loadTail()
	helpers.LoopUI(tail)
}

func loadTail() *tail.Tail {
	tail, err := tail.TailFile(
		helpers.LogFileLocation,
		tail.Config{
			Follow:    true,
			MustExist: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: os.SEEK_END,
			},
		},
	)
	if err != nil {
		log.Fatalf("Could not open log file at %s", helpers.LogFileLocation)
	}
	return tail
}
