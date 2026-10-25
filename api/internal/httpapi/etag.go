package httpapi

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

// personETag identifies the current API representation. It is informational;
// enforcing optimistic concurrency requires a schema-backed row version.
func personETag(id int64, weeklyHours float64) string {
	value := strconv.FormatFloat(weeklyHours, 'g', -1, 64)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", id, value)))
	return fmt.Sprintf(`W/"%x"`, hash[:8])
}
