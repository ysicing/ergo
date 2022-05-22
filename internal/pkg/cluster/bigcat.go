package cluster

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/imroc/req/v3"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/kube"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/version"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Cluster) InstallBigCat() error {
	log.Flog.Info("executing init bigcat logic...")
	ctx := context.Background()
	c, err := kube.NewClient("", "")
	if err != nil {
		return err
	}
	log.Flog.Debug("waiting for storage to be ready...")
	waitsc := time.Now()
	// wait.BackoffUntil TODO
	for {
		sc, _ := c.GetDefaultSC(ctx)
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

	_, err = c.CreateNamespace(ctx, common.DefaultSystem, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Flog.Warnf("namespace %s already exists", common.DefaultSystem)
	}
	log.Flog.Donef("init bigcat done")
	helmbin, err := p.loadLocalBin(common.HelmBinName)
	if err != nil {
		return err
	}
	output, err := exec.Command(helmbin, "repo", "add", "install", common.DefaultChartRepo).CombinedOutput()
	if err != nil {
		errmsg := string(output)
		if !strings.Contains(errmsg, "exists") {
			log.Flog.Errorf("add bigcat install repo failed: %s", string(output))
			return err
		}
		log.Flog.Warnf("bigcat install repo already exists")
	} else {
		log.Flog.Donef("add bigcat install repo done")
	}

	output, err = exec.Command(helmbin, "repo", "update").CombinedOutput()
	if err != nil {
		log.Flog.Errorf("update bigcat install repo failed: %s", string(output))
		return err
	}
	log.Flog.Donef("update bigcat install repo done")
	// // helm upgrade -i nginx-ingress-controller bitnami/nginx-ingress-controller -n kube-system
	// output, err = exec.Command(helmbin, "upgrade", "-i", common.DefaultName, common.DefaultChartName, "-n", common.DefaultSystem).CombinedOutput()
	// if err != nil {
	// 	log.Flog.Errorf("upgrade install bigcat web failed: %s", string(output))
	// 	return err
	// }
	// output, err = exec.Command(helmbin, "upgrade", "-i", common.DefaultCneAPIName, common.DefaultAPIChartName, "-n", common.DefaultSystem).CombinedOutput()
	// if err != nil {
	// 	log.Flog.Errorf("upgrade install bigcat api failed: %s", string(output))
	// 	return err
	// }
	log.Flog.Donef("install bigcat done")
	p.Ready()
	initfile := common.GetCustomConfig(common.InitFileName)
	if err := file.Writefile(initfile, "init done"); err != nil {
		log.Flog.Warnf("write init done file failed, reason: %v.\n\t please run: touch %s", err, initfile)
	}
	return nil
}

func (p *Cluster) Ready() {
	clusterWaitGroup, ctx := errgroup.WithContext(context.Background())
	clusterWaitGroup.Go(func() error {
		return p.ready(ctx)
	})
	if err := clusterWaitGroup.Wait(); err != nil {
		log.Flog.Error(err)
	}
}

func (p *Cluster) ready(ctx context.Context) error {
	t1 := time.Now()
	client := req.C().SetUserAgent(version.GetUG()).SetTimeout(time.Second * 1)
	log.Flog.StartWait("waiting for bigcat ready")
	status := false
	for {
		t2 := time.Now()
		if time.Duration(t2.Sub(t1).Seconds()) > time.Second*180 {
			log.Flog.Warnf("waiting for bigcat ready 3min timeout: check your network or storage. after install you can run: ergo kube status")
			break
		}
		_, err := client.R().Get(fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0]))
		if err == nil {
			status = true
			break
		}
		time.Sleep(time.Second * 10)
	}
	log.Flog.StopWait()
	if status {
		log.Flog.Donef("bigcat ready, cost: %v", time.Since(t1))
	}
	return nil
}
