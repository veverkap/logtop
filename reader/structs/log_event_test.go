package structs

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func generateTime(now time.Time) (time.Time, string) {
	const longForm = "02/Jan/2006:15:04:05 -0700"
	formattedDate := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%02d +0000", now.Day(), now.Month().String()[:3], now.Year(), now.Hour(), now.Minute(), now.Second())
	date, _ := time.Parse(longForm, formattedDate)
	return date, formattedDate
}

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

func TestGroupBySection(t *testing.T) {
	type args struct {
		vs []LogEvent
	}
	now := time.Now()
	section1 := LogEvent{Section: "/scuba", Path: "/scuba/doo", Error: false, Date: now}
	section2 := LogEvent{Section: "/scuba", Path: "/scuba/loo", Error: true, Date: now}
	section3 := LogEvent{Section: "/other", Path: "/other/loo", Error: false, Date: now}

	tests := []struct {
		name string
		args args
		want []SectionDetail
	}{
		{
			name: "sections",
			args: args{
				vs: append(make([]LogEvent, 0), section1, section2),
			},
			want: []SectionDetail{
				SectionDetail{
					Section: "/scuba",
					Events:  []LogEvent{section1, section2},
					Hits:    2,
					Errors:  1,
				},
			},
		},
		{
			name: "sections grouped",
			args: args{
				vs: append(make([]LogEvent, 0), section1, section2, section3),
			},
			want: []SectionDetail{
				SectionDetail{
					Section: "/scuba",
					Events: []LogEvent{
						section1,
						section2,
					},
					Hits:   2,
					Errors: 1,
				},
				SectionDetail{
					Section: "/other",
					Events: []LogEvent{
						section3,
					},
					Hits:   1,
					Errors: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupBySection(tt.args.vs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GOT  = %v\nWANT = %v", got, tt.want)
			}
		})
	}
}
