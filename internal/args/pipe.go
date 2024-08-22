package args

import "time"

type Standard struct {
	Directory  string        `mapstructure:"KS_DIR"`
	Kubeconfig string        `mapstructure:"KUBECONFIG"`
	Timeout    time.Duration `mapstructure:"KS_TIMEOUT"`
}
