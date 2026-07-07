package args

import "time"

type Standard struct {
	Directory  string        `mapstructure:"KS_DIR"`
	DataDir    string        `mapstructure:"KS_DATA_DIR"`
	Kubeconfig string        `mapstructure:"KUBECONFIG"`
	Timeout    time.Duration `mapstructure:"KS_TIMEOUT"`
}
