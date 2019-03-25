package structs

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestTrailingEvents(t *testing.T) {
	type args struct {
		logEvents   []LogEvent
		lastSeconds int64
	}

	within10 := LogEvent{
		Date: time.Now().Add(-15.0),
	}

	moreThan10 := LogEvent{
		Date: time.Date(
			2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}
	tests := []struct {
		name string
		args args
		want []LogEvent
	}{
		{
			name: "nothing begets nothing",
			args: args{
				lastSeconds: 10,
				logEvents:   make([]LogEvent, 0),
			},
			want: make([]LogEvent, 0),
		},
		{
			name: "get it if within 10",
			args: args{
				lastSeconds: 10,
				logEvents:   append(make([]LogEvent, 0), within10),
			},
			want: append(make([]LogEvent, 0), within10),
		},
		{
			name: "gets none if over 10",
			args: args{
				lastSeconds: 10,
				logEvents:   append(make([]LogEvent, 0), moreThan10),
			},
			want: make([]LogEvent, 0),
		},
		{
			name: "get only the one",
			args: args{
				lastSeconds: 10,
				logEvents:   append(append(make([]LogEvent, 0), within10), moreThan10),
			},
			want: append(make([]LogEvent, 0), within10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrailingEvents(tt.args.logEvents, tt.args.lastSeconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrailingEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLogEvent(t *testing.T) {
	// Setup a timestamp for the log line
	date, formattedDate := generateTime(time.Now())

	type args struct {
		line string
	}

	tests := []struct {
		name string
		args args
		want LogEvent
	}{
		{
			name: "Normal form",
			args: args{
				line: fmt.Sprintf("127.0.0.1 - frank [%s] \"DELETE /config/update HTTP/1.0\" 401 491", formattedDate),
			},
			want: LogEvent{
				Host:       "127.0.0.1",
				User:       "frank",
				Date:       date,
				Verb:       "DELETE",
				Section:    "/config",
				Path:       "/config/update",
				StatusCode: 401,
				ByteSize:   491,
				Error:      true,
			},
		},
		{
			name: "Success",
			args: args{
				line: fmt.Sprintf("127.0.0.1 - frank [%s] \"DELETE /config HTTP/1.0\" 200 491", formattedDate),
			},
			want: LogEvent{
				Host:       "127.0.0.1",
				User:       "frank",
				Date:       date,
				Verb:       "DELETE",
				Section:    "/config",
				Path:       "/config",
				StatusCode: 200,
				ByteSize:   491,
				Error:      false,
			},
		},
		{
			name: "Bad format on host",
			args: args{
				line: fmt.Sprintf("127.0.0 - frank [%s] \"DELETE /config HTTP/1.0\" 200 491", formattedDate),
			},
			want: LogEvent{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ParseLogEvent(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLogEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateTime(now time.Time) (time.Time, string) {
	const longForm = "02/Jan/2006:15:04:05 -0700"
	formattedDate := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%02d +0000", now.Day(), now.Month().String()[:3], now.Year(), now.Hour(), now.Minute(), now.Second())
	date, _ := time.Parse(longForm, formattedDate)
	return date, formattedDate
}

func TestFilter(t *testing.T) {
	type args struct {
		vs []LogEvent
		f  func(LogEvent) bool
	}
	tests := []struct {
		name string
		args args
		want []LogEvent
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.vs, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortBySection(t *testing.T) {
	type args struct {
		events []LogEvent
	}
	tests := []struct {
		name string
		args args
		want []LogEvent
	}{
		{
			name: "Normal",
			args: args{
				events: []LogEvent{
					LogEvent{Section: "/blah"},
					LogEvent{Section: "/apple"},
				},
			},
			want: []LogEvent{
				LogEvent{Section: "/apple"},
				LogEvent{Section: "/blah"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortBySection(tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortBySection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortByDateAsc(t *testing.T) {
	type args struct {
		events []LogEvent
	}
	tests := []struct {
		name string
		args args
		want []LogEvent
	}{
		{
			name: "Normal",
			args: args{
				events: []LogEvent{
					LogEvent{Date: time.Date(
						2009, 11, 17, 20, 34, 58, 651387237, time.UTC)},
					LogEvent{Date: time.Date(
						2007, 11, 17, 20, 34, 58, 651387237, time.UTC)},
				},
			},
			want: []LogEvent{
				LogEvent{Date: time.Date(
					2007, 11, 17, 20, 34, 58, 651387237, time.UTC)},
				LogEvent{Date: time.Date(
					2009, 11, 17, 20, 34, 58, 651387237, time.UTC)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortByDateAsc(tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByDateAsc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortByDateDesc(t *testing.T) {
	type args struct {
		events []LogEvent
	}
	tests := []struct {
		name string
		args args
		want []LogEvent
	}{
		{
			name: "Normal",
			args: args{
				events: []LogEvent{
					LogEvent{Date: time.Date(
						2009, 11, 17, 20, 34, 58, 651387237, time.UTC)},
					LogEvent{Date: time.Date(
						2007, 11, 17, 20, 34, 58, 651387237, time.UTC)},
				},
			},
			want: []LogEvent{
				LogEvent{Date: time.Date(
					2009, 11, 17, 20, 34, 58, 651387237, time.UTC)},
				LogEvent{Date: time.Date(
					2007, 11, 17, 20, 34, 58, 651387237, time.UTC)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortByDateDesc(tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByDateDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}
