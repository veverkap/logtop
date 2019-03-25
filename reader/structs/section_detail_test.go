package structs

import (
	"reflect"
	"testing"
)

func TestGroupBySection(t *testing.T) {
	type args struct {
		vs []LogEvent
	}
	section1 := LogEvent{Section: "/scuba", Path: "/scuba/doo"}
	section2 := LogEvent{Section: "/scuba", Path: "/scuba/loo"}
	section3 := LogEvent{Section: "/other", Path: "/other/loo"}

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
					Hits: 2,
				},
				SectionDetail{
					Section: "/other",
					Events: []LogEvent{
						section3,
					},
					Hits: 1,
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
