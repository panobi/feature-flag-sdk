package panobi

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	eventsURI string = "http://panobi.com/integrations/flags-sdk/events"

	attempts          int = 3
	backoffInitial    int = 1
	backoffMultiplier int = 2
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
	si, err := CalculateSignature(b, t.ki, nil)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", eventsURI, url.PathEscape(t.ki.ExternalID))
	backoff := backoffInitial
	i := 1

	for {
		err := func() error {
			req, err := http.NewRequest("POST", url, bytes.NewReader(b))
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
			case code == 408 || code == 429:
				if i == attempts {
					return fmt.Errorf("http error: %d", resp.StatusCode)
				}
				time.Sleep(getRetryAfter(resp, backoff))
				backoff = backoff * backoffMultiplier
				i++
				return nil
			default:
				return fmt.Errorf("http error: %d", resp.StatusCode)
			}
		}()

		if err != nil {
			return err
		}
	}
}

func (t *transport) getHeaders(si *SignatureInfo) http.Header {
	headers := make(http.Header)

	headers.Set("Content-Type", "application/json")
	headers.Set("X-Panobi-Signature", si.S)
	headers.Set("X-Panobi-Request-Timestamp", si.TS)

	return headers
}

func getRetryAfter(resp *http.Response, defaultRetryAfter int) time.Duration {
	headerVal := strings.TrimSpace(resp.Header.Get("Retry-After"))
	if headerVal == "" {
		return time.Duration(defaultRetryAfter) * time.Second
	}

	retryAfter, err := strconv.ParseInt(headerVal, 10, 64)
	if err != nil {
		return time.Duration(defaultRetryAfter) * time.Second
	}

	return time.Duration(retryAfter) * time.Second
}
