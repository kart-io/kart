package options

type SecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{}
}
