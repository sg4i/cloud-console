package provider

// Credential 表示一个凭证
type Credential struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Token     string `json:"token,omitempty"`
}

// SessionCredential 表示带有过期时间的会话凭证
type SessionCredential struct {
	Cred      *Credential `json:"cred"`
	ExpiredAt uint64      `json:"expired_at"`
}

// AssumeRoleOptions 定义 AssumeRole 的可选参数
type AssumeRoleOptions struct {
	RoleArn         string
	RoleSessionName string
	DurationSeconds uint64
}
