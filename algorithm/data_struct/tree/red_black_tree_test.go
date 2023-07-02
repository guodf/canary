package tree

import (
	"log"
	"testing"
)

type IntValue int

func (value IntValue) Compare(v interface{}) int {
	return int(value) - int(v.(IntValue))
}

func TestNewRBTree(t *testing.T) {
	tree := NewRBTree()
	arr := []IntValue{IntValue(16), IntValue(3), IntValue(7), IntValue(11), IntValue(9), IntValue(26), IntValue(18), IntValue(14), IntValue(15),IntValue(1),IntValue(4),IntValue(5)}
	for _, i := range arr {
		tree.Add(i)
	}
	log.Println(tree)
}
