package main

import (
	"fmt"
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
	p := fmt.Println

	if len(os.Args) > 1 {
		helpers.AccessLog = os.Args[1]

	}
	fmt.Printf("helpers.AccessLog = %s", helpers.AccessLog)

	// for _, event := range logEvents {
	// 	// print(event.Date)
	// 	//fmt.Printf("%+v\n", event)
	// 	then := event.Date
	// 	now := time.Now()
	// 	p(then)
	// 	p(now)

	// 	diff := now.Sub(then)
	// 	p(diff)
	// 	p("\n")
	// }

	for i, event := range structs.TrailingEvents(helpers.LogEvents, 10) {
		p(i)
		fmt.Printf("%+v\n", event)
	}

	// for _, event := range left {
	// 	// fmt.Printf("%+v\n", event)
	// }

}

func loader() {
	if len(os.Args) > 1 {
		helpers.AccessLog = os.Args[1]
	}

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
			}
			ui.Render(grid)

			// .Text = text
			tickerCount++
		}
	}
}
