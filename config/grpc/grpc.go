package grpc

import (
	"github.com/sg4i/cloud-console/internal/utils"
)

type GRPCConfig struct {
	Token string
}

func New() *GRPCConfig {
	return &GRPCConfig{
		Token: utils.GetString("rpc.token"),
	}
}

func (c *GRPCConfig) GetToken() string {
	return c.Token
}
