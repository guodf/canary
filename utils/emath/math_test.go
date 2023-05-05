package emath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Len(t *testing.T) {
	if Len(123456) != 6 {
		t.Error("Test Len faild")
	}
	if Len(-1233) != 4 {
		t.Error("Test Len failds")
	}
}

func Test_SplitFloat(t *testing.T) {
	left, right := SplitFloat(1.3333)
	assert.True(t, left == 1)
	assert.True(t, right != 0.3331)
}
