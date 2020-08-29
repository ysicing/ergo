// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package drone

import (
	api "github.com/drone/drone-go/drone"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"k8s.io/klog"
	"os"
	"sync"
)

var (
	Host  string
	Token string
)

var droneClient api.Client
var droneOnce sync.Once

type DroneService struct {
	Client api.Client
}

func InitDrone() {
	droneOnce.Do(func() {
		config := new(oauth2.Config)
		auth := config.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: viper.GetString("Drone.Token"),
		})
		droneClient = api.NewClient(viper.GetString("Drone.Host"), auth)
	})
}

func NewDrone() *DroneService {
	return &DroneService{Client: droneClient}
}

func (d DroneService) GetClient() {

}

func (d DroneService) Use() {
	user, err := d.Client.Self()
	if err != nil {
		klog.Error("get user err: %v", err.Error())
		os.Exit(-1)
	}
	klog.Info(user)
}
