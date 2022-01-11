module github.com/ysicing/ergo

go 1.16

require (
	github.com/ProtonMail/go-crypto v0.0.0-20211112122917-428f8eabeeb3 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1434
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/blang/semver v3.5.1+incompatible
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/containerd/continuity v0.2.2
	github.com/elazarl/goproxy v0.0.0-20210801061803-8e322dfb79c4 // indirect
	github.com/ergoapi/log v0.0.1
	github.com/ergoapi/util v0.1.4
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-logr/logr v1.2.1 // indirect
	github.com/go-ping/ping v0.0.0-20211130115550-779d1e919534
	github.com/gofrs/flock v0.8.1
	github.com/google/go-github/v39 v39.2.1-0.20211020014439-17a925b6f848
	github.com/gopasspw/gopass v1.13.0
	github.com/gosuri/uitable v0.0.4
	github.com/hashicorp/go-version v1.4.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/juju/ansiterm v0.0.0-20210929141451-8b71cc96ebdc // indirect
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213
	github.com/kardianos/service v1.2.1-0.20211104163826-b9d1d5b7279b
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-isatty v0.0.14
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6
	github.com/pkg/sftp v1.13.4
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shirou/gopsutil/v3 v3.21.12
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.327
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.327
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.328
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse v1.0.327
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr v1.0.327
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc v1.0.327
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/wangle201210/githubapi v0.0.0-20200804144924-cde7bbdc36ab
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/net v0.0.0-20220105145211-5b0dc2dfae98 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.7.2
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
	k8s.io/kubectl v0.23.1
	k8s.io/utils v0.0.0-20211208161948-7d6a63dca704 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.0 // indirect
	sigs.k8s.io/yaml v1.3.0
)

// replace github.com/google/go-github/v39 => ../go-github
// github.com/kardianos/service => ../service
replace github.com/kardianos/service v1.2.1-0.20211104163826-b9d1d5b7279b => github.com/BeidouCloudPlatform/service v1.2.1-0.20211108101837-1fc58f473016
