package main

import (
	"flag"
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

	// displayUI()
}

func displayUI() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	liveLog := widgets.NewList()
	liveLog.Title = "Live Log"
	liveLog.Rows = []string{}
	liveLog.WrapText = true
	liveLog.SetRect(0, 0, 25, 8)

	alerts := widgets.NewList()
	alerts.Title = "Alerts"
	alerts.Rows = []string{}
	alerts.WrapText = true
	alerts.SetRect(0, 0, 25, 8)

	statistics := widgets.NewTable()

	events := structs.TrailingEvents(helpers.LogEvents, 10)
	details := structs.GroupBySection(events)
	rows := [][]string{
		[]string{"Section", "Hits", "Errors"},
	}
	for _, detail := range details {
		rows = append(rows, []string{detail.Section, string(detail.Hits), string(detail.Errors)})
	}
	statistics.Rows = rows

	statistics.Title = "Statistics (Last 10 Seconds)"
	statistics.TextStyle = ui.NewStyle(ui.ColorWhite)
	statistics.SetRect(0, 0, 60, 10)

	allTimeStatistics := widgets.NewTable()
	allTimeStatistics.Rows = [][]string{
		[]string{"Section", "Hits", "Errors"},
		[]string{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "10", "11"},
	}
	allTimeStatistics.Title = "Statistics (All Time)"
	allTimeStatistics.TextStyle = ui.NewStyle(ui.ColorWhite)
	allTimeStatistics.SetRect(0, 0, 60, 10)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, statistics),
			ui.NewCol(1.0/2, alerts),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, allTimeStatistics),
			ui.NewCol(1.0/2, liveLog),
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
				liveLog.ScrollDown()
			case "k", "<Up>":
				liveLog.ScrollUp()
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-tenTicker:
			events := structs.TrailingEvents(helpers.LogEvents, 10)
			details := structs.GroupBySection(events)
			rows := [][]string{
				[]string{"Section", "Hits", "Errors"},
			}
			for _, detail := range details {
				rows = append(rows, []string{detail.Section, string(detail.Hits), string(detail.Errors)})
			}
			statistics.Rows = rows
			ui.Render(grid)
		case <-ticker:
			line, err := helpers.LogFileLastLine()

			if err == nil {
				liveLog.Rows = append(liveLog.Rows, strings.ReplaceAll(line, "\n", ""))
				liveLog.ScrollPageDown()
				event, err := structs.ParseLogEvent(line)
				if err == nil {
					helpers.LogEvents = append(helpers.LogEvents, event)
				}
			}

			ui.Render(grid)
		}
	}
}
