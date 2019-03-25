package main

import (
	"github.com/veverkap/logtop/reader/helpers"
)

func main() {
	helpers.ParseFlags()
	helpers.LoadExistingLogFile()
	helpers.LoopUI()
}
