package helpers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/hpcloud/tail"
	"github.com/veverkap/logtop/reader/structs"
)

// UIStartTime is when the ui started
var UIStartTime time.Time

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func loadDebugValues() [][]string {
	now := time.Now()
	diff := now.Sub(UIStartTime)
	seconds := int(diff.Seconds())

	return [][]string{
		[]string{"Program Duration", fmt.Sprintf("%d secs", seconds)},
		[]string{"Total Event Count", fmt.Sprintf("%d", len(LogEvents))},

		[]string{"AlertThresholdDuration", fmt.Sprintf("%d secs", AlertThresholdDuration)},
		[]string{"AlertThreshold", fmt.Sprintf("%d/sec", AlertThreshold)},
		[]string{fmt.Sprintf("Events in last %d secs", AlertThresholdDuration), fmt.Sprintf("%d", structs.EventCount)},
		[]string{fmt.Sprintf("Event rate for last %d secs", AlertThresholdDuration), fmt.Sprintf("%.2f/sec", structs.PerSecondRate)},
		[]string{"Current Alert State", fmt.Sprintf("%s", structs.CurrentErrorState)},
	}
}
func reloadStatistics(events []structs.LogEvent) [][]string {
	details := structs.GroupBySection(events)
	details = structs.SortByHitsDesc(details)

	rows := [][]string{
		[]string{"Section", "Hits", "Errors"},
	}
	for _, detail := range details {
		rows = append(rows, []string{detail.Section, strconv.Itoa(detail.Hits), strconv.Itoa(detail.Errors)})
	}
	return rows
}

// LoopUI loads the UI and then goes into loop
func LoopUI(tail *tail.Tail) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	UIStartTime = time.Now()
	termWidth, termHeight := ui.TerminalDimensions()

	debugTable := widgets.NewTable()
	debugTable.Rows = loadDebugValues()
	debugTable.Title = "Debug Output"

	liveLog := widgets.NewList()
	liveLog.Title = "Live Log"
	liveLog.Rows = []string{}
	liveLog.WrapText = true
	liveLog.SetRect(0, 0, termWidth/2, termHeight/2)

	alerts := widgets.NewList()
	alerts.Title = "Alerts"
	alerts.Rows = []string{}
	alerts.WrapText = true
	alerts.SetRect(0, 0, 25, 8)

	statistics := widgets.NewTable()
	statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))
	statistics.Title = "Statistics (Last 10 Seconds)"
	statistics.TextStyle = ui.NewStyle(ui.ColorWhite)
	statistics.SetRect(0, 0, 60, 10)

	grid := ui.NewGrid()

	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(1.0/2,
				ui.NewRow(1.0/2, alerts),
				ui.NewRow(1.0/2, statistics),
			),
			ui.NewCol(1.0/2,
				ui.NewRow(1.0/2, debugTable),
				ui.NewRow(1.0/2, liveLog),
			),
		),
	)

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(500 * time.Millisecond).C

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case line, _ := <-tail.Lines:
			event, err := structs.ParseLogEvent(line.Text)
			if err == nil {
				LogEvents = append(LogEvents, event)
				processErrorState(alerts)

				liveLog.Rows = append(liveLog.Rows, line.Text)
				liveLog.ScrollPageDown()

				statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))
				debugTable.Rows = loadDebugValues()
				ui.Render(grid)
			}
		case <-ticker:
			processErrorState(alerts)
			statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))
			debugTable.Rows = loadDebugValues()
			ui.Render(grid)
		}
	}
}

func processErrorState(alerts *widgets.List) {
	errorState := structs.CalculateErrorState(LogEvents, AlertThresholdDuration, AlertThreshold)

	switch errorState {
	case structs.Triggered:
		displayErrorState(alerts)
	case structs.Recovered:
		hideErrorState(alerts)
	}

}
func displayErrorState(alerts *widgets.List) {
	t := time.Now()
	alerts.Rows = append(
		alerts.Rows,
		fmt.Sprintf("High traffic generated an alert - hits = %.2f/sec, triggered at %02d/%s/%d:%02d:%02d:%02d +0000", structs.PerSecondRate, t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second()),
	)
	alerts.ScrollPageDown()
}

func hideErrorState(alerts *widgets.List) {
	t := time.Now()
	alerts.Rows = append(
		alerts.Rows,
		fmt.Sprintf("High traffic alert recovered - hits = %.2f/sec, triggered at %02d/%s/%d:%02d:%02d:%02d +0000", structs.PerSecondRate, t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second()),
	)
	alerts.ScrollPageDown()
}
