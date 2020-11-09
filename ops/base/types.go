// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package base

const dockersh = `

curl -fsSL https://gitee.com/godu/install/raw/master/docker/get-docker.sh | bash -s docker --mirror Aliyun
cat > /etc/docker/daemon.json <<EOF
{
  "registry-mirrors": ["https://reg-mirror.qiniu.com","https://dyucrs4l.mirror.aliyuncs.com"],
  "bip": "169.254.1.1/24",
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
docker pull registry.cn-beijing.aliyuncs.com/k7scn/tools
docker run --rm -v /usr/local/bin:/sysdir registry.cn-beijing.aliyuncs.com/k7scn/tools tar zxf /pkg.tgz -C /sysdir
`

const mysql = `
mkdir -pv ~/svc/db

[ -f ~/svc/db/docker-compose.yaml ] && exit 0

cat > ~/svc/db/docker-compose.yaml <<EOF
version: '2.1'
services:
  mariadb:
    image: 'registry.cn-beijing.aliyuncs.com/k7scn/mariadb:10.5-debian-10'
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

cat > ~/svc/etcd/docker-compose.yaml <<EOF
version: '2.1'

services:
  mariadb:
    image: 'registry.cn-beijing.aliyuncs.com/k7scn/adminer'
    ports:
      - '127.0.0.1:10000:8080'
EOF

docker-compose -f ~/svc/adminer/docker-compose.yaml up -d
`
