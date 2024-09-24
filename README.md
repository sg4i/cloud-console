# cloud console

- tencent cloud
- aws

## tencent cloud

### 设计

- 密钥及配置读取优先级：命令行参数，环境变量、配置文件

  - 命令行参数：--secretId、--secretKey、--token, --arn, --s-url
  - 环境变量：TENCENTCLOUD_SECRET_ID、TENCENTCLOUD_SECRET_KEY、TENCENTCLOUD_TOKEN、TENCENTCLOUD_ARN, TENCENTCLOUD_SURL
  - 配置文件搜索

    - ～/.cloud-console/default.credential /.tccli/default.configure

      ```json
        {
            "secretId": "AKIDoOOCKBllZ8gxxMUO2Qm07HBuMZWCbnmU",
            "secretKey": "YonxrTnu64mDOK3YoEj11nNiqmaDeoou",
            "token": ""
        }
        {
            "arn": "",
            "surl": "",
        },

      ```

API:

- UpdateRoleConsoleLogin

  ```json
  {
    "serviceType": "cam",
    "cmd": "UpdateRoleConsoleLogin",
    "data": {
      "Version": "2019-01-16",
      "RoleId": "4611686018440646060",
      "ConsoleLogin": 1,
      "Language": "zh-CN"
    },
    "regionId": 1
  }
  ```

### REF

- [角色免密登录控制台](https://cloud.tencent.com/document/product/598/45529)

## alibaba cloud

### 设计

- 密钥及配置读取优先级：命令行参数，环境变量、配置文件

- 命令行参数：--access-key-id、--access-key-secret、--security-token, --role-arn, --destination, --login-url

- 环境变量：ALIBABA_CLOUD_ACCESS_KEY_ID、ALIBABA_CLOUD_ACCESS_KEY_SECRET、ALIBABA_CLOUD_SECURITY_TOKEN
- 配置文件：

  ```yaml
  accessKeyId
  accessKeySecret
  ```

### note

[调用 Login 使用登录令牌登录阿里云控制台](https://help.aliyun.com/document_detail/91914.html)是阿里云集成转售解决方案，这个 API 是给虚商使用的，需要线上申请虚商伙伴业务，签订商务合同。参考文档的使用前提： [集成概述](https://help.aliyun.com/document_detail/91976.html)。经测试，sls4service.console.aliyun.com 服务针对账号类型没有要求，可以生成免密登录链接。

TicketType

- 若类型为normal
  - DMS域名： dms.aliyun.com
  - SLS域名：sls.console.aliyun.com
  - 数据库自治服务: hdm.console.aliyun.com

- 类型为mini，则一般应用于BID虚拟商
  - DMS域名为 dms-jst4service.aliyun.com、 dms-Itwo4service.aliyun.com
  - SLS域名：sls4service.console.aliyun.com
  - 数据库自治服务:  hdm4service.console.aliyun.com

### REF

- [调用 GetSigninToken 使用安全令牌获取登录令牌](https://help.aliyun.com/document_detail/91913.html)
- [调用 Login 使用登录令牌登录阿里云控制台](https://help.aliyun.com/document_detail/91914.html)
- [控制台内嵌及分享](https://www.alibabacloud.com/help/zh/sls/developer-reference/embed-console-pages-and-share-log-data)
- [Flink支持的登录方式](https://help.aliyun.com/zh/flink/user-guide/supported-logon-methods)
