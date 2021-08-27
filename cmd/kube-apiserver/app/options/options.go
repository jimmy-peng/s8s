package options

import (
	cliflag "s8s/component-base/cli/flag"
)

type ServerRunOptions struct {
	EnableLogsHandler bool
	MasterCount int
};

func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions {
		EnableLogsHandler: true,
		MasterCount: 1,
	}

	return &s
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.IntVar(&s.MasterCount, "apiserver-count", s.MasterCount,
		"The number of apiservers running in the cluster, must be a positive number. (In use when --endpoint-reconciler-type=master-count is enabled.)")
	fs.BoolVar(&s.EnableLogsHandler, "enable-logs-handler", s.EnableLogsHandler, "If true, install a /logs handler for the apiserver logs.")
	return fss
}