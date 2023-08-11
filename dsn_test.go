package dsn_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-go/naming/selector"
	dsn "trpc.group/trpc-go/trpc-selector-dsn"
)

func TestDsn(t *testing.T) {

	s := selector.Get("dsn")
	assert.NotNil(t, s)

	node, err := s.Select("user:passwd@tcp(ip:port)/db")
	assert.Nil(t, err)

	assert.Equal(t, "user:passwd@tcp(ip:port)/db", node.Address)

	_, err = s.Select("")
	assert.EqualErrorf(t, err, "dsn address can not be empty", "empty src test")

	err = s.Report(nil, time.Second, nil)
	assert.Nil(t, err)
}

type FakeSelector struct {
	Addrs map[string]string
}

func (s *FakeSelector) Select(serviceName string, opt ...selector.Option) (*registry.Node, error) {
	addr, ok := s.Addrs[serviceName]
	if !ok {
		return nil, fmt.Errorf("unknown service name %s", serviceName)
	}

	return &registry.Node{
		Address:     addr,
		ServiceName: serviceName,
	}, nil
}

func (s *FakeSelector) Report(node *registry.Node, cost time.Duration, success error) error {
	node.Metadata = map[string]interface{}{
		"reported": node,
	}
	return nil
}

type serviceNameExtractor struct {
}

func (e *serviceNameExtractor) Extract(dsn string) (int, int, error) {
	idx := strings.Index(dsn, "@")
	if idx+1 >= len(dsn) {
		return 0, 0, fmt.Errorf("extract service name failed, src is %s", dsn)
	}
	pos := idx + 1
	length := len(dsn) - pos
	return pos, length, nil
}

func TestResolvableSelector(t *testing.T) {
	assert := assert.New(t)

	selector.Register("empty", dsn.NewResolvableSelector("", &serviceNameExtractor{}))
	selector.Register("noextractor", dsn.NewResolvableSelector("test", nil))
	selector.Register("fake", &FakeSelector{Addrs: map[string]string{"abc": "127.0.0.1:8080"}})
	selector.Register("dsn+fake", dsn.NewResolvableSelector("fake", &serviceNameExtractor{}))

	s := selector.Get("dsn+fake")
	assert.NotNil(t, s)

	cases := []struct {
		src  string
		dest string
		err  string
	}{
		{
			src:  "@abc",
			dest: "@127.0.0.1:8080",
		},
		{
			src:  "abc@",
			dest: "",
			err:  "extract service name failed, src is abc@",
		},
		{
			src:  "@abc/",
			dest: "@",
			err:  "unknown service name abc/",
		},
	}

	for _, c := range cases {
		node, err := s.Select(c.src)
		if len(c.err) == 0 {
			assert.Nilf(err, "case: %+v", c)
			assert.Equalf(c.dest, node.Address, "case: %+v", c)
			assert.Nilf(s.Report(node, node.CostTime, nil), "case: %+v", c)
			resolved, _ := node.Metadata["resolved"].(*registry.Node)
			assert.Equalf(resolved.Metadata["reported"], resolved, "case: %+v", c)
		} else {
			assert.EqualErrorf(err, c.err, "case: %+v", c)
		}
	}
	// compare with no cache and cache, hit cache, same node
	src := "user@abc"
	dest := "user@127.0.0.1:8080"
	nodeWithNoCache, err := s.Select(src)
	assert.Nil(err)
	assert.Equal(dest, nodeWithNoCache.Address)
	nodeWithCache, err := s.Select(src)
	assert.Nil(err)
	assert.True(nodeWithCache == nodeWithNoCache)

	emptyS := selector.Get("empty")
	_, err = emptyS.Select("123")
	assert.EqualError(err, "resolver selector name can not be empty", "case: empty")

	noExtractorS := selector.Get("noextractor")
	_, err = noExtractorS.Select("123")
	assert.EqualError(err, "service name extractor can not be nil", "case: empty")
}
