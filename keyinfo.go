package panobi

import (
	"fmt"
	"strings"
)

const (
	idLen int = 22
)

// Holds information about a secret key.
type KeyInfo struct {
	K   string // actual key
	WID string // workspace ID
}

// Parses the given string, and returns a KeyInfo structure holding the
// component parts.
//
// You can find your key in your Panobi workspace's integration settings.
// Keys are in the format `WID-K`, where WID is the workspace ID, and K
// is actually the secret key generated for your integration.
func ParseKey(input string) (*KeyInfo, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf(errInvalidKey)
	}

	wid := strings.TrimSpace(parts[0])
	if len(wid) != idLen {
		return nil, fmt.Errorf(errInvalidKey)
	}

	k := strings.TrimSpace(parts[1])
	if k == "" {
		return nil, fmt.Errorf(errInvalidKey)
	}

	return &KeyInfo{
		K:   k,
		WID: wid,
	}, nil
}

// Test for equality against the given key information.
func (ki *KeyInfo) Equals(other *KeyInfo) bool {
	if other == nil {
		return false
	}

	return ki.K == other.K && ki.WID == other.WID
}
