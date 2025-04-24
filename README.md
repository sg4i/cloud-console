# Cloud Console

Cloud Console 是一个命令行工具和服务，用于简化多云环境下控制台的访问。它支持通过角色扮演（AssumeRole）方式访问各大云服务商的控制台，无需手动登录过程。

## 支持的云服务提供商

- 腾讯云 (Tencent Cloud)
- 阿里云 (Alibaba Cloud)
- AWS (Amazon Web Services)

## 功能特点

- 支持命令行直接访问云控制台
- 支持配置文件管理多个云账号
- 支持 RPC 服务和 HTTP API 网关
- 自动生成并打开云控制台登录 URL
- 支持角色扮演 (AssumeRole) 访问

## 安装

### 从二进制文件安装

从 [GitHub Releases](https://github.com/sg4i/cloud-console/releases) 下载适合您系统的二进制文件，并将其添加到 PATH 环境变量中。

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/sg4i/cloud-console.git
cd cloud-console

# 构建项目
make build
```

## 使用方法

### 命令行模式

#### 腾讯云

```bash
# 基本使用
cloudconsole tencent --secret-id YOUR_SECRET_ID --secret-key YOUR_SECRET_KEY --role-arn YOUR_ROLE_ARN
```

#### 阿里云

```bash
# 基本使用
cloudconsole alibaba --access-key-id YOUR_ACCESS_KEY_ID --access-key-secret YOUR_ACCESS_KEY_SECRET --role-arn YOUR_ROLE_ARN
```

#### AWS

```bash
# 基本使用
cloudconsole aws --access-key-id YOUR_ACCESS_KEY_ID --secret-access-key YOUR_SECRET_ACCESS_KEY --role-arn YOUR_ROLE_ARN
```

### 使用配置文件

创建 `config.yml` 文件，参考以下格式：

```yaml
provider:
  tencent:
    credential:
      secretId: YOUR_SECRET_ID
      secretKey: YOUR_SECRET_KEY
      token: YOUR_TOKEN  # 可选
    roleArn: YOUR_ROLE_ARN # 可选
    destination: "https://console.cloud.tencent.com"
  alibaba:
    credential:
      secretId: YOUR_ACCESS_KEY_ID
      secretKey: YOUR_ACCESS_KEY_SECRET
      token: YOUR_TOKEN  # 可选
    roleArn: YOUR_ROLE_ARN # 可选
    loginUrl: "https://signin.aliyun.com/federation"
    destination: "https://console.aliyun.com"
  aws:
    credential:
      secretId: YOUR_ACCESS_KEY_ID
      secretKey: YOUR_SECRET_ACCESS_KEY
      token: YOUR_TOKEN  # 可选
    roleArn: YOUR_ROLE_ARN # 可选
    loginUrl: "https://signin.aws.amazon.com/federation"
    destination: "https://console.aws.amazon.com"

grpc:
  token: YOUR_RPC_TOKEN  # 可选，用于 RPC 服务认证
  rpcAddress: ":50050"
  httpAddress: ":50080"
```

然后使用以下命令：

```bash
# 使用默认配置文件
cloudconsole tencent

# 指定配置文件路径
cloudconsole tencent --config /path/to/your/config.yml
```

### 服务模式

启动 RPC 服务和 HTTP 网关：

```bash
# 使用默认配置
cloudconsole server

# 自定义 RPC 和 HTTP 地址
cloudconsole server --rpc-address ":8080" --http-address ":8081"

# 使用自定义配置文件
cloudconsole server --config /path/to/your/config.yml

# 仅启动 RPC 服务，不启动 HTTP 网关
cloudconsole server --no-http

# 设置认证令牌
cloudconsole server --token "your-auth-token"
```

HTTP API 示例:

```bash
curl -X POST http://localhost:50080/api/v1/role_login \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "tencent",
    "secret_id": "YOUR_SECRET_ID",
    "secret_key": "YOUR_SECRET_KEY",
    "role_arn": "YOUR_ROLE_ARN"
  }'
```

## 开发

### 生成 Proto 文件

```bash
make proto
```

### 构建

```bash
make build
```

### 构建 Docker 镜像

```bash
make docker-build
```


## 贡献

欢迎提交问题和合并请求。

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](./LICENSE) 文件。

