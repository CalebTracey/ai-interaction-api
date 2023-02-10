package facade

import "testing"

func Test_remainder(t *testing.T) {
	tests := []struct {
		name        string
		requestNum  int
		requestSize int
		want        int
	}{
		{
			name:        "happy",
			requestNum:  15,
			requestSize: 12,
			want:        3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := remainder(tt.requestNum, tt.requestSize); got != tt.want {
				t.Errorf("remainder() = %v, want %v", got, tt.want)
			}
		})
	}
}
