package panobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type SignatureInfo struct {
	S  string
	TS string
}

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
