package mapper_test

import (
	"fmt"
	"testing"

	"e.coding.net/guodf/gopkg/goutil/mapper"
)

// HelloDto dto
type HelloDto struct {
	Field1 string
	field2 string
}

// HelloModel model
type HelloModel struct {
	Field1 string
	field2 string
}

func TestDtoToModel(t *testing.T) {
	var helloDto HelloDto
	mapper.Mapper(HelloDto{
		Field1: "field1",
	}, &helloDto)
	fmt.Println(helloDto)
}
