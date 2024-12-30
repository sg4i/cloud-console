package grpc

import (
	"github.com/sg4i/cloud-console/internal/utils"
)

type GRPCConfig struct {
	Token       string
	RPCAddress  string
	HTTPAddress string
}

func New() *GRPCConfig {
	return &GRPCConfig{
		Token:       utils.GetString("grpc.token"),
		RPCAddress:  utils.GetString("grpc.rpcAddress"),
		HTTPAddress: utils.GetString("grpc.httpAddress"),
	}
}

func (c *GRPCConfig) GetToken() string {
	return c.Token
}

func (c *GRPCConfig) GetRPCAddress() string {
	return c.RPCAddress
}

func (c *GRPCConfig) GetHTTPAddress() string {
	return c.HTTPAddress
}
