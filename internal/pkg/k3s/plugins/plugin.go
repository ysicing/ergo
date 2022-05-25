package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAll() ([]Meta, error) {
	var plugins []Meta
	pf := fmt.Sprintf("%s/manifests/plugins/plugins.json", common.GetDefaultDataDir())
	log.Flog.Debug("load local plugin config from", pf)
	content, err := ioutil.ReadFile(pf)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &plugins)
	if err != nil {
		log.Flog.Errorf("unmarshal plugin meta failed: %v", err)
		return nil, err
	}
	return plugins, nil
}

func GetMaps() (map[string]Meta, error) {
	plugins, err := GetAll()
	if err != nil {
		return nil, err
	}
	maps := make(map[string]Meta)
	for _, p := range plugins {
		maps[p.Type] = p
	}
	return maps, nil
}

func GetMeta(args ...string) (Item, error) {
	ps, err := GetMaps()
	if err != nil {
		return Item{}, err
	}
	t := args[0]
	name := ""
	if len(args) == 2 {
		name = args[1]
	} else if strings.Contains(t, "/") {
		ts := strings.Split(t, "/")
		t = ts[0]
		if len(ts) > 1 {
			name = ts[1]
		}
	}
	var plugin Item
	if v, ok := ps[t]; ok {
		if name == "" {
			name = v.Default
		}
		exist := false
		for _, item := range v.Item {
			if item.Name == name {
				exist = true
				plugin = item
				plugin.Type = v.Type
				break
			}
		}
		if !exist {
			log.Flog.Warnf("%s not found %s, will use default: %s", t, name, v.Default)
			return GetMeta(t, v.Default)
		}
		log.Flog.Infof("install %s plugin: %s", t, name)
		return plugin, nil
	}
	return Item{}, fmt.Errorf("plugin %s not found", t)
}

func (p *Item) UnInstall() error {
	pluginName := fmt.Sprintf("%s-%s", common.KubePluginPrefix, p.Type)
	_, err := p.Client.GetSecret(context.TODO(), common.DefaultSystem, pluginName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Flog.Warnf("plugin %s is already uninstalled", p.Type)
			return nil
		}
		log.Flog.Fatalf("get plugin secret failed: %v", err)
		return nil
	}
	// #nosec
	output, err := exec.Command(os.Args[0], "kubectl", "delete", "-f", filepath.Join(common.GetDefaultDataDir(), p.Path), "-n", common.DefaultSystem).CombinedOutput()
	if err != nil {
		log.Flog.Errorf("uninstall plugin %s failed: %s", p.Type, string(output))
		return err
	}
	log.Flog.Donef("uninstall %s plugin done", p.Type)
	p.Client.DeleteSecret(context.TODO(), common.DefaultSystem, pluginName, metav1.DeleteOptions{})
	return nil
}

func (p *Item) Install() error {
	pluginName := fmt.Sprintf("%s-%s", common.KubePluginPrefix, p.Type)
	_, err := p.Client.GetSecret(context.TODO(), common.DefaultSystem, pluginName, metav1.GetOptions{})
	if err == nil {
		log.Flog.Warnf("plugin %s is already installed", p.Type)
		return nil
	}
	if !errors.IsNotFound(err) {
		log.Flog.Debugf("get plugin secret failed: %v", err)
		return fmt.Errorf("plugin %s install failed", p.Name)
	}
	// #nosec
	applycmd := exec.Command(os.Args[0], "kubectl", "apply", "-f", fmt.Sprintf("%s/%s", common.GetDefaultDataDir(), p.Path), "-n", common.DefaultSystem)
	// qcexec.Trace(applycmd)
	output, err := applycmd.CombinedOutput()

	if err != nil {
		log.Flog.Errorf("install %s plugin %s failed: %s", p.Type, p.Name, string(output))
		return err
	}
	log.Flog.Donef("upgrade install %s plugin %s done", p.Type, p.Name)
	plugindata := map[string]string{
		"type":       p.Type,
		"name":       p.Name,
		"version":    p.Version,
		"cliversion": version.Version,
	}
	_, err = p.Client.CreateSecret(context.TODO(), common.DefaultSystem, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: pluginName,
		},
		StringData: plugindata,
	}, metav1.CreateOptions{})
	return err
}
