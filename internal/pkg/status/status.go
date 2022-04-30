package status

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ergoapi/util/color"
	"github.com/ysicing/ergo/pkg/util/output"
)

type MapCount map[string]int

type MapMapCount map[string]MapCount

type PodStateMap map[string]PodStateCount

type Status struct {
	output     string     `json:"-" yaml:"-"`
	KubeStatus KubeStatus `json:"k8s" yaml:"k8s"`
}

type KubeStatus struct {
	Version   string      `json:"version" yaml:"version"`
	Type      string      `json:"type" yaml:"type"`
	NodeCount MapCount    `json:"nodes" yaml:"nodes"`
	PodState  PodStateMap `json:"service,omitempty" yaml:"service,omitempty"`
}

type PodStateCount struct {
	// Type is the type of deployment ("Deployment", "DaemonSet", ...)
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Desired is the number of desired pods to be scheduled
	Desired int `json:"desired,omitempty" yaml:"desired,omitempty"`

	// Ready is the number of ready pods
	Ready int `json:"ready,omitempty" yaml:"ready,omitempty"`

	// Available is the number of available pods
	Available int `json:"available,omitempty" yaml:"available,omitempty"`

	// Unavailable is the number of unavailable pods
	Unavailable int `json:"unavailable,omitempty" yaml:"unavailable,omitempty"`

	Disabled bool `json:"disabled" yaml:"disabled"`
}

func newStatus(output string) *Status {
	return &Status{
		output: output,
		KubeStatus: KubeStatus{
			NodeCount: MapCount{},
			Version:   "unknow",
			PodState:  PodStateMap{},
		},
	}
}

func (s *Status) Format() error {
	switch strings.ToLower(s.output) {
	case "json":
		return output.EncodeJSON(os.Stdout, s)
	case "yaml":
		return output.EncodeYAML(os.Stdout, s)
	default:
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)
		// k8s
		fmt.Fprintf(w, "Cluster Status: \n")
		fmt.Fprintf(w, "%s\t%s\n", "version", s.KubeStatus.Version)
		fmt.Fprintf(w, "%s\t%s\n", "mode", s.KubeStatus.Type)
		if s.KubeStatus.NodeCount["ready"] > 0 {
			fmt.Fprintf(w, "%s\t%s\n", "status", color.SGreen("health"))
		} else {
			fmt.Fprintf(w, "%s\t%s\n", "status", color.SRed("unhealth"))
		}
		for name, state := range s.KubeStatus.PodState {
			if state.Disabled {
				fmt.Fprintf(w, "%s\t%s\n", name, color.SBlue("disabled"))
			} else {
				if state.Available > 0 {
					fmt.Fprintf(w, "%s\t%s\n", name, color.SGreen("ok"))
				} else {
					fmt.Fprintf(w, "%s\t%s\n", name, color.SRed("warn"))
				}
			}
		}
		fmt.Fprintf(w, "\n")
		w.Flush()
		return output.EncodeText(os.Stdout, buf.Bytes())
	}
}
