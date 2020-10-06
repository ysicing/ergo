// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package vm

const dockersh = `

curl -fsSL https://gitee.com/godu/install/raw/master/docker/get-docker.sh | bash -s docker --mirror Aliyun
cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://reg-mirror.qiniu.com","https://dyucrs4l.mirror.aliyuncs.com"],
  "bip": "169.254.0.1/24",
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "20m",
    "max-file": "2"
  },
  "storage-driver": "overlay2"
}
EOF
systemctl enable docker
systemctl daemon-reload
systemctl start docker
systemctl restart docker 
docker info -f "{{json .ServerVersion }}"
docker run --rm -v /usr/local/bin:/sysdir registry.cn-beijing.aliyuncs.com/k7scn/tools tar zxf /pkg.tgz -C /sysdir
`
