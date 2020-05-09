package filesys

import (
	"os"
	"testing"
)


func TestFindFilesInExt(t *testing.T) {
	type args struct {
		dirPath string
		exts    []string
		depth   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
	}
	files := []string{"d:/download/aaa.exe", "d:/download/bbb.zip"}
	for _, file := range files {
		os.Create(file)
		tests = append(tests, struct {
			name string
			args args
			want int
		}{name: file, args: args{
			"d:/download",
			[]string{".exe", ".zip"},
			0,
		}, want: len(files)})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindFilesInExt(tt.args.dirPath, tt.args.exts, tt.args.depth); len(got) < tt.want {
				t.Errorf("FindFilesInExt() = %v, want %v", got, tt.want)
			}
		})
	}
}