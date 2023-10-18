English | [中文](README.zh_CN.md)

# Data Source Name Selector

[![Go Reference](https://pkg.go.dev/badge/github.com/trpc-ecosystem/go-selector-dsn.svg)](https://pkg.go.dev/github.com/trpc-ecosystem/go-selector-dsn)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-naming-polarismesh)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-naming-polarismesh)
[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://github.com/trpc-ecosystem/go-selector-dsn/blob/main/LICENSE)
[![Releases](https://img.shields.io/github/release/trpc-ecosystem/go-selector-dsn.svg?style=flat-square)](https://github.com/trpc-ecosystem/go-selector-dsn/releases)
[![Tests](https://github.com/trpc-ecosystem/go-selector-dsn/actions/workflows/prc.yml/badge.svg)](https://github.com/trpc-ecosystem/go-selector-dsn/actions/workflows/prc.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-selector-dsn/branch/main/graph/badge.svg)](https://app.codecov.io/gh/trpc-ecosystem/go-selector-dsn/tree/main)

DSN(Data Source Name) Selector implements a selector for tRPC-Go, which uses the client's target as a data source name , and returns it in the Node's Address

## client
```
client:                                            # backend-config for client
  service:                                         # backend's config
    - name: trpc.dsn.xxx.xxx         
      target: dsn://user:passwd@tcp(ip:port)/db      # select retruns the address after "://"
    - name: trpc.dsn.xxx.xxx1         
      # dsn+polaris means that the host in target will be resolved by polaris, and the actual address will be replaced 
      # after the host, and the part after "://" will be returned
      # polaris is specified when registering the selector, and can also be other selectors
      target: dsn+polaris://user:passwd@tcp(host)/db
```

```
// register selector
func init() {
    // use target as data source name or uri directly
    selector.Register("dsn", dsn.DefaultSelector)

    // selector which supports address resolution, polaris is the name of the address resolution selector
    // dsn.URIHostExtractor{} is the extractor to extract the key of polaris service from target
    selector.Register("dsn+polaris", dsn.NewResolvableSelector("polaris", &dsn.URIHostExtractor{}))
}

```
