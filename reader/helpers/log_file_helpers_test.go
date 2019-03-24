package helpers

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/veverkap/logtop/reader/structs"
)

func Test_parseStructs(t *testing.T) {
	type args struct {
		file io.Reader
	}
	tests := []struct {
		name string
		args args
		want []structs.LogEvent
	}{
		{
			name: "Boo",
			args: args{
				file: strings.NewReader("this is it"),
			},
			want: make([]structs.LogEvent, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseStructs(tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStructs() = %v, want %v", got, tt.want)
			}
		})
	}
}
