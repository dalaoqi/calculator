package main

import (
	"calculator/tools"
	"testing"
)

func Test_median(t *testing.T) {
	tests := []struct {
		n    string
		name string
		want string
	}{
		{
			n:    "1 5 5 10 15 2 3\n",
			name: "test_median_1",
			want: "5\n",
		},
		{
			n:    "9 10 12 13 13 13 15 15 16 16 18 22 23 24 24 25\n",
			name: "test_median_2",
			want: "15.5000000000000\n",
		},
		{
			n:    "0\n",
			name: "test_median_3",
			want: "0\n",
		},
		{
			n:    "0000\n",
			name: "test_median_4",
			want: "0\n",
		},
		{
			n:    "1 2 3 4 5 1000000000 7 8 9 10 \n",
			name: "test_median_5",
			want: "6.0000000000000\n",
		},
		{
			n:    "\n",
			name: "test_median_6",
			want: "NAN\n",
		},
		{
			n:    "1 2 3 4 5 6 7 8 9 10 0\n",
			name: "test_median_7",
			want: "5\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := median(tools.Str2slice(tt.n)); got != tt.want {
				t.Errorf("median() = %v, want %v", got, tt.want)
			}
		})
	}
}
