package panobi

import (
	"fmt"
	"strings"
)

type KeyInfo struct {
	K   string
	WID string
}

func ParseKey(input string) (*KeyInfo, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf(errInvalidKey)
	}

	wid := strings.TrimSpace(parts[0])
	if wid == "" {
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
