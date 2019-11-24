package profiler

import (
	_ "net/http/pprof"
)

type Config struct {
	Enable      bool
	BindAddress string
	Password    string
}

type Profiler struct {
}

func New(config Config) *Profiler {
	return &Profiler{}
}
