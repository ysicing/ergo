package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Metadata struct {
	Name            string      `json:"name" yaml:"name"`
	Provider        string      `json:"provider" yaml:"provider"`
	EIP             string      `json:"eip,omitempty" yaml:"eip,omitempty"`
	TLSSans         StringArray `json:"tls-sans,omitempty" yaml:"tls-sans,omitempty"`
	ClusterCidr     string      `json:"cluster-cidr,omitempty" yaml:"cluster-cidr,omitempty"`
	ServiceCidr     string      `json:"service-cidr,omitempty" yaml:"service-cidr,omitempty"`
	MasterExtraArgs string      `json:"master-extra-args,omitempty" yaml:"master-extra-args,omitempty"`
	WorkerExtraArgs string      `json:"worker-extra-args,omitempty" yaml:"worker-extra-args,omitempty"`
	DataStore       string      `json:"datastore,omitempty" yaml:"datastore,omitempty"`
	Network         string      `json:"network,omitempty" yaml:"network,omitempty"`
	Plugins         StringArray `json:"plugins,omitempty" yaml:"plugins,omitempty"`
	Mode            string      `json:"mode,omitempty" yaml:"mode,omitempty"`
}

type Status struct {
	Status string `json:"status,omitempty"`
}

// Flag struct for flag.
type Flag struct {
	Name      string
	P         interface{}
	V         interface{}
	ShortHand string
	Usage     string
	Required  bool
	EnvVar    string
}

type StringArray []string

// Scan gorm Scan implement.
func (a *StringArray) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		if v != "" {
			*a = strings.Split(v, ",")
		}
	default:
		return fmt.Errorf("failed to scan array value %v", value)
	}
	return nil
}

// Value gorm Value implement.
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return strings.Join(a, ","), nil
}
