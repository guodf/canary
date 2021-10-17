package lfu

import (
	"github.com/guodf/goutil/cache"
	"reflect"
	"testing"
)

func Test_newLfuCache(t *testing.T) {
	tests := []struct {
		name string
		want cache.ICache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLfuCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLfuCache() = %v, want %v", got, tt.want)
			}
		})
	}
}