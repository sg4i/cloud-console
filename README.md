# cloud console

- tencent cloud
- aws

## tencent cloud

### 设计

- 密钥读取优先级：命令行参数，环境变量、配置文件

  - 环境变量：TENCENTCLOUD_SECRET_ID、TENCENTCLOUD_SECRET_KEY、TENCENTCLOUD_TOKEN
  - 命令行参数：--secretId、--secretKey、--token
  - 配置文件搜索

    - ～/.cloud-console/default.credential /.tccli/default.configure
    - ～/.tccli/default.credential /.tccli/default.configure

      ```json
        {
            "secretId": "AKIDoOOCKBllZ8gxxMUO2Qm07HBuMZWCbnmU",
            "secretKey": "YonxrTnu64mDOK3YoEj11nNiqmaDeoou"
        }
        {
            "platform": {
                "output": "json",
                "region": "ap-guangzhou",
            }
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
