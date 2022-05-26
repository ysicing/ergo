package cluster

import (
	"context"
	"os/exec"
	"time"

	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	binfile "github.com/ysicing/ergo/pkg/util/util"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Cluster) InstallBigCat() error {
	log.Flog.Info("executing init BigCat logic...")
	ctx := context.Background()
	log.Flog.Debug("waiting for storage to be ready...")
	waitsc := time.Now()
	// wait.BackoffUntil TODO
	for {
		sc, _ := p.client.GetDefaultSC(ctx)
		if sc != nil {
			log.Flog.Donef("default storage %s is ready", sc.Name)
			break
		}
		time.Sleep(time.Second * 5)
		trywaitsc := time.Now()
		if trywaitsc.Sub(waitsc) > time.Minute*3 {
			log.Flog.Warnf("wait storage %s ready, timeout: %v", sc.Name, trywaitsc.Sub(waitsc).Seconds())
			break
		}
	}

	_, err := p.client.CreateNamespace(ctx, common.DefaultSystem, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}
	log.Flog.Done("start init BigCat")
	getbin := binfile.Meta{}
	helmbin, err := getbin.LoadLocalBin(common.HelmBinName)
	if err != nil {
		return err
	}
	// helm upgrade -i nginx-ingress-controller bitnami/nginx-ingress-controller -n kube-system
	output, err := exec.Command(helmbin, "upgrade", "-i", "bigcat", common.DefaultChartName, "-n", common.DefaultSystem).CombinedOutput()
	if err != nil {
		log.Flog.Errorf("upgrade install BigCat web failed: %s", string(output))
		return err
	}
	log.Flog.Donef("install BigCat done")
	p.Ready()
	initfile := common.GetCustomConfig(common.InitFileName)
	if err := file.Writefile(initfile, "init done"); err != nil {
		log.Flog.Warnf("write init done file failed, reason: %v.\n\t please run: touch %s", err, initfile)
	}
	return nil
}
