# develop

## proto

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/tencent.proto
```

## TODO

- [x] 模块设计
- [] 自动release，多平台发布
- [x] 命令行生成option，调用服务层生成实例启动
- [] option校验
- [x] 自动登录跳转
- [x] 指定配置文件路径
- []生产环境同源策略

- feature:

  - [] 支持多账号角色扮演
