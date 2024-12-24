package grpc

import (
	"github.com/sg4i/cloud-console/internal/utils"
)

type RPCConfig struct {
	Token string
}

func New() *RPCConfig {
	utils.LoadConfig()

	return &RPCConfig{
		Token: utils.GetString("rpc.token"),
	}
}

func (c *RPCConfig) GetToken() string {
	return c.Token
}
