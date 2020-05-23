## vm

```bash
Usage:
  ergo vm [command]

Available Commands:
  create      创建debian virtualbox虚拟机
  init        初始化debian
  reinstall   重装debian

```

### vm create 创建debian环境

```bash
Usage:
  ergo vm create [flags]

Flags:
  -h, --help            help for create
      --path string     Vagrantfile所在目录, $HOME/vm
      --vmcpus string   虚拟机CPU数 (default "2")
      --vmmem string    虚拟机Mem MB数 (default "4096")
      --vmname string   虚拟机名
      --vmnum string    虚拟机副本数 (default "1")

```

### init debian

```bash
Usage:
  ergo vm init [flags]

Flags:
      --docker        是否安装docker
  -h, --help          help for init
      --ip strings    ssh ip (default [11.11.11.111])
      --pass string   管理员密码 (default "vagrant")
      --port string   ssh端口 (default "22")
      --user string   管理员用户 (default "root")
```

## reinstall debian

```bash
Usage:
  ergo vm reinstall [flags]

Flags:
  -h, --help            help for reinstall
      --ip strings      ssh ip (default [11.11.11.111])
      --local           本地安装
      --pass string     管理员密码
      --pk string       管理员私钥
      --redisk string   自定义硬盘,如/dev/sdb
      --repass string   默认重装密码 (default "vagrant")
      --user string     管理员用户

```