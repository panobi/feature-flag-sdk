package panobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Holds information about a signature.
type SignatureInfo struct {
	S  string // the signature itself, calculated from a payload
	TS string // unix milliseconds at which it was calculated
}

// Calculates a signature for the given byte payload, using the given key
// information. The events endpoint requires that you include the calculated
// signature and timestamp when making requests.
func CalculateSignature(b []byte, ki *KeyInfo) (*SignatureInfo, error) {
	if len(b) > maxInputBytes {
		return nil, fmt.Errorf(errMaxNumberSize, "input", maxInputBytes, "bytes")
	}

	ts := fmt.Sprint(time.Now().UnixMilli())

	message := fmt.Sprintf("%s:%s:%s", "v0", ts, b)
	mac := hmac.New(sha256.New, []byte(ki.K))
	mac.Write([]byte(message))
	signature := "v0=" + string(hex.EncodeToString(mac.Sum(nil)))

	return &SignatureInfo{
		S:  signature,
		TS: ts,
	}, nil
}
