// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package util

import (
	"github.com/ysicing/ergo/common"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Copier struct{}

// CopyFile copies the contents of src to dst atomically.
func (c *Copier) CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "copyfile")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Chmod(tmp.Name(), common.FileMode0755); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}

func Copy(dst, src string) error {
	var c Copier
	return c.CopyFile(dst, src)
}
