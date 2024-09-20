package config


type Credential struct {
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
	Token     string `json:"token"`
}
