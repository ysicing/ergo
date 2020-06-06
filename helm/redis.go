// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"fmt"
	"github.com/wonderivan/logger"
)

type RedisMeta struct {
	redismeta HelmMeta
}

func (redis RedisMeta) CheckHelm() {
	logger.Info("检查ns %v 是否存在", redis.redismeta.NameSpace)
	checkns := fmt.Sprintf("kubectl create ns %v || sleep 1", redis.redismeta.NameSpace)
	SSHConfig.Cmd(redis.redismeta.Host, checkns)
	logger.Info("检查 helmservice %v 是否存在", fmt.Sprintf("redis-%v", redis.redismeta.NameSpace))
	checkhs := fmt.Sprintf("helm list -n %v  | grep redis-%v && exit 1 || sleep 0", redis.redismeta.NameSpace, redis.redismeta.NameSpace)
	SSHConfig.Cmd(redis.redismeta.Host, checkhs)
	logger.Info("检查 helmservice %v 是否存在", fmt.Sprintf("redis-%v", redis.redismeta.NameSpace))
	checkhelminit := fmt.Sprintf("helminit || exit 1")
	SSHConfig.Cmd(redis.redismeta.Host, checkhelminit)
}

func (redis RedisMeta) InstallHelm() {
	cmd := fmt.Sprintf("helm install redis-%v bitnami/redis --set cluster.enabled=false --set password=qwerasdf --set metrics.enabled=true --set master.service.type=NodePort -n %v",
		redis.redismeta.NameSpace, redis.redismeta.NameSpace)
	SSHConfig.Cmd(redis.redismeta.Host, cmd)
}
