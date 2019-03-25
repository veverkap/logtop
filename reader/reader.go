package main

import (
	"github.com/veverkap/logtop/reader/helpers"
)

func main() {
	helpers.ParseFlags()
	helpers.LoadExistingLogFile()
	// structs.CountLast24Hours(helpers.LogEvents)
	helpers.LoopUI()
}
