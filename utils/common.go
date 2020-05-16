// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"k8s.io/klog"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//DirIsEmpty 验证目录是否为空
func DirIsEmpty(dir string) bool {
	infos, err := ioutil.ReadDir(dir)
	if len(infos) == 0 || err != nil {
		return true
	}
	return false
}

//OpenOrCreateFile open or create file
func OpenOrCreateFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
}

//FileExists check file exist
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

//SearchFileBody 搜索文件中是否含有指定字符串
func SearchFileBody(filename, searchStr string) bool {
	body, _ := ioutil.ReadFile(filename)
	return strings.Contains(string(body), searchStr)
}

//IsHaveFile 指定目录是否含有文件
//.开头文件除外
func IsHaveFile(path string) bool {
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			return true
		}
	}
	return false
}

//SearchFile 搜索指定目录是否有指定文件，指定搜索目录层数，-1为全目录搜索
func SearchFile(pathDir, name string, level int) bool {
	if level == 0 {
		return false
	}
	files, _ := ioutil.ReadDir(pathDir)
	var dirs []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
			continue
		}
		if file.Name() == name {
			return true
		}
	}
	if level == 1 {
		return false
	}
	for _, dir := range dirs {
		ok := SearchFile(path.Join(pathDir, dir.Name()), name, level-1)
		if ok {
			return ok
		}
	}
	return false
}

//FileExistsWithSuffix 指定目录是否含有指定后缀的文件
func FileExistsWithSuffix(pathDir, suffix string) bool {
	files, _ := ioutil.ReadDir(pathDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			return true
		}
	}
	return false
}

//
func RandomString(lenstr int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if lenstr == 0 {
		lenstr = 8
	}
	b := make([]rune, lenstr)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomStringv2() string {
	t := time.Now().Format("2020010215")
	h := md5.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
}

func String2Int(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

func SysCmpOk(a, b, c string) bool {
	if String2Int(a)*String2Int(b) >= String2Int(c) {
		klog.Info(String2Int(a), String2Int(b), String2Int(c))
		return false
	}
	return true
}

func ErgoExit(s string) {
	logger.Error(s)
	os.Exit(1)
}

func WarningOs() {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		logger.Warn("或许不支持: ", runtime.GOOS, runtime.GOARCH)
	}
}

func WarningDocker() {
	Cmd("which", "docker")
}
