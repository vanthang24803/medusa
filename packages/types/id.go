package types

import (
	"github.com/google/uuid"
)

// GenerateID creates a prefix_<UUIDv7> ID, e.g., prod_018f6f96-3c0f-7000-8000-000000000000
// Following Medusa conventions: prod_, variant_, cus_, ord_, cart_...
func GenerateID(prefix string) string {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to UUID v4 in case of a rare error
		return prefix + "_" + uuid.NewString()
	}
	return prefix + "_" + id.String()
}
