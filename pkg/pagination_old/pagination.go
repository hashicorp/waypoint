package pagination_old

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type PaginationRequestType int

const (
	FirstPage PaginationRequestType = iota
	PreviousPage
	NextPage
)

// Decodes and parses a base64 encoded string of the format 'key:value'
func DecodeAndParsePageToken(encodedPageToken string) (string, string, error) {
	var tokenKey string
	var tokenValue string
	tokenFormattingError := errors.New("Incorrectly formatted pagination token.")
	if encodedPageToken == "" {
		return "", "", nil
	}
	rawDecodedText, err := base64.StdEncoding.DecodeString(encodedPageToken)
	if err != nil {
		return "", "", tokenFormattingError // use generic formatting error to preserve opaque nature of pagination token
	}
	tokenList := strings.SplitN(string(rawDecodedText), ":", 2)
	if len(tokenList) != 2 {
		return tokenKey, tokenValue, tokenFormattingError
	}
	tokenKey = tokenList[0]
	tokenValue = tokenList[1]
	return tokenKey, tokenValue, nil
}

// base64 encodes a key and value
func EncodeAndSerializePageToken(key string, value string) (string, error) {
	if key == "" || value == "" {
		return "", nil
	}
	serializedPageToken := fmt.Sprintf("%s:%s", key, value)
	return base64.StdEncoding.EncodeToString([]byte(serializedPageToken)), nil
}
