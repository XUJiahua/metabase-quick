package util

import "testing"

func TestGetFilenameWithExt(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				filename: "abc",
			},
			want: "abc",
		},
		{
			args: args{
				filename: "abc.txt",
			},
			want: "abc",
		},
		{
			args: args{
				filename: "/path/to/abc.txt",
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilenameWithExt(tt.args.filename); got != tt.want {
				t.Errorf("GetFilenameWithExt() = %v, want %v", got, tt.want)
			}
		})
	}
}
