package panobi

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type transport struct {
	c  *http.Client
	ki *KeyInfo
}

func createTransport(ki *KeyInfo) *transport {
	return &transport{
		c:  &http.Client{},
		ki: ki,
	}
}

func (t *transport) post(b []byte) error {
	si, err := CalculateSignature(b, t.ki)
	if err != nil {
		return err
	}

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("POST", eventsURI, bytes.NewReader(b))
		if err != nil {
			return err
		}

		req.Header = t.getHeaders(si)

		resp, err := t.c.Do(req)
		if err != nil {
			return err
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "Error closing body:", err)
			}
		}()

		switch code := resp.StatusCode; {
		case code >= 200 && code < 300:
			return nil
		case code == 408:
			time.Sleep(time.Duration(i+1) * time.Second)
		case code == 429:
			time.Sleep(getRetryAfter(resp))
		default:
			return fmt.Errorf("http error: %d", resp.StatusCode)
		}
	}

	return fmt.Errorf("http error: too many retries")
}

func (t *transport) getHeaders(si *SignatureInfo) http.Header {
	headers := make(http.Header)

	headers.Set("Content-Type", "application/json")
	headers.Set("X-Panobi-Signature", si.S)
	headers.Set("X-Panobi-Request-Timestamp", si.TS)
	headers.Set("X-WID", t.ki.WID)

	return headers
}

func getRetryAfter(resp *http.Response) time.Duration {
	headerVal := strings.TrimSpace(resp.Header.Get("Retry-After"))
	if headerVal == "" {
		return defaultRetryAfter
	}

	retryAfter, err := strconv.ParseInt(headerVal, 10, 64)
	if err != nil {
		return defaultRetryAfter
	}

	return time.Duration(retryAfter) * time.Second
}
