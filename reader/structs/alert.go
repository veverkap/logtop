package structs

// ErrorState int
type ErrorState int

/* The idea here is that we have a few states in our system
1. We are not in alert mode
2. We are in alert mode and are waiting for recovery (we displayed our alert)
3. We have recovered and need to go into non-alert mode
*/
const (
	Default            ErrorState = iota
	Triggered          ErrorState = iota
	WaitingForRecovery ErrorState = iota
	Recovered          ErrorState = iota
)

// CurrentErrorState represents the state of the system
var CurrentErrorState ErrorState

// EventCount represents the number of events in the threshold period
var EventCount int

// PerSecondRate is the current rate
var PerSecondRate float64

func (state ErrorState) String() string {
	names :=
		[]string{"Default", "Triggered", "Waiting For Recovery", "Recovered"}

	if state < Default || state > Recovered {
		return "Unknown"
	}
	return names[state]
}

// CalculateErrorState checks the trailing events and calculates whether the threshold has been met
func CalculateErrorState(events []LogEvent, alertThresholdDuration int, alertThreshold int) ErrorState {
	EventCount = len(TrailingEvents(events, int64(alertThresholdDuration)))
	perSecond := EventCount / alertThresholdDuration

	PerSecondRate = float64(EventCount*1.0) / float64(alertThresholdDuration*1.0)

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
