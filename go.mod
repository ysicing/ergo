module github.com/ysicing/ergo

go 1.16

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210920160938-87db9fbc61c7 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1317
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/blang/semver v3.5.1+incompatible
	github.com/ergoapi/log v0.0.0-20211027064103-103783bd0168
	github.com/ergoapi/util v0.0.9
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-ping/ping v0.0.0-20211014180314-6e2b003bffdd
	github.com/gofrs/flock v0.8.1
	github.com/google/go-github/v39 v39.2.1-0.20211020014439-17a925b6f848
	github.com/gopasspw/gopass v1.12.8
	github.com/gosuri/uitable v0.0.4
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/juju/ansiterm v0.0.0-20210929141451-8b71cc96ebdc // indirect
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6
	github.com/pkg/sftp v1.13.4
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shirou/gopsutil/v3 v3.21.10
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.281
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.281
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.281
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse v1.0.281
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr v1.0.281
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc v1.0.281
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/wangle201210/githubapi v0.0.0-20200804144924-cde7bbdc36ab
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	helm.sh/helm/v3 v3.7.1
	k8s.io/klog/v2 v2.30.0 // indirect
	k8s.io/kubectl v0.22.3
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b // indirect
	sigs.k8s.io/yaml v1.3.0
)

// replace github.com/google/go-github/v39 => ../go-github
