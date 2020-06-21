package hashing

import "testing"

func TestMurmur3_128(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Murmur3_128(tt.args.str); got != tt.want {
				t.Errorf("Murmur3_128() = %v, want %v", got, tt.want)
			}
		})
	}
}