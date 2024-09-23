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

### REF

- [调用 GetSigninToken 使用安全令牌获取登录令牌](https://help.aliyun.com/document_detail/91913.html)
- [调用 Login 使用登录令牌登录阿里云控制台](https://help.aliyun.com/document_detail/91914.html)
- [控制台内嵌及分享](https://www.alibabacloud.com/help/zh/sls/developer-reference/embed-console-pages-and-share-log-data)
