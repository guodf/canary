package mapper_test

import (
	"fmt"
	"testing"

	"e.coding.net/guodf/gopkg/goutil/mapper"
)

type SubDto struct {
	Field1 string
	Field2 int
}

// HelloDto dto
type HelloDto struct {
	Field1 string
	Field2 int
	Field3 bool
	Field4 byte
	Field5 []byte
	Field6 SubDto
	Field7 *SubDto
	Field8 []*SubDto
	Field9 map[string]interface{}
}

// HelloModel model
type HelloModel struct {
	Field1 string
	Field2 int
	Field3 bool
	Field4 byte
	Field5 []byte
	Field6 SubDto
	Field7 *SubDto
	Field8 []*SubDto
	Field9 map[string]interface{}
}

func TestDtoToModel(t *testing.T) {
	var helloDto HelloDto
	mapper.Mapper(HelloDto{
		Field1: "field1",
	}, &helloDto)
	fmt.Println(helloDto)
}
