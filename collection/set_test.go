package collection

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_Set
func Test_Set(t *testing.T) {
	set := Set[string]{}
	set.Add("1")
	set.Add("2")
	str := set.String()
	assert.True(t, len(str) == 5)
	for _, v := range set {
		log.Println(v)
	}
	_, err := json.Marshal(set)
	assert.True(t, err == nil)
	json.Unmarshal([]byte(`["5","2","3","4"]`), &set)
	assert.True(t, set.Len() == 5)
}
