package structs

import (
	"reflect"
	"testing"
	"time"
)

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
					Section:            "/scuba",
					Events:             []LogEvent{section1, section2},
					Hits:               2,
					Errors:             1,
					HitsLastTenSeconds: 2,
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
					Hits:               2,
					HitsLastTenSeconds: 2,
					Errors:             1,
				},
				SectionDetail{
					Section: "/other",
					Events: []LogEvent{
						section3,
					},
					Hits:               1,
					HitsLastTenSeconds: 1,
					Errors:             0,
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
