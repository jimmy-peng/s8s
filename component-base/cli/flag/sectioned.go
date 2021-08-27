package flag

import (
	"fmt"

	"github.com/spf13/pflag"
)

type NamedFlagSets struct {
	Order             []string
	FlagSets          map[string]*pflag.FlagSet
	NormalizeNameFunc func(f *pflag.FlagSet, name string) pflag.NormalizedName
}

// FlagSet returns the flag set with the given name and adds it to the
// ordered name list if it is not in there yet.
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}
	if _, ok := nfs.FlagSets[name]; !ok {
		flagSet := pflag.NewFlagSet(name, pflag.ExitOnError)
		flagSet.SetNormalizeFunc(pflag.CommandLine.GetNormalizeFunc())
		if nfs.NormalizeNameFunc != nil {
			flagSet.SetNormalizeFunc(nfs.NormalizeNameFunc)
		}
		nfs.FlagSets[name] = flagSet
		nfs.Order = append(nfs.Order, name)
	}
	return nfs.FlagSets[name]
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Println("Flag: --%s=%q", flag.Name, flag.Value)
	})
}
