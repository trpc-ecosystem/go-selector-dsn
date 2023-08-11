# Data Source Name 寻址方式

[English](./README.md) | 简体中文

实现了 tRPC-Go 的 selector 接口，将 client 的 target 当做 data source name 来使用，并在 Node 的 Address 中返回

## client
```
client:                                            # 客户端调用的后端配置
  service:                                         # 针对后端的配置
    - name: trpc.dsn.xxx.xxx         
      target: dsn://user:passwd@tcp(ip:port)/db      # select 返回的 address 为"://"之后的部分
    - name: trpc.dsn.xxx.xxx1         
      # dsn+polaris 表示 target 中的 host 会进行北极星解析，然后用实际地址替换 host 后返回“://”之后的部分
      # polaris 是在注册 selector 时指定的，也可以为其它的 selector
      target: dsn+polaris://user:passwd@tcp(host)/db
```

```
// 注册 selector
func init() {
    // 直接用 target 作为 data source name 或 uri
    selector.Register("dsn", dsn.DefaultSelector)

    // 支持地址解析的 selector，polaris 为地址解析的 selector 名称，
    // dsn.URIHostExtractor{}为从 target 中提取 polaris 服务的 key 的提取器
    selector.Register("dsn+polaris", dsn.NewResolvableSelector("polaris", &dsn.URIHostExtractor{}))
}

```
