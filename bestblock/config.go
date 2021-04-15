package bestblock

import "time"

type Config struct {
	HTTP         string
	WS           string
	PollInterval time.Duration `mapstructure:"poll-interval"`
}
