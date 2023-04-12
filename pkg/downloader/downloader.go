// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ysicing/ergo/version"

	"github.com/cheggaaa/pb/v3"
	"github.com/containerd/continuity/fs"
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/validation"
	"github.com/ergoapi/util/zos"
	"github.com/mattn/go-isatty"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/util/log"
)

type Status = string

const (
	StatusUnknown    Status = ""
	StatusDownloaded Status = "downloaded"
	StatusSkipped    Status = "skipped"
	StatusUsedCache  Status = "used-cache"
)

type Result struct {
	Status Status
}

func Download(remote, local string) (*Result, error) {
	dlog := log.GetInstance()
	localPath, err := canonicalLocalPath(local)
	if err != nil {
		return nil, err
	}
	if fileinfo, err := os.Stat(localPath); err == nil {
		if fileinfo.Size() > 0 {
			dlog.Debugf("file %q already exists, skipping downloading from %q", localPath, remote)
			res := &Result{
				Status: StatusSkipped,
			}
			return res, nil
		}
	}
	localPathDir := filepath.Dir(localPath)
	if err := os.MkdirAll(localPathDir, common.FileMode0755); err != nil {
		return nil, err
	}
	if validation.IsLocal(remote) {
		if err := CopyLocal(localPath, remote); err != nil {
			return nil, err
		}
		res := &Result{
			Status: StatusDownloaded,
		}
		return res, nil
	}
	temp, _ := os.CreateTemp("", "")
	defer func() {
		os.Remove(temp.Name())
	}()
	if code, err := downloadHTTP(temp.Name(), remote, dlog); err != nil {
		if code != 403 {
			return nil, err
		}
		dlog.Warn("cdn 403, try use wget or curl")
		if err := downloadBySystem(temp.Name(), remote, dlog); err != nil {
			return nil, err
		}
	}
	if err := CopyLocal(localPath, temp.Name()); err != nil {
		return nil, err
	}
	res := &Result{
		Status: StatusDownloaded,
	}
	return res, nil
}

func CopyLocal(dst, src string) error {
	srcPath, err := canonicalLocalPath(src)
	if err != nil {
		return err
	}
	if dst == "" {
		// empty dst means caching-only mode
		return nil
	}
	dstPath, err := canonicalLocalPath(dst)
	if err != nil {
		return err
	}
	return fs.CopyFile(dstPath, srcPath)
}

func canonicalLocalPath(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("got empty path")
	}
	if !validation.IsLocal(s) {
		return "", fmt.Errorf("got non-local path: %q", s)
	}
	if strings.HasPrefix(s, "file://") {
		res := strings.TrimPrefix(s, "file://")
		if !filepath.IsAbs(res) {
			return "", fmt.Errorf("got non-absolute path %q", res)
		}
		return res, nil
	}
	return zos.HomeExpand(s)
}

// downloadByWget used wget
func downloadBySystem(localPath, url string, dlog log.Logger) error {
	if localPath == "" {
		return fmt.Errorf("downloadHTTP: got empty localPath")
	}
	if err := os.RemoveAll(localPath); err != nil {
		return err
	}
	url = proxyurl(url)
	dlog.Debugf("downloading %q into %q", url, localPath)
	wgetbin, err := exec.LookPath("wget")
	if err != nil {
		curlbin, err := exec.LookPath("curl")
		if err != nil {
			return fmt.Errorf("not found curl/wget")
		}
		curlargs := []string{url, "-s", "-L", "--user-agent", version.GetUG(), "-o", localPath}
		dlog.Debugf("curl %v", curlargs)
		if _, err := exec.Command(curlbin, curlargs...).CombinedOutput(); err != nil {
			return err
		}
		return nil
	}
	wgetagent := fmt.Sprintf("--user-agent=ERGO/%v", common.Version)
	wgetargs := []string{"-O", localPath, "-q", "-t=3", "-c", wgetagent, url}
	dlog.Debugf("wget %v", wgetargs)
	if _, err := exec.Command(wgetbin, wgetargs...).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

func downloadHTTP(localPath, url string, dlog log.Logger) (int, error) {
	if localPath == "" {
		return 0, fmt.Errorf("downloadHTTP: got empty localPath")
	}
	localPathTmp := localPath + ".tmp"
	if err := os.RemoveAll(localPathTmp); err != nil {
		return 0, err
	}
	fileWriter, err := os.Create(localPathTmp)
	if err != nil {
		return 0, err
	}
	defer fileWriter.Close()

	url = proxyurl(url)

	dlog.Debugf("downloading %q into %q", url, localPath)
	// resp, err := http.Get(url)
	client := &http.Client{
		Timeout: time.Minute * 5,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", version.GetUG())
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("expected HTTP status %d, got %s", http.StatusOK, resp.Status)
	}
	bar, err := createBar(resp.ContentLength)
	if err != nil {
		return 0, err
	}

	writers := []io.Writer{fileWriter}
	multiWriter := io.MultiWriter(writers...)

	bar.Start()
	if _, err := io.Copy(multiWriter, bar.NewProxyReader(resp.Body)); err != nil {
		return 0, err
	}
	bar.Finish()
	if err := fileWriter.Sync(); err != nil {
		return 0, err
	}
	if err := fileWriter.Close(); err != nil {
		return 0, err
	}
	if err := os.RemoveAll(localPath); err != nil {
		return 0, err
	}
	if err := os.Rename(localPathTmp, localPath); err != nil {
		return 0, err
	}

	return 0, nil
}

func createBar(size int64) (*pb.ProgressBar, error) {
	bar := pb.New64(size)

	bar.Set(pb.Bytes, true)
	if isatty.IsTerminal(os.Stdout.Fd()) {
		bar.SetTemplateString(`{{counters . }} {{bar . | green }} {{percent .}} {{speed . "%s/s"}}`)
		bar.SetRefreshRate(200 * time.Millisecond)
	} else {
		bar.Set(pb.Terminal, false)
		bar.Set(pb.ReturnSymbol, "\n")
		bar.SetTemplateString(`{{counters . }} ({{percent .}}) {{speed . "%s/s"}}`)
		bar.SetRefreshRate(5 * time.Second)
	}
	bar.SetWidth(80)
	if err := bar.Err(); err != nil {
		return nil, err
	}

	return bar, nil
}

func proxyurl(url string) string {
	if strings.Contains(url, "github") && environ.GetEnv("NO_MIRROR") == "" && !strings.HasPrefix(url, common.PluginGithubJiasu) {
		url = fmt.Sprintf("%v/%v", common.PluginGithubJiasu, url)
	}
	return url
}
