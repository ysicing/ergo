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
}
