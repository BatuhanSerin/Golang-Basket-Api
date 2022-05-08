package check

import "testing"

func TestCheckDigits(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				s: "1234567890",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDigits(tt.args.s); got != tt.want {
				t.Errorf("CheckDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}
