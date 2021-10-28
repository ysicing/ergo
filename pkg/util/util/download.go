// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package util

import (
	"github.com/ysicing/ergo/common"
	"io/ioutil"
	"net/http"
)

func HttpGet(url, indexFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(indexFile, data, common.FileMode0600)
}
