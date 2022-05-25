package main

import (
	"os"

	"github.com/BeidouCloudPlatform/go-bindata/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	os.Unsetenv("GOPATH")
	bc := &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "manifests/bin",
				Recursive: true,
			},
		},
		Package:    "data",
		Arch:       "arm64",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/data/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "manifests/bin",
				Recursive: true,
			},
		},
		Package:    "data",
		Arch:       "amd64",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/data/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:       "manifests/scripts",
				Recursive:  true,
				FileSuffix: ".sh",
			},
		},
		Package:    "scripts",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/scripts/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "manifests/plugins",
				Recursive: true,
			},
		},
		Package:    "plugins",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/plugins/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
}
