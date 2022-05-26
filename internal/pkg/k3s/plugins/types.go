package plugins

import "github.com/ysicing/ergo/internal/kube"

type Meta struct {
	Type    string `json:"type"`
	Default string `json:"default"`
	Item    []Item `json:"item"`
}

type Item struct {
	Client      *kube.Client `json:"-"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Version     string       `json:"version"`
	Home        string       `json:"home"`
	Appversion  string       `json:"appversion"`
	Type        string       `json:"type"`
	Path        string       `json:"path"`
	Tool        string       `json:"tool"`
	Namespace   string       `json:"namespace"`
}

type List []Meta
