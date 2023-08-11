package dsn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHost(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		uri  string // The uri has been preprocessed, no longer contains the strings preceding :// and ://.
		host string
		err  string
	}{
		{
			uri:  "localhost",
			host: "localhost",
		},
		{
			uri:  "admin:123456@localhost/",
			host: "localhost",
		},
		{
			uri:  "admin:123456@localhost",
			host: "localhost",
		},
		{
			uri:  "example1.com:27017,example2.com:27017",
			host: "example1.com:27017,example2.com:27017",
		},
		{
			uri:  "host1,host2,host3/?slaveOk=true",
			host: "host1,host2,host3",
		},
		{
			uri: "admin:123456@localhost?",
			err: "parse host from uri: must have a / before the query ?",
		},
		{
			uri:  "user:secret@localhost:6379/0?foo=bar&qux=baz", // redis示例
			host: "localhost:6379",
		},
		{
			uri:  "user:secretWith@secretWith@localhost:6379/0?foo=bar&qux=baz", // 密码包含@符号示例
			host: "localhost:6379",
		},
		{
			uri:  "user:secret@tcp(localhost:6379)/database?timeout=1s&interpolateParams=true", // mysql示例
			host: "localhost:6379",
		},
	}

	extractor := new(URIHostExtractor)
	for _, c := range cases {
		pos, length, err := extractor.Extract(c.uri)
		if len(c.err) != 0 {
			assert.EqualErrorf(err, c.err, "case: %+v ", c)
		} else {
			assert.Equalf(c.host, c.uri[pos:pos+length], "case: %+v", c)
		}
	}
}
