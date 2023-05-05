package strutil

import (
	"testing"
)

func TestValidatePhoneNumber(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "19211111111_test",
			args: args{
				phone: "18761739725",
			},
			want: true,
		},
		{
			name: "19611111111_test",
			args: args{
				phone: "19611111111",
			},
			want: true,
		},
		{
			name: "19711111111_test",
			args: args{
				phone: "19711111111",
			},
			want: true,
		},
		{
			name: "19911111111_test",
			args: args{
				phone: "19911111111",
			},
			want: true,
		},
		{
			name: "19811111111_test",
			args: args{
				phone: "19811111111",
			},
			want: true,
		},
		{
			name: "15511111111_test",
			args: args{
				phone: "15511111111",
			},
			want: true,
		},
		{
			name: "15611111111_test",
			args: args{
				phone: "15611111111",
			},
			want: true,
		},
		{
			name: "16611111111_test",
			args: args{
				phone: "16611111111",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePhoneNumber(tt.args.phone); got != tt.want {
				t.Errorf("ValidatePhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateEmailNumber(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				email: "a@123.com",
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				email: "b",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEmailNumber(tt.args.email); got != tt.want {
				t.Errorf("ValidateEmailNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateStrongPassword(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				s: "Aa123456",
			},
			want: true,
		},
		{
			name: "test2",
			args: args{
				s: "aa123456",
			},
			want: false,
		},
		{
			name: "test3",
			args: args{
				s: "asdfghjk",
			},
			want: false,
		},
		{
			name: "test4",
			args: args{
				s: "12345678",
			},
			want: false,
		},
		{
			name: "test5",
			args: args{
				s: "Aa12345",
			},
			want: false,
		},
		{
			name: "test6",
			args: args{
				s: "yuhu8888",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStrongPassword(tt.args.s); got != tt.want {
				t.Errorf("ValidateStrongPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
