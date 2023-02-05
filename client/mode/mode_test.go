package main

import (
	"calculator/tools"
	"testing"
)

func Test_mode(t *testing.T) {
	tests := []struct {
		n    string
		name string
		want string
	}{
		{
			n:    "1 5 5 10 15 2 3\n",
			name: "test_mode_1",
			want: "5\n",
		},
		{
			n:    "9 10 12 13 13 13 15 15 16 16 18 22 23 24 24 25\n",
			name: "test_mode_2",
			want: "13\n",
		},
		{
			n:    "0\n",
			name: "test_mode_3",
			want: "0\n",
		},
		{
			n:    "0000\n",
			name: "test_mode_4",
			want: "0\n",
		},
		{
			n:    "1 2 3 4 5 1000000000 7 8 9 10 \n",
			name: "test_mode_5",
			want: "1 2 3 4 5 7 8 9 10 1000000000\n",
		},
		{
			n:    "\n",
			name: "test_mode_6",
			want: "NAN\n",
		},
		{
			n:    "1 2 3 4 5 6 7 8 9 10 0\n",
			name: "test_mode_7",
			want: "0 1 2 3 4 5 6 7 8 9 10\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mode(tools.Str2slice(tt.n)); got != tt.want {
				t.Errorf("mode() = %v, want %v", got, tt.want)
			}
		})
	}
}
