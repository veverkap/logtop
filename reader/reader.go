package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

var previousOffset int64

func readLastLine(filename string) (string, error) {

	file, err := os.Open(filename)
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

	fileInfo, err := os.Stat(filename)

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
		previousOffset = offset
		return string(buffer), nil
	}
	return "", errors.New("No new line")

}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Live Log"
	l.Rows = []string{}
	l.WrapText = true
	l.SetRect(0, 0, 25, 8)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, l),
		),
	)

	ui.Render(grid)
	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(500 * time.Millisecond).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "j", "<Down>":
				l.ScrollDown()
			case "k", "<Up>":
				l.ScrollUp()
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			line, err := readLastLine("/tmp/access.log")
			if err == nil {
				line = strings.TrimSuffix(line, "\n")
				l.Rows = append(l.Rows, line)
			}

			ui.Render(grid)

			// .Text = text
			tickerCount++
		}
	}
}
