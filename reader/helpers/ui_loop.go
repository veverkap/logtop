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

// LogEvents is a slice of LogEvents (representations of lines from the log)
var LogEvents = make([]structs.LogEvent, 0)

// UIStartTime is when the ui started
var UIStartTime time.Time

// loadDebugValues generates a table of debug values
func loadDebugValues() [][]string {
	now := time.Now()
	diff := now.Sub(UIStartTime)
	seconds := int(diff.Seconds())

	return [][]string{
		[]string{"Program Duration", fmt.Sprintf("%d secs", seconds)},
		[]string{"Total Event Count", fmt.Sprintf("%d", len(LogEvents))},

		[]string{"AlertThresholdDuration", fmt.Sprintf("%d secs", AlertThresholdDuration)},
		[]string{"AlertThreshold", fmt.Sprintf("%d/sec", AlertThreshold)},
		[]string{fmt.Sprintf("Events in last %d secs", AlertThresholdDuration), fmt.Sprintf("%d", ThresholdEventCount)},
		[]string{fmt.Sprintf("Event rate for last %d secs", AlertThresholdDuration), fmt.Sprintf("%.2f/sec", ThresholdRate)},
		[]string{"Current Alert State", fmt.Sprintf("%s", CurrentErrorState)},
	}
}

// reloadStatistics generates a table of statistics
func reloadStatistics(events []structs.LogEvent) [][]string {
	details := structs.GroupBySection(events)
	details = structs.SortSectionDetailsByHitsDesc(details)

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
	UIStartTime = time.Now()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	// this is our debug table
	debugTable := widgets.NewTable()
	debugTable.Rows = loadDebugValues()
	debugTable.Title = "Debug Output"

	// this will include the log (an echo)
	liveLog := widgets.NewList()
	liveLog.Title = "Live Log"
	liveLog.Rows = []string{}
	liveLog.WrapText = true
	liveLog.SetRect(0, 0, termWidth/2, termHeight/2)

	// holder for any alerts
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
			// we receive a message in the tail file chan
			event, err := structs.ParseLogEvent(line.Text)
			if err == nil {
				// we were able to parse this line and need to add it to our LogEvents slice
				LogEvents = append(LogEvents, event)

				// add this line to our liveLog
				liveLog.Rows = append(liveLog.Rows, line.Text)
				liveLog.ScrollPageDown()

				// let's check if these changes triggered an alert
				processErrorState(alerts)

				// recalculate statistics for the last 10 seconds
				statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))

				// load debug values and display
				debugTable.Rows = loadDebugValues()
				ui.Render(grid)
			}
		case <-ticker:
			// it's been 500 ms, let's see if we are in alert
			processErrorState(alerts)

			// recalculate statistics for the last 10 seconds
			statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))

			// load debug values and display
			debugTable.Rows = loadDebugValues()
			ui.Render(grid)
		}
	}
}

// processErrorState calls the structs.Alert.CalculateErrorState and adds an Alert when appropriate
func processErrorState(alerts *widgets.List) {
	errorState := CalculateErrorState(LogEvents, AlertThresholdDuration, AlertThreshold)

	switch errorState {
	case Triggered:
		displayErrorState(alerts)
	case Recovered:
		hideErrorState(alerts)
	}

}

// displayErrorState adds a text notification to the list that we generated an alert
func displayErrorState(alerts *widgets.List) {
	t := time.Now()
	alerts.Rows = append(
		alerts.Rows,
		fmt.Sprintf("High traffic generated an alert - hits = %.2f/sec, triggered at %02d/%s/%d:%02d:%02d:%02d +0000", ThresholdRate, t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second()),
	)
	alerts.ScrollPageDown()
}

// displayErrorState adds a text notification to the list that we have recovered from our alert
func hideErrorState(alerts *widgets.List) {
	t := time.Now()
	alerts.Rows = append(
		alerts.Rows,
		fmt.Sprintf("High traffic alert recovered - hits = %.2f/sec, triggered at %02d/%s/%d:%02d:%02d:%02d +0000", ThresholdRate, t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second()),
	)
	alerts.ScrollPageDown()
}
