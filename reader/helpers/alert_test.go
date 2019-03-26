package helpers

import (
	"fmt"
	"testing"
	"time"

	"github.com/veverkap/logtop/reader/structs"
)

func generateLogEventsSlice(count int) []structs.LogEvent {
	events := make([]structs.LogEvent, count)
	for index := 0; index < count; index++ {
		events[index] = structs.LogEvent{Date: time.Now()}
	}
	return events
}

func TestCalculateErrorStateFromDefault(t *testing.T) {
	type args struct {
		events                 []structs.LogEvent
		alertThresholdDuration int
		alertThreshold         int
	}
	// Testing Default ErrorState
	tests := []struct {
		name string
		args args
		want ErrorState
	}{
		{
			name: "Default - no events",
			args: args{
				events:                 generateLogEventsSlice(0),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Default,
		},
		{
			name: "Default - 1 event < 10",
			args: args{
				events:                 generateLogEventsSlice(1),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Default,
		},
		{
			name: "Default - 5 events violate",
			args: args{
				events:                 generateLogEventsSlice(5),
				alertThresholdDuration: 1,
				alertThreshold:         5,
			},
			want: Triggered,
		},
		{
			name: "Default - 6 events violate",
			args: args{
				events:                 generateLogEventsSlice(6),
				alertThresholdDuration: 1,
				alertThreshold:         5,
			},
			want: Triggered,
		},
	}
	for _, tt := range tests {
		fmt.Printf("Running %s\n", tt.name)
		CurrentErrorState = Default // need to reset this
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateErrorState(tt.args.events, tt.args.alertThresholdDuration, tt.args.alertThreshold); got != tt.want {
				t.Errorf("CalculateErrorState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateErrorStateFromTriggered(t *testing.T) {
	type args struct {
		events                 []structs.LogEvent
		alertThresholdDuration int
		alertThreshold         int
	}
	// Testing Default ErrorState
	tests := []struct {
		name string
		args args
		want ErrorState
	}{
		{
			name: "Triggered - no events should go to Recovered",
			args: args{
				events:                 generateLogEventsSlice(0),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Recovered,
		},
		{
			name: "Triggered - 1 event < 10",
			args: args{
				events:                 generateLogEventsSlice(1),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Recovered,
		},
		{
			name: "Triggered - 5 events violate",
			args: args{
				events:                 generateLogEventsSlice(5),
				alertThresholdDuration: 1,
				alertThreshold:         5,
			},
			want: WaitingForRecovery,
		},
		{
			name: "Triggered - 6 events violate",
			args: args{
				events:                 generateLogEventsSlice(6),
				alertThresholdDuration: 1,
				alertThreshold:         5,
			},
			want: WaitingForRecovery,
		},
	}
	for _, tt := range tests {
		fmt.Printf("Running %s\n", tt.name)
		CurrentErrorState = Triggered // need to reset this
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateErrorState(tt.args.events, tt.args.alertThresholdDuration, tt.args.alertThreshold); got != tt.want {
				t.Errorf("CalculateErrorState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateErrorStateFromWaitingForRecovery(t *testing.T) {
	type args struct {
		events                 []structs.LogEvent
		alertThresholdDuration int
		alertThreshold         int
	}
	// Testing Default ErrorState
	tests := []struct {
		name string
		args args
		want ErrorState
	}{
		{
			name: "WaitingForRecovery - no events should go to Recovered",
			args: args{
				events:                 generateLogEventsSlice(0),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Recovered,
		},
		{
			name: "WaitingForRecovery - 1 event < 10",
			args: args{
				events:                 generateLogEventsSlice(1),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Recovered,
		},
		{
			name: "WaitingForRecovery - 5 events violate",
			args: args{
				events:                 generateLogEventsSlice(5),
				alertThresholdDuration: 1,
				alertThreshold:         5,
			},
			want: WaitingForRecovery,
		},
		{
			name: "WaitingForRecovery - 6 events violate",
			args: args{
				events:                 generateLogEventsSlice(6),
				alertThresholdDuration: 1,
				alertThreshold:         5,
			},
			want: WaitingForRecovery,
		},
	}
	for _, tt := range tests {
		fmt.Printf("Running %s\n", tt.name)
		CurrentErrorState = WaitingForRecovery // need to reset this
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateErrorState(tt.args.events, tt.args.alertThresholdDuration, tt.args.alertThreshold); got != tt.want {
				t.Errorf("CalculateErrorState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateErrorStateFromRecovered(t *testing.T) {
	type args struct {
		events                 []structs.LogEvent
		alertThresholdDuration int
		alertThreshold         int
	}
	// Testing Default ErrorState
	tests := []struct {
		name string
		args args
		want ErrorState
	}{
		{
			name: "Recovered - no events should go to Recovered",
			args: args{
				events:                 generateLogEventsSlice(0),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Default,
		},
		{
			name: "Recovered - 1 event < 10",
			args: args{
				events:                 generateLogEventsSlice(1),
				alertThresholdDuration: 1,
				alertThreshold:         10,
			},
			want: Default,
		},
	}
	for _, tt := range tests {
		fmt.Printf("Running %s\n", tt.name)
		CurrentErrorState = Recovered // need to reset this
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateErrorState(tt.args.events, tt.args.alertThresholdDuration, tt.args.alertThreshold); got != tt.want {
				t.Errorf("CalculateErrorState() = %v, want %v", got, tt.want)
			}
		})
	}
}
