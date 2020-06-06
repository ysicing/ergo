// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package helm

import (
	"fmt"
	"github.com/wonderivan/logger"
)

type MysqlMeta struct {
	mysqlmeta HelmMeta
}

func (mysql MysqlMeta) CheckHelm() {
	logger.Info("检查ns %v 是否存在", mysql.mysqlmeta.NameSpace)
	checkns := fmt.Sprintf("kubectl create ns %v || sleep 1", mysql.mysqlmeta.NameSpace)
	SSHConfig.Cmd(mysql.mysqlmeta.Host, checkns)
	logger.Info("检查 helmservice %v 是否存在", fmt.Sprintf("mysql-%v", mysql.mysqlmeta.NameSpace))
	checkhs := fmt.Sprintf("helm list -n %v  | grep mysql-%v && exit 1 || sleep 0", mysql.mysqlmeta.NameSpace, mysql.mysqlmeta.NameSpace)
	SSHConfig.Cmd(mysql.mysqlmeta.Host, checkhs)
	logger.Info("检查 helmservice %v 是否存在", fmt.Sprintf("mysql-%v", mysql.mysqlmeta.NameSpace))
	checkhelminit := fmt.Sprintf("helminit || exit 1")
	SSHConfig.Cmd(mysql.mysqlmeta.Host, checkhelminit)
}

func (mysql MysqlMeta) InstallHelm() {
	cmd := fmt.Sprintf("helm install mysql-%v bitnami/mariadb --set master.persistence.enabled=false --set replication.enabled=false --set rootUser.password=qwerasdf --set metrics.enabled=true --set service.type=NodePort -n %v",
		mysql.mysqlmeta.NameSpace, mysql.mysqlmeta.NameSpace)
	SSHConfig.Cmd(mysql.mysqlmeta.Host, cmd)
}
