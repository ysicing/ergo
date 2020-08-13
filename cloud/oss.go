// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cloud

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"k8s.io/klog"
	"strings"
)

func (ali CloudConfig) AliOss() {
	region := ali.AliRegionID[0]
	if !strings.HasSuffix(region, "aliyuncs.com") {
		region = fmt.Sprintf("oss-%v.aliyuncs.com", region)
	}
	client, err := oss.New(ali.AliRegionID[0], ali.AliKey, ali.AliSecret)
	if err != nil {
		klog.Exit("create oss client err: ", err)
	}
	bucket, err := client.Bucket(ali.OssBucket.Bucket)
	if err != nil {
		klog.Exit("get bucket err:", err)
	}
	if err = bucket.PutObjectFromFile(ali.OssBucket.Remote, ali.OssBucket.Local); err != nil {
		klog.Exit("upload file err:", err)
	}
	klog.Info("upload done")
}

func AliOssUpload() {
	oss := CloudConfig{
		AliKey:      AliKey,
		AliSecret:   AliSecret,
		AliRegionID: AliRegionID,
		OssBucket: Ossbucket{
			Bucket: OssBucket,
			Local:  OssLocal,
			Remote: OssRemote,
		},
	}
	oss.AliOss()
}
