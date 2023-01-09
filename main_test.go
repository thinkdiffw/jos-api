package main

import "testing"

func Test_getApiList(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{id: 573}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getApiList(tt.args.id)
		})
	}
}
