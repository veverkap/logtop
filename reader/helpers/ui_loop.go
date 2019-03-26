package helpers

import (
	"log"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/hpcloud/tail"
	"github.com/veverkap/logtop/reader/structs"
)

var isErrorState bool

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

func checkErrorState(alerts *widgets.List, alertThresholdDuration int, alertThreshold int) {
	eventCount := len(structs.TrailingEvents(LogEvents, int64(alertThresholdDuration)))
	perSecond := eventCount / alertThresholdDuration

	if isErrorState {
		// We are already in an error state, let's check if we should get OUT of error
		if perSecond < alertThreshold {
			// We can get out of error state
			alerts.Rows = append(alerts.Rows, "We are out of trouble.")
			isErrorState = false
		} else {
			// we can't do anything
			isErrorState = true
		}
	} else {
		// We are not currently in error state, let's check
		if perSecond > alertThreshold {
			// rate := fmt.Sprintf(
			// 	"%v seconds / %v hits / %d rate",
			// 	alertThresholdDuration,
			// 	eventCount,
			// 	perSecond,
			// )
			alerts.Rows = append(alerts.Rows, "We are in trouble")
			alerts.ScrollPageDown()
			isErrorState = true
		} else {
			isErrorState = false
		}
	}

}

// LoopUI loads the UI and then goes into loop
func LoopUI(tail *tail.Tail) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

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
				ui.NewRow(2.0/2, liveLog),
			),
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
				liveLog.Rows = append(liveLog.Rows, line.Text)
				liveLog.ScrollPageDown()

				statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))

				ui.Render(grid)
			}
		case <-tenTicker:
			statistics.Rows = reloadStatistics(structs.TrailingEvents(LogEvents, 10))
			checkErrorState(alerts, AlertThresholdDuration, AlertThreshold)
			ui.Render(grid)
		case <-ticker:
			ui.Render(grid)
		}
	}
}
