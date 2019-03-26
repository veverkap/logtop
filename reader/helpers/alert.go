package helpers

import "github.com/veverkap/logtop/reader/structs"

// ErrorState int for "enum"
type ErrorState int

/*
The idea here is that we have a few states in our system
	1. We are not in alert mode (Default)
	2. We received an alert but we have not shown the UI indicator (Triggered)
	3. We displayed the UI indicator and are waiting for recovery (WaitingForRecovery)
	4. We have recovered and displayed a UI indicator that we have recovered (Recovered)
*/
const (
	Default            ErrorState = iota
	Triggered          ErrorState = iota
	WaitingForRecovery ErrorState = iota
	Recovered          ErrorState = iota
)

// String converts the ErrorState to a string representation
func (state ErrorState) String() string {
	names := []string{"Default", "Triggered", "Waiting For Recovery", "Recovered"}

	if state < Default || state > Recovered {
		return "Unknown"
	}
	return names[state]
}

// CurrentErrorState represents the CURRENT state of the system
var CurrentErrorState ErrorState

// ThresholdEventCount represents the number of events in the threshold period
var ThresholdEventCount int

// ThresholdRate is the current rate (hits / sec) in the current threshold (used for UI)
var ThresholdRate float64

// CalculateErrorState checks the trailing events and calculates whether the threshold has been met
func CalculateErrorState(events []structs.LogEvent, alertThresholdDuration int, alertThreshold int) ErrorState {
	ThresholdEventCount = len(structs.TrailingEvents(events, int64(alertThresholdDuration)))
	perSecond := ThresholdEventCount / alertThresholdDuration

	ThresholdRate = float64(ThresholdEventCount*1.0) / float64(alertThresholdDuration*1.0)

	if CurrentErrorState == Default {
		// We are in default state
		if perSecond >= alertThreshold {
			// We need to alert on this
			CurrentErrorState = Triggered
			return Triggered
		}
		// We continue in default
		return Default
	}

	if CurrentErrorState == Triggered {
		// We are in triggered mode, looking to recover
		if perSecond < alertThreshold {
			// We can recover at this point
			CurrentErrorState = Recovered
			return Recovered
		}
		// We continue in triggered
		return WaitingForRecovery
	}

	if CurrentErrorState == WaitingForRecovery {
		// We are in WaitingForRecovery mode, looking to recover
		if perSecond < alertThreshold {
			// We can recover at this point
			CurrentErrorState = Recovered
			return Recovered
		}
		// We continue in triggered
		return WaitingForRecovery
	}

	if CurrentErrorState == Recovered {
		// We are in recovered mode and just need to transition back to Default
		CurrentErrorState = Default
		return Default
	}
	return Default
}
