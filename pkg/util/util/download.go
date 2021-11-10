// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/environ"
	"github.com/ysicing/ergo/common"
)

func HTTPGet(url, indexFile string) error {
	exlog := log.GetInstance()
	if strings.Contains(url, "github") && environ.GetEnv("NO_MIRROR") == "" {
		url = fmt.Sprintf("%v/%v", common.PluginGithubJiasu, url)
	}
	exlog.Debugf("url: %v, path: %v", url, indexFile)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	indexFiletmp, err := os.OpenFile(indexFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, common.FileMode0600)
	if err != nil {
		return fmt.Errorf("failed to get file info: %s: %s", indexFile, err)
	}
	defer func() { _ = indexFiletmp.Close() }()
	_, err = io.Copy(indexFiletmp, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %s: %s", indexFile, err)
	}
	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// return ioutil.WriteFile(indexFile, data, common.FileMode0600)
	return nil
}
