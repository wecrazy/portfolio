// Package cryptoutil provides lightweight cryptographic helpers built on the
// Go standard library. It has no dependencies outside stdlib and is safe to
// import from any package in the project.
package cryptoutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// RandomHex returns a cryptographically secure random hex string produced from
// nBytes of random data (i.e. the resulting string is 2*nBytes characters long).
// Falls back to a time-based string on the unlikely event that crypto/rand
// fails, so callers never need to handle an error.
func RandomHex(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		// Extremely rare; crypto/rand only fails on severely broken systems.
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
