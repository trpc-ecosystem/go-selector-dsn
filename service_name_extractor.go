package dsn

import (
	"errors"
	"strings"
)

// URIHostExtractor extracts host from URI, used for ip resolve(like get ip from polaris), work with ResolvableSelector
type URIHostExtractor struct {
}

// Extract extracts host from uri.
// Note: The uri has been preprocessed, no longer contains the strings preceding :// and ://.
func (e *URIHostExtractor) Extract(uri string) (int, int, error) {
	// mongodb+polaris://user:pswd@xxx.mongodb.com
	offset := 0

	// beginning of the host
	if idx := strings.LastIndex(uri, "@"); idx != -1 {
		uri = uri[idx+1:]
		offset += idx + 1
	}

	//  resolve end-part of the host
	begin := offset
	length, err := dealHostEndPart(uri)
	if err != nil {
		return 0, 0, err
	}
	uri = uri[0:length]

	return e.dealProtocolToken(uri, begin, length)
}

func dealHostEndPart(uri string) (int, error) {
	length := len(uri)
	if idx := strings.IndexAny(uri, "/?@"); idx != -1 {
		if uri[idx] == '@' {
			return 0, errors.New("parse host from uri: unescaped @ sign in user info")
		}
		if uri[idx] == '?' {
			return 0, errors.New("parse host from uri: must have a / before the query ?")
		}
		length = idx
	}
	return length, nil
}

func (e *URIHostExtractor) dealProtocolToken(uri string, begin, length int) (int, int, error) {
	begin, length = dealProtocolPrefix(uri, begin, length)
	length = dealProtocolSuffix(uri, length)
	return begin, length, nil
}

func dealProtocolPrefix(uri string, begin, length int) (int, int) {
	if strings.HasPrefix(uri, "tcp(") {
		return begin + 4, length - 4
	}
	return begin, length
}

func dealProtocolSuffix(uri string, length int) int {
	if strings.HasSuffix(uri, ")") {
		return length - 1
	}
	return length
}
