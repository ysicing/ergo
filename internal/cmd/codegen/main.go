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
				Path:      "hack/bin",
				Recursive: true,
			},
		},
		Package:    "static",
		Arch:       "arm64",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/pkg/static/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "hack/bin",
				Recursive: true,
			},
		},
		Package:    "static",
		Arch:       "amd64",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/pkg/static/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:       "hack/scripts",
				Recursive:  true,
				FileSuffix: ".sh",
			},
		},
		Package:    "static",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/pkg/static/zz_generated_scriptdata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
}
