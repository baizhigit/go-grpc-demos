package config

// `{
//	"methodConfig": [{
//		"name": [{"service": "config.ConfigService"}],
//		"timeout": "10s",
//		"retryPolicy": {
//		  "maxAttempts": 4,
//		  "initialBackoff": "0.1s",
//		  "maxBackoff": "1s",
//		  "backoffMultiplier": 2,
//		  "retryableStatusCodes": [
//			"INTERNAL", "UNAVAILABLE"
//		  ],
//	}]
//}`

type Config struct {
	MethodConfig []MethodConfig `json:"methodConfig,omitempty"`
}

type MethodConfig struct {
	Name        []NameConfig `json:"name,omitempty"`
	Timeout     string       `json:"timeout,omitempty"`
	RetryPolicy *RetryPolicy `json:"retryPolicy,omitempty"`
}

type RetryPolicy struct {
	MaxAttempts          int      `json:"maxAttempts"`
	InitialBackoff       string   `json:"initialBackoff"`
	MaxBackoff           string   `json:"maxBackoff"`
	BackoffMultiplier    float64  `json:"backoffMultiplier"`
	RetryableStatusCodes []string `json:"retryableStatusCodes"`
}

type NameConfig struct {
	Service string `json:"service,omitempty"`
	Method  string `json:"method,omitempty"`
}
