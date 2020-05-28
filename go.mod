module github.com/ysicing/ergo

go 1.14

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.174 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.1.0+incompatible
	github.com/cuisongliu/sshcmd v1.5.2
	github.com/ghodss/yaml v1.0.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/koding/logging v0.0.0-20160720134017-8b5a689ed69b // indirect
	github.com/koding/vagrantutil v0.0.0-20180710063911-70827343f116
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shirou/gopsutil v2.20.3+incompatible
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/wonderivan/logger v1.0.0
	k8s.io/klog v1.0.0
)

replace github.com/cuisongliu/sshcmd v1.5.2 => github.com/kunnos/sshcmd v1.6.0
