package main

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/veverkap/logtop/reader/helpers"
	"github.com/veverkap/logtop/reader/structs"
)

func main() {
	helpers.LoadExistingLogFile()

	if len(os.Args) > 1 {
		helpers.AccessLog = os.Args[1]
	}
	// fmt.Printf("Loaded %d logevents from %s", len(helpers.LogEvents), helpers.AccessLog)
	displayUI()
}

func displayUI() {
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
			line, err := helpers.LogFileLastLine()

			if err == nil {
				l.Rows = append(l.Rows, line)
				event, err := structs.ParseLogEvent(line)
				if err == nil {
					helpers.LogEvents = append(helpers.LogEvents, event)
				}
			}

			ui.Render(grid)

			// .Text = text
			tickerCount++
		}
	}
}
