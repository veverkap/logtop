package helpers

import "flag"

// AlertThreshold represents "Number of requests per second maximum for alert"
var AlertThreshold int

// AlertThresholdDuration represents "Duration in seconds of sampling period for alerts"
var AlertThresholdDuration int

// LogFileLocation represents "Location of log file to parse"
var LogFileLocation string

// ParseFlags loads the flags passed at the command line or sets defaults
func ParseFlags() {
	flag.IntVar(&AlertThreshold, "threshold", 10, "Number of requests per second maximum for alert")
	flag.IntVar(&AlertThresholdDuration, "thresholdDuration", 120, "Duration in seconds of sampling period for alerts")
	flag.StringVar(&LogFileLocation, "logFileLocation", "/tmp/access.log", "Location of log file to parse")
	flag.Parse()
}
