# alibaba

## note

[调用 Login 使用登录令牌登录阿里云控制台](https://help.aliyun.com/document_detail/91914.html)是阿里云集成转售解决方案，这个 API 是给虚商使用的，需要线上申请虚商伙伴业务，签订商务合同。参考文档的使用前提： [集成概述](https://help.aliyun.com/document_detail/91976.html)。经测试，sls4service.console.aliyun.com 服务针对账号类型没有要求，可以生成免密登录链接。

TicketType

- 若类型为 normal

  - DMS 域名： dms.aliyun.com
  - SLS 域名：sls.console.aliyun.com
  - 数据库自治服务: hdm.console.aliyun.com

- 类型为 mini，则一般应用于 BID 虚拟商
  - DMS 域名为 dms-jst4service.aliyun.com、 dms-Itwo4service.aliyun.com
  - SLS 域名：sls4service.console.aliyun.com
  - 数据库自治服务: hdm4service.console.aliyun.com

## REF

- [调用 GetSigninToken 使用安全令牌获取登录令牌](https://help.aliyun.com/document_detail/91913.html)
- [调用 Login 使用登录令牌登录阿里云控制台](https://help.aliyun.com/document_detail/91914.html)
- [控制台内嵌及分享](https://www.alibabacloud.com/help/zh/sls/developer-reference/embed-console-pages-and-share-log-data)
- [Flink 支持的登录方式](https://help.aliyun.com/zh/flink/user-guide/supported-logon-methods)
