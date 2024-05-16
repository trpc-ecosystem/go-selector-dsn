//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

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
	length := len(uri)
	if idx := strings.IndexAny(uri, "/?@"); idx != -1 {
		if uri[idx] == '@' {
			return 0, 0, errors.New("parse host from uri: unescaped @ sign in user info")
		}
		if uri[idx] == '?' {
			return 0, 0, errors.New("parse host from uri: must have a / before the query ?")
		}
		length = idx
	}
	uri = uri[0:length]

	return e.dealProtocolToken(uri, begin, length)
}

func (e *URIHostExtractor) dealProtocolToken(uri string, begin, length int) (int, int, error) {
	if strings.HasPrefix(uri, "tcp(") {
		begin += 4
		length -= 4
	}
	if strings.HasSuffix(uri, ")") {
		length--
	}
	return begin, length, nil
}
