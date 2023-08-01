package kart_grpc

type GrpcConfig struct {
	Port string `yaml:"port" json:"port" toml:"port"`
	Addr string `yaml:"addr" json:"addr" toml:"addr"`
}
