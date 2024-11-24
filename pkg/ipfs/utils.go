package ipfs

import (
	"errors"
	"strings"
)

// ValidateCID checks if the given CID looks valid
func ValidateCID(cid string) error {
	if len(cid) == 0 {
		return errors.New("CID cannot be empty")
	}
	if strings.ContainsAny(cid, " \t\n") {
		return errors.New("CID contains invalid characters")
	}
	return nil
}
