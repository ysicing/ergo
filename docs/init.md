## init

### init vm

```bash
Usage:
  ergo init vm [flags]

Flags:
  -h, --help            help for vm
      --path string     Vagrantfile所在目录, $HOME/vm
      --vmcpus string   虚拟机CPU数 (default "2")
      --vmmem string    虚拟机Mem MB数 (default "4096")
      --vmname string   虚拟机名
      --vmnum string    虚拟机副本数 (default "1")
```

### init debian

```bash
Usage:
  ergo init debian [flags]

Flags:
      --docker        是否安装docker
  -h, --help          help for debian
      --ip string     ssh ip (default "192.168.100.101")
      --pass string   管理员密码 (default "vagrant")
      --port string   ssh端口 (default "22")
      --user string   管理员用户 (default "root")
```