package lfu

import (
	"reflect"
	"testing"

	"e.coding.net/guodf/gopkg/goutil/cache"
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
