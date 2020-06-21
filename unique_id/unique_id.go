package unique_id

import (
	"github.com/google/uuid"
	"strings"
)

// NewUUID 生成新的UUID
func NewUUID(separated ...string) string {
	uID, _ := uuid.NewUUID()

	uIDStr := uID.String()
	if separated != nil {
		return strings.ReplaceAll(uIDStr, "-", separated[0])
	}
	return uIDStr
}