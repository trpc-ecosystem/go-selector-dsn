# Data Source Name Selector

English | [简体中文](./README-zh.md)

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
