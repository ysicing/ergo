// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package base

const dockersh = `

curl -fsSL https://gitee.com/godu/install/raw/master/docker/get-docker.sh | bash -s docker --mirror Aliyun
mkdir /etc/docker
cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://dyucrs4l.mirror.aliyuncs.com"],
  "bip": "169.254.123.1/24",
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "30m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
EOF
systemctl enable docker
systemctl daemon-reload
systemctl start docker
systemctl restart docker 
docker info -f "{{json .ServerVersion }}"
docker pull registry.cn-beijing.aliyuncs.com/k7scn/tools
docker run --rm -v /usr/local/bin:/sysdir registry.cn-beijing.aliyuncs.com/k7scn/tools tar zxf /pkg.tgz -C /sysdir
# 腾讯云使用bip有问题 169.254.123.1/24, 请使用172.30.42.1/16
`

const mysql = `
mkdir -pv ~/svc/db

[ -f ~/svc/db/docker-compose.yaml ] && exit 0

cat > ~/svc/db/docker-compose.yaml <<EOF
version: '2.1'
services:
  mariadb:
    image: 'registry.cn-beijing.aliyuncs.com/k7scn/mariadb:10.5-debian-10'
    container_name: mariadb
    ports:
      - '3306:3306'
    volumes:
      - 'mariadb_data:/bitnami/mariadb'
    environment:
      # ALLOW_EMPTY_PASSWORD is recommended only for development.
      # - ALLOW_EMPTY_PASSWORD=yes
      - MARIADB_EXTRA_FLAGS=--max-connect-errors=1000 --max_connections=155
      - MARIADB_ROOT_PASSWORD=Eetohchi7aoGe8yaingai2eetahgoo9L
      - MARIADB_DATABASE=mydb
    healthcheck:
      test: ['CMD', '/opt/bitnami/scripts/mariadb/healthcheck.sh']
      interval: 15s
      timeout: 5s
      retries: 6

volumes:
  mariadb_data:
    driver: local
EOF

docker-compose -f ~/svc/db/docker-compose.yaml up -d
`

const redis = `
mkdir -pv ~/svc/redis

[ -f ~/svc/redis/docker-compose.yaml ] && exit 0

cat > ~/svc/redis/docker-compose.yaml <<EOF
version: '2'

services:
  redis:
    image: 'registry.cn-beijing.aliyuncs.com/k7scn/redis:6.0-debian-10'
    container_name: redis
    environment:
      # ALLOW_EMPTY_PASSWORD is recommended only for development.
      - REDIS_PASSWORD=ahphu9nah9iuheid1aew2eiPei6Ach
      - REDIS_DISABLE_COMMANDS=FLUSHDB,FLUSHALL
    ports:
      - '6379:6379'
    volumes:
      - 'redis_data:/bitnami/redis/data'

volumes:
  redis_data:
    driver: local
EOF

docker-compose -f ~/svc/redis/docker-compose.yaml up -d
`

const etcd = `
mkdir -pv ~/svc/etcd

[ -f ~/svc/etcd/docker-compose.yaml ] && exit 0

cat > ~/svc/etcd/docker-compose.yaml <<EOF
version: '2'

services:
  etcd:
    image: registry.cn-beijing.aliyuncs.com/k7scn/etcd:3-debian-10
    container_name: etcd
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
    ports:
      - 2379:2379
      - 2380:2380
    volumes:
      - etcd_data:/bitnami/etcd
volumes:
  etcd_data:
    driver: local
EOF

docker-compose -f ~/svc/etcd/docker-compose.yaml up -d
`

const adminer = `
mkdir -pv ~/svc/adminer

[ -f ~/svc/adminer/docker-compose.yaml ] && exit 0

cat > ~/svc/adminer/docker-compose.yaml <<EOF
version: '2.1'

services:
  adminer:
    container_name: adminer
    image: 'registry.cn-beijing.aliyuncs.com/k7scn/adminer'
    ports:
      - '127.0.0.1:10000:8080'
EOF

docker-compose -f ~/svc/adminer/docker-compose.yaml up -d
`

const prom = `
mkdir -pv ~/svc/prom

[ -f ~/svc/prom/docker-compose.yaml ] && exit 0
cat > ~/svc/prom/prometheus.yml <<EOF
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.

rule_files:

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
EOF

cat > ~/svc/prom/docker-compose.yaml <<EOF
version: '2'
services:
  prometheus:
    image: registry.cn-beijing.aliyuncs.com/k7scn/prometheus:2-debian-10
    container_name: prometheus
    volumes:
    - ./prometheus.yml:/opt/bitnami/prometheus/conf/prometheus.yml
    network_mode: host
EOF

docker-compose -f ~/svc/prom/docker-compose.yaml up -d
`

const grafana = `
mkdir -pv ~/svc/grafana

[ -f ~/svc/grafana/docker-compose.yaml ] && exit 0

cat > ~/svc/grafana/docker-compose.yaml <<EOF
version: '2'

services:
  grafana:
    image: registry.cn-beijing.aliyuncs.com/k7scn/grafana:7-debian-10
    container_name: grafana
    ports:
      - '3000:3000'
    environment:
      - 'GF_SECURITY_ADMIN_PASSWORD=admin'
      - 'GF_INSTALL_PLUGINS=grafana-clock-panel:1.1.0,grafana-kubernetes-app,farski-blendstat-panel,grafana-simple-json-datasource,yesoreyeram-boomtheme-panel'
    volumes:
      - grafana_data:/opt/bitnami/grafana/data
volumes:
  grafana_data:
    driver: local
EOF

docker-compose -f ~/svc/grafana/docker-compose.yaml up -d
`

const goscript = `
#!/bin/bash

go::install(){
    pushd /tmp
    # 下载
    wget https://dl.google.com/go/go1.16.6.linux-amd64.tar.gz
    # 解压
    tar -C /usr/local -xzf go1.16.6.linux-amd64.tar.gz
    popd
}

go::config(){
    cat >> /root/.bashrc <<EOF
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
export GOPATH="/root/go"
export GOBIN="$GOPATH/bin"
export PATH=$PATH:$GOBIN:/usr/local/go/bin
EOF

    source /root/.bashrc
}

go::test(){
    go env
}

go::install
go::config
go::test
`

const nodeexpoter = `
#!/bin/bash

set -e

dl(){
    pushd /tmp
    # 下载
	wget https://github.com/prometheus/node_exporter/releases/download/v1.2.0/node_exporter-1.2.0.linux-amd64.tar.gz
    # 解压
    tar -xzf node_exporter-1.2.0.linux-amd64.tar.gz
	rm -rf node_exporter-1.2.0.linux-amd64.tar.gz
	mv node_exporter-1.2.0.linux-amd64/node_exporter /usr/local/bin/
	chmod +x /usr/local/bin/node_exporter
    popd
}

systemd(){

cat > /etc/systemd/system/node-exporter.service <<EOF
[Unit]
Description=Node Exporter

[Service]
User=root
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
EOF

systemctl enable node-exporter
systemctl start node-exporter
}

dl
systemd
`
