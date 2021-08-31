package options

import (
	"fmt"
	cliflag "s8s/component-base/cli/flag"
	kubeoptions "s8s/pkg/kubeapiserver/options"
)

type ServerRunOptions struct {
	EnableLogsHandler bool
	MasterCount       int
	CloudProvider     *kubeoptions.CloudProviderOptions
}

func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{
		CloudProvider:     kubeoptions.NewCloudProviderOptions(),
		EnableLogsHandler: true,
		MasterCount:       1,
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

func (s *ServerRunOptions) Validate() []error {
	var errs []error
	if s.MasterCount <= 0 {
		errs = append(errs, fmt.Errorf("--apiserver-count should be a positive number"))
	}
	errs = append(errs, s.CloudProvider.Validates())
	return errs
}
