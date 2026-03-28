package config

// `{
//	"methodConfig": [{
//		"name": [{"service": "config.ConfigService", "method": "LongRunning"}],
//		"timeout": "10s"
//	}]
//}`

type Config struct {
	MethodConfig []MethodConfig `json:"methodConfig,omitempty"`
}

type MethodConfig struct {
	Name    []NameConfig `json:"name,omitempty"`
	Timeout string       `json:"timeout,omitempty"`
}

type NameConfig struct {
	Service string `json:"service,omitempty"`
	Method  string `json:"method,omitempty"`
}
