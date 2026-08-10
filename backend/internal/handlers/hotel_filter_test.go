package handlers

import (
	"reflect"
	"testing"
)

func TestParsePositiveIDs(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    []int
		wantErr bool
	}{
		{name: "empty", raw: "", want: []int{}},
		{name: "multiple", raw: "3, 1,3,2", want: []int{3, 1, 2}},
		{name: "invalid text", raw: "1,a", wantErr: true},
		{name: "invalid zero", raw: "0", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parsePositiveIDs(test.raw)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil || !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %v, %v; want %v", got, err, test.want)
			}
		})
	}
}
