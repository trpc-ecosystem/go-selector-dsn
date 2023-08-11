// Package dsn data layer storage address method, mainly used for mysql mongodb and other database address
package dsn

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-go/naming/selector"
)

// SeletorName dsn selector name
var SeletorName string = "dsn"

func init() {
	selector.Register(SeletorName, DefaultSelector)
}

// DefaultSelector dsn default selector
var DefaultSelector = &DsnSelector{dsns: make(map[string]*registry.Node)}

// DsnSelector returns original service name node, with memory cache
type DsnSelector struct {
	dsns map[string]*registry.Node
	lock sync.RWMutex
}

// Select selects address from dsn://user:passwd@tcp(ip:port)/db
func (s *DsnSelector) Select(serviceName string, opt ...selector.Option) (*registry.Node, error) {
	if len(serviceName) == 0 {
		return nil, errors.New("dsn address can not be empty")
	}
	s.lock.RLock()
	node, ok := s.dsns[serviceName]
	s.lock.RUnlock()
	if ok {
		return node, nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	node, ok = s.dsns[serviceName]
	if ok {
		return node, nil
	}
	node = &registry.Node{
		ServiceName: serviceName,
		Address:     serviceName,
	}
	s.dsns[serviceName] = node
	return node, nil
}

// Report dsn selector no need to report
func (s *DsnSelector) Report(node *registry.Node, cost time.Duration, err error) error {
	return nil
}

// ResolvableSelector dsn-selector with address resolver
type ResolvableSelector struct {
	dsns                 map[string]*registry.Node
	lock                 sync.RWMutex
	resolverSelectorName string
	extractor            ServiceNameExtractor
}

// ServiceNameExtractor extracts the part of the service name in the dsn, and return the starting position and length
type ServiceNameExtractor interface {
	Extract(string) (int, int, error)
}

// NewResolvableSelector selector contains address resolution, implemented by other selectors. selectorName is the
// selector name for address resolution extractor is the func to extract the service name from the dsn,
// the extracted service name is the parameter of the selector used for address resolution
// eg：
// target: mongodb+polaris://user:passwd@poloars_name
// extractor extract polaris_name from target
// polaris selector will resolve polaris_name to address
func NewResolvableSelector(selectorName string, extractor ServiceNameExtractor) selector.Selector {
	return &ResolvableSelector{
		dsns:                 make(map[string]*registry.Node),
		resolverSelectorName: selectorName,
		extractor:            extractor,
	}
}

// Select selects address from mongodb+polaris://user:passwd@poloars_name/db
func (s *ResolvableSelector) Select(serviceName string, opt ...selector.Option) (*registry.Node, error) {
	// resolve serviceName from dsn
	pos, length, err := s.extractService(serviceName)
	if err != nil {
		return nil, err
	}
	extractedServiceName := serviceName[pos : pos+length]
	// selector select a available node
	resolvedNode, err := s.resolveAddress(extractedServiceName, opt...)
	if err != nil {
		return nil, err
	}
	address := serviceName[:pos] + resolvedNode.Address + serviceName[pos+length:]
	if len(address) == 0 {
		return nil, errors.New("dsn address can not be empty")
	}
	return s.dsnRW(address, serviceName, resolvedNode), nil
}

func (s *ResolvableSelector) extractService(serviceName string) (int, int, error) {
	if len(s.resolverSelectorName) == 0 {
		return 0, 0, errors.New("resolver selector name can not be empty")
	}
	if s.extractor == nil {
		return 0, 0, errors.New("service name extractor can not be nil")
	}
	pos, length, err := s.extractor.Extract(serviceName)
	if err != nil {
		return 0, 0, err
	}
	if length == 0 {
		return 0, 0, fmt.Errorf("the extracted service name is empty and the dsn is %s", serviceName)
	}
	return pos, length, nil
}

func (s *ResolvableSelector) resolveAddress(serviceName string, opt ...selector.Option) (*registry.Node, error) {
	resolver := selector.Get(s.resolverSelectorName)
	if resolver == nil {
		return nil, errors.New("unknown selector name " + s.resolverSelectorName)
	}
	return resolver.Select(serviceName, opt...)
}

func (s *ResolvableSelector) dsnRW(address, serviceName string, resolvedNode *registry.Node) *registry.Node {
	s.lock.RLock()
	node, ok := s.dsns[address]
	s.lock.RUnlock()
	if ok {
		return node
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	node, ok = s.dsns[address]
	if ok {
		return node
	}

	node = &registry.Node{
		ServiceName: serviceName,
		Address:     address,
		Metadata: map[string]interface{}{
			"resolved": resolvedNode,
		},
	}
	s.dsns[address] = node
	return node
}

// Report dsn selector does not need to report, but the embedded selector needs to report. this func is only executed
// after Select() is successful, so the node returned by Select() needs to be checked before report
func (s *ResolvableSelector) Report(node *registry.Node, cost time.Duration, err error) error {
	if node.Metadata == nil {
		return errors.New("metadata can not be nil")
	}

	resolved, ok := node.Metadata["resolved"]
	if !ok {
		return errors.New("the resolved in the metadata can not be nil")
	}

	resolvedNode, ok := resolved.(*registry.Node)
	if !ok {
		return errors.New("the resolved in the metadata is illegal")
	}

	resolver := selector.Get(s.resolverSelectorName)
	if resolver == nil {
		return errors.New("unknown selector name " + s.resolverSelectorName)
	}

	return resolver.Report(resolvedNode, cost, err)
}
