package main

import (
	"calculator/tools"
	"testing"
)

var (
	n string
)

func Test_Mean(t *testing.T) {
	tests := []struct {
		n    string
		name string
		want string
	}{
		{
			n:    "1 5 5 10 15 2 3\n",
			name: "test_mean_1",
			want: "5.8571428571429\n",
		},
		{
			n:    "9 10 12 13 13 13 15 15 16 16 18 22 23 24 24 25\n",
			name: "test_mean_2",
			want: "16.7500000000000\n",
		},
		{
			n:    "0\n",
			name: "test_mean_3",
			want: "0\n",
		},
		{
			n:    "0000\n",
			name: "test_mean_4",
			want: "0\n",
		},
		{
			n:    "1 2 3 4 5 6 7 8 9 10 1000000000\n",
			name: "test_mean_5",
			want: "90909095.9090909063816\n",
		},
		{
			n:    "1 2 3 4 5 6 7 8 9 10 0\n",
			name: "test_mean_6",
			want: "5.0000000000000\n",
		},
		{
			n:    "\n",
			name: "test_mean_6",
			want: "NAN\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mean(tools.Str2slice(tt.n)); got != tt.want {
				t.Errorf("mean() = %v, want %v", got, tt.want)
			}
		})
	}
}
