package flags

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ysicing/ergo/internal/pkg/k3s/types"
)

// BashCompEnvVarFlag cobra flag's annotation used for bind env to flag.
const BashCompEnvVarFlag = "cobra_annotation_bash_env_var_flag"

// FlagHackLookup hack lookup function.
// Borrowed from https://github.com/docker/machine/blob/master/commands/create.go#L267.
func FlagHackLookup(flagName string) string {
	// i.e. "-d" for "--driver"
	flagPrefix := flagName[1:3]

	// TODO: Should we support -flag-name (single hyphen) syntax as well?
	for i, arg := range os.Args {
		if strings.Contains(arg, flagPrefix) {
			// format '--driver foo' or '-d foo'
			if arg == flagPrefix || arg == flagName {
				if i+1 < len(os.Args) {
					return os.Args[i+1]
				}
			}

			// format '--driver=foo' or '-d=foo'
			if strings.HasPrefix(arg, flagPrefix+"=") || strings.HasPrefix(arg, flagName+"=") {
				return strings.Split(arg, "=")[1]
			}
		}
	}

	return ""
}

// ConvertFlags change autok3s flags to FlagSet, will mark required annotation if possible.
func ConvertFlags(cmd *cobra.Command, fs []types.Flag) *pflag.FlagSet {
	for _, f := range fs {
		if f.ShortHand == "" {
			if cmd.Flags().Lookup(f.Name) == nil {
				pf := cmd.Flags()
				switch t := f.V.(type) {
				case bool:
					pf.BoolVar(f.P.(*bool), f.Name, t, f.Usage)
				case string:
					pf.StringVar(f.P.(*string), f.Name, t, f.Usage)
				case map[string]string:
					pf.StringToStringVar(f.P.(*map[string]string), f.Name, t, f.Usage)
				case []string:
					pf.StringArrayVar(f.P.(*[]string), f.Name, t, f.Usage)
				case types.StringArray:
					pf.Var(newStringArrayValue(t, f.P.(*types.StringArray)), f.Name, f.Usage)
				case int:
					pf.IntVar(f.P.(*int), f.Name, t, f.Usage)
				default:
					continue
				}
				if f.Required {
					_ = cobra.MarkFlagRequired(pf, f.Name)
				}
			}
		} else {
			if cmd.Flags().Lookup(f.Name) == nil {
				pf := cmd.Flags()
				switch t := f.V.(type) {
				case bool:
					pf.BoolVarP(f.P.(*bool), f.Name, f.ShortHand, t, f.Usage)
				case string:
					pf.StringVarP(f.P.(*string), f.Name, f.ShortHand, t, f.Usage)
				case map[string]string:
					pf.StringToStringVarP(f.P.(*map[string]string), f.Name, f.ShortHand, t, f.Usage)
				case []string:
					pf.StringArrayVarP(f.P.(*[]string), f.Name, f.ShortHand, t, f.Usage)
				case types.StringArray:
					pf.VarP(newStringArrayValue(t, f.P.(*types.StringArray)), f.Name, f.ShortHand, f.Usage)
				default:
					continue
				}
				if f.Required {
					_ = cobra.MarkFlagRequired(pf, f.Name)
				}
			}
		}

		if f.EnvVar != "" {
			_ = cmd.Flags().SetAnnotation(f.Name, BashCompEnvVarFlag, []string{f.EnvVar})
		}
	}

	return cmd.Flags()
}

type stringArrayValue struct {
	value   *types.StringArray
	changed bool
}

func newStringArrayValue(val []string, p *types.StringArray) *stringArrayValue {
	ssv := new(stringArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringArrayValue) Set(val string) error {
	if !s.changed {
		*s.value = []string{val}
		s.changed = true
	} else {
		*s.value = append(*s.value, val)
	}
	return nil
}

func (s *stringArrayValue) Append(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringArrayValue) Replace(val []string) error {
	out := make([]string, len(val))
	for i, d := range val {
		out[i] = d
	}
	*s.value = out
	return nil
}

func (s *stringArrayValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = d
	}
	return out
}

func (s *stringArrayValue) Type() string {
	return "stringArray"
}

func (s *stringArrayValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

func writeAsCSV(ss []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(ss)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}
