# aws

## note

- 临时安全凭证

  - AssumeRole\*
  - GetFederationToken

推荐 AssumeRole

- 登录 URL 有效期

URL 自创建之日起 15 分钟内有效

- GovCloud

```golang
// Different URIs for GovCloud
if strings.HasPrefix(region, "us-gov-") {
  signinURI = "https://signin.amazonaws-us-gov.com/federation"
  consoleURI = "https://console.amazonaws-us-gov.com"
  signoutURI = "https://signin.amazonaws-us-gov.com/oauth?Action=logout&redirect_uri=https://amazonaws-us-gov.com"
}

```

## REF

- [允许自定义身份代理访问 AWS 控制台](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html)
