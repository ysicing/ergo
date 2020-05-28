// MIT License
// Copyright (c) 2019 ysicing <i@ysicing.me>

package cloud

var (
	AliRegionID []string
	AliKey      string
	AliSecret   string
	OssBucket   string
	OssRemote   string
	OssLocal    string
)

type CloudConfig struct {
	AliRegionID []string
	AliKey      string
	AliSecret   string
	OssBucket   Ossbucket
}

type Ossbucket struct {
	Bucket string
	Remote string
	Local  string
}
