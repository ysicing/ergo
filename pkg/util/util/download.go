// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ergoapi/util/environ"
	"github.com/ysicing/ergo/common"
)

func HTTPGet(url, indexFile string) error {
	if strings.Contains(url, "github") && environ.GetEnv("NO_MIRROR") == "" {
		url = fmt.Sprintf("%v/%v", common.PluginGithubJiasu, url)
	}
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
