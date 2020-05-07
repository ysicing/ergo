## install

### docker,go,tools

安装docker

```bash
Usage:
  ergo install docker/go/tools [flags]

Flags:
  -h, --help   help for docker

Global Flags:
      --config string   config file (default is $HOME/.doge/config.yaml)
      --ip strings      需要安装节点ip (default [192.168.100.101])
      --pass string     管理员密码 (default "vagrant")
      --pk string       管理员私钥
      --user string     管理员 (default "root")
```

### k8s

```bash
Usage:
  ergo install k8s [flags]

Flags:
      --enableingress    k8s启用ingress (default true)
      --enablenfs        k8s启用nfs sc
      --exnfs string     外部nfs地址, 若无则为空
  -h, --help             help for k8s
      --mip string       管理节点ip,eg ip或者ip-ip (default "11.11.11.111")
      --nfspath string   nfs路径 (default "/k8sdata")
      --nfssc string     默认nfs storageclass (default "nfs-data")
      --wip string       计算节点ip,eg ip或者ip-ip

Global Flags:
      --config string   config file (default is $HOME/.doge/config.yaml)
      --ip strings      需要安装节点ip (default [11.11.11.111])
      --pass string     管理员密码 (default "vagrant")
      --pk string       管理员私钥
      --user string     管理员 (default "root")

```

