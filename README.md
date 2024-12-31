# cloud console

- tencent cloud
- alibaba cloud
- aws

```shell
# 基本使用
$ cloud-console tencent --secret-id YOUR_SECRET_ID --secret-key YOUR_SECRET_KEY --role-arn YOUR_ROLE_ARN

# 指定自定义目标 URL
$ cloud-console tencent --secret-id YOUR_SECRET_ID --secret-key YOUR_SECRET_KEY --role-arn YOUR_ROLE_ARN --destination "https://custom.console.url"
```

配置读取优先级：
命令行参数优先
环境变量
配置文件
默认值
