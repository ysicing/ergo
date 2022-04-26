/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package experimental

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ergoapi/exgin"
	"github.com/gin-gonic/gin"
	"github.com/ysicing/ergo/pkg/util/log"
)

var (
	head = `<html>
	<form action="/newdir/" method="POST">
		<input type="text" name="dirname">
		<input type="submit" value="新建目录">
	</form>
	<form action="/" enctype="multipart/form-data" method="POST">
		<input type="file" name="files" multiple="multiple" />
		<input type="submit" value="上传" />
	</form>
	<script>
		window.onload = function() {
			var prelist = document.getElementsByTagName("pre");
			if (prelist.length > 0) {
				var pre = prelist[0];
				var alist = pre.getElementsByTagName("a");
				if (alist.length > 0) {
					for (var i = 0; i < alist.length; i++) {
						var a = alist[i];
						var input = document.createElement("input");
						input.type = "checkbox";
						input.name = "list[]";
						input.value = a.text;
						pre.insertBefore(input, a);
					}
					var input = document.createElement("input");
					input.type = "submit";
					input.value = "删除选中";
					pre.parentNode.append(input);
				}
			}
		}
    </script>
	<form action="/del/" method="POST">
`
	errFmt = `<font color="red">%s</font><br />
	<input type="button" value="返回" onclick="history.back()">`
)

type Files struct {
	List []string `form:"list[]"`
}

// SimpleFile 简单文件服务
func (exp *Options) SimpleFile() {
	g := exgin.Init(exp.SimpleFileCfg.Debug)
	g.Use(gin.Logger(), gin.Recovery(), exgin.ExCors())
	var authGroup *gin.RouterGroup
	if exp.SimpleFileCfg.User != "" {
		authGroup = g.Group("/", gin.BasicAuth(gin.Accounts{
			exp.SimpleFileCfg.User: exp.SimpleFileCfg.Pass,
		}))
	} else {
		authGroup = g.Group("/")
	}
	getGroup := authGroup.Group("/", writeHead)
	log.Flog.Debug("dir: ", exp.SimpleFileCfg.Dir)
	getGroup.StaticFS("/", gin.Dir(exp.SimpleFileCfg.Dir, true))
	authGroup.POST("/", exp.uploadFile)
	authGroup.POST("/del/", exp.delFile)
	authGroup.POST("/newdir/", exp.newDir)
	g.Run(":" + exp.SimpleFileCfg.Port)
}

func writeHead(c *gin.Context) {
	if strings.HasSuffix(c.Request.URL.Path, "/") {
		c.Writer.WriteString(head)
	}
}

func (exp *Options) uploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, fmt.Sprintf(errFmt, "Error: "+err.Error()))
		return
	}

	p := getPath(c.GetHeader("referer"))
	files := form.File["files"]
	for _, file := range files {
		err = c.SaveUploadedFile(file, path.Join(exp.SimpleFileCfg.Dir, p, path.Clean("/"+file.Filename)))
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, fmt.Sprintf(errFmt, file.Filename+" 上传失败<br />Error: "+err.Error()))
			return
		}
	}
	c.Redirect(302, p)
}

func (exp *Options) delFile(c *gin.Context) {
	var files Files
	err := c.ShouldBind(&files)
	if err != nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, fmt.Sprintf(errFmt, "Error: "+err.Error()))
		return
	}

	p := getPath(c.GetHeader("referer"))
	for _, file := range files.List {
		err = os.RemoveAll(path.Join(exp.SimpleFileCfg.Dir, p, path.Clean("/"+file)))
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, fmt.Sprintf(errFmt, file+" 删除失败<br />Error: "+err.Error()))
			return
		}
	}
	c.Redirect(302, p)
}

func (exp *Options) newDir(c *gin.Context) {
	p := getPath(c.GetHeader("referer"))
	dirname := c.PostForm("dirname")
	if dirname != "" {
		err := os.MkdirAll(path.Join(exp.SimpleFileCfg.Dir, p, path.Clean("/"+dirname)), os.ModeDir)
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(200, fmt.Sprintf(errFmt, "Error: "+err.Error()))
			return
		}
	}
	c.Redirect(302, p)
}

func getPath(referer string) string {
	var p = "/"
	u, err := url.ParseRequestURI(referer)
	if err == nil {
		// remove ../
		p = path.Clean("/" + u.Path)
		if p != "/" {
			p += "/"
		}
	}
	return p
}
