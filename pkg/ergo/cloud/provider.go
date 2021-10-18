// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package cloud

type ProviderType string

func (s ProviderType) Value() string {
	return string(s)
}

const (
	ProviderAliyun = ProviderType("aliyun")
	ProviderQcloud = ProviderType("qcloud")
)

var CloudType = []struct {
	Key   string
	Value string
}{
	{
		Key:   "aliyun",
		Value: "阿里云",
	}, {
		Key:   "qcloud",
		Value: "腾讯云",
	},
	{
		Key:   "all",
		Value: "所有云服务商",
	},
}
