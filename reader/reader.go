package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/veverkap/logtop/reader/helpers"
	"github.com/veverkap/logtop/reader/structs"
)

var threshold int
var thresholdDuration int
var logFileLocation string

func main() {
	flag.IntVar(&threshold, "threshold", 10, "Number of requests per second maximum for alert")
	flag.IntVar(&thresholdDuration, "thresholdDuration", 120, "Duration in seconds of sampling period for alerts")
	flag.StringVar(&logFileLocation, "logFileLocation", "/tmp/access.log", "Location of log file to parse")
	flag.Parse()

	helpers.AccessLog = logFileLocation
	helpers.LoadExistingLogFile()
	fmt.Printf("There are currently %d events\n", len(helpers.LogEvents))

	// for {
	// 	line, err := helpers.LogFileLastLine()
	// 	if err == nil {
	// 		structs.ParseLogEvent(line)
	// 		fmt.Printf("There are currently %d events\n", len(helpers.LogEvents))
	// 		hits := structs.TrailingEvents(helpers.LogEvents, 10)
	// 		fmt.Printf("There have been %d events in last 10\n", len(hits))
	// 	}
	// }
	// displayUI()
	// twoD := make([][]int, 0)
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

	trailingEvents := structs.TrailingEvents(helpers.LogEvents, 10)
	hits := len(trailingEvents)

	p := widgets.NewParagraph()
	p.Title = "Hits/Sec For Last 10 Seconds"
	p.Text = "LOK"
	p.SetRect(0, 0, 10, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	table1 := widgets.NewTable()
	table1.Rows = [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "10", "11"},
	}
	table1.TextStyle = ui.NewStyle(ui.ColorWhite)
	table1.SetRect(0, 0, 60, 10)
	hitrate := widgets.NewParagraph()
	hitrate.Title = "Hits/Sec For Last 10 Seconds"

	hitrate.Text = fmt.Sprintf("\n    %d req/sec", hits)
	hitrate.SetRect(0, 0, 10, 5)
	hitrate.TextStyle.Fg = ui.ColorWhite
	hitrate.BorderStyle.Fg = ui.ColorCyan

	p2 := widgets.NewParagraph()
	p2.Title = "ALAMRS"
	p2.Text = "DUDE"

	alarms := widgets.NewParagraph()
	alarms.Text = "<> This row has 3 columns\n<- Widgets can be stacked up like left side\n<- Stacked widgets are treated as a single widget"
	alarms.Title = "Demonstration"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/4, p),
			ui.NewCol(1.0/4,
				ui.NewRow(.5/3, hitrate),
				ui.NewRow(.9/3, p),
				ui.NewRow(1.2/3, p2),
			),
			ui.NewCol(1.0/2, alarms),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, table1),
			ui.NewCol(1.0/2, l),
		),
	)

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(500 * time.Millisecond).C
	tenTicker := time.NewTicker(10 * time.Second).C

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
		case <-tenTicker:
			// hits := len(structs.TrailingEvents(helpers.LogEvents, 10))
			p2.Text += "\n10"
			ui.Render(grid)
		case <-ticker:
			line, err := helpers.LogFileLastLine()

			if err == nil {
				l.Rows = append(l.Rows, strings.ReplaceAll(line, "\n", ""))
				l.ScrollPageDown()
				event, err := structs.ParseLogEvent(line)
				if err == nil {
					helpers.LogEvents = append(helpers.LogEvents, event)
				}
				hits := len(structs.TrailingEvents(helpers.LogEvents, 10))
				p.Text = fmt.Sprintf("    %d req/sec", hits)
				hitrate.Text = fmt.Sprintf("\n    %d req/sec", hits)
			}

			ui.Render(grid)
		}
	}
}
